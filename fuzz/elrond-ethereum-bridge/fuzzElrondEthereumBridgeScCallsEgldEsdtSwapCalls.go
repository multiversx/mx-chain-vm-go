package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
)

func (fe *fuzzExecutor) wrapEgld(userAddress string, amount *big.Int) error {
	wrappedEgldTokenId := string(fe.interpretExpr(fe.data.wrappedEgldTokenId))

	egldBalanceBefore := fe.getBalance(userAddress)
	esdtBalanceBefore := fe.getEsdtBalance(userAddress, wrappedEgldTokenId)

	fe.performSmartContractCall(
		userAddress,
		fe.data.actorAddresses.egldEsdtSwap,
		amount,
		"wrapEgld",
		[]string{},
		true,
		"",
		[]string{},
	)

	actualEgldBalanceAfter := fe.getBalance(userAddress)
	actualEsdtBalanceAfter := fe.getEsdtBalance(userAddress, wrappedEgldTokenId)

	expectedEgldBalanceAfter := big.NewInt(0)
	expectedEgldBalanceAfter.Sub(egldBalanceBefore, amount)

	expectedEsdtBalanceAfter := big.NewInt(0)
	expectedEsdtBalanceAfter.Add(esdtBalanceBefore, amount)

	if actualEgldBalanceAfter.Cmp(expectedEgldBalanceAfter) != 0 {
		return fmt.Errorf("Wrong EGLD balance after wrapping, Expected: %s, Have: %s",
			expectedEgldBalanceAfter.String(),
			actualEgldBalanceAfter.String())
	}
	if actualEsdtBalanceAfter.Cmp(expectedEsdtBalanceAfter) != 0 {
		return fmt.Errorf("Wrong ESDT balance after wrapping, Expected: %s, Have: %s",
			expectedEsdtBalanceAfter.String(),
			actualEsdtBalanceAfter.String())
	}

	return nil
}

func (fe *fuzzExecutor) unwrapEgld(userAddress string, amount *big.Int) error {
	scEgldBalance := fe.getBalance(fe.data.actorAddresses.egldEsdtSwap)
	if scEgldBalance.Cmp(amount) < 0 {
		return fmt.Errorf("EgldEsdtSwap does not have enough EGLD balance. Have: %s, Need: %s",
			scEgldBalance.String(),
			amount.String(),
		)
	}

	wrappedEgldTokenId := string(fe.interpretExpr(fe.data.wrappedEgldTokenId))

	egldBalanceBefore := fe.getBalance(userAddress)
	esdtBalanceBefore := fe.getEsdtBalance(userAddress, wrappedEgldTokenId)

	fe.performSmartContractCall(
		userAddress,
		fe.data.actorAddresses.egldEsdtSwap,
		amount,
		"ESDTTransfer",
		[]string{wrappedEgldTokenId, amount.String(), "str:unwrapEgld"},
		true,
		"",
		[]string{},
	)

	actualEgldBalanceAfter := fe.getBalance(userAddress)
	actualEsdtBalanceAfter := fe.getEsdtBalance(userAddress, wrappedEgldTokenId)

	expectedEgldBalanceAfter := big.NewInt(0)
	expectedEgldBalanceAfter.Add(egldBalanceBefore, amount)

	expectedEsdtBalanceAfter := big.NewInt(0)
	expectedEsdtBalanceAfter.Sub(esdtBalanceBefore, amount)

	if actualEgldBalanceAfter.Cmp(expectedEgldBalanceAfter) != 0 {
		return fmt.Errorf("Wrong EGLD balance after unwrapping, Expected: %s, Have: %s",
			expectedEgldBalanceAfter.String(),
			actualEgldBalanceAfter.String())
	}
	if actualEsdtBalanceAfter.Cmp(expectedEsdtBalanceAfter) != 0 {
		return fmt.Errorf("Wrong ESDT balance after unwrapping, Expected: %s, Have: %s",
			expectedEsdtBalanceAfter.String(),
			actualEsdtBalanceAfter.String())
	}

	return nil
}
