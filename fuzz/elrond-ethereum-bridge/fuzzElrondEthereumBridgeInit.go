package elrond_ethereum_bridge

import "math/big"

type PriceAggregatorInitArgs struct {
	paymentToken       string
	oracleAddresses    []string
	submissionCount    int
	decimals           int
	queryPaymentAmount *big.Int
}

type MultisigInitArgs struct {
	requiredStake *big.Int
	slashAmount   *big.Int
	quorum        int
	boardMembers  []string
}

type DeployChildContractsArgs struct {
	egldEsdtSwapCode       []byte
	multiTransferEsdtCode  []byte
	ethereumFeePrepayCode  []byte
	esdtSafeCode           []byte
	priceAggregatorAddress string
	wrappedEgldTokenId     string
	wrappedEthTokenId      string
	tokenWhitelist         []string
}
