package vmjsonintegrationtest

import (
	"testing"
)

func TestDelegation_v0_2(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("delegation/v0_2").
		Run().
		CheckNoError()
}

func TestDelegation_v0_3(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("delegation/v0_3").
		Exclude("delegation/v0_3/test/integration/genesis/genesis.scen.json").
		Run().
		CheckNoError()
}

func TestDelegation_v0_4_genesis(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("delegation/v0_4_genesis").
		Run().
		CheckNoError()
}

func TestDelegation_v0_5_latest_full(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("delegation/v0_5_latest_full").
		Run().
		CheckNoError()
}

func TestDelegation_v0_5_latest_update(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("delegation/v0_5_latest_update").
		Run().
		CheckNoError()

}
