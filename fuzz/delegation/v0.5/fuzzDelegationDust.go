package delegation

import (
	"fmt"
	"math/big"
)

func (pfe *fuzzDelegationExecutor) dustCleanup(dustLimit *big.Int) error {
	for pfe.hasDustItemsWaitingList(dustLimit) {
		err := pfe.dustCleanupWaitingList(dustLimit)
		if err != nil {
			return err
		}
	}

	for pfe.hasDustItemsActive(dustLimit) {
		err := pfe.dustCleanupActive(dustLimit)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pfe *fuzzDelegationExecutor) hasDustItemsWaitingList(dustLimit *big.Int) bool {
	count, err := pfe.querySingleResult("countDustItemsWaitingList", fmt.Sprintf(`"%d"`, dustLimit))
	if err != nil {
		panic(err)
	}
	pfe.log("dust (waiting). Limit: %d, Count: %d", dustLimit, count)
	return count.Sign() > 0
}

func (pfe *fuzzDelegationExecutor) hasDustItemsActive(dustLimit *big.Int) bool {
	count, err := pfe.querySingleResult("countDustItemsActive", fmt.Sprintf(`"%d"`, dustLimit))
	if err != nil {
		panic(err)
	}
	pfe.log("dust (active).  Limit: %d, Count: %d", dustLimit, count)
	return count.Sign() > 0
}

func (pfe *fuzzDelegationExecutor) dustCleanupWaitingList(dustLimit *big.Int) error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "%d",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "dustCleanupWaitingList",
				"arguments": [
					"%d"
				],
				"gasLimit": "300,000,000",
				"gasPrice": "0"
			},
			"expect": {
				"out": [],
				"status": "",
				"refund": "*"
			}
		}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		dustLimit,
	))
	return err
}

func (pfe *fuzzDelegationExecutor) dustCleanupActive(dustLimit *big.Int) error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "%d",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "dustCleanupActive",
				"arguments": [
					"%d"
				],
				"gasLimit": "300,000,000",
				"gasPrice": "0"
			},
			"expect": {
				"out": [],
				"status": "",
				"refund": "*"
			}
		}`,
		pfe.nextTxIndex(),
		pfe.ownerAddress,
		pfe.delegationContractAddress,
		dustLimit,
	))
	return err
}
