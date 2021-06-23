package elrond_ethereum_bridge

import "math/big"

type TransactionStatus int

const (
	None TransactionStatus = iota
	Pending
	InProgress
	Executed
	Rejected
)

type ActorAddresses struct {
	owner             string
	relayers          []string
	users             []string
	multisig          string
	priceAggregator   string
	egldEsdtSwap      string
	esdtSafe          string
	ethereumFeePrepay string
	multiTransferEsdt string
}

type Transaction struct {
	blockNonce int
	nonce      int
	from       string
	to         string
	tokenId    string
	amount     *big.Int
}

type SimpleTransfer struct {
	to      string
	tokenId string
	amount  *big.Int
}

type EgldEsdtSwapState struct {
	wrappedEgldTokenId string
	egldBalance        *big.Int
}

type EsdtSafeState struct {
	tokenWhitelist []string
	transactions   []*Transaction
	txStatus       []TransactionStatus
	balances       map[string]*big.Int
}

type EthereumFeePrepayState struct {
	wrappedEthTokenId string
	deposits          map[string]map[string]*big.Int
}

type MultiTransferEsdtState struct {
	tokenWhitelist []string
}

type MultisigState struct {
	owner                   string
	boardMembers            []string
	requiredStake           int
	stakedAmounts           map[string]*big.Int
	quorum                  int
	paused                  bool
	currentTransactionBatch []*Transaction
	lastValidActionId       int
	lastExecutedActionId    int
}
