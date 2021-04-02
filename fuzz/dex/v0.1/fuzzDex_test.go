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
			wegldTokenId: "WEGLD-abcdef",
			numUsers: 10,
			numTokens: 5,
		},
	)
	require.Nil(t, err)

	// Creating Pairs is done by users; but we'll do it ourselves,
	// since is not a matter of fuzzing (crashing or stuck funds).
	// Testing about pair creation and lp token issuing is done via mandos.
	err = pfe.createPairs()
	require.Nil(t, err)

	//Pais are created. Set fee on for each pair that has WEGLD-abcdef as a token.
	err = pfe.setFeeOn()
	require.Nil(t, err)

	err = pfe.increaseBlockNonce(r.Intn(10000))
	require.Nil(t, err)

	re := fuzzutil.NewRandomEventProvider()
	for stepIndex := 0; stepIndex < 2500; stepIndex++ {
		generateRandomEvent(t, pfe, r, re)
	}
}


func generateRandomEvent(
	t *testing.T,
	pfe *fuzzDexExecutor,
	r *rand.Rand,
	re *fuzzutil.RandomEventProvider,
) {
	re.Reset()

	tokenA := ""
	tokenB := ""

	tokenAIndex := r.Intn(pfe.numTokens * 2)
	if tokenAIndex > pfe.numTokens {
		tokenA = pfe.wegldTokenId
	} else {
		tokenA = pfe.tokenTicker(tokenAIndex)
	}
	tokenBIndex := r.Intn(pfe.numTokens)
	tokenB = pfe.tokenTicker(tokenBIndex)

	userId := r.Intn(pfe.numUsers) + 1
	user := string(pfe.userAddress(userId))

	switch {
		//swap
		case re.WithProbability(0.4):

			fixedInput := false
			fromAtoB := false
			amountA := 0
			amountB := 0

			fromAtoB = r.Intn(2) != 0
			if fromAtoB == false {
				aux := tokenA
				tokenA = tokenB
				tokenB = aux
			}

			fixedInput = r.Intn(2) != 0
			seed := r.Intn(1000000000000)
			amountA = seed
			amountB = seed / 10000

			if fixedInput {
				err := pfe.swapFixedInput(user, tokenA, amountA, tokenB, amountB)
				require.Nil(t, err)
			} else {
				err := pfe.swapFixedOutput(user, tokenA, amountA, tokenB, amountB)
				require.Nil(t, err)
			}

		//add liquidity
		case re.WithProbability(0.2):

			seed := r.Intn(1000000000000)
			amountA := seed
			amountB := seed
			amountAmin := seed / 10000
			amountBmin := seed / 10000

			err := pfe.addLiquidity(user, tokenA, tokenB, amountA, amountB, amountAmin, amountBmin)
			require.Nil(t, err)

		//remove liquidity
		case re.WithProbability(0.2):

			seed := r.Intn(10000000000)
			amount := seed
			amountAmin := seed / 10000
			amountBmin := seed / 10000

			err := pfe.removeLiquidity(user, tokenA, tokenB, amount, amountAmin, amountBmin)
			require.Nil(t, err)

	default:
	}
}
