package host

import (
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
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
	parentCompilationCost_SameCtx = uint64(len(arwen.GetTestSCCode("exec-same-ctx-parent", "../../")))
	childCompilationCost_SameCtx = uint64(len(arwen.GetTestSCCode("exec-same-ctx-child", "../../")))

	parentCompilationCost_DestCtx = uint64(len(arwen.GetTestSCCode("exec-dest-ctx-parent", "../../")))
	childCompilationCost_DestCtx = uint64(len(arwen.GetTestSCCode("exec-dest-ctx-child", "../../")))
}

func expectedVMOutput_SameCtx_Prepare(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)
	arwen.AddFinishData(vmOutput, []byte("succ"))

	expectedExecutionCost := uint64(137)
	gas := gasProvided
	gas -= parentCompilationCost_SameCtx
	gas -= expectedExecutionCost
	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_WrongContractCalled(code []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare(code)

	arwen.AddFinishData(vmOutput, []byte("fail"))

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
	vmOutput := arwen.MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)

	arwen.AddFinishData(vmOutput, []byte("fail"))

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

func expectedVMOutput_SameCtx_Simple(parentCode []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	arwen.AddFinishData(vmOutput, []byte("child"))
	arwen.AddFinishData(vmOutput, []byte{})
	for i := 1; i < 100; i++ {
		arwen.AddFinishData(vmOutput, []byte{byte(i)})
	}
	arwen.AddFinishData(vmOutput, []byte{})
	arwen.AddFinishData(vmOutput, []byte("child"))
	arwen.AddFinishData(vmOutput, []byte{})
	for i := 1; i < 100; i++ {
		arwen.AddFinishData(vmOutput, []byte{byte(i)})
	}
	arwen.AddFinishData(vmOutput, []byte{})
	arwen.AddFinishData(vmOutput, []byte("parent"))

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-198,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		198,
		nil,
	)
	childAccount.GasUsed = 3434

	parentGasUsed := uint64(3107)
	childGasUsed := childAccount.GasUsed
	executionCost := parentGasUsed + childGasUsed

	gas := gasProvided
	gas -= uint64(len(parentCode))
	gas -= uint64(len(childCode))
	gas -= executionCost

	vmOutput.GasRemaining = gas

	return vmOutput
}

func expectedVMOutput_SameCtx_SuccessfulChildCall(parentCode []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := expectedVMOutput_SameCtx_Prepare(parentCode)

	parentAccount := vmOutput.OutputAccounts[string(parentAddress)]
	parentAccount.BalanceDelta = big.NewInt(-141)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		3,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	executeAPICost := uint64(39)
	childExecutionCost := uint64(436)
	childAccount.GasUsed = childCompilationCost_SameCtx + childExecutionCost

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		96,
		[]byte("qwerty"),
	)

	// The child SC has stored data on the parent's storage
	arwen.SetStorageUpdate(parentAccount, childKey, childData)

	// The called child SC will output some arbitrary data, but also data that it
	// has read from the storage of the parent.
	arwen.AddFinishData(vmOutput, childFinish)
	arwen.AddFinishData(vmOutput, parentDataA)
	for _, c := range parentDataA {
		arwen.AddFinishData(vmOutput, []byte{c})
	}
	arwen.AddFinishData(vmOutput, parentDataB)
	for _, c := range parentDataB {
		arwen.AddFinishData(vmOutput, []byte{c})
	}
	arwen.AddFinishData(vmOutput, []byte("child ok"))
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

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
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	// parentAccount.BalanceDelta = big.NewInt(-99)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)
	childExecutionCost := uint64(107)
	childAccount.GasUsed = childCompilationCost_SameCtx + childExecutionCost

	// The child SC will output "child ok" if it could read some expected Big
	// Ints directly from the parent's context.
	arwen.AddFinishData(vmOutput, []byte("child ok"))
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

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
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
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
		arwen.AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		arwen.SetStorageUpdateStrings(account, key, value)
	}

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_Direct_ErrMaxInstances(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	finishString := fmt.Sprintf("Rfinish%03d", recursiveCalls)
	arwen.AddFinishData(vmOutput, []byte(finishString))

	key := fmt.Sprintf("Rkey%03d.........................", recursiveCalls)
	value := fmt.Sprintf("Rvalue%03d", recursiveCalls)
	arwen.SetStorageUpdateStrings(account, key, value)

	arwen.AddFinishData(vmOutput, []byte("fail"))
	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(1)})

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_MutualMethods(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	account.GasUsed = 25412

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	arwen.AddFinishData(vmOutput, []byte("start recursive mutual calls"))

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
		arwen.SetStorageUpdateStrings(account, key, value)
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	arwen.AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutput_SameCtx_Recursive_MutualSCs(_ []byte, _ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.GasUsed = 3650

	childAccount := arwen.AddNewOutputAccount(
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
		arwen.SetStorageUpdateStrings(parentAccount, key, value)
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	arwen.SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(recursiveCalls+1)).Bytes())

	return vmOutput
}

func expectedVMOutput_SameCtx_BuiltinFunctions_1(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		42,
		nil,
	)
	account.Balance = big.NewInt(1000)

	arwen.AddFinishData(vmOutput, []byte("succ"))
	gasConsumed_builtinClaim := 100
	vmOutput.GasRemaining = uint64(98504 - gasConsumed_builtinClaim)

	return vmOutput
}

func expectedVMOutput_SameCtx_BuiltinFunctions_2(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0)

	arwen.AddFinishData(vmOutput, []byte("succ"))

	gasConsumed_builtinDoSomething := 0
	vmOutput.GasRemaining = uint64(98500 - gasConsumed_builtinDoSomething)

	return vmOutput
}

func expectedVMOutput_SameCtx_BuiltinFunctions_3(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	arwen.AddFinishData(vmOutput, []byte("fail"))
	vmOutput.GasRemaining = 98000

	return vmOutput
}

func expectedVMOutput_DestCtx_Prepare(_ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-parentTransferValue,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		parentTransferReceiver,
		parentTransferValue,
		parentTransferData,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)

	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)
	arwen.AddFinishData(vmOutput, []byte("succ"))

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

	arwen.AddFinishData(vmOutput, []byte("fail"))

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
	vmOutput := arwen.MakeVMOutput()

	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)

	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)

	arwen.AddFinishData(vmOutput, []byte("fail"))

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

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		99-12,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		childTransferReceiver,
		12,
		[]byte("Second sentence."),
	)

	arwen.SetStorageUpdate(childAccount, childKey, childData)

	arwen.AddFinishData(vmOutput, childFinish)
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

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
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.Ok

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-99,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		99,
		nil,
	)

	// The child SC will output "child ok" if it could NOT read the Big Ints from
	// the parent's context.
	arwen.AddFinishData(vmOutput, []byte("child ok"))
	arwen.AddFinishData(vmOutput, []byte("succ"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

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
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	for i := recursiveCalls; i >= 0; i-- {
		finishString := fmt.Sprintf("Rfinish%03d", i)
		arwen.AddFinishData(vmOutput, []byte(finishString))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	for i := 0; i <= recursiveCalls; i++ {
		key := fmt.Sprintf("Rkey%03d.........................", i)
		value := fmt.Sprintf("Rvalue%03d", i)
		arwen.SetStorageUpdateStrings(account, key, value)
	}

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_MutualMethods(_ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	account := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	account.Balance = big.NewInt(1000)
	account.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))

	arwen.SetStorageUpdate(account, recursiveIterationCounterKey, []byte{byte(recursiveCalls + 1)})
	arwen.SetStorageUpdate(account, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())

	arwen.AddFinishData(vmOutput, []byte("start recursive mutual calls"))

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
		arwen.SetStorageUpdateStrings(account, key, value)
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	arwen.AddFinishData(vmOutput, []byte("end recursive mutual calls"))

	return vmOutput
}

func expectedVMOutput_DestCtx_Recursive_MutualSCs(_ []byte, _ []byte, recursiveCalls int) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentIterations := (recursiveCalls / 2) + (recursiveCalls % 2)
	childIterations := recursiveCalls - parentIterations

	balanceDelta := int64(5*parentIterations - 3*childIterations)

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.BalanceDelta = big.NewInt(-balanceDelta)

	childAccount := arwen.AddNewOutputAccount(
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
			arwen.SetStorageUpdateStrings(parentAccount, key, value)
		} else {
			finishData = fmt.Sprintf("Cfinish%03d", iteration)
			key = fmt.Sprintf("Ckey%03d.........................", iteration)
			value = fmt.Sprintf("Cvalue%03d", iteration)
			arwen.SetStorageUpdateStrings(childAccount, key, value)
		}
		arwen.AddFinishData(vmOutput, []byte(finishData))
	}

	for i := recursiveCalls - 1; i >= 0; i-- {
		arwen.AddFinishData(vmOutput, []byte("succ"))
	}

	counterValue := (recursiveCalls + recursiveCalls%2) / 2
	arwen.SetStorageUpdate(parentAccount, recursiveIterationCounterKey, []byte{byte(counterValue + 1)})
	arwen.SetStorageUpdate(childAccount, recursiveIterationCounterKey, []byte{byte(counterValue)})
	if recursiveCalls%2 == 0 {
		arwen.SetStorageUpdate(parentAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	} else {
		arwen.SetStorageUpdate(childAccount, recursiveIterationBigCounterKey, big.NewInt(int64(1)).Bytes())
	}

	return vmOutput
}

func expectedVMOutput_AsyncCall(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)
	outTransfer := vmcommon.OutputTransfer{Data: []byte(" there"), Value: big.NewInt(3)}
	outAcc := vmOutput.OutputAccounts[string(thirdPartyAddress)]
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer)
	outAcc.BalanceDelta = big.NewInt(6)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	arwen.SetStorageUpdate(childAccount, childKey, childData)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	arwen.AddFinishData(vmOutput, []byte{0})
	arwen.AddFinishData(vmOutput, []byte("thirdparty"))
	arwen.AddFinishData(vmOutput, []byte("vault"))
	arwen.AddFinishData(vmOutput, []byte{0})
	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutput_AsyncCall_ChildFails(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-7,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutput_AsyncCall_CallBackFails(_ []byte, _ []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	vmOutput.ReturnCode = vmcommon.UserError

	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-10,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	arwen.SetStorageUpdate(parentAccount, parentKeyA, parentDataA)
	arwen.SetStorageUpdate(parentAccount, parentKeyB, parentDataB)
	arwen.AddFinishData(vmOutput, parentFinishA)
	arwen.AddFinishData(vmOutput, parentFinishB)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		thirdPartyAddress,
		3,
		[]byte("hello"),
	)
	outTransfer2 := vmcommon.OutputTransfer{Value: big.NewInt(3), Data: []byte(" there")}
	outAcc := vmOutput.OutputAccounts[string(thirdPartyAddress)]
	outAcc.OutputTransfers = append(outAcc.OutputTransfers, outTransfer2)
	outAcc.BalanceDelta = big.NewInt(6)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		childAddress,
		0,
		nil,
	)
	childAccount.Balance = big.NewInt(0)
	childAccount.BalanceDelta = big.NewInt(0).Sub(big.NewInt(1), big.NewInt(1))
	arwen.SetStorageUpdate(childAccount, childKey, childData)

	_ = arwen.AddNewOutputAccount(
		vmOutput,
		vaultAddress,
		4,
		[]byte{},
	)

	arwen.AddFinishData(vmOutput, []byte{3})
	arwen.AddFinishData(vmOutput, []byte("thirdparty"))
	arwen.AddFinishData(vmOutput, []byte("vault"))
	arwen.AddFinishData(vmOutput, []byte("user error"))
	arwen.AddFinishData(vmOutput, []byte("txhash"))

	vmOutput.ReturnMessage = "callBack error"

	return vmOutput
}

func expectedVMOutput_CreateNewContract_Success(_ []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		-42,
		nil,
	)
	parentAccount.Balance = big.NewInt(1000)
	parentAccount.Nonce = 1
	arwen.SetStorageUpdate(parentAccount, []byte{'A'}, childCode)

	childAccount := arwen.AddNewOutputAccount(
		vmOutput,
		[]byte("newAddress"),
		42,
		nil,
	)
	childAccount.Code = childCode
	childAccount.CodeMetadata = []byte{1, 0}
	childAccount.CodeDeployerAddress = parentAddress

	l := len(childCode)
	arwen.AddFinishData(vmOutput, []byte{byte(l / 256), byte(l % 256)})
	arwen.AddFinishData(vmOutput, []byte("init successful"))
	arwen.AddFinishData(vmOutput, []byte("succ"))

	return vmOutput
}

func expectedVMOutput_CreateNewContract_Fail(_ []byte, childCode []byte) *vmcommon.VMOutput {
	vmOutput := arwen.MakeVMOutput()
	parentAccount := arwen.AddNewOutputAccount(
		vmOutput,
		parentAddress,
		0,
		nil,
	)
	parentAccount.Nonce = 0
	arwen.SetStorageUpdate(parentAccount, []byte{'A'}, childCode)

	l := len(childCode)
	arwen.AddFinishData(vmOutput, []byte{byte(l / 256), byte(l % 256)})
	arwen.AddFinishData(vmOutput, []byte("fail"))

	return vmOutput
}
