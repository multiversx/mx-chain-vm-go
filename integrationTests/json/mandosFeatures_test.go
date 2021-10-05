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
	runAllTestsInFolder(t, "features/composability/mandos")
}

func TestRustLegacyComposability(t *testing.T) {
	runAllTestsInFolder(t, "features/composability/mandos-legacy")
}

func TestTimelocks(t *testing.T) {
	runAllTestsInFolder(t, "timelocks")
}
