package host

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

var parentCompilationCostSameCtx uint64
var childCompilationCostSameCtx uint64

var parentCompilationCostDestCtx uint64
var childCompilationCostDestCtx uint64

var vaultAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fvaultAddress..........")
var thirdPartyAddress = []byte("\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x0fthirdPartyAddress.....")

func init() {
	parentCompilationCostSameCtx = uint64(len(GetTestSCCode("exec-same-ctx-parent", "../../")))
	childCompilationCostSameCtx = uint64(len(GetTestSCCode("exec-same-ctx-child", "../../")))

	parentCompilationCostDestCtx = uint64(len(GetTestSCCode("exec-dest-ctx-parent", "../../")))
	childCompilationCostDestCtx = uint64(len(GetTestSCCode("exec-dest-ctx-child", "../../")))
}
