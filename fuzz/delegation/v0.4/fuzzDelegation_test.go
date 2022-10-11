package delegation

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/wasm-vm-v1_4/fuzz/util"
	mc "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/controller"
	"github.com/stretchr/testify/require"
)

var fuzz = flag.Bool("fuzz", false, "fuzz")

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzDelegationExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"delegation.wasm",
			filepath.Join(getTestRoot(), "delegation/v0_4_genesis/output/delegation.wasm")).
		ReplacePath(
			"auction-mock.wasm",
			filepath.Join(getTestRoot(), "delegation/auction-mock/output/auction-mock.wasm"))

	pfe, err := newFuzzDelegationExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDelegation_v0_4_0_genesis(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()

	err := pfe.init(&fuzzDelegationExecutorInitArgs{
		serviceFee:                  r.Intn(10000),
		ownerMinStake:               0,
		numBlocksBeforeForceUnstake: r.Intn(1000),
		numBlocksBeforeUnbond:       r.Intn(1000),
		numDelegators:               10,
		stakePerNode:                big.NewInt(1000000000),
	})
	require.Nil(t, err)

	maxStake := big.NewInt(0).Mul(pfe.stakePerNode, big.NewInt(2))
	maxSystemReward := big.NewInt(10)

	// add nodes
	numNodesAtGenesis := 10 + r.Intn(190)
	_ = pfe.addNodes(numNodesAtGenesis)

	// stake genesis
	genesisStake := big.NewInt(0).Mul(pfe.stakePerNode, big.NewInt(int64(numNodesAtGenesis)))
	for stepIndex := 0; stepIndex < 30; stepIndex++ {
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		stake := big.NewInt(0).Rand(r, maxStake)
		nextTotalStake := big.NewInt(0).Add(pfe.totalStakeAdded, stake)
		if nextTotalStake.Cmp(genesisStake) <= 0 {
			err = pfe.stakeGenesis(delegatorIdx, stake)
			require.Nil(t, err)
		}
	}

	if pfe.totalStakeAdded.Cmp(genesisStake) < 0 {
		// stake the remainder
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		remainingStake := big.NewInt(0).Sub(genesisStake, pfe.totalStakeAdded)
		err = pfe.stakeGenesis(delegatorIdx, remainingStake)
		require.Nil(t, err)
	}

	// activate genesis
	err = pfe.activateGenesis()
	require.Nil(t, err)

	// after genesis
	_ = pfe.increaseBlockNonce(r.Intn(10000))
	re := fuzzutil.NewRandomEventProvider(r)

	for stepIndex := 0; stepIndex < 500; stepIndex++ {
		re.Reset()
		switch {
		case re.WithProbability(0.2):
			// add system rewards
			rewards := big.NewInt(0).Rand(r, maxSystemReward)
			err = pfe.addRewards(rewards)
			require.Nil(t, err)
		case re.WithProbability(0.2):
			// claim rewards
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			err = pfe.claimRewards(delegatorIdx)
			require.Nil(t, err)
		case re.WithProbability(0.1):
			// computeAllRewards
			err = pfe.computeAllRewards()
			require.Nil(t, err)

		default:
		}
	}

	// all delegators (incl. owner) claim all rewards
	err = pfe.computeAllRewards()
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)
	}

	require.True(t, pfe.getContractBalance().Sign() == 0)

}
