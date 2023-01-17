package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRustAllocFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/alloc-features/scenarios")
}

func TestRustBasicFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runTestsInFolder(t, "features/basic-features/scenarios", []string{
		"features/basic-features/scenarios/storage_mapper_fungible_token.scen.json"})
}

func TestRustBasicFeaturesNoSmallIntApi(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features-no-small-int-api/scenarios")
}

// Backwards compatibility.
func TestRustBasicFeaturesLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/basic-features-legacy/scenarios")
}

func TestRustBigFloatFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/big-float-features/scenarios")
}

func TestRustPayableFeaturesLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "features/payable-features/scenarios")
}

func TestRustComposability(t *testing.T) {
	runAllTestsInFolder(t, "features/composability/scenarios")
}

func TestRustFormattedMessageFeatures(t *testing.T) {
	runAllTestsInFolder(t, "features/formatted-message-features/scenarios")
}

func TestRustLegacyComposability(t *testing.T) {
	runAllTestsInFolder(t, "features/composability/scenarios-legacy")
}

func TestTimelocks(t *testing.T) {
	runAllTestsInFolder(t, "timelocks")
}

func TestIndividualScenarios(t *testing.T) {
	err := runSingleTestReturnError("features/composability/scenarios", "forw_raw_contract_upgrade_self.scen.json")
	require.Nil(t, err)
}
