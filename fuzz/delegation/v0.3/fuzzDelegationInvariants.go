package delegation

import (
	"fmt"
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
			"should not have unexpected balance in the fuzzer. Unexpected balance: %d",
			unexpectedBalance)
	}

	return nil
}
