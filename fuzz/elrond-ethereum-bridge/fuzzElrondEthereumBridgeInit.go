package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
)

type MultisigInitArgs struct {
	requiredStake *big.Int
	slashAmount   *big.Int
	quorum        int
	boardMembers  []Address
}

type DeployChildContractsArgs struct {
	egldEsdtSwapCode       []byte
	multiTransferEsdtCode  []byte
	ethereumFeePrepayCode  []byte
	esdtSafeCode           []byte
	priceAggregatorAddress Address
	wrappedEgldTokenId     TokenIdentifier
	wrappedEthTokenId      TokenIdentifier
	tokenWhitelist         []TokenIdentifier
}

func (fe *fuzzExecutor) initData() error {
	fe.data = &fuzzData{
		actorAddresses: &ActorAddresses{
			accounts:          []Address{"address:owner"},
			multisig:          "sc:multisig",
			priceAggregator:   "sc:price_aggregator",
			egldEsdtSwap:      "sc:egld_esdt_swap",
			esdtSafe:          "sc:esdt_safe",
			ethereumFeePrepay: "sc:ethereum_fee_prepay",
			multiTransferEsdt: "sc:multi_transfer_esdt",
		},
		egldEsdtSwapState:      nil,
		esdtSafeState:          nil,
		ethereumFeePrepayState: nil,
		multiTransferEsdtState: nil,
		multisigState:          nil,
	}
	fe.world.Clear()

	return nil
}

func (fe *fuzzExecutor) setup(
	multisigInitArgs *MultisigInitArgs,
	deployChildContractsArgs *DeployChildContractsArgs) error {

	err := fe.setupAggregator()
	if err != nil {
		return err
	}

	/*
		err = fe.executeStep(fmt.Sprintf(`
			{
				"step": "setState",
				"accounts": {
					"%s": {
						"nonce": "0",
						"balance": "0",
						"storage": {}
					}
				}
			}`,
			fe.data.actorAddresses.owner,
		))
		if err != nil {
			return err
		}
	*/

	return nil
}

func (fe *fuzzExecutor) setupAggregator() error {
	err := fe.executeStep(fmt.Sprintf(`
		{
			"step": "externalSteps",
			"path": "%s"
		}`,
		"/home/elrond/arwen-wasm-vm/test/elrond-ethereum-bridge/price-aggregator/mandos/oracle_submit.scen.json",
	))
	if err != nil {
		return err
	}

	return nil
}
