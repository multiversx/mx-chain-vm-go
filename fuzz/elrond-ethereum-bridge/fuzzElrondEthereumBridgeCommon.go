package elrond_ethereum_bridge

import "math/big"

type Address string
type EthAddress string
type TokenIdentifier string
type TransactionStatus int

const (
	None TransactionStatus = iota
	Pending
	InProgress
	Executed
	Rejected
)

type ActorAddresses struct {
	accounts          []Address
	multisig          Address
	priceAggregator   Address
	egldEsdtSwap      Address
	esdtSafe          Address
	ethereumFeePrepay Address
	multiTransferEsdt Address
}

type Transaction struct {
	blockNonce int
	nonce      int
	from       Address
	to         EthAddress
	tokenId    TokenIdentifier
	amount     *big.Int
}

type SimpleTransfer struct {
	to      Address
	tokenId TokenIdentifier
	amount  *big.Int
}

type EgldEsdtSwapState struct {
	wrappedEgldTokenId TokenIdentifier
	egldBalance        *big.Int
}

type EsdtSafeState struct {
	tokenWhitelist []TokenIdentifier
	transactions   []*Transaction
	txStatus       []TransactionStatus
	balances       map[TokenIdentifier]*big.Int
}

type EthereumFeePrepayState struct {
	wrappedEthTokenId TokenIdentifier
	deposits          map[Address]map[TokenIdentifier]*big.Int
}

type MultiTransferEsdtState struct {
	tokenWhitelist []TokenIdentifier
}

type MultisigState struct {
	owner                   Address
	boardMembers            []Address
	requiredStake           int
	stakedAmounts           map[Address]*big.Int
	quorum                  int
	paused                  bool
	currentTransactionBatch []*Transaction
	lastValidActionId       int
	lastExecutedActionId    int
}
