package arwen

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"unsafe"

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

var (
	vmContextCounter uint8
	vmContextMap     = map[uint8]HostContext{}
	vmContextMapMu   sync.Mutex
)

func addHostContext(ctx HostContext) int {
	vmContextMapMu.Lock()
	id := vmContextCounter
	vmContextCounter++
	vmContextMap[id] = ctx
	vmContextMapMu.Unlock()
	return int(id)
}

func removeHostContext(idx int) {
	vmContextMapMu.Lock()
	delete(vmContextMap, uint8(idx))
	vmContextMapMu.Unlock()
}

func getHostContext(pointer unsafe.Pointer) HostContext {
	var idx = *(*int)(pointer)

	vmContextMapMu.Lock()
	ctx := vmContextMap[uint8(idx)]
	vmContextMapMu.Unlock()

	return ctx
}

type logTopicsData struct {
	topics [][]byte
	data   []byte
}

// vmContext implements evmc.HostContext interface.
type vmContext struct {
	blockChainHook vmcommon.BlockchainHook
	cryptoHook     vmcommon.CryptoHook
	imports        *wasmer.Imports

	vmInput vmcommon.VMInput

	vmType       []byte
	callFunction string
	scAddress    []byte

	logs          map[string]logTopicsData
	storageUpdate map[string](map[string][]byte)

	outputAccounts map[string]*vmcommon.OutputAccount

	output     []byte
	returnCode vmcommon.ReturnCode

	selfDestruct map[string][]byte
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

	instance, err := wasmer.NewInstanceWithImports(input.ContractCode, host.imports)
	if err != nil {
		fmt.Println("arwen Error: ", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	idContext := addHostContext(host)
	instance.SetContextData(unsafe.Pointer(&idContext))

	result := make([]byte, 0)
	init := instance.Exports["init"]
	if init != nil {
		out, err := init()
		if err != nil {
			fmt.Println("arwen Error", err.Error())
			return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature), nil
		}
		if out.GetType() != wasmer.TypeVoid {
			result = []byte(out.String())
		}
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

	removeHostContext(idContext)

	vmOutput := host.createVMOutput(result, gasLeft)

	return vmOutput, err
}

var ErrInitFuncCalledInRun = errors.New("it is not allowed to call init in run")

func (host *vmContext) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	host.initInternalValues()
	host.vmInput = input.VMInput
	host.scAddress = input.RecipientAddr
	host.callFunction = input.Function

	host.addTxValueToSmartContract(input.CallValue, input.RecipientAddr)

	gasLeft := input.GasProvided.Int64()
	contract := host.GetCode(host.scAddress)

	instance, err := wasmer.NewInstanceWithImports(contract, host.imports)
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.ContractInvalid), nil
	}

	defer instance.Close()

	idContext := addHostContext(host)
	instance.SetContextData(unsafe.Pointer(&idContext))

	if host.callFunction == "init" {
		fmt.Println("arwen Error", ErrInitFuncCalledInRun.Error())
		return host.createVMOutputInCaseOfError(vmcommon.UserError), nil
	}

	function := instance.Exports[host.callFunction]
	result, err := function()
	if err != nil {
		fmt.Println("arwen Error", err.Error())
		return host.createVMOutputInCaseOfError(vmcommon.FunctionWrongSignature), nil
	}

	addOutput := make([]byte, 0)
	if result.GetType() != wasmer.TypeVoid {
		addOutput = []byte(result.String())
	}

	vmOutput := host.createVMOutput(addOutput, gasLeft)
	globalDebuggingTrace.PutVMOutput(host.scAddress, vmOutput)
	displayVMOutput(vmOutput)

	return vmOutput, nil
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

	output = append(output, host.output...)
	vmOutput.ReturnData = append(vmOutput.ReturnData, big.NewInt(0).SetBytes(output))
	vmOutput.GasRemaining = big.NewInt(gasLeft)
	vmOutput.GasRefund = big.NewInt(0)
	vmOutput.ReturnCode = host.returnCode

	return vmOutput
}

func displayVMOutput(output *vmcommon.VMOutput) {
	fmt.Println("=============Resulted VM Output=============")
	fmt.Println("RetunCode: ", output.ReturnCode)
	fmt.Println("ReturnData: ", output.ReturnData)
	fmt.Println("GasRemaining: ", output.GasRemaining)
	fmt.Println("GasRefund: ", output.GasRefund)

	for id, touchedAccount := range output.TouchedAccounts {
		fmt.Println("Touched account ", id, ": "+hex.EncodeToString(touchedAccount))
	}

	for id, deletedAccount := range output.DeletedAccounts {
		fmt.Println("Deleted account ", id, ": "+hex.EncodeToString(deletedAccount))
	}

	for id, outputAccount := range output.OutputAccounts {
		fmt.Println("Output account ", id, ": "+hex.EncodeToString(outputAccount.Address))
		if outputAccount.BalanceDelta != nil {
			fmt.Println("           Balance change with : ", outputAccount.BalanceDelta)
		}
		if outputAccount.Nonce != 0 {
			fmt.Println("           Nonce change to : ", outputAccount.Nonce)
		}
		if len(outputAccount.Code) > 0 {
			fmt.Println("           Code change to : ", outputAccount.Code)
		}

		for _, storageUpdate := range outputAccount.StorageUpdates {
			fmt.Println("           Storage update key: "+string(storageUpdate.Offset)+" value: ", big.NewInt(0).SetBytes(storageUpdate.Data))
		}
	}

	for _, log := range output.Logs {
		fmt.Println("Log address: " + hex.EncodeToString(log.Address) + " data: " + string(log.Data))
		fmt.Println("Topics started: ")
		for _, topic := range log.Topics {
			fmt.Print(topic, " ")
		}
		fmt.Println("Topics end")
	}
	fmt.Println("============================================")
}

func (host *vmContext) initInternalValues() {
	host.storageUpdate = make(map[string]map[string][]byte, 0)
	host.logs = make(map[string]logTopicsData, 0)
	host.selfDestruct = make(map[string][]byte)
	host.vmInput = vmcommon.VMInput{}
	host.outputAccounts = make(map[string]*vmcommon.OutputAccount, 0)
	host.output = make([]byte, 0)
	host.scAddress = make([]byte, 0)
	host.callFunction = ""
	host.returnCode = vmcommon.Ok
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

func NewArwenVM(
	blockChainHook vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
	vmType []byte,
) (*vmContext, error) {

	imports, err := ElrondEImports()
	if err != nil {
		return nil, err
	}

	context := &vmContext{
		blockChainHook: blockChainHook,
		cryptoHook:     cryptoHook,
		vmType:         vmType,
		imports:        imports,
	}

	context.initInternalValues()

	return context, nil
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

func (host *vmContext) SelfDestruct(addr []byte, beneficiary []byte) {
	panic("implement me")
}

func (host *vmContext) GetSCAddress() []byte {
	return host.scAddress
}

func (host *vmContext) Finish(data []byte) {
	host.output = append(host.output, data...)
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
		return host.outputAccounts[strAdr].Balance.Bytes()
	}

	balance, err := host.blockChainHook.GetBalance(addr)
	if err != nil {
		fmt.Printf("GetBalance returned with error %s \n", err.Error())
		return nil
	}

	host.outputAccounts[strAdr] = &vmcommon.OutputAccount{Balance: big.NewInt(0).Set(balance), Address: addr}

	return balance.Bytes()
}

func (host *vmContext) GetCodeSize(addr []byte) int {
	code, err := host.blockChainHook.GetCode(addr)
	if err != nil {
		fmt.Printf("GetCodeSize returned with error %s \n", err.Error())
	}

	return len(code)
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

func (host *vmContext) Selfdestruct(addr []byte, beneficiary []byte) {
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

var ErrInvalidTransfer = errors.New("invalid sender")

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
