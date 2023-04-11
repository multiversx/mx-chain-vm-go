package vmjsonintegrationtest

import (
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
	// TODO The two excluded tests perform async calls from within async calls,
	// which are unsupported by the legacy async calls on which the forwarder is
	// currently based. The new AsyncContext will block multi-level async calls
	// anyway in its first release.
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		Exclude("features/composability/scenarios/forwarder_send_twice_egld.scen.json").
		Exclude("features/composability/scenarios/forwarder_send_twice_esdt.scen.json").
		Run().
		CheckNoError()
}

func TestRustPromisesFeatures(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios-promises").
		Run().
		CheckNoError()
}

// TODO: debug, then delete
func TestRustPromisesFeaturesDebug(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios-promises/promises_call_async_retrieve_egld.scen.json").
		Run().
		CheckNoError()
}

func TestRustFormattedMessageFeatures(t *testing.T) {
	ScenariosTest(t).
		Folder("features/formatted-message-features/scenarios").
		Run().
		CheckNoError()
}

// New contracts no longer the older, unmanaged hooks.
// We have older contracts that just do regression checking.
func TestRustLegacyComposability(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability-legacy/scenarios-legacy").
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
