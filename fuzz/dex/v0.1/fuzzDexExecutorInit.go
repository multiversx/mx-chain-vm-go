package dex

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
)

func (pfe *fuzzDexExecutor) init(args *fuzzDexExecutorInitArgs) error {
	pfe.wegldTokenId = args.wegldTokenId
	pfe.mexTokenId = args.mexTokenId
	pfe.numTokens = args.numTokens
	pfe.numUsers = args.numUsers
	pfe.numEvents = args.numEvents
	pfe.removeLiquidityProb = args.removeLiquidityProb
	pfe.addLiquidityProb = args.addLiquidityProb
	pfe.swapProb = args.swapProb
	pfe.queryPairsProb = args.queryPairsProb
	pfe.enterFarmProb = args.enterFarmProb
	pfe.exitFarmProb = args.exitFarmProb
	pfe.unbondProb = args.unbondProb
	pfe.increaseEpochProb = args.increaseEpochProb
	pfe.removeLiquidityMaxValue = args.removeLiquidityMaxValue
	pfe.addLiquidityMaxValue = args.addLiquidityMaxValue
	pfe.swapMaxValue = args.swapMaxValue
	pfe.enterFarmMaxValue = args.enterFarmMaxValue
	pfe.exitFarmMaxValue = args.exitFarmMaxValue
	pfe.blockEpochIncrease = args.blockEpochIncrease
	pfe.tokensCheckFrequency = args.tokensCheckFrequency
	pfe.farmers = make(map[int]FarmerInfo)

	pfe.world.Clear()

	pfe.ownerAddress = []byte("fuzz_owner_addr_______________s1")
	pfe.routerAddress = []byte("fuzz_dex_router_addr__________s1")
	pfe.wegldFarmingAddress = []byte("fuzz_dex_wegld_farming_addr___s1")
	pfe.mexFarmingAddress = []byte("fuzz_dex_mex_farming_addr_____s1")

	// users
	esdtString := pfe.fullOfEsdtWalletString()
	for i := 1; i <= args.numUsers; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"''%s": {
					"nonce": "0",
					"balance": "0",
					"storage": {},
					"esdt": {
						%s
					},
					"code": ""
				}
			}
		}`,
			string(pfe.userAddress(i)),
			esdtString,
		))
		if err != nil {
			return err
		}
	}

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"''%s": {
				"nonce": "0",
				"balance": "1,000,000,000,000,000,000,000,000,000,000",
				"storage": {},
				"code": ""
			}
		},
		"newAddresses": [
			{
				"creatorAddress": "''%s",
				"creatorNonce": "0",
				"newAddress": "''%s"
			},
			{
				"creatorAddress": "''%s",
				"creatorNonce": "1",
				"newAddress": "''%s"
			},
			{
				"creatorAddress": "''%s",
				"creatorNonce": "2",
				"newAddress": "''%s"
			}
		]
	}`,
		string(pfe.ownerAddress),
		string(pfe.ownerAddress),
		string(pfe.routerAddress),
		string(pfe.ownerAddress),
		string(pfe.wegldFarmingAddress),
		string(pfe.ownerAddress),
		string(pfe.mexFarmingAddress),
	))
	if err != nil {
		return err
	}

	// deploy router
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scDeploy",
		"txId": "-router-deploy-",
		"tx": {
			"from": "''%s",
			"value": "0",
			"contractCode": "file:elrond_dex_router.wasm",
			"arguments": [
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
	))
	if err != nil {
		return err
	}

	// deploy wegld farming
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scDeploy",
		"txId": "-farm-deploy-",
		"tx": {
			"from": "''%s",
			"value": "0",
			"contractCode": "file:elrond_dex_farm.deprecated.wasm",
			"arguments": [
				"str:%s",
				"''%s",
				"1"
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
		pfe.wegldTokenId,
		string(pfe.routerAddress),
	))
	if err != nil {
		return err
	}

	// deploy mex farming
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scDeploy",
		"txId": "-farm-deploy-",
		"tx": {
			"from": "''%s",
			"value": "0",
			"contractCode": "file:elrond_dex_farm.deprecated.wasm",
			"arguments": [
				"str:%s",
				"''%s",
				"1"
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
		pfe.mexTokenId,
		string(pfe.routerAddress),
	))
	if err != nil {
		return err
	}

	// setup pair code
	fileBytes, err := ioutil.ReadFile("../../../test/dex/v0_2/output/elrond_dex_pair.wasm")
	if err != nil {
		fmt.Print(err)
	}

	pairCode := hex.EncodeToString(fileBytes)
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "start-pair-code-construction",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "startPairCodeConstruction",
			"arguments": [],
			"gasLimit": "10,000,000",
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
		string(pfe.routerAddress),
	))

	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "append-pair-code",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "appendPairCode",
			"arguments": [
				"0x%s"
			],
			"gasLimit": "10,000,000",
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
		string(pfe.routerAddress),
		pairCode,
	))
	if err != nil {
		return err
	}

	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "end-pair-code-construction",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "endPairCodeConstruction",
			"arguments": [],
			"gasLimit": "10,000,000",
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
		string(pfe.routerAddress),
	))
	if err != nil {
		return err
	}

	// issue stake token
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "issue-stake-token",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "5,000,000,000,000,000,000",
			"function": "issueFarmToken",
			"arguments": [
				"0x53656d6946756e6769626c65",
				"0x53454d4946554e47"
			],
			"gasLimit": "10,000,000",
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
		string(pfe.wegldFarmingAddress),
	))
	if err != nil {
		return err
	}

	// set local roles
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "set-local-roles-staking",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"function": "setLocalRolesFarmToken",
			"arguments": [],
			"gasLimit": "10,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "4",
			"message": "*",
			"gas": "*",
			"refund": "*"
		}
	}`,
		string(pfe.ownerAddress),
		string(pfe.wegldFarmingAddress),
	))
	if err != nil {
		return err
	}

	pfe.log("init ok")
	return nil
}
