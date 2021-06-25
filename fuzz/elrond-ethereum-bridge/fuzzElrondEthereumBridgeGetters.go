package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
)

func (fe *fuzzExecutor) interpretExpr(expression string) []byte {
	bytes, err := fe.mandosParser.ExprInterpreter.InterpretString(expression)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (fe *fuzzExecutor) relayerAddress(index int) string {
	return fmt.Sprintf("address:relayer-%05d", index)
}

func (fe *fuzzExecutor) userAddress(index int) string {
	return fmt.Sprintf("address:user-%05d", index)
}

func (fe *fuzzExecutor) getNonce(address string) int {
	acct := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	return int(acct.Nonce)
}

func (fe *fuzzExecutor) getBalance(address string) *big.Int {
	acct := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	return acct.Balance
}

func (fe *fuzzExecutor) getEsdtBalance(address string, tokenId string) *big.Int {
	acct := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	balance, err := acct.GetTokenBalanceByName(tokenId)

	if err != nil {
		return big.NewInt(0)
	}
	return balance
}

func (fe *fuzzExecutor) nextTxIndex() int {
	fe.txIndex++
	return fe.txIndex
}

func (fe *fuzzExecutor) getRandomUser() string {
	index := fe.randSource.Intn(len(fe.data.actorAddresses.users))
	return fe.data.actorAddresses.users[index]
}
