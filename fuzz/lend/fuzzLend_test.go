package lendFuzz

import (
	"flag"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	mandosController "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/controller"
)

var fuzz = flag.Bool("fuzz", false, "enable lend fuzzer")

var seedFlag = flag.Int64("seed", 0, "random seed use it to replay fuzz scenarios")

func newExecutor() *fuzzLendExecutor {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fileResolver := mandosController.NewDefaultFileResolver().
		ReplacePath("liquidity_pool.wasm", filepath.Join(pwd, "wasm/liquidity_pool.wasm"))

	executor, err := newFuzzLendExecutor(fileResolver)
	if err != nil {
		panic(err)
	}

	return executor
}

func TestFuzzLend(t *testing.T) {
	if !*fuzz {
		t.Skip("skipping test - missing --fuzz arg")
	}

	executor := newExecutor()
	defer executor.saveGeneratedScenario()

	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}

	log.Println("random seed: ", seed)
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	err := executor.init(&fuzzLendExecutorArgs{
		wegldTokenID:  wegld,
		lwegldTokenID: lwegld,
		bwegldTokenID: bwegld,
		busdTokenID:   busd,
		lbusdTokenID:  lbusd,
		bbusdTokenID:  bbusd,
		numUsers:      10,
		numEvents:     500,
	})
	require.Nil(t, err)
}

func generateRandomEvent(
	t *testing.T,
	executor *fuzzLendExecutor,
	r *rand.Rand,
	stats *statistics,
) {

}
