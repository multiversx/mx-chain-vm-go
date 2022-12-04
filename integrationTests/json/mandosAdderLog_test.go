package vmjsonintegrationtest

import (
	"testing"
)

const expectedAdderLog = `starting log:
GetFunctionNames: [add callBack getSum init]
ValidateVoidFunction(add): true
ValidateVoidFunction(callBack): true
ValidateVoidFunction(getSum): true
ValidateVoidFunction(init): true
GetFunctionNames: [add callBack getSum init]
GetFunctionNames: [add callBack getSum init]
HasFunction(init): true
CallFunction(init):
VM hook begin: CheckNoPayment() points used: 0
VM hook end:   CheckNoPayment() points used: 100
VM hook begin: GetNumArguments() points used: 100
VM hook end:   GetNumArguments() points used: 200
VM hook begin: BigIntGetUnsignedArgument(0, -101) points used: 200
VM hook end:   BigIntGetUnsignedArgument(0, -101) points used: 1200
VM hook begin: MBufferSetBytes(-102, 1048601, 3) points used: 1200
VM hook end:   MBufferSetBytes(-102, 1048601, 3) points used: 6200
VM hook begin: MBufferFromBigIntUnsigned(-103, -101) points used: 6200
VM hook end:   MBufferFromBigIntUnsigned(-103, -101) points used: 10200
VM hook begin: MBufferStorageStore(-102, -103) points used: 10200
VM hook end:   MBufferStorageStore(-102, -103) points used: 135200
Reset: true
SetGasLimit: 18446744073708343115
SetBreakpointValue: 0
HasFunction(getSum): true
CallFunction(getSum):
VM hook begin: CheckNoPayment() points used: 0
VM hook end:   CheckNoPayment() points used: 100
VM hook begin: GetNumArguments() points used: 100
VM hook end:   GetNumArguments() points used: 200
VM hook begin: MBufferSetBytes(-101, 1048601, 3) points used: 200
VM hook end:   MBufferSetBytes(-101, 1048601, 3) points used: 5200
VM hook begin: MBufferStorageLoad(-101, -102) points used: 5200
VM hook end:   MBufferStorageLoad(-101, -102) points used: 56200
VM hook begin: MBufferToBigIntUnsigned(-102, -103) points used: 56200
VM hook end:   MBufferToBigIntUnsigned(-102, -103) points used: 60200
VM hook begin: BigIntFinishUnsigned(-103) points used: 60200
VM hook end:   BigIntFinishUnsigned(-103) points used: 71200
Reset: true
SetGasLimit: 3791500
SetBreakpointValue: 0
HasFunction(add): true
CallFunction(add):
VM hook begin: CheckNoPayment() points used: 0
VM hook end:   CheckNoPayment() points used: 100
VM hook begin: GetNumArguments() points used: 100
VM hook end:   GetNumArguments() points used: 200
VM hook begin: BigIntGetUnsignedArgument(0, -101) points used: 200
VM hook end:   BigIntGetUnsignedArgument(0, -101) points used: 1200
VM hook begin: MBufferSetBytes(-102, 1048601, 3) points used: 1200
VM hook end:   MBufferSetBytes(-102, 1048601, 3) points used: 6200
VM hook begin: MBufferStorageLoad(-102, -103) points used: 6200
VM hook end:   MBufferStorageLoad(-102, -103) points used: 57200
VM hook begin: MBufferToBigIntUnsigned(-103, -104) points used: 57200
VM hook end:   MBufferToBigIntUnsigned(-103, -104) points used: 61200
VM hook begin: BigIntAdd(-104, -104, -101) points used: 61200
VM hook end:   BigIntAdd(-104, -104, -101) points used: 63200
VM hook begin: MBufferFromBigIntUnsigned(-105, -104) points used: 63200
VM hook end:   MBufferFromBigIntUnsigned(-105, -104) points used: 67200
VM hook begin: MBufferStorageStore(-102, -105) points used: 67200
VM hook end:   MBufferStorageStore(-102, -105) points used: 142200
`

func TestRustAdderLog(t *testing.T) {
	MandosTest(t).
		Folder("adder/mandos").
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expectedAdderLog)
}
