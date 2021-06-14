package vmjsonintegrationtest

import (
	"testing"
)

func TestDelegation_v0_2(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runSingleTest(t, "delegation/v0_3/test", "fuzz_gen.scen.json")
}

func TestDelegation_v0_3(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runTestsInFolder(t, "delegation/v0_3", []string{
		"delegation/v0_3/test/integration/genesis/genesis.scen.json",
	})
}

func TestDelegation_v0_4_genesis(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_4_genesis")
}

func TestDelegation_v0_5_latest_full(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_5_latest_full")
}

func TestDelegation_v0_5_latest_update(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "delegation/v0_5_latest_update")
}
