package fuzzForwarder

import (
	"flag"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/fuzz/util"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/controller"
	"github.com/stretchr/testify/require"
)

var fuzz = flag.Bool("fuzz", false, "Enable fuzz test")

var seedFlag = flag.Int64("seed", 0, "Random seed, use it to replay fuzz scenarios")

var iterationsFlag = flag.Int("iterations", 10, "Number of iterations")

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"forwarder.wasm",
			filepath.Join(getTestRoot(), "features/async/forwarder/output/forwarder.wasm"))

	pfe, err := newFuzzExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return pfe
}

func TestFuzzForwarder(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	pfe := newExecutorWithPaths()
	defer pfe.saveGeneratedScenario()

	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}
	pfe.log("Random seed: %d\n", seed)
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	err := pfe.initData()
	require.Nil(t, err)

	err = pfe.setUp()
	require.Nil(t, err)

	for stepIndex := 0; stepIndex < *iterationsFlag; stepIndex++ {
		fromIndex := r.Intn(pfe.data.numForwarders) + 1
		toIndex := r.Intn(pfe.data.numForwarders) + 1
		err = pfe.programCall(fromIndex, toIndex, "EGLD", 0, "1000")
		require.Nil(t, err)
	}

	// err = pfe.programCall(1, 1, "EGLD", 0, "1000")
	// require.Nil(t, err)

	// err = pfe.programCall(2, 3, "EGLD", 0, "1000")
	// require.Nil(t, err)

	// err = pfe.programCall(3, 4, "EGLD", 0, "1000")
	// require.Nil(t, err)

	err = pfe.executeCall(1)
	require.Nil(t, err)
}

func generateRandomEvent(
	t *testing.T,
	pfe *fuzzExecutor,
	r *rand.Rand,
	re *fuzzutil.RandomEventProvider,
	maxDelegationCap *big.Int,
) {

	re.Reset()

	switch {
	case re.WithProbability(0.9):
		// increment block nonce
		// err := pfe.increaseBlockNonce(r.Intn(1000))
		// require.Nil(t, err)

		// pfe.checkInvariants(t)
	default:
	}
}
