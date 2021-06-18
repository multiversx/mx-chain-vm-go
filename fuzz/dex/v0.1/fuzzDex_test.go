package dex

import (
	"errors"
	"flag"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/controller"
	"github.com/stretchr/testify/require"
)

func newExecutorWithPaths() *fuzzDexExecutor {
	pwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"elrond_dex_router.wasm",
			filepath.Join(pwd, "wasms/elrond_dex_router.wasm")).
		ReplacePath(
			"elrond_dex_pair.wasm",
			filepath.Join(pwd, "wasms/elrond_dex_pair.wasm")).
		ReplacePath(
			"elrond_dex_farm.wasm",
			filepath.Join(pwd, "wasms/elrond_dex_farm.wasm"))

	pfe, err := newFuzzDexExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzDex_v0_1(t *testing.T) {
	if !*flag.Bool("fuzz", false, "fuzz") {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	err := pfe.init(
		&fuzzDexExecutorInitArgs{
			wegldTokenId:            "WEGLD-abcdef",
			mexTokenId:              "MEX-abcdef",
			busdTokenId:             "BUSD-abcdef",
			wemeLpTokenId:           "WEMELP-abcdef",
			webuLpTokenId:           "WEBULP-abcdef",
			wemeFarmTokenId:         "WEMEFARM-abcdef",
			webuFarmTokenId:         "WEBUFARM-abcdef",
			mexFarmTokenId:          "MEXFARM-abcdef",
			numUsers:                10,
			numEvents:               500,
			removeLiquidityProb:     5,
			addLiquidityProb:        20,
			swapProb:                25,
			queryPairsProb:          5,
			enterFarmProb:           18,
			exitFarmProb:            6,
			claimRewardsProb:        20,
			removeLiquidityMaxValue: 1000000000,
			addLiquidityMaxValue:    1000000000,
			swapMaxValue:            10000000,
			enterFarmMaxValue:       100000000,
			exitFarmMaxValue:        100000000,
			claimRewardsMaxValue:    10000000,
			blockNonceIncrease:      1,
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
		claimRewardsHits:           0,
		claimRewardsMisses:         0,
		claimRewardsWithRewards:    0,
	}

	for stepIndex := 0; stepIndex < pfe.numEvents; stepIndex++ {
		generateRandomEvent(t, pfe, r, &stats)
		err := pfe.increaseBlockNonce(1)
		require.Nil(t, err)
	}

	printStatistics(&stats, pfe)
}

func generateRandomEvent(
	t *testing.T,
	pfe *fuzzDexExecutor,
	r *rand.Rand,
	statistics *eventsStatistics,
) {
	events := map[string]int{
		"removeLiquidity": pfe.removeLiquidityProb,
		"addLiquidity":    pfe.addLiquidityProb,
		"swap":            pfe.swapProb,
		"checkPairViews":  pfe.queryPairsProb,
		"enterFarm":       pfe.enterFarmProb,
		"exitFarm":        pfe.exitFarmProb,
		"claimRewards":    pfe.claimRewardsProb,
	}

	event, err := weighted_random_choice(r, events)
	assert.Nil(t, err)

	switch event {
	case "removeLiquidity":
		err = pfe.removeLiquidity(r, statistics)
		assert.Nil(t, err)
	case "addLiquidity":
		err = pfe.addLiquidity(r, statistics)
		assert.Nil(t, err)
	case "swap":
		err = pfe.swap(r, statistics)
		assert.Nil(t, err)
	case "checkPairViews":
		err = pfe.checkPairViews(r, statistics)
		assert.Nil(t, err)
	case "enterFarm":
		err = pfe.enterFarm(r, statistics)
		assert.Nil(t, err)
	case "exitFarm":
		err = pfe.exitFarm(r, statistics)
		assert.Nil(t, err)
	case "claimRewards":
		err = pfe.claimRewards(r, statistics)
		assert.Nil(t, err)
	}
}

func weighted_random_choice(r *rand.Rand, choices map[string]int) (string, error) {
	sumOfWeight := 0
	for _, v := range choices {
		sumOfWeight = sumOfWeight + v
	}

	pick := r.Intn(sumOfWeight)

	current := 0
	for k, v := range choices {
		current += v
		if current > pick {
			return k, nil
		}
	}

	return "", errors.New("no event")
}

func printStatistics(statistics *eventsStatistics, pfe *fuzzDexExecutor) {
	pfe.log("\nStatistics:")
	pfe.log("\tswapFixedInputHits			%d", statistics.swapFixedInputHits)
	pfe.log("\tswapFixedInputMisses			%d", statistics.swapFixedInputMisses)
	pfe.log("")
	pfe.log("\tswapFixedOutputHits			%d", statistics.swapFixedOutputHits)
	pfe.log("\tswapFixedOutputMissed		%d", statistics.swapFixedOutputMisses)
	pfe.log("")
	pfe.log("\taddLiquidityHits				%d", statistics.addLiquidityHits)
	pfe.log("\taddLiquidityMisses			%d", statistics.addLiquidityMisses)
	pfe.log("\taddLiquidityPriceChecks 		%d", statistics.addLiquidityPriceChecks)
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
	pfe.log("\texitFarmHits					%d", statistics.exitFarmHits)
	pfe.log("\texitFarmMisses				%d", statistics.exitFarmMisses)
	pfe.log("\texitFarmWithRewards			%d", statistics.exitFarmWithRewards)
	pfe.log("")
	pfe.log("\tclaimRewardsHits				%d", statistics.claimRewardsHits)
	pfe.log("\tclaimRewardsMisses			%d", statistics.claimRewardsMisses)
	pfe.log("\tclaimRewardsWithRewards		%d", statistics.claimRewardsWithRewards)
	pfe.log("")
}
