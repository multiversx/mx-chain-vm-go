package delegation

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

var fuzz = flag.Bool("fuzz", false, "fuzz")

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzDelegationExecutor {
	fileResolver := ij.NewDefaultFileResolver().
		ReplacePath(
			"delegation.wasm",
			filepath.Join(getTestRoot(), "delegation/delegation.wasm")).
		ReplacePath(
			"auction-mock.wasm",
			filepath.Join(getTestRoot(), "delegation/auction-mock.wasm"))

	executor, err := newFuzzDelegationExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return executor
}

func TestFuzzDelegation(t *testing.T) {
	// if !*fuzz {
	// 	t.Skip("skipping test; only run with --fuzz argument")
	// }

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	executor := newExecutorWithPaths()

	err := executor.init(&fuzzDelegationExecutorInitArgs{
		nodeShare:              r.Intn(10000),
		timeBeforeForceUnstake: 0,
		numDelegators:          r.Intn(49) + 1,
		numNodes:               r.Intn(29) + 1,
		stakePerNode:           big.NewInt(1000000000),
	})
	if err != nil {
		panic(err)
	}

	for stepIndex := 0; stepIndex < 1000 && !executor.active; stepIndex++ {
		switch {
		case rand.Float32() < 0.3:
			err = executor.maybeActivate()
			if err != nil {
				panic(err)
			}
		case rand.Float32() < 0.01:
			// finish staking, activate
			delegatorIdx := r.Intn(executor.numDelegators)
			err = executor.stakeTheRest(delegatorIdx)
			if err != nil {
				panic(err)
			}
			err = executor.mustActivate()
			if err != nil {
				panic(err)
			}
		default:
			delegatorIdx := r.Intn(executor.numDelegators)
			stake := big.NewInt(0).Rand(r, executor.expectedStake)
			err = executor.tryStake(delegatorIdx, stake)
			if err != nil {
				panic(err)
			}
		}
	}

}
