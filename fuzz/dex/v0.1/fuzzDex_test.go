package dex

import (
	"flag"
	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/fuzz/util"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/controller"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
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
			filepath.Join(getTestRoot(), "dex/v0_1/output/elrond_dex_router.wasm")).
		ReplacePath(
			"elrond_dex_pair.wasm",
			filepath.Join(getTestRoot(), "dex/v0_1/output/elrond_dex_pair.wasm")).
		ReplacePath(
			"elrond_dex_staking.wasm",
			filepath.Join(getTestRoot(), "dex/v0_1/output/elrond_dex_staking.wasm"))

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
			wegldTokenId:				"WEGLD-abcdef",
			numUsers:					10,
			numTokens:					3,
			numEvents:					15000,
			removeLiquidityProb:		0.1,
			addLiquidityProb:			0.3,
			swapProb:					0.5,
			queryPairsProb:				0.1,
			removeLiquidityMaxValue:	1000000000,
			addLiquidityMaxValue: 		1000000000,
			swapMaxValue: 				10000000,
			tokensCheckFrequency:		100,
		},
	)
	require.Nil(t, err)

	// Creating Pairs is done by users; but we'll do it ourselves,
	// since is not a matter of fuzzing (crashing or stuck funds).
	// Testing about pair creation and lp token issuing is done via mandos.
	err = pfe.createPairs()
	require.Nil(t, err)

	err = pfe.doHackishSteps()
	require.Nil(t, err)

	//Pais are created. Set fee on for each pair that has WEGLD-abcdef as a token.
	err = pfe.setFeeOn()
	require.Nil(t, err)

	stats := eventsStatistics{
		swapFixedInputHits:				0,
		swapFixedInputMisses:			0,
		swapFixedOutputHits:			0,
		swapFixedOutputMisses:			0,
		addLiquidityHits:				0,
		addLiquidityMisses:				0,
		addLiquidityPriceChecks:		0,
		removeLiquidityHits:			0,
		removeLiquidityMisses:			0,
		removeLiquidityPriceChecks: 	0,
		queryPairsHits:					0,
		queryPairsMisses:				0,
	}

	re := fuzzutil.NewRandomEventProvider()
	for stepIndex := 0; stepIndex < pfe.numEvents; stepIndex++ {
		generateRandomEvent(t, pfe, r, re, &stats)

		if stepIndex != 0 && stepIndex % pfe.tokensCheckFrequency == 0 {
			err = pfe.checkTokens()
			require.Nil(t, err)
		}
	}

	err = pfe.checkTokens()
	require.Nil(t, err)

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

	tokenA := ""
	tokenB := ""

	tokenAIndex := r.Intn(pfe.numTokens * 2) + 1
	if tokenAIndex > pfe.numTokens {
		tokenA = pfe.wegldTokenId
	} else {
		tokenA = pfe.tokenTicker(tokenAIndex)
	}
	tokenBIndex := r.Intn(pfe.numTokens) + 1
	tokenB = pfe.tokenTicker(tokenBIndex)

	userId := r.Intn(pfe.numUsers) + 1
	user := string(pfe.userAddress(userId))

	fromAtoB := r.Intn(2) != 0
	if fromAtoB == false {
		aux := tokenA
		tokenA = tokenB
		tokenB = aux
	}

	switch {
		//remove liquidity
		case re.WithProbability(pfe.removeLiquidityProb):

			seed := r.Intn(pfe.removeLiquidityMaxValue) + 1
			amount := seed
			amountAmin := seed / 100
			amountBmin := seed / 100

			err := pfe.removeLiquidity(user, tokenA, tokenB, amount, amountAmin, amountBmin, statistics)
			require.Nil(t, err)

		//add liquidity
		case re.WithProbability(pfe.addLiquidityProb):

			seed := r.Intn(pfe.addLiquidityMaxValue) + 1
			amountA := seed
			amountB := seed
			amountAmin := seed / 100
			amountBmin := seed / 100

			err := pfe.addLiquidity(user, tokenA, tokenB, amountA, amountB, amountAmin, amountBmin, statistics)
			require.Nil(t, err)

		//swap
		case re.WithProbability(pfe.swapProb):

			fixedInput := false
			amountA := 0
			amountB := 0

			fixedInput = r.Intn(2) != 0
			seed := r.Intn(pfe.swapMaxValue) + 1
			amountA = seed
			amountB = seed / 100

			if fixedInput {
				err := pfe.swapFixedInput(user, tokenA, amountA, tokenB, amountB, statistics)
				require.Nil(t, err)
			} else {
				err := pfe.swapFixedOutput(user, tokenA, amountA, tokenB, amountB, statistics)
				require.Nil(t, err)
			}

		// pair views
		case re.WithProbability(pfe.queryPairsProb):

			err := pfe.checkPairViews(user, tokenA, tokenB, statistics)
			require.Nil(t, err)

	default:
	}
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
}
