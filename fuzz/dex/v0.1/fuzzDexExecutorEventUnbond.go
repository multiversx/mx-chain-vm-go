package dex

import (
	"errors"
	"fmt"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"math/big"
	"math/rand"
)


func (pfe *fuzzDexExecutor) unbond(amountMax int, statistics *eventsStatistics, rand *rand.Rand) error {
	unstakersLen := len(pfe.unstakers)
	if unstakersLen == 0 {
		return nil
	}

	nonce := rand.Intn(unstakersLen)
	user := pfe.unstakers[nonce].user
	amount := pfe.unstakers[nonce].value
	if pfe.unstakers[nonce].value == 0 {
		return nil
	}

	unbondAmount := int64(amountMax)
	if int64(amountMax) > amount {
		unbondAmount = amount
	} else {
		unbondAmount = int64(amountMax)
	}
	lpToken := pfe.unstakers[nonce].lpToken
	pfe.unstakers[nonce] = UnstakeInfo{
		value: amount - unbondAmount,
		user: user,
		lpToken: lpToken,
	}

	lpBefore, errGet := pfe.getTokens([]byte(user), lpToken)
	if errGet != nil {
		return nil
	}

	output, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "unbond",
		"tx": {
			"from": "''%s",
			"to": "''%s",
			"value": "0",
			"esdt": {
				"tokenIdentifier": "str:%s",
				"value": "%d",
				"nonce": "%d"
			},
			"function": "unbond",
			"arguments": [],
			"gasLimit": "100,000,000",
			"gasPrice": "0"
		}
	}`,
		user,
		string(pfe.wegldStakingAddress),
		"UNSTAK-abcdef",
		unbondAmount,
		nonce,
	))
	if output == nil {
		return err
	}

	success := output.ReturnCode == vmi.Ok
	if success {
		statistics.unbondHits += 1

		lpAfter, errAfter := pfe.getTokens([]byte(user), lpToken)
		if errAfter != nil {
			return nil
		}

		if lpAfter.Cmp(big.NewInt(0).Add(big.NewInt(unbondAmount), lpBefore)) != 0 {
			return errors.New("UNBOND failed")
		}
	} else {
		statistics.unbondMisses += 1
		pfe.log("unbond")
		pfe.log("could not unbond because %s", output.ReturnMessage)

		if output.ReturnMessage == "insufficient funds" {
			return errors.New(output.ReturnMessage)
		}
	}

	return nil
}