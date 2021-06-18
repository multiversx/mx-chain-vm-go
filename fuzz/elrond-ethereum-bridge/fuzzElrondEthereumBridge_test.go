package elrond_ethereum_bridge

import (
	"flag"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	fuzzutil "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/fuzz/util"
	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/mandos-go/controller"
)

var fuzz = flag.Bool("fuzz", true, "Enable fuzz test")

var seedFlag = flag.Int64("seed", 0, "Random seed, use it to replay fuzz scenarios")

var iterationsFlag = flag.Int("iterations", 1000, "Number of iterations")

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../../test")
	return arwenTestRoot
}

func newExecutorWithPaths() *fuzzExecutor {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"price-aggregator.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/price-aggregator/price-aggregator.wasm"))

	fe, err := newFuzzExecutor(fileResolver)
	if err != nil {
		panic(err)
	}
	return fe
}

func TestElrondEthereumBridge(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test; only run with --fuzz argument")
	}

	fe := newExecutorWithPaths()
	defer fe.saveGeneratedScenario()

	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}
	fe.log("Random seed: %d\n", seed)
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	re := fuzzutil.NewRandomEventProvider(r)
	re.Reset()

	err := fe.initData()
	if err != nil {
		t.Error(err)
	}

	err = fe.setup(nil, nil)
	if err != nil {
		t.Error(err)
	}
}
