package dex

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/fuzz/util"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
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

func newExecutorWithPaths() *fuzzDexExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"elrond_dex_router.wasm",
			filepath.Join(getTestRoot(), "dex/v0_2/output/elrond_dex_router.wasm")).
		ReplacePath(
			"elrond_dex_pair.wasm",
			filepath.Join(getTestRoot(), "dex/v0_2/output/elrond_dex_pair.wasm")).
		ReplacePath(
			"elrond_dex_farm.wasm",
			filepath.Join(getTestRoot(), "dex/v0_2/output/elrond_dex_farm.wasm"))

	pfe, err := newFuzzDexExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDelegation_v0_5(t *testing.T) {
	//if !*fuzz {
	//	t.Skip("skipping test; only run with --fuzz argument")
	//}

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	err := pfe.init(
		&fuzzDexExecutorInitArgs{
			wegldTokenId:            "WEGLD-abcdef",
			mexTokenId:              "MEX-abcdef",
			busdTokenId:			 "BUSD-abcdef",
			wemeLpTokenId:           "WEMELP-abcdef",
			webuLpTokenId:           "WEBULP-abcdef",
			wemeFarmTokenId:		 "WEMEFARM-abcdef",
			webuFarmTokenId:		 "WEBUFARM-abcdef",
			mexFarmTokenId: 		 "MEXFARM-abcdef",
			numUsers:                10,
			numEvents:               1000,
			removeLiquidityProb:     0.05,
			addLiquidityProb:        0.20,
			swapProb:                0.35,
			queryPairsProb:          0.05,
			enterFarmProb:           0.20,
			exitFarmProb:            0.12,
			increaseEpochProb:       0.02,
			removeLiquidityMaxValue: 1000000000,
			addLiquidityMaxValue:    1000000000,
			swapMaxValue:            10000000,
			enterFarmMaxValue:       100000000,
			exitFarmMaxValue:        100000000,
			blockEpochIncrease:      10,
		},
	)
	require.Nil(t, err)

	stats := eventsStatistics{
		swapFixedInputHits:         0,
		swapFixedInputMisses:       0,
		swapFixedOutputHits:        0,
		swapFixedOutputMisses:      0,
		addLiquidityHits:           0,
		addLiquidityMisses:         0,
		addLiquidityPriceChecks:    0,
		removeLiquidityHits:        0,
		removeLiquidityMisses:      0,
		removeLiquidityPriceChecks: 0,
		queryPairsHits:             0,
		queryPairsMisses:           0,
		enterFarmHits:              0,
		enterFarmMisses:            0,
		exitFarmHits:               0,
		exitFarmMisses:             0,
		exitFarmWithRewards:        0,
	}

	re := fuzzutil.NewRandomEventProvider(r)
	for stepIndex := 0; stepIndex < pfe.numEvents; stepIndex++ {
		generateRandomEvent(t, pfe, r, re, &stats)
	}

	printStatistics(&stats, pfe)
}

func generateRandomEvent(
	t *testing.T,
	pfe *fuzzDexExecutor,
	r *rand.Rand,
	re *fuzzutil.RandomEventProvider,
	statistics *eventsStatistics,
) {
	re.Reset()

	userId := r.Intn(pfe.numUsers) + 1
	user := pfe.userAddress(userId)

	switch {
	//remove liquidity
	case re.WithProbability(pfe.removeLiquidityProb):

		swapPair := getRandomSwapPair(r, pfe)
		seed := r.Intn(pfe.removeLiquidityMaxValue) + 1
		amount := seed
		amountAmin := seed / 100
		amountBmin := seed / 100

		err := pfe.removeLiquidity(user, swapPair, amount, amountAmin, amountBmin, statistics)
		require.Nil(t, err)

	//add liquidity
	case re.WithProbability(pfe.addLiquidityProb):

		swapPair := getRandomSwapPair(r, pfe)
		seed := r.Intn(pfe.addLiquidityMaxValue) + 1
		amountA := seed
		amountB := seed
		amountAmin := seed / 100
		amountBmin := seed / 100

		err := pfe.addLiquidity(user, swapPair, amountA, amountB, amountAmin, amountBmin, statistics)
		require.Nil(t, err)

	//swap
	case re.WithProbability(pfe.swapProb):

		swapPair := getRandomSwapPair(r, pfe)
		fixedInput := false
		amountA := 0
		amountB := 0
		fixedInput = r.Intn(2) != 0
		seed := r.Intn(pfe.swapMaxValue) + 1
		amountA = seed
		amountB = seed / 100

		if fixedInput {
			err := pfe.swapFixedInput(user, swapPair, amountA, amountB, statistics)
			require.Nil(t, err)
		} else {
			err := pfe.swapFixedOutput(user, swapPair, amountA, amountB, statistics)
			require.Nil(t, err)
		}

	// pair views
	case re.WithProbability(pfe.queryPairsProb):

		swapPair := getRandomSwapPair(r, pfe)
		err := pfe.checkPairViews(user, swapPair, statistics)
		require.Nil(t, err)

	// enterFarm
	case re.WithProbability(pfe.enterFarmProb):

		amount := r.Intn(pfe.enterFarmMaxValue) + 1
		farm := getRandomFarm(r, pfe)
		err := pfe.enterFarm(user, farm, amount, statistics)
		require.Nil(t, err)

	// exitFarm
	case re.WithProbability(pfe.exitFarmProb):

		amount := r.Intn(pfe.removeLiquidityMaxValue) + 1
		err := pfe.exitFarm(amount, statistics, r)
		require.Nil(t, err)

	// increase block epoch. required for unbond
	case re.WithProbability(pfe.increaseEpochProb):

		err := pfe.increaseBlockEpoch(pfe.blockEpochIncrease)
		require.Nil(t, err)
	default:
	}
}

func getRandomSwapPair(r *rand.Rand, pfe *fuzzDexExecutor) SwapPair {
	seed := r.Intn(2)
	randomPair := SwapPair{}

	if seed == 0 {
		randomPair.address = pfe.wemeSwapAddress
		randomPair.lpToken = pfe.wemeLpTokenId
		randomPair.firstToken = pfe.wegldTokenId
		randomPair.secondToken = pfe.mexTokenId
	} else {
		randomPair.address = pfe.webuSwapAddress
		randomPair.lpToken = pfe.webuLpTokenId
		randomPair.firstToken = pfe.wegldTokenId
		randomPair.secondToken = pfe.busdTokenId
	}

	return randomPair
}

func getRandomFarm(r *rand.Rand, pfe *fuzzDexExecutor) Farm {
	seed := r.Intn(3)
	randomFarm := Farm{}

	if seed == 0 {
		randomFarm.address = pfe.wemeFarmAddress
		randomFarm.farmToken = pfe.wemeFarmTokenId
		randomFarm.farmingToken = pfe.wemeLpTokenId
		randomFarm.rewardToken = pfe.mexTokenId
	} else if seed == 1 {
		randomFarm.address = pfe.webuFarmAddress
		randomFarm.farmToken = pfe.webuFarmTokenId
		randomFarm.farmingToken = pfe.webuLpTokenId
		randomFarm.rewardToken = pfe.mexTokenId
	} else {
		randomFarm.address = pfe.mexFarmAddress
		randomFarm.farmToken = pfe.mexFarmTokenId
		randomFarm.farmingToken = pfe.mexTokenId
		randomFarm.rewardToken = pfe.mexTokenId
	}

	return randomFarm
}

func printStatistics(statistics *eventsStatistics, pfe *fuzzDexExecutor) {
	pfe.log("\nStatistics:")
	pfe.log("\tswapFixedInputHits			%d", statistics.swapFixedInputHits)
	pfe.log("\tswapFixedInputMisses		%d", statistics.swapFixedInputMisses)
	pfe.log("")
	pfe.log("\tswapFixedOutputHits			%d", statistics.swapFixedOutputHits)
	pfe.log("\tswapFixedOutputMissed		%d", statistics.swapFixedOutputMisses)
	pfe.log("")
	pfe.log("\taddLiquidityHits			%d", statistics.addLiquidityHits)
	pfe.log("\taddLiquidityMisses			%d", statistics.addLiquidityMisses)
	pfe.log("\taddLiquidityPriceChecks 	%d", statistics.addLiquidityPriceChecks)
	pfe.log("")
	pfe.log("\tremoveLiquidityHits			%d", statistics.removeLiquidityHits)
	pfe.log("\tremoveLiquidityMisses		%d", statistics.removeLiquidityMisses)
	pfe.log("\tremoveLiquidityPriceChecks	%d", statistics.removeLiquidityPriceChecks)
	pfe.log("")
	pfe.log("\tqueryPairHits				%d", statistics.queryPairsHits)
	pfe.log("\tqueryPairMisses				%d", statistics.queryPairsMisses)
	pfe.log("")
	pfe.log("\tenterFarmHits				%d", statistics.enterFarmHits)
	pfe.log("\tenterFarmMisses				%d", statistics.enterFarmMisses)
	pfe.log("")
	pfe.log("\texitFarmHits				%d", statistics.exitFarmHits)
	pfe.log("\texitFarmMisses				%d", statistics.exitFarmMisses)
	pfe.log("\texitFarmWithRewards			%d", statistics.exitFarmWithRewards)
	pfe.log("")
}
