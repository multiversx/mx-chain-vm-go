package lendFuzz

import "fmt"

const (
	setAccountEsdts = `
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
	}`

	setAccountBalance = `
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
	}`
)

func (e *executor) init(args *executorArgs) error {
	e.wegldTokenID = args.wegldTokenID
	e.lwegldTokenID = args.lwegldTokenID
	e.bwegldTokenID = args.bwegldTokenID
	e.busdTokenID = args.busdTokenID
	e.lbusdTokenID = args.lbusdTokenID
	e.bbusdTokenID = args.bbusdTokenID

	e.ownerAddress = "address:owner"
	e.wegldLPAddress = "sc:wegld_liq_pool"
	e.busdLPAddress = "sc:busd_liq_pool"

	e.numUsers = args.numUsers
	e.numTokens = args.numTokens
	e.numEvents = args.numEvents

	err := e.mintWallets()
	if err != nil {
		return err
	}

	return nil
}

func (e *executor) setupLiquidityPool(lpAddress, tokenID, owner string) error {
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
					"code": "file:elrond_dex_pair.wasm",
					"owner": "%s"
				}
			}
		}`, lpAddress, tokenID, owner))
}

func (e *executor) setupLendingPool() error {
	return nil
}

func (e *executor) mintWallets() error {
	esdts := e.esdtBalances()
	for i := 1; i <= e.numUsers; i++ {
		err := e.executeStep(fmt.Sprintf(setAccountEsdts, e.userAddress(i), esdts))
		if err != nil {
			return err
		}
	}

	return e.executeStep(fmt.Sprintf(setAccountBalance, e.ownerAddress))
}

func (e *executor) userAddress(index int) string {
	return fmt.Sprintf("address:user%06d", index)
}

func (e *executor) esdtBalances() string {
	s := ""

	s += fmt.Sprintf(`"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, e.wegldTokenID)
	s += fmt.Sprintf(`"str:%s": "1,000,000,000,000,000,000,000,000,000,000"`, e.busdTokenID)

	return s
}
