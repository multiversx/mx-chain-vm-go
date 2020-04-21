package arwenjsontest

import (
	"bytes"
	"fmt"
	"math/big"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func checkTxResults(
	txIndex int,
	blResult *ij.TransactionResult,
	checkGas bool,
	output *vmi.VMOutput,
) error {

	expectedStatus := 0
	if blResult.Status.Value != nil {
		expectedStatus = int(blResult.Status.Value.Int64())
	}
	if expectedStatus != int(output.ReturnCode) {
		return fmt.Errorf("result code mismatch. Tx #%d. Want: %d. Have: %d (%s). Message: %s",
			txIndex, expectedStatus, int(output.ReturnCode), output.ReturnCode.String(), output.ReturnMessage)
	}

	if output.ReturnMessage != blResult.Message {
		return fmt.Errorf("result message mismatch. Tx #%d. Want: %s. Have: %s",
			txIndex, blResult.Message, output.ReturnMessage)
	}

	// check result
	if len(output.ReturnData) != len(blResult.Out) {
		return fmt.Errorf("result length mismatch. Tx #%d. Want: %s. Have: %s",
			txIndex,
			ij.ResultAsString(ij.JSONBytesValues(blResult.Out)),
			ij.ResultAsString(output.ReturnData))
	}
	for i, expected := range blResult.Out {
		if !ij.ResultEqual(expected, output.ReturnData[i]) {
			return fmt.Errorf("result mismatch. Tx #%d. Want: %s. Have: %s",
				txIndex,
				ij.ResultAsString(ij.JSONBytesValues(blResult.Out)),
				ij.ResultAsString(output.ReturnData))
		}
	}

	// check refund
	if !blResult.Refund.IsStar {
		if blResult.Refund.Value.Cmp(output.GasRefund) != 0 {
			return fmt.Errorf("result gas refund mismatch. Want: 0x%x. Have: 0x%x",
				blResult.Refund.Value, output.GasRefund)
		}
	}

	// check gas
	if checkGas && !blResult.Gas.IsStar {
		if blResult.Gas.Value != output.GasRemaining {
			return fmt.Errorf("result gas mismatch. Want: %d (0x%x). Got: %d (0x%x)",
				blResult.Gas.Value,
				blResult.Gas.Value,
				output.GasRemaining,
				output.GasRemaining)
		}
	}

	// "logs": "*" means any value is accepted, log check ignored
	if blResult.IgnoreLogs {
		return nil
	}

	// this is the real log check
	if len(blResult.Logs) != len(output.Logs) {
		return fmt.Errorf("wrong number of logs. Want:%d. Got:%d",
			len(blResult.Logs), len(output.Logs))
	}
	for i, outLog := range output.Logs {
		testLog := blResult.Logs[i]
		if !bytes.Equal(outLog.Address, testLog.Address.Value) {
			return fmt.Errorf("bad log address. Want:\n%s\nGot:\n%s",
				ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
		}
		if !bytes.Equal(outLog.Identifier, testLog.Identifier.Value) {
			return fmt.Errorf("bad log identifier. Want:\n%s\nGot:\n%s",
				ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
		}
		if len(outLog.Topics) != len(testLog.Topics) {
			return fmt.Errorf("wrong number of log topics. Want:\n%s\nGot:\n%s",
				ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
		}
		for ti := range outLog.Topics {
			if !bytes.Equal(outLog.Topics[ti], testLog.Topics[ti].Value) {
				return fmt.Errorf("bad log topic. Want:\n%s\nGot:\n%s",
					ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
			}
		}
		if big.NewInt(0).SetBytes(outLog.Data).Cmp(big.NewInt(0).SetBytes(testLog.Data.Value)) != 0 {
			return fmt.Errorf("bad log data. Want:\n%s\nGot:\n%s",
				ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
		}
	}

	return nil
}
