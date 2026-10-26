package vmjsonintegrationtest

import (
	"testing"
)

const expectedAdderLog = `starting log:
GetFunctionNames: [add callBack getSum init upgrade]
ValidateFunctionArities: true
GetFunctionNames: [add callBack getSum init upgrade]
HasFunction(init): true
CallFunction(init):
VM hook begin: CheckNoPayment()
GetPointsUsed: 3
GetPointsUsed: 3
SetPointsUsed: 103
VM hook end:   CheckNoPayment()
VM hook begin: GetNumArguments()
GetPointsUsed: 110
GetPointsUsed: 110
SetPointsUsed: 210
VM hook end:   GetNumArguments()
VM hook begin: BigIntGetUnsignedArgument(0, -201)
GetPointsUsed: 249
GetPointsUsed: 249
SetPointsUsed: 1249
GetPointsUsed: 1249
GetPointsUsed: 1249
SetPointsUsed: 2249
VM hook end:   BigIntGetUnsignedArgument(0, -201)
VM hook begin: MBufferSetBytes(-202, 131097, 3)
GetPointsUsed: 2289
GetPointsUsed: 2289
SetPointsUsed: 4289
GetPointsUsed: 4289
GetPointsUsed: 4289
SetPointsUsed: 7289
VM hook end:   MBufferSetBytes(-202, 131097, 3)
VM hook begin: MBufferFromBigIntUnsigned(-203, -201)
GetPointsUsed: 7333
GetPointsUsed: 7333
SetPointsUsed: 11333
VM hook end:   MBufferFromBigIntUnsigned(-203, -201)
VM hook begin: MBufferStorageStore(-202, -203)
GetPointsUsed: 11345
GetPointsUsed: 11345
SetPointsUsed: 86345
GetPointsUsed: 86345
GetPointsUsed: 86345
SetPointsUsed: 86345
GetPointsUsed: 86345
GetPointsUsed: 86345
SetPointsUsed: 136345
VM hook end:   MBufferStorageStore(-202, -203)
GetPointsUsed: 136352
GetPointsUsed: 136352
GetPointsUsed: 136352
GetPointsUsed: 136352
GetPointsUsed: 136352
GetPointsUsed: 136352
GetPointsUsed: 136352
Reset: true
SetPointsUsed: 0
SetGasLimit: 9223372036853566107
SetBreakpointValue: 0
HasFunction(getSum): true
CallFunction(getSum):
VM hook begin: CheckNoPayment()
GetPointsUsed: 3
GetPointsUsed: 3
SetPointsUsed: 103
VM hook end:   CheckNoPayment()
VM hook begin: GetNumArguments()
GetPointsUsed: 110
GetPointsUsed: 110
SetPointsUsed: 210
VM hook end:   GetNumArguments()
VM hook begin: MBufferSetBytes(-201, 131097, 3)
GetPointsUsed: 250
GetPointsUsed: 250
SetPointsUsed: 2250
GetPointsUsed: 2250
GetPointsUsed: 2250
SetPointsUsed: 5250
VM hook end:   MBufferSetBytes(-201, 131097, 3)
VM hook begin: MBufferStorageLoad(-201, -202)
GetPointsUsed: 5291
GetPointsUsed: 5291
SetPointsUsed: 6291
GetPointsUsed: 6291
GetPointsUsed: 6291
SetPointsUsed: 21578
GetPointsUsed: 21578
GetPointsUsed: 21578
SetPointsUsed: 22578
VM hook end:   MBufferStorageLoad(-201, -202)
VM hook begin: MBufferToBigIntUnsigned(-202, -203)
GetPointsUsed: 22611
GetPointsUsed: 22611
SetPointsUsed: 26611
GetPointsUsed: 26611
GetPointsUsed: 26611
SetPointsUsed: 27611
VM hook end:   MBufferToBigIntUnsigned(-202, -203)
VM hook begin: BigIntFinishUnsigned(-203)
GetPointsUsed: 27622
GetPointsUsed: 27622
SetPointsUsed: 28622
GetPointsUsed: 28622
GetPointsUsed: 28622
SetPointsUsed: 38622
VM hook end:   BigIntFinishUnsigned(-203)
GetPointsUsed: 38624
GetPointsUsed: 38624
GetPointsUsed: 38624
GetPointsUsed: 38624
GetPointsUsed: 38624
GetPointsUsed: 38624
GetPointsUsed: 38624
GetPointsUsed: 38624
Reset: true
SetPointsUsed: 0
SetGasLimit: 3790300
SetBreakpointValue: 0
HasFunction(add): true
CallFunction(add):
VM hook begin: CheckNoPayment()
GetPointsUsed: 3
GetPointsUsed: 3
SetPointsUsed: 103
VM hook end:   CheckNoPayment()
VM hook begin: GetNumArguments()
GetPointsUsed: 110
GetPointsUsed: 110
SetPointsUsed: 210
VM hook end:   GetNumArguments()
VM hook begin: BigIntGetUnsignedArgument(0, -201)
GetPointsUsed: 249
GetPointsUsed: 249
SetPointsUsed: 1249
GetPointsUsed: 1249
GetPointsUsed: 1249
SetPointsUsed: 2249
VM hook end:   BigIntGetUnsignedArgument(0, -201)
VM hook begin: MBufferSetBytes(-202, 131097, 3)
GetPointsUsed: 2289
GetPointsUsed: 2289
SetPointsUsed: 4289
GetPointsUsed: 4289
GetPointsUsed: 4289
SetPointsUsed: 7289
VM hook end:   MBufferSetBytes(-202, 131097, 3)
VM hook begin: MBufferStorageLoad(-202, -203)
GetPointsUsed: 7333
GetPointsUsed: 7333
SetPointsUsed: 8333
GetPointsUsed: 8333
GetPointsUsed: 8333
SetPointsUsed: 23620
GetPointsUsed: 23620
GetPointsUsed: 23620
SetPointsUsed: 24620
VM hook end:   MBufferStorageLoad(-202, -203)
VM hook begin: MBufferToBigIntUnsigned(-203, -204)
GetPointsUsed: 24653
GetPointsUsed: 24653
SetPointsUsed: 28653
GetPointsUsed: 28653
GetPointsUsed: 28653
SetPointsUsed: 29653
VM hook end:   MBufferToBigIntUnsigned(-203, -204)
VM hook begin: BigIntAdd(-204, -204, -201)
GetPointsUsed: 29673
GetPointsUsed: 29673
SetPointsUsed: 31673
VM hook end:   BigIntAdd(-204, -204, -201)
VM hook begin: MBufferFromBigIntUnsigned(-205, -204)
GetPointsUsed: 31712
GetPointsUsed: 31712
SetPointsUsed: 35712
VM hook end:   MBufferFromBigIntUnsigned(-205, -204)
VM hook begin: MBufferStorageStore(-202, -205)
GetPointsUsed: 35724
GetPointsUsed: 35724
SetPointsUsed: 110724
GetPointsUsed: 110724
GetPointsUsed: 110724
SetPointsUsed: 110724
GetPointsUsed: 110724
GetPointsUsed: 110724
SetPointsUsed: 110724
VM hook end:   MBufferStorageStore(-202, -205)
GetPointsUsed: 110731
GetPointsUsed: 110731
GetPointsUsed: 110731
GetPointsUsed: 110731
GetPointsUsed: 110731
GetPointsUsed: 110731
GetPointsUsed: 110731
GetPointsUsed: 110731
Clean: true
`

func TestRustAdderLog(t *testing.T) {
	ScenariosTest(t).
		Folder("adder/scenarios/adder.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expectedAdderLog)
}
