package delegation

import (
	"flag"
	"fmt"
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

	pfe, err := newFuzzDelegationExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDelegation(t *testing.T) {
	// if !*fuzz {
	// 	t.Skip("skipping test; only run with --fuzz argument")
	// }

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()

	err := pfe.init(&fuzzDelegationExecutorInitArgs{
		serviceFee:                  r.Intn(10000),
		numBlocksBeforeForceUnstake: 0,
		numBlocksBeforeUnbond:       0,
		numDelegators:               100,
		stakePerNode:                big.NewInt(1000000000),
	})
	require.Nil(t, err)
	pfe.enableAutoActivation()

	maxStake := big.NewInt(0).Mul(pfe.stakePerNode, big.NewInt(2))
	maxSystemReward := big.NewInt(1000000000)

	re := newRandomEventProvider()
	for stepIndex := 0; stepIndex < 1000; stepIndex++ {
		re.reset()
		switch {
		case re.withProbability(0.1):
			// add nodes
			err = pfe.addNodes(r.Intn(3))
			require.Nil(t, err)
		case re.withProbability(0.2):
			// stake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = pfe.stake(delegatorIdx, stake)
			require.Nil(t, err)
		case re.withProbability(0.1):
			// withdraw inactive stake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			stake := big.NewInt(0).Rand(r, maxStake)
			err = pfe.withdrawInactiveStake(delegatorIdx, stake)
			require.Nil(t, err)
		case re.withProbability(0.2):
			// add system rewards
			rewards := big.NewInt(0).Rand(r, maxSystemReward)
			err = pfe.addRewards(rewards)
			require.Nil(t, err)
		case re.withProbability(0.15):
			// claim rewards
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			err = pfe.claimRewards(delegatorIdx)
			require.Nil(t, err)
		case re.withProbability(0.05):
			// computeAllRewards
			err = pfe.computeAllRewards()
			require.Nil(t, err)
		case re.withProbability(0.1):
			// announceUnStake
			delegatorIdx := r.Intn(pfe.numDelegators + 1)
			amount := big.NewInt(0).Rand(r, maxStake)
			err = pfe.announceUnStake(delegatorIdx, amount)
			require.Nil(t, err)
		default:
		}
	}

	err = pfe.checkContractBalanceVsState()
	if err != nil {
		fmt.Println(err)
		return
	}

	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.withdrawAllInactiveStake(delegatorIdx)
		require.Nil(t, err)
	}

	fmt.Println(pfe.getAllDelegatorsBalance())

	err = pfe.computeAllRewards()
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)

		err = pfe.checkContractBalanceVsState()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	totalDelegatorBalance := pfe.getAllDelegatorsBalance()

	fmt.Println(pfe.getContractBalance())
	fmt.Println()
	fmt.Println(pfe.getAuctionBalance())
	fmt.Println(pfe.getWithdrawTargetBalance())
	require.True(t, pfe.totalRewards.Cmp(totalDelegatorBalance) == 0,
		"Rewards don't match. Total rewards: %d. Total delegator balance: %d.",
		pfe.totalRewards, totalDelegatorBalance)

}
