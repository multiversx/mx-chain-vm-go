package lendFuzz

import (
	"fmt"
	"log"
)

func (e *fuzzLendExecutor) init(args *fuzzLendExecutorArgs) error {
	e.wegldTokenID = args.wegldTokenID
	e.lwegldTokenID = args.lwegldTokenID
	e.bwegldTokenID = args.bwegldTokenID
	e.busdTokenID = args.busdTokenID
	e.lbusdTokenID = args.lbusdTokenID
	e.bbusdTokenID = args.bbusdTokenID

	e.ownerAddress = "address:owner"
	e.lendPoolAddress = "sc:lend_pool"
	e.wegldLPAddress = "sc:wegld_liq_pool"
	e.busdLPAddress = "sc:busd_liq_pool"

	e.numUsers = args.numUsers
	e.numTokens = args.numTokens
	e.numEvents = args.numEvents

	err := e.initAccounts()
	if err != nil {
		return err
	}

	log.Println("init successful")

	return nil
}

func (e *fuzzLendExecutor) initAccounts() error {
	esdtBalances := e.esdtBalances()
	for i := 1; i <= e.numUsers; i++ {
		err := e.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
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
			e.userAddress(i),
			esdtBalances,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *fuzzLendExecutor) initOwner() error {
	return e.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"%s": {
				"nonce": "0",
				"balance": "1,000,000,000,000,000,000,000,000,000,000",
				"storage": {},
				"code": ""
			}
		}
	}`,
		e.ownerAddress,
	))
}

func (e *fuzzLendExecutor) initContracts() error {
	return nil
}

func (e *fuzzLendExecutor) setupLiquidityPool(lpAddress, tokenID, owner string) error {
	return e.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0",
					"esdt": {
						"str:%s": {
							"roles": [
								"ESDTRoleLocalMint",
								"ESDTRoleLocalBurn"
							]
						}
					},
					"storage": {
					},
					"code": "file:liquidity_pool.wasm",
					"owner": "%s"
				}
			}
		}`, lpAddress, tokenID, owner))
}

func (e *fuzzLendExecutor) setupLendingPool() error {
	return e.executeStep(fmt.Sprintf(`
	
`))
}

func (e *fuzzLendExecutor) userAddress(index int) string {
	return fmt.Sprintf("address:user%06d", index)
}

func (e *fuzzLendExecutor) esdtBalances() string {
	s := ""

	s += fmt.Sprintf(`"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, e.wegldTokenID)
	s += fmt.Sprintf(`"str:%s": "1,000,000,000,000,000,000,000,000,000,000"`, e.busdTokenID)

	return s
}
