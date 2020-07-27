package delegation

import (
	"fmt"
	"math/big"
)

func (pfe *fuzzDelegationExecutor) checkNoUnexpectedBalance() error {
	err := pfe.computeAllRewards()
	if err != nil {
		return err
	}

	unexpectedBalance, err := pfe.simpleQuery("getUnexpectedBalance")
	if err != nil {
		return err
	}

	if unexpectedBalance.Sign() > 0 {
		return fmt.Errorf(
			"Should not have unexpected balance in the fuzzer. Unexpected balance: %d",
			unexpectedBalance)
	}

	return nil
}

// Currently outdated, not used.
func (pfe *fuzzDelegationExecutor) checkContractBalanceVsState() error {
	err := pfe.computeAllRewards()
	if err != nil {
		return err
	}

	totalInactive, err := pfe.simpleQuery("getTotalInactiveStake")
	if err != nil {
		return err
	}

	totalUnclRewards, err := pfe.simpleQuery("getTotalUnclaimedRewards")
	if err != nil {
		return err
	}

	expectedBalance := big.NewInt(0).Add(totalInactive, totalUnclRewards)

	contractBalance := pfe.getContractBalance()
	if contractBalance.Cmp(expectedBalance) != 0 {
		return fmt.Errorf(
			"bad contract balance.\nWant: %d (inactive stake) + %d (unclaimed rewards) = %d\nHave: %d",
			totalInactive, totalUnclRewards, expectedBalance,
			contractBalance)
	}
	return nil
}
