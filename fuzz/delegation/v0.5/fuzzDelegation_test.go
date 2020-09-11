package delegation

import (
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/fuzz/util"
	mc "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/controller"
	"github.com/stretchr/testify/require"
)

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
			filepath.Join(getTestRoot(), "delegation/v0_5/output/delegation.wasm")).
		ReplacePath(
			"auction-mock.wasm",
			filepath.Join(getTestRoot(), "delegation/auction-mock/output/auction-mock.wasm"))

	pfe, err := newFuzzDelegationExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDelegation_v0_5(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()

	err := pfe.init(
		&fuzzDelegationExecutorInitArgs{
			serviceFee:                  r.Intn(10000),
			ownerMinStake:               0,
			minStake:                    r.Intn(1000000),
			numBlocksBeforeForceUnstake: r.Intn(1000),
			numBlocksBeforeUnbond:       r.Intn(1000),
			numDelegators:               10,
			stakePerNode:                big.NewInt(1000000000),
		},
	)
	require.Nil(t, err)

	err = pfe.increaseBlockNonce(r.Intn(10000))
	require.Nil(t, err)

	re := fuzzutil.NewRandomEventProvider()
	for stepIndex := 0; stepIndex < 500; stepIndex++ {
		generateRandomEvent(t, pfe, r, re)
	}

	// all delegators (incl. owner) withdraw all inactive stake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err := pfe.withdrawAllInactiveStake(delegatorIdx)
		require.Nil(t, err)
	}

	// all delegators (incl. owner) claim all rewards
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err := pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)
	}

	// check that delegators got all rewards out
	totalDelegatorBalance := pfe.getAllDelegatorsBalance()
	require.True(t, pfe.totalRewards.Cmp(totalDelegatorBalance) == 0,
		"Rewards don't match. Total rewards: %d. Total delegator balance: %d.",
		pfe.totalRewards, totalDelegatorBalance)

	err = pfe.increaseBlockNonce(pfe.numBlocksBeforeForceUnstake + 1)
	require.Nil(t, err)

	// all delegators (incl. owner) unStake a part of stake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		stake := big.NewInt(0).Rand(r, pfe.stakePerNode)
		err = pfe.unStake(delegatorIdx, stake)
		require.Nil(t, err)
	}

	err = pfe.increaseBlockNonce(pfe.numBlocksBeforeUnbond + 1)
	require.Nil(t, err)

	// all delegators (incl. owner) unBond
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.unBond(delegatorIdx)
		require.Nil(t, err)
	}

	// auction SC should have no more funds
	auctionBalanceAfterUnbond := pfe.getAuctionBalance()
	require.True(t, auctionBalanceAfterUnbond.Sign() == 0,
		"Auction still has balance after full unbond: %d",
		auctionBalanceAfterUnbond)

	// all delegators (incl. owner) withdraw all inactive stake
	for delegatorIdx := 0; delegatorIdx <= pfe.numDelegators; delegatorIdx++ {
		err = pfe.withdrawAllInactiveStake(delegatorIdx)
		require.Nil(t, err)
	}

	withdrawnAtTheEnd := pfe.getWithdrawTargetBalance()
	require.True(t, withdrawnAtTheEnd.Cmp(pfe.totalStakeAdded) == 0,
		"Stake added and withdrawn doesn't match. Staked: %d. Withdrawn: %d. Off by: %d",
		pfe.totalStakeAdded, withdrawnAtTheEnd,
		big.NewInt(0).Sub(pfe.totalStakeAdded, withdrawnAtTheEnd))
}

func generateRandomEvent(
	t *testing.T,
	pfe *fuzzDelegationExecutor,
	r *rand.Rand,
	re *fuzzutil.RandomEventProvider,
) {
	maxStake := big.NewInt(0).Mul(pfe.stakePerNode, big.NewInt(2))
	maxSystemReward := big.NewInt(1000000000)
	re.Reset()

	switch {
	case re.WithProbability(0.05):
		// increment block nonce
		err := pfe.increaseBlockNonce(r.Intn(1000))
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// add nodes
		err := pfe.addNodes(r.Intn(3))
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// stake
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		stake := big.NewInt(0).Rand(r, maxStake)

		err := pfe.stake(delegatorIdx, stake)
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// withdraw inactive stake
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		stake := big.NewInt(0).Rand(r, maxStake)

		err := pfe.withdrawInactiveStake(delegatorIdx, stake)
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// add system rewards
		rewards := big.NewInt(0).Rand(r, maxSystemReward)

		err := pfe.addRewards(rewards)
		require.Nil(t, err)
	case re.WithProbability(0.2):
		// claim rewards
		delegatorIdx := r.Intn(pfe.numDelegators + 1)

		err := pfe.claimRewards(delegatorIdx)
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// unStake
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		stake := big.NewInt(0).Rand(r, maxStake)

		err := pfe.unStake(delegatorIdx, stake)
		require.Nil(t, err)
	case re.WithProbability(0.05):
		// unBond
		delegatorIdx := r.Intn(pfe.numDelegators + 1)
		err := pfe.unBond(delegatorIdx)
		require.Nil(t, err)
	default:
	}
}
