package delegation

import (
	"fmt"
	"math/big"
)

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
