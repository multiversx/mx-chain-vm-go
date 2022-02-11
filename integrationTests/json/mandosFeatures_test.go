package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRustBasicFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features/mandos")
}

func TestRustBasicFeaturesNoSmallIntApi(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features-no-small-int-api/mandos")
}

// Backwards compatibility.
func TestRustBasicFeaturesLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features-legacy/mandos")
}

func TestRustPayableFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/payable-features/mandos")
}

func TestRustComposability(t *testing.T) {
	runAllTestsInFolder(t, "features/composability/mandos")
}

// For debugging:
// func TestESDTMultiTransferOnCallback(t *testing.T) {
// 	err := runSingleTestReturnError(
// 		"features/composability/mandos",
// 		"forw_raw_call_async_retrieve_multi_transfer.scen.json")
// 	require.Nil(t, err)
// }

// func TestESDTMultiTransferOnCallAndCallback(t *testing.T) {
// 	err := runSingleTestReturnError(
// 		"features/composability/mandos",
// 		"forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json")
// 	require.Nil(t, err)
// }

// func TestExecOnDestByCallerAndNFTCreate(t *testing.T) {
// 	err := runSingleTestReturnError(
// 		"features/composability/mandos",
// 		"forwarder_builtin_nft_create_by_caller.scen.json")
// 	require.Nil(t, err)
// }

func TestRustLegacyComposability(t *testing.T) {
	runAllTestsInFolder(t, "features/composability/mandos-legacy")
}

func TestTimelocks(t *testing.T) {
	runAllTestsInFolder(t, "timelocks")
}

// func TestPromises(t *testing.T) {
// 	executor, err := am.NewArwenTestExecutor()
// 	require.Nil(t, err)
// 	runner := mc.NewScenarioRunner(
// 		executor,
// 		mc.NewDefaultFileResolver(),
// 	)
// 	err = runner.RunAllJSONScenariosInDirectory(
// 		getTestRoot(),
// 		"promises",
// 		".scen.json",
// 		[]string{})

// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestForwarderTransfExec(t *testing.T) {
	err := runSingleTestReturnError("features/composability/mandos", "forwarder_call_transf_exec_nft_reject.scen.json")
	require.Nil(t, err)
}

func TestForwarderTransfExecMultiReject(t *testing.T) {
	err := runSingleTestReturnError("features/composability/mandos", "forwarder_call_transf_exec_multi_transfer_reject.scen.json")
	require.Nil(t, err)
}
