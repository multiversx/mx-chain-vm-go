package vmjsonintegrationtest

import (
	"runtime"
	"testing"
)

func TestDelegation_v0_2(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}
	if runtime.GOARCH == "arm64" {
		t.Skip("skipping test on arm64")
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
	if runtime.GOARCH == "arm64" {
		t.Skip("skipping test on arm64")
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
	if runtime.GOARCH == "arm64" {
		t.Skip("skipping test on arm64")
	}

	ScenariosTest(t).
		Folder("delegation/v0_4_genesis").
		Run().
		CheckNoError()
}

func TestDelegation_v0_5_latest(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}
	if runtime.GOARCH == "arm64" {
		t.Skip("skipping test on arm64")
	}

	ScenariosTest(t).
		Folder("delegation/v0_5_latest").
		Run().
		CheckNoError()
}
