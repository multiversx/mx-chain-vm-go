package delegation

import (
	"fmt"
)

func (pfe *fuzzDelegationExecutor) init(args *fuzzDelegationExecutorInitArgs) error {
	pfe.serviceFee = args.serviceFee
	pfe.numBlocksBeforeForceUnstake = args.numBlocksBeforeForceUnstake
	pfe.numBlocksBeforeUnbond = args.numBlocksBeforeUnbond
	pfe.numDelegators = args.numDelegators
	pfe.stakePerNode = args.stakePerNode

	pfe.world.Clear()

	pfe.ownerAddress = []byte("fuzz_owner_addr_______________s1")
	pfe.delegationContractAddress = []byte("fuzz_sc_delegation_addr_______s1")
	pfe.auctionMockAddress = []byte("fuzz_sc_auction_mock_addr_____s1")
	pfe.faucetAddress = []byte("endless_sack_of_erd___________s1")
	pfe.withdrawTargetAddress = []byte("withdraw_target_______________s1")
	pfe.stakePurchaseForwardAddress = []byte("stake_purchase_forwarded______s1")

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"''%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {},
				"code": ""
			},
			"''%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {
					"''stake_per_node": "%d"
				},
				"code": "file:auction-mock.wasm"
			},
			"''%s": {
				"nonce": "0",
				"balance": "1,000,000,000,000,000,000,000,000,000,000",
				"storage": {},
				"code": ""
			},
			"''%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {},
				"code": ""
			},
			"''%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {},
				"code": ""
			}
		},
		"newAddresses": [
			{
				"creatorAddress": "''%s",
				"creatorNonce": "0",
				"newAddress": "''%s"
			}
		]
	}`,
		string(pfe.ownerAddress),
		string(pfe.auctionMockAddress),
		pfe.stakePerNode,
		string(pfe.faucetAddress),
		string(pfe.withdrawTargetAddress),
		string(pfe.stakePurchaseForwardAddress),
		string(pfe.ownerAddress),
		string(pfe.delegationContractAddress),
	))
	if err != nil {
		return err
	}

	// delegators
	for i := 1; i <= args.numDelegators; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"''%s": {
					"nonce": "0",
					"balance": "0",
					"storage": {},
					"code": ""
				}
			}
		}`,
			string(pfe.delegatorAddress(i)),
		))
		if err != nil {
			return err
		}
	}

	// deploy delegation
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scDeploy",
		"txId": "-deploy-",
		"tx": {
			"from": "''%s",
			"value": "0",
			"contractCode": "file:delegation.wasm",
			"arguments": [
				"''%s",
				"%d",
				"%d",
				"%d",
				"%d",
				"%d"
			],
			"gasLimit": "1,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.auctionMockAddress),
		args.serviceFee,
		args.ownerMinStake,
		args.numBlocksBeforeUnbond,
		args.minStake,
		args.totalDelegationCap,
	))
	if err != nil {
		return err
	}

	pfe.log("init ok")
	return nil
}
