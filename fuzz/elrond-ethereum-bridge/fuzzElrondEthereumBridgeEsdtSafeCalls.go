package elrond_ethereum_bridge

import (
	"encoding/binary"
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

	transaction := &Transaction{
		from:    userAddress,
		tokenId: tokenId,
		amount:  amount,
	}
	fe.data.multisigState.allEsdtSafeTransactions = append(fe.data.multisigState.allEsdtSafeTransactions, transaction)

	return nil
}

func (fe *fuzzExecutor) getNextTransactionBatch() error {
	// call has to be done through the multisig SC
	output, err := fe.performSmartContractCall(
		fe.getRandomRelayer(),
		fe.data.actorAddresses.multisig,
		big.NewInt(0),
		"getNextTransactionBatch",
		[]string{},
	)
	if err != nil {
		return err
	}

	// no transactions were fetched
	if len(fe.data.multisigState.allEsdtSafeTransactions) == 0 {
		return nil
	}

	// output[0] is the whole serialized batch, so we ignore that
	// This is due to a limitation in the current executeOnDestContext implementation
	// SCResults from child contracts propagate to the original caller

	// append zeroes in front to fill 8 bytes (needed for Uint64 conversion)
	batchIdBytes := output[1]
	for len(batchIdBytes) != 8 {
		batchIdBytes = append([]byte{0}, batchIdBytes...)
	}

	batchId := int(binary.BigEndian.Uint64(batchIdBytes))
	nrStructFields := 6
	nrTransactionsInBatch := (len(output) - 2) / nrStructFields

	if nrTransactionsInBatch > 0 {
		fe.data.multisigState.currentEsdtSafeBatchId = batchId
		fe.data.multisigState.currentEsdtSafeTransactionBatch = fe.data.multisigState.allEsdtSafeTransactions[:nrTransactionsInBatch]

		if nrTransactionsInBatch == len(fe.data.multisigState.allEsdtSafeTransactions) {
			fe.data.multisigState.allEsdtSafeTransactions = []*Transaction{}
		} else {
			fe.data.multisigState.allEsdtSafeTransactions = fe.data.multisigState.allEsdtSafeTransactions[nrTransactionsInBatch:]
		}
	}

	return nil
}
