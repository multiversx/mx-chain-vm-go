package context

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/crypto"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen/ethapi"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

type StorageStatus int

const (
	StorageUnchanged StorageStatus = 0
	StorageModified  StorageStatus = 1
	StorageAdded     StorageStatus = 3
	StorageDeleted   StorageStatus = 4
)

type logTopicsData struct {
	topics [][]byte
	data   []byte
}

// vmContext implements HostContext interface.
type vmContext struct {
	BigIntContainer
	blockChainHook vmcommon.BlockchainHook
	cryptoHook     vmcommon.CryptoHook
	imports        *wasmer.Imports
	instance       wasmer.Instance

	vmInput vmcommon.VMInput

	vmType       []byte
	callFunction string
	scAddress    []byte

	logs           map[string]logTopicsData
	storageUpdate  map[string](map[string][]byte)
	outputAccounts map[string]*vmcommon.OutputAccount
	returnData     []*big.Int
	returnCode     vmcommon.ReturnCode

	selfDestruct  map[string][]byte
	ethInput      []byte
	blockGasLimit uint64

	gasCostConfig *config.GasCost
}

func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
	vmType []byte,
	blockGasLimit uint64,
	gasSchedule map[string]uint64,
) (*vmContext, error) {

	imports, err := elrondapi.ElrondEImports()
	if err != nil {
		return nil, err
	}

	imports, err = elrondapi.BigIntImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = ethapi.EthereumImports(imports)
	if err != nil {
		return nil, err
	}

	imports, err = crypto.CryptoImports(imports)
	if err != nil {
		return nil, err
	}

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	if err != nil {
		return nil, err
	}

	context := &vmContext{
		BigIntContainer: NewBigIntContainer(),
		blockChainHook:  blockChainHook,
		cryptoHook:      cryptoHook,
		vmType:          vmType,
		imports:         imports,
		blockGasLimit:   blockGasLimit,
		gasCostConfig:   gasCostConfig,
	}

	context.initInternalValues()

	return context, nil
}

func (host *vmContext) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	host.initInternalValues()
	host.vmInput = input.VMInput

	nonce, err := host.blockChainHook.GetNonce(input.CallerAddr)
	if err != nil {
		return nil, err
	}

	if nonce > 0 {
		nonce -= 1
	}

	address, err := host.blockChainHook.NewAddress(input.CallerAddr, nonce, host.vmType)
	if err != nil {
		return nil, err
	}

	host.scAddress = address
	host.addTxValueToSmartContract(input.CallValue, address)

	host.instance, err = wasmer.NewInstanceWithImports(input.ContractCode, host.imports)
	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	defer host.instance.Close()

	idContext := arwen.AddHostContext(host)
	host.instance.SetContextData(unsafe.Pointer(&idContext))

	var result []byte
	init := host.instance.Exports["init"]
	if init != nil {
		out, err := init()
		if err != nil {
			fmt.Println("arwen Error", err.Error())
			return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature), nil
		}
		convertedResult := arwen.ConvertReturnValue(out)
		result = convertedResult.Bytes()
	}

	gasLeft := input.GasProvided.Int64()

	// take out contract creation gas
	gasLeft = gasLeft - int64(len(input.ContractCode))

	newSCAcc, ok := host.outputAccounts[string(address)]
	if !ok {
		host.outputAccounts[string(address)] = &vmcommon.OutputAccount{
			Address:        address,
			Nonce:          0,
			BalanceDelta:   big.NewInt(0),
			StorageUpdates: nil,
			Code:           input.ContractCode,
		}
	} else {
		newSCAcc.Code = input.ContractCode
	}

	arwen.RemoveHostContext(idContext)

	vmOutput := host.createVMOutput(result, gasLeft)

	return vmOutput, err
}

func (host *vmContext) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.initInternalValues()
	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	host.addTxValueToSmartContract(input.CallValue, input.RecipientAddr)

	gasLeft := input.GasProvided.Int64()
	contract := host.GetCode(host.scAddress)

	var err error
	opcode_costs := opcode_costs_uniform_value(2)
	host.instance, err = wasmer.NewMeteredInstanceWithImports(contract, host.imports, uint64(gasLeft), opcode_costs)
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	defer host.instance.Close()

	idContext := arwen.AddHostContext(host)
	host.instance.SetContextData(unsafe.Pointer(&idContext))

	if host.callFunction == "init" {
		fmt.Println("arwen Error", ErrInitFuncCalledInRun.Error())
		return host.createVMOutputInCaseOfError(vmcommon.UserError), nil
	}

	function, ok := host.instance.Exports[host.callFunction]
	if !ok {
		fmt.Println("arwen Error", "Function not found")
		return host.createVMOutputInCaseOfError(vmcommon.FunctionNotFound), nil
	}

	result, err := function()
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature), nil
	}

	if host.returnCode != vmcommon.Ok {
		// user error: signalError()
		return host.createVMOutputInCaseOfError(host.returnCode), nil
	}

	convertedResult := arwen.ConvertReturnValue(result)
	gasLeft = gasLeft - int64(host.instance.GetPointsUsed())
	vmOutput := host.createVMOutput(convertedResult.Bytes(), gasLeft)

	return vmOutput, nil
}

func opcode_costs_uniform_value(value uint32) *[wasmer.OPCODE_COUNT]uint32 {
	opcode_costs := [wasmer.OPCODE_COUNT]uint32{}
	for i := 0; i < wasmer.OPCODE_COUNT; i++ {
		opcode_costs[i] = value
	}
	return &opcode_costs
}

func (host *vmContext) createVMOutputInCaseOfError(errCode vmcommon.ReturnCode) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{GasRemaining: big.NewInt(0), GasRefund: big.NewInt(0)}
	vmOutput.ReturnCode = errCode
	return vmOutput
}

// adapt vm output and all saved data from sc run into VM Output
func (host *vmContext) createVMOutput(output []byte, gasLeft int64) *vmcommon.VMOutput {
	vmOutput := &vmcommon.VMOutput{}
	// save storage updates
	outAccs := make(map[string]*vmcommon.OutputAccount, 0)
	for addr, updates := range host.storageUpdate {
		if _, ok := outAccs[addr]; !ok {
			outAccs[addr] = &vmcommon.OutputAccount{Address: []byte(addr)}
		}

		for key, value := range updates {
			storageUpdate := &vmcommon.StorageUpdate{
				Offset: []byte(key),
				Data:   value,
			}

			outAccs[addr].StorageUpdates = append(outAccs[addr].StorageUpdates, storageUpdate)
		}
	}

	// add balances
	for addr, outAcc := range host.outputAccounts {
		if _, ok := outAccs[addr]; !ok {
			outAccs[addr] = &vmcommon.OutputAccount{}
		}

		outAccs[addr].Address = outAcc.Address
		outAccs[addr].BalanceDelta = outAcc.BalanceDelta

		if len(outAcc.Code) > 0 {
			outAccs[addr].Code = outAcc.Code
		}
		if outAcc.Nonce > 0 {
			outAccs[addr].Nonce = outAcc.Nonce
		}
	}

	// save to the output finally
	for _, outAcc := range outAccs {
		vmOutput.OutputAccounts = append(vmOutput.OutputAccounts, outAcc)
	}

	// save logs
	for addr, value := range host.logs {
		logEntry := &vmcommon.LogEntry{
			Address: []byte(addr),
			Data:    value.data,
		}

		topics := make([]*big.Int, len(value.topics))
		for i := 0; i < len(value.topics); i++ {
			topics[i] = big.NewInt(0).SetBytes(value.topics[i])
		}

		logEntry.Topics = topics
		vmOutput.Logs = append(vmOutput.Logs, logEntry)
	}

	if len(host.returnData) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, host.returnData...)
	}
	if len(output) > 0 {
		vmOutput.ReturnData = append(vmOutput.ReturnData, big.NewInt(0).SetBytes(output))
	}

	vmOutput.GasRemaining = big.NewInt(gasLeft)
	vmOutput.GasRefund = big.NewInt(0)
	vmOutput.ReturnCode = host.returnCode

	return vmOutput
}

func (host *vmContext) initInternalValues() {
	host.Clean()
	host.storageUpdate = make(map[string]map[string][]byte, 0)
	host.logs = make(map[string]logTopicsData, 0)
	host.selfDestruct = make(map[string][]byte)
	host.vmInput = vmcommon.VMInput{}
	host.outputAccounts = make(map[string]*vmcommon.OutputAccount, 0)
	host.scAddress = make([]byte, 0)
	host.callFunction = ""
	host.returnData = nil
	host.returnCode = vmcommon.Ok
	host.ethInput = nil
}

func (host *vmContext) addTxValueToSmartContract(value *big.Int, scAddress []byte) {
	destAcc, ok := host.outputAccounts[string(scAddress)]
	if !ok {
		destAcc = &vmcommon.OutputAccount{
			Address:      scAddress,
			BalanceDelta: big.NewInt(0),
		}
		host.outputAccounts[string(destAcc.Address)] = destAcc
	}

	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)
}

func (host *vmContext) GasSchedule() *config.GasCost {
	return host.gasCostConfig
}

func (host *vmContext) EthContext() arwen.EthContext {
	return host
}

func (host *vmContext) CoreContext() arwen.HostContext {
	return host
}

func (host *vmContext) BigInContext() arwen.BigIntContext {
	return host
}

func (host *vmContext) CryptoContext() arwen.CryptoContext {
	return host
}

func (host *vmContext) CryptoHooks() vmcommon.CryptoHook {
	return host.cryptoHook
}

func (host *vmContext) Finish(data []byte) {
	host.returnData = append(host.returnData, big.NewInt(0).SetBytes(data))
}

func (host *vmContext) SignalUserError() {
	host.returnCode = vmcommon.UserError
}

func (host *vmContext) Arguments() []*big.Int {
	return host.vmInput.Arguments
}

func (host *vmContext) Function() string {
	return host.callFunction
}

func (host *vmContext) GetSCAddress() []byte {
	return host.scAddress
}

func (host *vmContext) AccountExists(addr []byte) bool {
	exists, err := host.blockChainHook.AccountExists(addr)
	if err != nil {
		fmt.Printf("Account exsits returned with error %s \n", err.Error())
	}
	return exists
}

func (host *vmContext) GetStorage(addr []byte, key []byte) []byte {
	strAdr := string(addr)
	if _, ok := host.storageUpdate[strAdr]; ok {
		if value, ok := host.storageUpdate[strAdr][string(key)]; ok {
			return value
		}
	}

	hash, err := host.blockChainHook.GetStorageData(addr, key)
	if err != nil {
		fmt.Printf("GetStorage returned with error %s \n", err.Error())
	}

	return hash
}

func (host *vmContext) SetStorage(addr []byte, key []byte, value []byte) int32 {
	strAdr := string(addr)

	if _, ok := host.storageUpdate[strAdr]; !ok {
		host.storageUpdate[strAdr] = make(map[string][]byte, 0)
	}

	oldValue := host.storageUpdate[strAdr][string(key)]
	length := len(value)
	host.storageUpdate[strAdr][string(key)] = make([]byte, length)
	copy(host.storageUpdate[strAdr][string(key)][:length], value[:length])

	if bytes.Equal(oldValue, value) {
		return int32(StorageUnchanged)
	}

	zero := []byte{}
	if bytes.Equal(oldValue, zero) {
		return int32(StorageAdded)
	}
	if bytes.Equal(value, zero) {
		return int32(StorageDeleted)
	}

	return int32(StorageModified)
}

func (host *vmContext) getBalanceFromBlockChain(addr []byte) *big.Int {
	balance, err := host.blockChainHook.GetBalance(addr)

	if err != nil {
		fmt.Printf("GetBalance returned with error %s \n", err.Error())
		return big.NewInt(0)
	}

	return balance
}

func (host *vmContext) GetBalance(addr []byte) []byte {
	strAdr := string(addr)
	if _, ok := host.outputAccounts[strAdr]; ok {
		balance := host.outputAccounts[strAdr].Balance
		return balance.Bytes()
	}

	balance, err := host.blockChainHook.GetBalance(addr)
	if err != nil {
		fmt.Printf("GetBalance returned with error %s \n", err.Error())
		return big.NewInt(0).Bytes()
	}

	host.outputAccounts[strAdr] = &vmcommon.OutputAccount{Balance: big.NewInt(0).Set(balance), Address: addr}
	return balance.Bytes()
}

func (host *vmContext) GetCodeSize(addr []byte) int32 {
	code, err := host.blockChainHook.GetCode(addr)
	if err != nil {
		fmt.Printf("GetCodeSize returned with error %s \n", err.Error())
	}

	return int32(len(code))
}

func (host *vmContext) GetCodeHash(addr []byte) []byte {
	code, err := host.blockChainHook.GetCode(addr)
	if err != nil {
		fmt.Printf("GetCodeSize returned with error %s \n", err.Error())
	}

	codeHash, err := host.cryptoHook.Keccak256(string(code))
	if err != nil {
		fmt.Printf("GetCodeSize returned with error %s \n", err.Error())
	}

	return []byte(codeHash)
}

func (host *vmContext) GetCode(addr []byte) []byte {
	code, err := host.blockChainHook.GetCode(addr)
	if err != nil {
		fmt.Printf("GetCodeSize returned with error %s \n", err.Error())
	}

	return code
}

func (host *vmContext) SelfDestruct(addr []byte, beneficiary []byte) {
	host.selfDestruct[string(addr)] = beneficiary
}

func (host *vmContext) GetVMInput() vmcommon.VMInput {
	return host.vmInput
}

func (host *vmContext) BlockHash(number int64) []byte {
	block, err := host.blockChainHook.GetBlockhash(big.NewInt(number))

	if err != nil {
		fmt.Printf("GetBlockHash returned with error %s \n", err.Error())
		return nil
	}

	return block
}

func (host *vmContext) WriteLog(addr []byte, topics [][]byte, data []byte) {
	if len(topics) != len(data) {
		return
	}

	strAdr := string(addr)

	if _, ok := host.logs[strAdr]; !ok {
		host.logs[strAdr] = logTopicsData{
			topics: make([][]byte, 0),
			data:   make([]byte, 0),
		}
	}

	currLogs := host.logs[strAdr]
	for i := 0; i < len(topics); i++ {
		currLogs.topics = append(currLogs.topics, topics[i])
	}
	currLogs.data = append(currLogs.data, data...)

	host.logs[strAdr] = currLogs
}

// Transfer handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (host *vmContext) Transfer(destination []byte, sender []byte, value *big.Int, input []byte, gas int64,
) (gasLeft int64, err error) {
	//TODO: should this be kept, or there are other use cases where a sender can be somebody else
	if !bytes.Equal(sender, host.GetSCAddress()) {
		return 0, ErrInvalidTransfer
	}

	senderAcc, ok := host.outputAccounts[string(sender)]
	if !ok {
		senderAcc = &vmcommon.OutputAccount{
			Address:      sender,
			BalanceDelta: big.NewInt(0),
		}
		host.outputAccounts[string(senderAcc.Address)] = senderAcc
	}

	destAcc, ok := host.outputAccounts[string(destination)]
	if !ok {
		destAcc = &vmcommon.OutputAccount{
			Address:      destination,
			BalanceDelta: big.NewInt(0),
		}
		host.outputAccounts[string(destAcc.Address)] = destAcc
	}

	senderAcc.BalanceDelta = big.NewInt(0).Sub(senderAcc.BalanceDelta, value)
	destAcc.BalanceDelta = big.NewInt(0).Add(destAcc.BalanceDelta, value)

	//TODO: support this - smart contract which call other smart contracts - not supported yet
	// if destination.HasCode ...

	gasLeft = gas - 1

	return gasLeft, err
}

func (host *vmContext) CallData() []byte {
	if host.ethInput == nil {
		host.ethInput = host.createETHCallInput()
	}
	return host.ethInput
}

func (host *vmContext) UseGas(gas uint64) {
	currGas := host.instance.GetPointsUsed()
	if currGas > gas {
		currGas = currGas - gas
	} else {
		currGas = 0
	}

	host.instance.SetPointsUsed(currGas)
}

func (host *vmContext) GasLeft() uint64 {
	return host.instance.GetPointsUsed()
}

func (host *vmContext) BlockGasLimit() uint64 {
	return host.blockGasLimit
}

func (host *vmContext) BlockChainHook() vmcommon.BlockchainHook {
	return host.blockChainHook
}

// The first four bytes is the method selector. The rest of the input data are method arguments in chunks of 32 bytes.
// The method selector is the kecccak256 hash of the method signature.
func (host *vmContext) createETHCallInput() []byte {
	newInput := make([]byte, 0)

	if len(host.callFunction) > 0 {
		hashOfFunction, err := host.cryptoHook.Keccak256(host.callFunction)
		if err != nil {
			return nil
		}

		methodSelectors, err := hex.DecodeString(hashOfFunction)
		if err != nil {
			return nil
		}

		newInput = append(newInput, methodSelectors[0:4]...)
	}

	for _, arg := range host.vmInput.Arguments {
		currInput := make([]byte, arwen.HashLen)
		copy(currInput[arwen.HashLen-len(arg.Bytes()):], arg.Bytes())

		newInput = append(newInput, currInput...)
	}

	return newInput
}
