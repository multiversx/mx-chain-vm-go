package testcommon

// ParentKeyA value exposed for test usage
var ParentKeyA = []byte("parentKeyA......................")

// ParentKeyB value exposed for test usage
var ParentKeyB = []byte("parentKeyB......................")

// ParentDataA value exposed for test usage
var ParentDataA = []byte("parentDataA")

// ParentDataB value exposed for test usage
var ParentDataB = []byte("parentDataB")

// ChildKey value exposed for test usage
var ChildKey = []byte("childKey........................")

// ChildData value exposed for test usage
var ChildData = []byte("childData")

// ParentFinishA value exposed for test usage
var ParentFinishA = []byte("parentFinishA")

// ParentFinishB value exposed for test usage
var ParentFinishB = []byte("parentFinishB")

// ChildFinish value exposed for test usage
var ChildFinish = []byte("childFinish")

// ParentTransferReceiver value exposed for test usage
var ParentTransferReceiver = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fparentTransferReceiver")

// ChildTransferReceiver value exposed for test usage
var ChildTransferReceiver = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fchildTransferReceiver.")

// ParentTransferValue value exposed for test usage
var ParentTransferValue = int64(42)

// ParentTransferData value exposed for test usage
var ParentTransferData = []byte("parentTransferData")

// RecursiveIterationCounterKey value exposed for test usage
var RecursiveIterationCounterKey = []byte("recursiveIterationCounter.......")

// RecursiveIterationBigCounterKey value exposed for test usage
var RecursiveIterationBigCounterKey = []byte("recursiveIterationBigCounter....")

// GasProvided value exposed for test usage
var GasProvided = uint64(1000000)

// ParentCompilationCostSameCtx value exposed for test usage
var ParentCompilationCostSameCtx uint64

// ChildCompilationCostSameCtx value exposed for test usage
var ChildCompilationCostSameCtx uint64

// ParentCompilationCostDestCtx value exposed for test usage
var ParentCompilationCostDestCtx uint64

// ChildCompilationCostDestCtx value exposed for test usage
var ChildCompilationCostDestCtx uint64

// VaultAddress value exposed for test usage
var VaultAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fvaultAddress..........")

// ThirdPartyAddress value exposed for test usage
var ThirdPartyAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fthirdPartyAddress.....")
