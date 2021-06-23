package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
)

func (pfe *fuzzExecutor) interpretExpr(expression string) []byte {
	bytes, err := pfe.mandosParser.ExprInterpreter.InterpretString(expression)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (pfe *fuzzExecutor) relayerAddress(index int) string {
	return fmt.Sprintf("address:relayer-%05d", index)
}

func (pfe *fuzzExecutor) userAddress(index int) string {
	return fmt.Sprintf("address:user-%05d", index)
}

func (pfe *fuzzExecutor) getNonce(address string) int {
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(address))
	return int(acct.Nonce)
}

func (pfe *fuzzExecutor) getBalance(address string) *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(address))
	return acct.Balance
}

func (pfe *fuzzExecutor) getEsdtBalance(address string, tokenId string) *big.Int {
	acct := pfe.world.AcctMap.GetAccount(pfe.interpretExpr(address))
	balance, err := acct.GetTokenBalanceByName(tokenId)

	if err != nil {
		return big.NewInt(0)
	}
	return balance
}
