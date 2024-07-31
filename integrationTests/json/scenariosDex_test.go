package vmjsonintegrationtest

import (
	"testing"
)

func TestDexClaimProgress(t *testing.T) {
	ScenariosTest(t).
		Folder("dex/scenarios/trace1.scen.json").
		WithConsoleExecutorLogs().
		WithExecutorLogs().
		Run().
		CheckNoError()
}

func TestDexClaimProgressFromUser(t *testing.T) {
	ScenariosTest(t).
		Folder("dex/scenarios/trace2.scen.json").
		WithConsoleExecutorLogs().
		WithExecutorLogs().
		Run().
		CheckNoError()
}
