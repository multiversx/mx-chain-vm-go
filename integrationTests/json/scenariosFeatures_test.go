package vmjsonintegrationtest

import (
	"path/filepath"
	"testing"
)

func TestRustAllocFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/alloc-features/scenarios").
		Run().
		CheckNoError()
}

func TestRustBasicFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/basic-features/scenarios").
		Exclude("features/basic-features/scenarios/storage_mapper_fungible_token.scen.json").
		Exclude("features/basic-features/scenarios/get_shard_of_address.scen.json").
		Run().
		CheckNoError()
}

func TestRustBasicFeaturesNoSmallIntApi(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/basic-features-no-small-int-api/scenarios").
		Run().
		CheckNoError()
}

// Backwards compatibility.
func TestRustBasicFeaturesLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/basic-features-legacy/scenarios").
		Run().
		CheckNoError()
}

func TestRustBigFloatFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/big-float-features/scenarios").
		Run().
		CheckNoError()
}

func TestRustManagedMapBenchmark(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/managed-map-benchmark/scenarios").
		Run().
		CheckNoError()
}

func TestRustManagedMapFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/managed-map-features/scenarios").
		Run().
		CheckNoError()
}

func TestRustPayableFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/payable-features/scenarios").
		Run().
		CheckNoError()
}

func TestRustPayableFeaturesBarnard(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	// we run the exact same tests, but on a different compiled contract,
	// this contract uses the new "managedGetAllTransfersCallValue" VM hook,
	// which is only available in Barnard
	// this will become the only test

	ScenariosTest(t).
		Folder("features/payable-features/scenarios").
		ReplacePath(
			"../output/payable-features.mxsc.json",
			filepath.Join(getTestRoot(), "features/payable-features/output/payable-features-barnard.mxsc.json"),
		).
		Run().
		CheckNoError()
}

func TestRustPanicMessageFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/panic-message-features/scenarios").
		Run().
		CheckNoError()
}

func TestRustComposability(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		Run().
		CheckNoError()
}

func TestRustFormattedMessageFeatures(t *testing.T) {
	ScenariosTest(t).
		Folder("features/formatted-message-features/scenarios").
		Run().
		CheckNoError()
}

func TestTimelocks(t *testing.T) {
	ScenariosTest(t).
		Folder("timelocks").
		Run().
		CheckNoError()
}

func TestForwarderTransfExecFallibleMultiReject(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_transf_exec_fallible_multi_egld_reject.scen.json").
		Run().
		CheckNoError()
}

func TestBalanceAfterGet(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forwarder_call_sync_retrieve_bt_multi.scen.json").
		Run().
		CheckNoError()
}

func TestNFTDecode(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forwarder_nft_decode_complex_attributes.scen.json").
		Run().
		CheckNoError()
}
