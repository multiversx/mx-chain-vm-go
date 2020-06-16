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
	"github.com/stretchr/testify/require"
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
	defer executor.saveGeneratedScenario()

	err := executor.init(&fuzzDelegationExecutorInitArgs{
		serviceFee:                  r.Intn(10000),
		numBlocksBeforeForceUnstake: 0,
		numBlocksBeforeUnbond:       0,
		numDelegators:               r.Intn(50) + 1,
		stakePerNode:                big.NewInt(1000000000),
	})
	require.Nil(t, err)
	executor.enableAutoActivation()

	maxStake := big.NewInt(0).Mul(executor.stakePerNode, big.NewInt(2))
	maxSystemReward := big.NewInt(1000000000)

	re := newRandomEventProvider()
	for stepIndex := 0; stepIndex < 10000; stepIndex++ {
		re.reset()
		switch {
		case re.withProbability(0.1):
			// add nodes
			err = executor.addNodes(r.Intn(3))
			require.Nil(t, err)
		case re.withProbability(0.3):
			// stake
			delegatorIdx := r.Intn(executor.numDelegators)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = executor.stake(delegatorIdx, stake)
			require.Nil(t, err)
		case re.withProbability(0.1):
			// withdraw inactive stake
			delegatorIdx := r.Intn(executor.numDelegators)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = executor.withdrawInactiveStake(delegatorIdx, stake)
			require.Nil(t, err)
		case re.withProbability(0.2):
			// add system rewards
			rewards := big.NewInt(0).Rand(r, maxSystemReward)
			err = executor.addRewards(rewards)
			require.Nil(t, err)
		case re.withProbability(0.1):
			// claim rewards
			delegatorIdx := r.Intn(executor.numDelegators)
			err = executor.claimRewards(delegatorIdx)
			require.Nil(t, err)
		case re.withProbability(0.2):
			// computeAllRewards
			err = executor.computeAllRewards()
		default:
		}
	}

}
