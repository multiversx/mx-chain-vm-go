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
	boardMembers      []string
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
