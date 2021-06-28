package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
)

func (fe *fuzzExecutor) createEsdtSafeTransaction(userAddress string,
	tokenId string, amount *big.Int, destEthAddress string) error {

	esdtBalanceBefore := fe.getEsdtBalance(userAddress, tokenId)

	_, err := fe.performEsdtTransferSmartContractCall(
		userAddress,
		fe.data.actorAddresses.esdtSafe,
		tokenId,
		amount,
		"createTransaction",
		[]string{destEthAddress},
		true,
		"",
		[]string{},
	)
	if err != nil {
		return err
	}

	actualEsdtBalanceAfter := fe.getEsdtBalance(userAddress, tokenId)

	expectedEsdtBalanceAfter := big.NewInt(0)
	expectedEsdtBalanceAfter.Sub(esdtBalanceBefore, amount)

	if actualEsdtBalanceAfter.Cmp(expectedEsdtBalanceAfter) != 0 {
		return fmt.Errorf("Wrong ESDT balance after creating EsdtSafe transaction, Expected: %s, Have: %s",
			expectedEsdtBalanceAfter.String(),
			actualEsdtBalanceAfter.String())
	}

	return nil
}
