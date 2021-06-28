package elrond_ethereum_bridge

import (
	"encoding/binary"
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
	balance, err := acct.GetTokenBalanceByName(string(fe.interpretExpr(tokenId)))

	if err != nil {
		return big.NewInt(0)
	}
	return balance
}

func (fe *fuzzExecutor) generateValidRandomEsdtPayment(address string) (string, *big.Int, error) {
	acc := fe.world.AcctMap.GetAccount(fe.interpretExpr(address))
	allEsdts, err := acc.GetFullMockESDTData()
	if err != nil {
		return "", nil, err
	}

	// map key order is not guaranteed, so this is "random"
	for tokenId := range allEsdts {
		tokenIdMandosFormat := "str:" + tokenId
		balance := fe.getEsdtBalance(address, tokenIdMandosFormat)
		amount := fe.getRandomBigInt(balance)

		return tokenIdMandosFormat, amount, nil
	}

	return "", nil, fmt.Errorf("Account has no ESDT")
}

func (fe *fuzzExecutor) nextTxIndex() int {
	fe.txIndex++
	return fe.txIndex
}

// Do NOT use with bigInts that don't fit [0, max(int64))
// maxValue is inclusive
func (fe *fuzzExecutor) getRandomBigInt(maxValue *big.Int) *big.Int {
	return big.NewInt(int64(fe.randSource.Intn(int(maxValue.Int64())) + 1))
}

func (fe *fuzzExecutor) getRandomUser() string {
	index := fe.randSource.Intn(len(fe.data.actorAddresses.users))
	return fe.data.actorAddresses.users[index]
}

func (fe *fuzzExecutor) getRandomRelayer() string {
	index := fe.randSource.Intn(len(fe.data.actorAddresses.relayers))
	return fe.data.actorAddresses.relayers[index]
}

func (fe *fuzzExecutor) getEthAddress() string {
	return "0x0102030405060708091011121314151617181920"
}

func (fe *fuzzExecutor) bytesToInt(bytes []byte) int {
	intBytes := bytes
	for len(intBytes) != 8 {
		intBytes = append([]byte{0}, intBytes...)
	}

	return int(binary.BigEndian.Uint64(intBytes))
}
