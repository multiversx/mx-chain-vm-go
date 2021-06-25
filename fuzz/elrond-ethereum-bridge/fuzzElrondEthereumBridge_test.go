package elrond_ethereum_bridge

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
)

var fuzz = flag.Bool("fuzz", true, "Enable fuzz test")

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
			"price-aggregator.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/price-aggregator/price-aggregator.wasm")).
		ReplacePath(
			"multisig.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/multisig/multisig.wasm")).
		ReplacePath(
			"egld-esdt-swap.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/egld-esdt-swap/egld-esdt-swap.wasm")).
		ReplacePath(
			"esdt-safe.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/esdt-safe/esdt-safe.wasm")).
		ReplacePath(
			"multi-transfer-esdt.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/multi-transfer-esdt/multi-transfer-esdt.wasm")).
		ReplacePath(
			"ethereum-fee-prepay.wasm",
			filepath.Join(getTestRoot(), "elrond-ethereum-bridge/ethereum-fee-prepay/ethereum-fee-prepay.wasm"))

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

	err := fe.initData()
	if err != nil {
		t.Error(err)
	}

	// TODO: Uncomment once aggregator is integrated
	// The current version doesn't have relayer incentives and user fees

	/*
		err = fe.setupAggregator()
		if err != nil {
			t.Error(err)
		}
	*/

	nrRelayers := 2
	nrUsers := 2
	initialBalance := big.NewInt(INIT_BALANCE)
	err = fe.initAccounts(nrRelayers, nrUsers, initialBalance)
	if err != nil {
		t.Error(err)
	}

	multisigInitArgs := MultisigInitArgs{
		requiredStake: big.NewInt(1000),
		slashAmount:   big.NewInt(500),
		quorum:        len(fe.data.actorAddresses.relayers) / 2,
		boardMembers:  fe.data.actorAddresses.relayers,
	}
	err = fe.deployMultisig(&multisigInitArgs)
	if err != nil {
		t.Error(err)
	}

	deployChildContractsArgs := DeployChildContractsArgs{
		egldEsdtSwapCodePath:      "file:egld-esdt-swap.wasm",
		multiTransferEsdtCodePath: "file:multi-transfer-esdt.wasm",
		ethereumFeePrepayCodePath: "file:ethereum-fee-prepay.wasm",
		esdtSafeCodePath:          "file:esdt-safe.wasm",
		priceAggregatorAddress:    "sc:price-aggregator",
		wrappedEgldTokenId:        "str:WEGLD-123456",
		wrappedEthTokenId:         "str:WETH-abcdef",
		tokenWhitelist:            []string{},
	}
	err = fe.setupChildContracts(&deployChildContractsArgs)
	if err != nil {
		t.Error(err)
	}

	var seed int64
	if *seedFlag == 0 {
		seed = time.Now().UnixNano()
	} else {
		seed = *seedFlag
	}
	fe.log("Random seed: %d\n", seed)
	r := rand.New(rand.NewSource(seed))
	r.Seed(seed)

	fe.randSource = *r

	re := fuzzutil.NewRandomEventProvider(r)
	for stepIndex := 0; stepIndex < *iterationsFlag; stepIndex++ {
		re.Reset()

		switch {
		case re.WithProbability(0.75):
			userAcc := fe.getRandomUser()
			wrapAmount := big.NewInt(int64(fe.randSource.Intn(100) + 1))

			err = fe.wrapEgld(userAcc, wrapAmount)
			if err != nil {
				t.Error(err)
			}
		case re.WithProbability(0.25):
			userAcc := fe.getRandomUser()
			userWrappedEgldBalance := fe.getEsdtBalance(userAcc, string(fe.interpretExpr(fe.data.wrappedEgldTokenId)))

			if userWrappedEgldBalance.Cmp(big.NewInt(0)) == 0 {
				continue
			}

			unwrapAmount := big.NewInt(int64(fe.randSource.Intn(int(userWrappedEgldBalance.Int64())) + 1))

			err = fe.wrapEgld(userAcc, unwrapAmount)
			if err != nil {
				t.Error(err)
			}
		default:
		}
	}
}
