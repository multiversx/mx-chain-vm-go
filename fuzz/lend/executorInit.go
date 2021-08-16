package lendFuzz

import "fmt"

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

	return nil
}

func (e *executor) mintWallets() error {
	esdts := e.esdtBalances()
}

func (e *executor) esdtBalances() string {
	s := ""

	s += fmt.Sprintf(`"str:%s": "1,000,000,000,000,000,000,000,000,000,000",`, e.wegldTokenID)
	s += fmt.Sprintf(`"str:%s": "1,000,000,000,000,000,000,000,000,000,000"`, e.busdTokenID)

	return s
}
