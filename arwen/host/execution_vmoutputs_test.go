package host

import (
	"fmt"
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
var parentTransferReceiver = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fparentTransferReceiver")
var childTransferReceiver = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fchildTransferReceiver.")
var parentTransferValue = int64(42)
var parentTransferData = []byte("parentTransferData")

var recursiveIterationCounterKey = []byte("recursiveIterationCounter.......")
var recursiveIterationBigCounterKey = []byte("recursiveIterationBigCounter....")

var gasProvided = uint64(1000000)

var parentCompilationCost_SameCtx uint64
var childCompilationCost_SameCtx uint64

var parentCompilationCost_DestCtx uint64
var childCompilationCost_DestCtx uint64

var vaultAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fvaultAddress..........")
var thirdPartyAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fthirdPartyAddress.....")

func init() {
	parentCompilationCost_SameCtx = uint64(len(GetTestSCCode("exec-same-ctx-parent", "../../")))
	childCompilationCost_SameCtx = uint64(len(GetTestSCCode("exec-same-ctx-child", "../../")))

	parentCompilationCost_DestCtx = uint64(len(GetTestSCCode("exec-dest-ctx-parent", "../../")))
	childCompilationCost_DestCtx = uint64(len(GetTestSCCode("exec-dest-ctx-child", "../../")))
}

func expectedVMOutput_SameCtx_Prepare(_ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)
	AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(137)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_WrongContractCalled(code []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare(code)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(39)
	gasLostOnFailure := uint64(50000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_OutOfGas(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	AddFinishData(vmOutput, parentFinishA)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_SuccessfulChildCall(parentCode []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		3,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	executeAPICost := uint64(39)
	childExecutionCost := uint64(436)
	childAccount.GasUsed = childCompilationCost_SameCtx + childExecutionCost

	_ = AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		96,
		[]byte("qwerty"),
	)

	// The child SC has stored data on the parent's storage
	SetStorageUpdate(parentAccount, childKey, childData)

	// The called child SC will output some arbitrary data, but also data that it
	// has read from the storage of the parent.
	AddFinishData(vmOutput, childFinish)
	AddFinishData(vmOutput, parentDataA)
	for _, c := range parentDataA {
		AddFinishData(vmOutput, []byte{c})
	}
	AddFinishData(vmOutput, parentDataB)
	for _, c := range parentDataB {
		AddFinishData(vmOutput, []byte{c})
	}
	AddFinishData(vmOutput, []byte("child ok"))
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(171)
	finalCost := uint64(134)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_SameCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_SuccessfulChildCall_BigInts(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	// parentAccount.BalanceDelta = big.NewInt(-99)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)
	childExecutionCost := uint64(107)
	childAccount.GasUsed = childCompilationCost_SameCtx + childExecutionCost

	// The child SC will output "child ok" if it could read some expected Big
	// Ints directly from the parent's context.
	AddFinishData(vmOutput, []byte("child ok"))
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(113)
	executeAPICost := uint64(13)
	finalCost := uint64(67)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_SameCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_Direct(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 21187

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		SetStorageUpdateStrings(account, key, value)
	}

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_Direct_ErrMaxInstances(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	finishString := fmt.Sprintf("Rfinish%03d", recursiveCalls)
	AddFinishData(vmOutput, []byte(finishString))

	key := fmt.Sprintf("Rkey%03d.........................", recursiveCalls)
	value := fmt.Sprintf("Rvalue%03d", recursiveCalls)
	SetStorageUpdateStrings(account, key, value)

	AddFinishData(vmOutput, []byte("fail"))
	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(1)})

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_MutualMethods(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 25412

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	AddFinishData(vmOutput, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		SetStorageUpdateStrings(account, key, value)
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_MutualSCs(_ []byte, _ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3650

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	childAccount.GasUsed = 5437

	if recursiveCalls%2 == 1 {
		parentAccount.BalanceDelta = big.NewInt(-5)
		childAccount.BalanceDelta = big.NewInt(5)
	} else {
		parentAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
		childAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	}

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
		}
		SetStorageUpdateStrings(parentAccount, key, value)
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutput_SameCtx_BuiltinFunctions_1(_ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		42,
		nil,
	)
	account.Balance = big.NewInt(1000)

	AddFinishData(vmOutput, []byte("succ"))
	gasConsumed_builtinClaim := 100
	vmOutput.GasRemaining = uint64(98504 - gasConsumed_builtinClaim)

	return vmOutput
}

func expectedVMOutput_SameCtx_BuiltinFunctions_2(_ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0)

	AddFinishData(vmOutput, []byte("succ"))

	gasConsumed_builtinDoSomething := 0
	vmOutput.GasRemaining = uint64(98500 - gasConsumed_builtinDoSomething)

	return vmOutput
}

func expectedVMOutput_SameCtx_BuiltinFunctions_3(_ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	_ = AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	AddFinishData(vmOutput, []byte("fail"))
	vmOutput.GasRemaining = 98000

	return vmOutput
}

func expectedVMOutput_DestCtx_Prepare(_ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)
	AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(137)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_DestCtx_WrongContractCalled(parentCode []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-42)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(156)
	executeAPICost := uint64(42)
	gasLostOnFailure := uint64(10000)
	finalCost := uint64(44)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_DestCtx_OutOfGas(_ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	AddFinishData(vmOutput, parentFinishA)

	AddFinishData(vmOutput, []byte("fail"))

	executionCostBeforeExecuteAPI := uint64(90)
	executeAPICost := uint64(1)
	gasLostOnFailure := uint64(3500)
	finalCost := uint64(54)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= executionCostBeforeExecuteAPI
	gas -= executeAPICost
	gas -= gasLostOnFailure
	gas -= finalCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_DestCtx_SuccessfulChildCall(parentCode []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		99-12,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		12,
		[]byte("Second sentence."),
	)

	SetStorageUpdate(childAccount, childKey, childData)

	AddFinishData(vmOutput, childFinish)
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(166)
	executeAPICost := uint64(42)
	childExecutionCost := uint64(91)
	finalCost := uint64(65)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_DestCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutput_DestCtx_SuccessfulChildCall_BigInts(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)

	// The child SC will output "child ok" if it could NOT read the Big Ints from
	// the parent's context.
	AddFinishData(vmOutput, []byte("child ok"))
	AddFinishData(vmOutput, []byte("succ"))
	AddFinishData(vmOutput, []byte("succ"))

	parentGasBeforeExecuteAPI := uint64(113)
	executeAPICost := uint64(13)
	childExecutionCost := uint64(101)
	finalCost := uint64(68)
	gas := gasProvided
	gas -= parentCompilationCost_DestCtx
	gas -= parentGasBeforeExecuteAPI
	gas -= executeAPICost
	gas -= childCompilationCost_DestCtx
	gas -= childExecutionCost
	gas -= finalCost
	vmOutput.GasRemaining = gas
	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_Direct(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		SetStorageUpdateStrings(account, key, value)
	}

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_MutualMethods(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	account := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	AddFinishData(vmOutput, []byte("start recursive mutual calls"))

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Afinish%03d", iteration)
			key = fmt.Sprintf("Akey%03d.........................", iteration)
			value = fmt.Sprintf("Avalue%03d", iteration)
		} else {
			finishData = fmt.Sprintf("Bfinish%03d", iteration)
			key = fmt.Sprintf("Bkey%03d.........................", iteration)
			value = fmt.Sprintf("Bvalue%03d", iteration)
		}
		SetStorageUpdateStrings(account, key, value)
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_MutualSCs(_ []byte, _ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentIterations := (recursiveCalls / 2) + (recursiveCalls % 2)
	childIterations := recursiveCalls - parentIterations

	balanceDelta := int64(5*parentIterations - 3*childIterations)

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.BalanceDelta = big.NewInt(-balanceDelta)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	childAccount.BalanceDelta = big.NewInt(balanceDelta)

	for i := 0; i <= recursiveCalls; i++ {
		var finishData string
		var key string
		var value string
		iteration := recursiveCalls - i
		if i%2 == 0 {
			finishData = fmt.Sprintf("Pfinish%03d", iteration)
			key = fmt.Sprintf("Pkey%03d.........................", iteration)
			value = fmt.Sprintf("Pvalue%03d", iteration)
			SetStorageUpdateStrings(parentAccount, key, value)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
			SetStorageUpdateStrings(childAccount, key, value)
		}
		AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		AddFinishData(vmOutput, []byte("succ"))
	}

	counterValue := (recursiveCalls + recursiveCalls%2) / 2
	SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(counterValue + 1)})
	SetStorageUpdate(childAccount, recursiveIterationCounterKey, []byte{byte(counterValue)})
	if recursiveCalls%2 == 0 {
		SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	} else {
		SetStorageUpdate(childAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	}

	return vmOutput
}

func expectedVMOutput_AsyncCall(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)

	_ = AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)
	outTransfer := vmcommon.OutputTransfer{Data: []byte(" there"), Value: big.NewInt(3)}
	outAcc := vmOutput.OutputAccounts[string(thirdPartyAddress)]
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer)
	outAcc.BalanceDelta = big.NewInt(6)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	SetStorageUpdate(childAccount, childKey, childData)

	_ = AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	AddFinishData(vmOutput, []byte{0})
	AddFinishData(vmOutput, []byte("thirdparty"))
	AddFinishData(vmOutput, []byte("vault"))
	AddFinishData(vmOutput, []byte{0})
	AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutput_AsyncCall_ChildFails(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-7,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)

	_ = AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutput_AsyncCall_CallBackFails(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()

	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	AddFinishData(vmOutput, parentFinishA)
	AddFinishData(vmOutput, parentFinishB)

	_ = AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)
	outTransfer2 := vmcommon.OutputTransfer{Value: big.NewInt(3), Data: []byte(" there")}
	outAcc := vmOutput.OutputAccounts[string(thirdPartyAddress)]
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer2)
	outAcc.BalanceDelta = big.NewInt(6)

	childAccount := AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	childAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	SetStorageUpdate(childAccount, childKey, childData)

	_ = AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	AddFinishData(vmOutput, []byte{3})
	AddFinishData(vmOutput, []byte("thirdparty"))
	AddFinishData(vmOutput, []byte("vault"))
	AddFinishData(vmOutput, []byte("user error"))
	AddFinishData(vmOutput, []byte("txhash"))

	vmOutput.ReturnMessage = "callBack error"

	return vmOutput
}

func expectedVMOutput_CreateNewContract_Success(_ []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-42,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.Nonce = 1
	SetStorageUpdate(parentAccount, []byte{'A'}, childCode)

	childAccount := AddNewOutputAccount(
		vmOutput,
		[]byte("newAddress"),
		42,
		nil,
	)
	childAccount.Code = childCode
	childAccount.CodeMetadata = []byte{1, 0}
	childAccount.CodeDeployerAddress = parentAddress

	l := len(childCode)
	AddFinishData(vmOutput, []byte{byte(l / 256), byte(l % 256)})
	AddFinishData(vmOutput, []byte("init successful"))
	AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutput_CreateNewContract_Fail(_ []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := MakeVMOutput()
	parentAccount := AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Nonce = 0
	SetStorageUpdate(parentAccount, []byte{'A'}, childCode)

	l := len(childCode)
	AddFinishData(vmOutput, []byte{byte(l / 256), byte(l % 256)})
	AddFinishData(vmOutput, []byte("fail"))

	return vmOutput
}
