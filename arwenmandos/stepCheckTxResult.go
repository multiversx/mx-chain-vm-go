package arwenmandos

import (
	"bytes"
	"fmt"
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	vmi "github.com/ElrondNetwork/elrond-go/core/vmcommon"
)

func checkTxResults(
	txIndex string,
	blResult *mj.TransactionResult,
	checkGas bool,
	output *vmi.VMOutput,
) error {

	if !blResult.Status.Check(big.NewInt(int64(output.ReturnCode))) {
		return fmt.Errorf("result code mismatch. Tx %s. Want: %s. Have: %d (%s). Message: %s",
			txIndex, blResult.Status.Original, int(output.ReturnCode), output.ReturnCode.String(), output.ReturnMessage)
	}

	if !blResult.Message.Check([]byte(output.ReturnMessage)) {
		return fmt.Errorf("result message mismatch. Tx %s. Want: %s. Have: %s",
			txIndex, blResult.Message.Original, output.ReturnMessage)
	}

	// check result
	if len(output.ReturnData) != len(blResult.Out) {
		return fmt.Errorf("result length mismatch. Tx %s. Want: %s. Have: %s",
			txIndex,
			checkBytesListPretty(blResult.Out),
			mj.ResultAsString(output.ReturnData))
	}
	for i, expected := range blResult.Out {
		if !expected.Check(output.ReturnData[i]) {
			return fmt.Errorf("result mismatch. Tx %s. Want: %s. Have: %s",
				txIndex,
				checkBytesListPretty(blResult.Out),
				mj.ResultAsString(output.ReturnData))
		}
	}

	// check refund
	if !blResult.Refund.Check(output.GasRefund) {
		return fmt.Errorf("result gas refund mismatch. Tx %s. Want: %s. Have: 0x%x",
			txIndex, blResult.Refund.Original, output.GasRefund)
	}

	// check gas
	if checkGas && !blResult.Gas.Check(output.GasRemaining) {
		return fmt.Errorf("result gas mismatch. Tx %s. Want: %s. Got: %d (0x%x)",
			txIndex,
			blResult.Gas.Original,
			output.GasRemaining,
			output.GasRemaining)
	}

	// "logs": "*" means any value is accepted, log check ignored
	if blResult.IgnoreLogs {
		return nil
	}

	// this is the real log check
	if len(blResult.Logs) != len(output.Logs) {
		return fmt.Errorf("wrong number of logs. Tx %s. Want:%d. Got:%d",
			txIndex,
			len(blResult.Logs),
			len(output.Logs))
	}
	for i, outLog := range output.Logs {
		testLog := blResult.Logs[i]
		if !bytes.Equal(outLog.Address, testLog.Address.Value) {
			return fmt.Errorf("bad log address. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
		if !bytes.Equal(outLog.Identifier, testLog.Identifier.Value) {
			return fmt.Errorf("bad log identifier. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
		if len(outLog.Topics) != len(testLog.Topics) {
			return fmt.Errorf("wrong number of log topics. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
		for ti := range outLog.Topics {
			if !bytes.Equal(outLog.Topics[ti], testLog.Topics[ti].Value) {
				return fmt.Errorf("bad log topic. Tx %s. Want:\n%s\nGot:\n%s",
					txIndex,
					mjwrite.LogToString(testLog),
					mjwrite.LogToString(convertLogToTestFormat(outLog)))
			}
		}
		if big.NewInt(0).SetBytes(outLog.Data).Cmp(big.NewInt(0).SetBytes(testLog.Data.Value)) != 0 {
			return fmt.Errorf("bad log data. Tx %s. Want:\n%s\nGot:\n%s",
				txIndex,
				mjwrite.LogToString(testLog),
				mjwrite.LogToString(convertLogToTestFormat(outLog)))
		}
	}

	return nil
}
