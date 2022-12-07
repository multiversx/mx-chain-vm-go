package vmjsonintegrationtest

import (
	"testing"
)

func TestRustAllocFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("features/alloc-features/mandos").
		Run().
		CheckNoError()
}

func TestRustBasicFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("features/basic-features/mandos").
		Exclude("features/basic-features/mandos/storage_mapper_fungible_token.scen.json").
		Run().
		CheckNoError()
}

func TestRustBasicFeaturesNoSmallIntApi(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("features/basic-features-no-small-int-api/mandos").
		Run().
		CheckNoError()
}

// Backwards compatibility.
func TestRustBasicFeaturesLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("features/basic-features-legacy/mandos").
		Run().
		CheckNoError()
}

func TestRustBigFloatFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("features/big-float-features/mandos").
		Run().
		CheckNoError()
}

func TestRustPayableFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("features/payable-features/mandos").
		Run().
		CheckNoError()
}

func TestRustComposability(t *testing.T) {
	// TODO The two excluded tests perform async calls from within async calls,
	// which are unsupported by the legacy async calls on which the forwarder is
	// currently based. The new AsyncContext will block multi-level async calls
	// anyway in its first release.
	MandosTest(t).
		Folder("features/composability/mandos").
		Exclude("features/composability/mandos/forwarder_send_twice_egld.scen.json").
		Exclude("features/composability/mandos/forwarder_send_twice_esdt.scen.json").
		Run().
		CheckNoError()
}

func TestRustPromisesFeatures(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos-promises").
		Run().
		CheckNoError()
}

// TODO: debug, then delete
func TestRustPromisesFeaturesDebug(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos-promises/promises_call_async_retrieve_egld.scen.json").
		Run().
		CheckNoError()
}

func TestRustFormattedMessageFeatures(t *testing.T) {
	MandosTest(t).
		Folder("features/formatted-message-features/mandos").
		Run().
		CheckNoError()
}

func TestRustLegacyComposability(t *testing.T) {
	// TODO The two excluded tests perform async calls from within async calls,
	// which are unsupported by the legacy async calls on which the forwarder is
	// currently based. The new AsyncContext will block multi-level async calls
	// anyway in its first release.
	MandosTest(t).
		Folder("features/composability/mandos-legacy").
		Exclude("features/composability/mandos-legacy/l_forwarder_send_twice_egld.scen.json").
		Exclude("features/composability/mandos-legacy/l_forwarder_send_twice_esdt.scen.json").
		Run().
		CheckNoError()

}

func TestTimelocks(t *testing.T) {
	MandosTest(t).
		Folder("timelocks").
		Run().
		CheckNoError()
}

func TestForwarderTransfExec(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos").
		File("forwarder_call_transf_exec_reject_nft.scen.json").
		Run().
		CheckNoError()
}

func TestForwarderTransfExecMultiReject(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos").
		File("forwarder_call_transf_exec_reject_multi_transfer.scen.json").
		Run().
		CheckNoError()
}
