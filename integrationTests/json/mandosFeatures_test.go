package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRustAllocFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/alloc-features/mandos")
}

func TestRustBasicFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runTestsInFolder(t, "features/basic-features/mandos", []string{
		"features/basic-features/mandos/storage_mapper_fungible_token.scen.json"})
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

func TestRustBigFloatFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/big-float-features/mandos")
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

func TestRustFormattedMessageFeatures(t *testing.T) {
	runAllTestsInFolder(t, "features/formatted-message-features/mandos")
}

func TestRustLegacyComposability(t *testing.T) {
	runAllTestsInFolder(t, "features/composability/mandos-legacy")
}

func TestTimelocks(t *testing.T) {
	runAllTestsInFolder(t, "timelocks")
}

func TestIndividualScenarios(t *testing.T) {
	var err error
	err = runSingleTestReturnError("features/composability/mandos", "forw_raw_contract_upgrade_self.scen.json")
	require.Nil(t, err)
}
