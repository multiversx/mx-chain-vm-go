package vmjsonintegrationtest

import (
	"testing"
)

const expectedAdderLog = `starting log:
GetFunctionNames: [add callBack getSum init]
ValidateFunctionArities: true
GetFunctionNames: [add callBack getSum init]
GetFunctionNames: [add callBack getSum init]
HasFunction(init): true
CallFunction(init):
VM hook begin: CheckNoPayment() points used: 3
VM hook end:   CheckNoPayment() points used: 103
VM hook begin: GetNumArguments() points used: 110
VM hook end:   GetNumArguments() points used: 210
VM hook begin: BigIntGetUnsignedArgument(0, -101) points used: 249
VM hook end:   BigIntGetUnsignedArgument(0, -101) points used: 1249
VM hook begin: MBufferSetBytes(-102, 1048576, 3) points used: 1289
VM hook end:   MBufferSetBytes(-102, 1048576, 3) points used: 6289
VM hook begin: MBufferFromBigIntUnsigned(-103, -101) points used: 6333
VM hook end:   MBufferFromBigIntUnsigned(-103, -101) points used: 10333
VM hook begin: MBufferStorageStore(-102, -103) points used: 10345
VM hook end:   MBufferStorageStore(-102, -103) points used: 135345
Reset: true
SetGasLimit: 18446744073708343115
SetBreakpointValue: 0
HasFunction(getSum): true
CallFunction(getSum):
VM hook begin: CheckNoPayment() points used: 3
VM hook end:   CheckNoPayment() points used: 103
VM hook begin: GetNumArguments() points used: 110
VM hook end:   GetNumArguments() points used: 210
VM hook begin: MBufferSetBytes(-101, 1048576, 3) points used: 250
VM hook end:   MBufferSetBytes(-101, 1048576, 3) points used: 5250
VM hook begin: MBufferStorageLoad(-101, -102) points used: 5291
VM hook end:   MBufferStorageLoad(-101, -102) points used: 56291
VM hook begin: MBufferToBigIntUnsigned(-102, -103) points used: 56324
VM hook end:   MBufferToBigIntUnsigned(-102, -103) points used: 60324
VM hook begin: BigIntFinishUnsigned(-103) points used: 60335
VM hook end:   BigIntFinishUnsigned(-103) points used: 71335
Reset: true
SetGasLimit: 3791500
SetBreakpointValue: 0
HasFunction(add): true
CallFunction(add):
VM hook begin: CheckNoPayment() points used: 3
VM hook end:   CheckNoPayment() points used: 103
VM hook begin: GetNumArguments() points used: 110
VM hook end:   GetNumArguments() points used: 210
VM hook begin: BigIntGetUnsignedArgument(0, -101) points used: 249
VM hook end:   BigIntGetUnsignedArgument(0, -101) points used: 1249
VM hook begin: MBufferSetBytes(-102, 1048576, 3) points used: 1289
VM hook end:   MBufferSetBytes(-102, 1048576, 3) points used: 6289
VM hook begin: MBufferStorageLoad(-102, -103) points used: 6333
VM hook end:   MBufferStorageLoad(-102, -103) points used: 57333
VM hook begin: MBufferToBigIntUnsigned(-103, -104) points used: 57366
VM hook end:   MBufferToBigIntUnsigned(-103, -104) points used: 61366
VM hook begin: BigIntAdd(-104, -104, -101) points used: 61386
VM hook end:   BigIntAdd(-104, -104, -101) points used: 63386
VM hook begin: MBufferFromBigIntUnsigned(-105, -104) points used: 63425
VM hook end:   MBufferFromBigIntUnsigned(-105, -104) points used: 67425
VM hook begin: MBufferStorageStore(-102, -105) points used: 67437
VM hook end:   MBufferStorageStore(-102, -105) points used: 142437
Clean: true
`

func TestRustAdderLog(t *testing.T) {
	ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expectedAdderLog)
}
