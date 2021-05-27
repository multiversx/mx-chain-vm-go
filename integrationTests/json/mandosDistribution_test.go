package vmjsonintegrationtest

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"testing"
)

func TestDistribution_v0_1(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "distribution/v0_1")
}

func TestDistribution_v0_1_single(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	arwen.SetLoggingForTests()
	runSingleTest(t, "distribution/v0_1/mandos", "claim_mex_rewards_proxy_after_mint_rewards.scen.json")
}
