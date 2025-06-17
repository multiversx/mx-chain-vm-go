package vmjsonintegrationtest

import (
	logger "github.com/multiversx/mx-chain-logger-go"
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

func TestRustBarnardFeatures(t *testing.T) {
	// TODO: will get merged into basic-features after barnard mainnet release
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/barnard-features/scenarios").
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

func TestRustPayableFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("features/payable-features/scenarios").
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

func TestForwarderTransfExec(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forwarder_call_transf_exec_reject_nft.scen.json").
		Run().
		CheckNoError()
}

func TestForwarderTransfExecMultiReject(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forwarder_call_transf_exec_reject_multi_transfer.scen.json").
		Run().
		CheckNoError()
}

func TestForwarderTransfExecFallibleMultiReject(t *testing.T) {
	_ = logger.SetLogLevel("*:TRACE")
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_transf_exec_fallible_multi_egld_reject.scen.json").
		Run().
		CheckNoError()
}
