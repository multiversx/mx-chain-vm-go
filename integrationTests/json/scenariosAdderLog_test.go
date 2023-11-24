package vmjsonintegrationtest

import (
	"testing"
)

const expectedAdderLog = `starting log:
GetFunctionNames: [add callBack getSum init upgrade]
ValidateFunctionArities: true
GetFunctionNames: [add callBack getSum init upgrade]
GetFunctionNames: [add callBack getSum init upgrade]
HasFunction(init): true
CallFunction(init):
VM hook begin: CheckNoPayment()
GetPointsUsed: 3
SetPointsUsed: 103
VM hook end:   CheckNoPayment()
VM hook begin: GetNumArguments()
GetPointsUsed: 110
SetPointsUsed: 210
VM hook end:   GetNumArguments()
VM hook begin: BigIntGetUnsignedArgument(0, -101)
GetPointsUsed: 249
SetPointsUsed: 1249
VM hook end:   BigIntGetUnsignedArgument(0, -101)
VM hook begin: MBufferSetBytes(-102, 131097, 3)
GetPointsUsed: 1289
SetPointsUsed: 3289
GetPointsUsed: 3289
SetPointsUsed: 6289
VM hook end:   MBufferSetBytes(-102, 131097, 3)
VM hook begin: MBufferFromBigIntUnsigned(-103, -101)
GetPointsUsed: 6333
SetPointsUsed: 10333
VM hook end:   MBufferFromBigIntUnsigned(-103, -101)
VM hook begin: MBufferStorageStore(-102, -103)
GetPointsUsed: 10345
SetPointsUsed: 85345
GetPointsUsed: 85345
GetPointsUsed: 85345
SetPointsUsed: 85345
GetPointsUsed: 85345
GetPointsUsed: 85345
SetPointsUsed: 135345
VM hook end:   MBufferStorageStore(-102, -103)
GetPointsUsed: 135352
GetPointsUsed: 135352
GetPointsUsed: 135352
GetPointsUsed: 135352
GetPointsUsed: 135352
GetPointsUsed: 135352
GetPointsUsed: 135352
Reset: true
SetPointsUsed: 0
SetGasLimit: 9223372036853566707
SetBreakpointValue: 0
HasFunction(getSum): true
CallFunction(getSum):
VM hook begin: CheckNoPayment()
GetPointsUsed: 3
SetPointsUsed: 103
VM hook end:   CheckNoPayment()
VM hook begin: GetNumArguments()
GetPointsUsed: 110
SetPointsUsed: 210
VM hook end:   GetNumArguments()
VM hook begin: MBufferSetBytes(-101, 131097, 3)
GetPointsUsed: 250
SetPointsUsed: 2250
GetPointsUsed: 2250
SetPointsUsed: 5250
VM hook end:   MBufferSetBytes(-101, 131097, 3)
VM hook begin: MBufferStorageLoad(-101, -102)
GetPointsUsed: 5291
SetPointsUsed: 6291
GetPointsUsed: 6291
GetPointsUsed: 6291
SetPointsUsed: 56291
VM hook end:   MBufferStorageLoad(-101, -102)
VM hook begin: MBufferToBigIntUnsigned(-102, -103)
GetPointsUsed: 56324
SetPointsUsed: 60324
VM hook end:   MBufferToBigIntUnsigned(-102, -103)
VM hook begin: BigIntFinishUnsigned(-103)
GetPointsUsed: 60335
SetPointsUsed: 61335
GetPointsUsed: 61335
SetPointsUsed: 71335
VM hook end:   BigIntFinishUnsigned(-103)
GetPointsUsed: 71337
GetPointsUsed: 71337
GetPointsUsed: 71337
GetPointsUsed: 71337
GetPointsUsed: 71337
GetPointsUsed: 71337
GetPointsUsed: 71337
GetPointsUsed: 71337
Reset: true
SetPointsUsed: 0
SetGasLimit: 3790900
SetBreakpointValue: 0
HasFunction(add): true
CallFunction(add):
VM hook begin: CheckNoPayment()
GetPointsUsed: 3
SetPointsUsed: 103
VM hook end:   CheckNoPayment()
VM hook begin: GetNumArguments()
GetPointsUsed: 110
SetPointsUsed: 210
VM hook end:   GetNumArguments()
VM hook begin: BigIntGetUnsignedArgument(0, -101)
GetPointsUsed: 249
SetPointsUsed: 1249
VM hook end:   BigIntGetUnsignedArgument(0, -101)
VM hook begin: MBufferSetBytes(-102, 131097, 3)
GetPointsUsed: 1289
SetPointsUsed: 3289
GetPointsUsed: 3289
SetPointsUsed: 6289
VM hook end:   MBufferSetBytes(-102, 131097, 3)
VM hook begin: MBufferStorageLoad(-102, -103)
GetPointsUsed: 6333
SetPointsUsed: 7333
GetPointsUsed: 7333
GetPointsUsed: 7333
SetPointsUsed: 57333
VM hook end:   MBufferStorageLoad(-102, -103)
VM hook begin: MBufferToBigIntUnsigned(-103, -104)
GetPointsUsed: 57366
SetPointsUsed: 61366
VM hook end:   MBufferToBigIntUnsigned(-103, -104)
VM hook begin: BigIntAdd(-104, -104, -101)
GetPointsUsed: 61386
SetPointsUsed: 63386
VM hook end:   BigIntAdd(-104, -104, -101)
VM hook begin: MBufferFromBigIntUnsigned(-105, -104)
GetPointsUsed: 63425
SetPointsUsed: 67425
VM hook end:   MBufferFromBigIntUnsigned(-105, -104)
VM hook begin: MBufferStorageStore(-102, -105)
GetPointsUsed: 67437
SetPointsUsed: 142437
GetPointsUsed: 142437
GetPointsUsed: 142437
SetPointsUsed: 142437
GetPointsUsed: 142437
GetPointsUsed: 142437
SetPointsUsed: 142437
VM hook end:   MBufferStorageStore(-102, -105)
GetPointsUsed: 142444
GetPointsUsed: 142444
GetPointsUsed: 142444
GetPointsUsed: 142444
GetPointsUsed: 142444
GetPointsUsed: 142444
GetPointsUsed: 142444
GetPointsUsed: 142444
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
