package delegation

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	mc "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/controller"
	"github.com/stretchr/testify/assert"
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
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"delegation.wasm",
			filepath.Join(getTestRoot(), "delegation_v0.3/delegation.wasm")).
		ReplacePath(
			"auction-mock.wasm",
			filepath.Join(getTestRoot(), "delegation_v0.3/auction-mock.wasm"))

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
	// defer executor.saveGeneratedScenario()

	err := executor.init(&fuzzDelegationExecutorInitArgs{
		serviceFee:                  r.Intn(10000),
		numBlocksBeforeForceUnstake: 0,
		numBlocksBeforeUnbond:       0,
		numDelegators:               r.Intn(10) + 1,
		stakePerNode:                big.NewInt(1000000000),
	})
	assert.Nil(t, err)
	executor.enableAutoActivation()

	maxStake := big.NewInt(0).Mul(executor.stakePerNode, big.NewInt(2))

	re := newRandomEventProvider()
	for stepIndex := 0; stepIndex < 100; stepIndex++ {
		re.reset()
		switch {
		case re.withProbability(0.2):
			// add nodes
			// err = executor.addNodes(r.Intn(3))
			// assert.Nil(t, err)
		case re.withProbability(0.5):
			// stake
			// delegatorIdx := r.Intn(executor.numDelegators)
			// stake := big.NewInt(0).Rand(r, maxStake)
			// err = executor.stake(delegatorIdx, stake)
			// assert.Nil(t, err)
		case re.withProbability(0.3):
			// unstake
			delegatorIdx := r.Intn(executor.numDelegators)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = executor.withdrawInactiveStake(delegatorIdx, stake)
			assert.Nil(t, err)
		default:
		}
		// executor.saveGeneratedScenario()
	}

}
