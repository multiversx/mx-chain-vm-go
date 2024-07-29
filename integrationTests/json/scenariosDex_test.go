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
