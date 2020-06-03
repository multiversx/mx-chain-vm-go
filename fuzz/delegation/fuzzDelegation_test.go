package delegation

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandosjson"
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
	fileResolver := mj.NewDefaultFileResolver().
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
		numDelegators:          r.Intn(4) + 1,
		numNodes:               r.Intn(9) + 1,
		stakePerNode:           big.NewInt(1000000000),
	})
	if err != nil {
		panic(err)
	}

	re := newRandomEventProvider()
	for stepIndex := 0; stepIndex < 1000 && !executor.active; stepIndex++ {
		re.reset()
		switch {
		case re.withProbability(0.3):
			err = executor.maybeActivate()
			if err != nil {
				panic(err)
			}
		case re.withProbability(0.01):
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
		case re.withProbability(0.5):
			// stake
			delegatorIdx := r.Intn(executor.numDelegators)
			stake := big.NewInt(0).Rand(r, executor.expectedStake)
			err = executor.tryStake(delegatorIdx, stake)
			if err != nil {
				panic(err)
			}
		case re.withProbability(0.1):
			// unstake
			delegatorIdx := r.Intn(executor.numDelegators)
			stake := big.NewInt(0).Rand(r, executor.expectedStake)
			err = executor.tryUnstake(delegatorIdx, stake)
			if err != nil {
				panic(err)
			}
		default:
		}
	}

}
