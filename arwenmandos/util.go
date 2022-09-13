package arwenmandos

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm/mandos-go/esdtconvert"
	er "github.com/ElrondNetwork/wasm-vm/mandos-go/expression/reconstructor"
	mj "github.com/ElrondNetwork/wasm-vm/mandos-go/model"
	worldmock "github.com/ElrondNetwork/wasm-vm/mock/world"
)

func convertAccount(testAcct *mj.Account, world *worldmock.MockWorld) (*worldmock.Account, error) {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}

	err := esdtconvert.WriteMandosESDTToStorage(testAcct.ESDTData, storage)
	if err != nil {
		return nil, err
	}

	if len(testAcct.Address.Value) != 32 {
		return nil, errors.New("bad test: account address should be 32 bytes long")
	}

	account := &worldmock.Account{
		Address:         testAcct.Address.Value,
		Nonce:           testAcct.Nonce.Value,
		Balance:         big.NewInt(0).Set(testAcct.Balance.Value),
		BalanceDelta:    big.NewInt(0),
		DeveloperReward: big.NewInt(0),
		Username:        testAcct.Username.Value,
		Storage:         storage,
		Code:            testAcct.Code.Value,
		OwnerAddress:    testAcct.Owner.Value,
		AsyncCallData:   testAcct.AsyncCallData,
		ShardID:         uint32(testAcct.Shard.Value),
		IsSmartContract: len(testAcct.Code.Value) > 0,
		CodeMetadata: (&vmcommon.CodeMetadata{
			Payable:     true,
			Upgradeable: true,
			Readable:    true,
		}).ToBytes(), // TODO: add explicit fields in mandos json
		MockWorld: world,
	}

	return account, nil
}

func validateSetStateAccount(mandosAccount *mj.Account, converted *worldmock.Account) error {
	err := converted.Validate()
	if err != nil {
		return fmt.Errorf(
			`"setState" step validation failed for account "%s": %w`,
			mandosAccount.Address.Original,
			err)
	}
	return nil
}

func validateNewAddressMocks(testNAMs []*mj.NewAddressMock) error {
	for _, testNAM := range testNAMs {
		if !worldmock.IsSmartContractAddress(testNAM.NewAddress.Value) {
			return fmt.Errorf(
				`address in "setState" "newAddresses" field should have SC format: %s`,
				testNAM.NewAddress.Original)
		}
	}
	return nil
}

func convertNewAddressMocks(testNAMs []*mj.NewAddressMock) []*worldmock.NewAddressMock {
	var result []*worldmock.NewAddressMock
	for _, testNAM := range testNAMs {
		result = append(result, &worldmock.NewAddressMock{
			CreatorAddress: testNAM.CreatorAddress.Value,
			CreatorNonce:   testNAM.CreatorNonce.Value,
			NewAddress:     testNAM.NewAddress.Value,
		})
	}
	return result
}

func convertBlockInfo(testBlockInfo *mj.BlockInfo, currentInfo *worldmock.BlockInfo) *worldmock.BlockInfo {
	if testBlockInfo == nil {
		return currentInfo
	}

	if currentInfo == nil {
		currentInfo = &worldmock.BlockInfo{
			BlockTimestamp: 0,
			BlockNonce:     0,
			BlockRound:     0,
			BlockEpoch:     0,
			RandomSeed:     nil,
		}
	}

	if !testBlockInfo.BlockTimestamp.OriginalEmpty() {
		currentInfo.BlockTimestamp = testBlockInfo.BlockTimestamp.Value

	}

	if !testBlockInfo.BlockNonce.OriginalEmpty() {
		currentInfo.BlockNonce = testBlockInfo.BlockNonce.Value
	}

	if !testBlockInfo.BlockRound.OriginalEmpty() {
		currentInfo.BlockRound = testBlockInfo.BlockRound.Value
	}

	if !testBlockInfo.BlockEpoch.OriginalEmpty() {
		currentInfo.BlockEpoch = uint32(testBlockInfo.BlockEpoch.Value)
	}

	if testBlockInfo.BlockRandomSeed != nil && !testBlockInfo.BlockRandomSeed.OriginalEmpty() {
		var randomsSeed [48]byte
		copy(randomsSeed[:], testBlockInfo.BlockRandomSeed.Value)
		currentInfo.RandomSeed = &randomsSeed

	}

	return currentInfo
}

// this is a small hack, so we can reuse mandos's JSON printing in error messages
func (ae *ArwenTestExecutor) convertLogToTestFormat(outputLog *vmcommon.LogEntry) *mj.LogEntry {
	topics := mj.JSONCheckValueList{
		Values: make([]mj.JSONCheckBytes, len(outputLog.Topics)),
	}
	for i, topic := range outputLog.Topics {
		topics.Values[i] = mj.JSONCheckBytesReconstructed(
			topic,
			ae.exprReconstructor.Reconstruct(topic,
				er.NoHint))
	}
	testLog := mj.LogEntry{
		Address: mj.JSONCheckBytesReconstructed(
			outputLog.Address,
			ae.exprReconstructor.Reconstruct(outputLog.Address,
				er.AddressHint)),
		Endpoint: mj.JSONCheckBytesReconstructed(
			outputLog.Identifier,
			ae.exprReconstructor.Reconstruct(outputLog.Identifier,
				er.StrHint)),
		Data:   mj.JSONCheckBytesReconstructed(outputLog.Data, ""),
		Topics: topics,
	}

	return &testLog
}

func generateTxHash(txIndex string) []byte {
	txIndexBytes := []byte(txIndex)
	if len(txIndexBytes) > 32 {
		return txIndexBytes[:32]
	}
	for i := len(txIndexBytes); i < 32; i++ {
		txIndexBytes = append(txIndexBytes, '.')
	}
	return txIndexBytes
}

func addESDTToVMInput(esdtData []*mj.ESDTTxData, vmInput *vmcommon.VMInput) {
	esdtDataLen := len(esdtData)

	if esdtDataLen > 0 {
		vmInput.ESDTTransfers = make([]*vmcommon.ESDTTransfer, esdtDataLen)
		for i := 0; i < esdtDataLen; i++ {
			vmInput.ESDTTransfers[i] = &vmcommon.ESDTTransfer{}
			vmInput.ESDTTransfers[i].ESDTTokenName = esdtData[i].TokenIdentifier.Value
			vmInput.ESDTTransfers[i].ESDTValue = esdtData[i].Value.Value
			vmInput.ESDTTransfers[i].ESDTTokenNonce = esdtData[i].Nonce.Value
			if vmInput.ESDTTransfers[i].ESDTTokenNonce != 0 {
				vmInput.ESDTTransfers[i].ESDTTokenType = uint32(core.NonFungible)
			} else {
				vmInput.ESDTTransfers[i].ESDTTokenType = uint32(core.Fungible)
			}
		}
	}
}

func logGasTrace(ae *ArwenTestExecutor) {
	if ae.PeekTraceGas() {
		metering := ae.GetVMHost().Metering()
		scGasTrace := metering.GetGasTrace()
		totalGasUsedByAPIs := 0
		for scAddress, gasTrace := range scGasTrace {
			fmt.Println("Gas Trace for: ", "SC Address", scAddress)
			for functionName, value := range gasTrace {
				totalGasUsed := uint64(0)
				for _, usedGas := range value {
					totalGasUsed += usedGas
				}
				fmt.Println("GasTrace: functionName:", functionName, ",  totalGasUsed:", totalGasUsed, ", numberOfCalls:", len(value))
				totalGasUsedByAPIs += int(totalGasUsed)
			}
			fmt.Println("TotalGasUsedByAPIs: ", totalGasUsedByAPIs)
		}
	}
}

func setGasTraceInMetering(ae *ArwenTestExecutor, enable bool) {
	metering := ae.GetVMHost().Metering()
	if enable && ae.PeekTraceGas() {
		metering.SetGasTracing(true)
	} else {
		metering.SetGasTracing(false)
	}
}

func setExternalStepGasTracing(ae *ArwenTestExecutor, step *mj.ExternalStepsStep) {
	switch step.TraceGas.ToInt() {
	case mj.Undefined.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, ae.PeekTraceGas())
	case mj.TrueValue.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, true)
	case mj.FalseValue.ToInt():
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, false)
	}
}

func resetGasTracesIfNewTest(ae *ArwenTestExecutor, scenario *mj.Scenario) {
	if ae.vm == nil || scenario.IsNewTest {
		ae.scenarioTraceGas = make([]bool, 0)
		ae.scenarioTraceGas = append(ae.scenarioTraceGas, scenario.TraceGas)
		scenario.IsNewTest = false
	}
}
