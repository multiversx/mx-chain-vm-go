package host

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var parentKeyA = []byte("parentKeyA......................")
var parentKeyB = []byte("parentKeyB......................")
var childKey = []byte("childKey........................")
var parentDataA = []byte("parentDataA")
var parentDataB = []byte("parentDataB")
var childData = []byte("childData")
var parentFinishA = []byte("parentFinishA")
var parentFinishB = []byte("parentFinishB")
var childFinish = []byte("childFinish")
var parentTransferReceiver = []byte("parentTransferReceiver..........")
var childTransferReceiver = []byte("childTransferReceiver...........")
var parentTransferValue = int64(42)
var parentTransferData = []byte("parentTransferData")

var gasProvided = uint64(1000000)
var parentCompilationCost = uint64(2282)
var childCompilationCost = uint64(3020)

func expectedVMOutput_Prepare() *vmcommon.VMOutput {
	expectedVMOutput := MakeVMOutput()
	expectedVMOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		expectedVMOutput,
		parentSCAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		expectedVMOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	AddFinishData(expectedVMOutput, parentFinishA)
	AddFinishData(expectedVMOutput, parentFinishB)
	AddFinishData(expectedVMOutput, []byte("succ"))

	expectedExecutionCost := uint64(135)
	gas := gasProvided
	gas -= parentCompilationCost
	gas -= expectedExecutionCost
	expectedVMOutput.GasRemaining = gas

	return expectedVMOutput
}

func expectedVMOutput_WrongContractCalled() *vmcommon.VMOutput {
	expectedVMOutput := expectedVMOutput_Prepare()

	parentAccount := expectedVMOutput.OutputAccounts[string(parentSCAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	_ = AddNewOutputAccount(
		expectedVMOutput,
		[]byte("wrongSC........................."),
		99,
		nil,
	)

	AddFinishData(expectedVMOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(180)
	executeAPICost := uint64(39)
	gasLostOnFailure := uint64(50000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCost
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	expectedVMOutput.GasRemaining = gas

	return expectedVMOutput
}

func expectedVMOutput_SuccessfulChildCall_SameCtx() *vmcommon.VMOutput {
	expectedVMOutput := expectedVMOutput_Prepare()

	parentAccount := expectedVMOutput.OutputAccounts[string(parentSCAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	childAccount := AddNewOutputAccount(
		expectedVMOutput,
		childSCAddress,
		3,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = AddNewOutputAccount(
		expectedVMOutput,
		childTransferReceiver,
		96,
		[]byte("qwerty"),
	)

	// The child SC has stored data on the parent's storage
	SetStorageUpdate(parentAccount, childKey, childData)

	// The called child SC will output some arbitrary data, but also data that it
	// has read from the storage of the parent.
	AddFinishData(expectedVMOutput, childFinish)
	AddFinishData(expectedVMOutput, parentDataA)
	for _, c := range parentDataA {
		AddFinishData(expectedVMOutput, []byte{c})
	}
	AddFinishData(expectedVMOutput, parentDataB)
	for _, c := range parentDataB {
		AddFinishData(expectedVMOutput, []byte{c})
	}
	AddFinishData(expectedVMOutput, []byte("child ok"))
	AddFinishData(expectedVMOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(188)
	executeAPICost := uint64(39)
	childExecutionCost := uint64(441)
	finalCost := uint64(36)
	gas := gasProvided
	gas -= parentCompilationCost
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost
	gas -= childExecutionCost
	gas -= finalCost
	expectedVMOutput.GasRemaining = gas

	return expectedVMOutput
}

func expectedVMOutput_SuccessfulChildCall_BigInts_SameCtx() *vmcommon.VMOutput {
	expectedVMOutput := MakeVMOutput()
	expectedVMOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		expectedVMOutput,
		parentSCAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.BalanceDelta = big.NewInt(-99)

	_ = AddNewOutputAccount(
		expectedVMOutput,
		childSCAddress,
		99,
		nil,
	)

	// The child SC will output "child ok" if it could read some expected Big
	// Ints directly from the parent's context.
	AddFinishData(expectedVMOutput, []byte("child ok"))
	AddFinishData(expectedVMOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(143)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(90)
	finalCost := uint64(36)
	gas := gasProvided
	gas -= parentCompilationCost
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost
	gas -= childExecutionCost
	gas -= finalCost
	expectedVMOutput.GasRemaining = gas
	return expectedVMOutput
}

func expectedVMOutput_SuccessfulChildCall_DestCtx() *vmcommon.VMOutput {
	expectedVMOutput := expectedVMOutput_Prepare()

	parentAccount := expectedVMOutput.OutputAccounts[string(parentSCAddress)]
	parentAccount.BalanceDelta = big.NewInt(-75)

	childAccount := AddNewOutputAccount(
		expectedVMOutput,
		childSCAddress,
		33,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = AddNewOutputAccount(
		expectedVMOutput,
		childTransferReceiver,
		12,
		[]byte("Second sentence."),
	)

	SetStorageUpdate(childAccount, childKey, childData)

	AddFinishData(expectedVMOutput, childFinish)
	AddFinishData(expectedVMOutput, []byte("succ"))

	return expectedVMOutput
}
