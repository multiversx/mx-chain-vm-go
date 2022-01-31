package vmjsonintegrationtest

import (
	"testing"
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
	// TODO The two excluded tests perform async calls from within async calls,
	// which are unsupported by the legacy async calls on which the forwarder is
	// currently based. The new AsyncContext will block multi-level async calls
	// anyway in its first release.
	runTestsInFolder(t, "features/composability/mandos", []string{
		"features/composability/mandos/forwarder_send_twice_egld.scen.json",
		"features/composability/mandos/forwarder_send_twice_esdt.scen.json",
	})
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
	// TODO The two excluded tests perform async calls from within async calls,
	// which are unsupported by the legacy async calls on which the forwarder is
	// currently based. The new AsyncContext will block multi-level async calls
	// anyway in its first release.
	runTestsInFolder(t, "features/composability/mandos-legacy", []string{
		"features/composability/mandos-legacy/l_forwarder_send_twice_egld.scen.json",
		"features/composability/mandos-legacy/l_forwarder_send_twice_esdt.scen.json",
	})
}

func TestTimelocks(t *testing.T) {
	runAllTestsInFolder(t, "timelocks")
}
