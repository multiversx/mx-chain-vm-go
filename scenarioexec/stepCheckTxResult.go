package scenarioexec

import (
	"fmt"
	"math/big"

	er "github.com/multiversx/mx-chain-scenario-go/expression/reconstructor"
	mjwrite "github.com/multiversx/mx-chain-scenario-go/json/write"
	mj "github.com/multiversx/mx-chain-scenario-go/model"
	oj "github.com/multiversx/mx-chain-scenario-go/orderedjson"
	vmi "github.com/multiversx/mx-chain-vm-common-go"
)

func (ae *VMTestExecutor) checkTxResults(
	txIndex string,
	blResult *mj.TransactionResult,
	checkGas bool,
	output *vmi.VMOutput,
) error {

	if !blResult.Status.Check(big.NewInt(int64(output.ReturnCode))) {
		return fmt.Errorf("result code mismatch. Tx '%s'. Want: %s. Have: %d (%s). Message: %s",
			txIndex, blResult.Status.Original, int(output.ReturnCode), output.ReturnCode.String(), output.ReturnMessage)
	}

	if !blResult.Message.Check([]byte(output.ReturnMessage)) {
		return fmt.Errorf("result message mismatch. Tx '%s'. Want: %s. Have: %s",
			txIndex, blResult.Message.Original, output.ReturnMessage)
	}

	// check result
	if !blResult.Out.CheckList(output.ReturnData) {
		return fmt.Errorf("result mismatch. Tx '%s'. Want: %s. Have: %s",
			txIndex,
			checkBytesListPretty(blResult.Out),
			ae.exprReconstructor.ReconstructList(output.ReturnData, er.NoHint))
	}

	// check refund
	if !blResult.Refund.Check(output.GasRefund) {
		return fmt.Errorf("result gas refund mismatch. Tx '%s'. Want: %s. Have: 0x%x",
			txIndex, blResult.Refund.Original, output.GasRefund)
	}

	// check gas
	// unlike other checks, if unspecified the remaining gas check is ignored
	if checkGas && !blResult.Gas.IsUnspecified() && !blResult.Gas.Check(output.GasRemaining) {
		return fmt.Errorf("result gas mismatch. Tx '%s'. Want: %s. Got: %d (0x%x)",
			txIndex,
			blResult.Gas.Original,
			output.GasRemaining,
			output.GasRemaining)
	}

	return ae.checkTxLogs(txIndex, blResult.Logs, output.Logs)
}

func (ae *VMTestExecutor) checkTxLogs(
	txIndex string,
	expectedLogs mj.LogList,
	actualLogs []*vmi.LogEntry,
) error {
	// "logs": "*" means any value is accepted, log check ignored
	if expectedLogs.IsStar {
		return nil
	}

	// this is the real log check
	if len(actualLogs) < len(expectedLogs.List) {
		return fmt.Errorf("too few logs. Tx '%s'. Want:%d. Got:%d",
			txIndex,
			len(expectedLogs.List),
			len(actualLogs))
	}

	for i, actualLog := range actualLogs {
		if i < len(expectedLogs.List) {
			testLog := expectedLogs.List[i]
			err := ae.checkTxLog(txIndex, i, testLog, actualLog)
			if err != nil {
				return err
			}
		} else if !expectedLogs.MoreAllowedAtEnd {
			return fmt.Errorf("unexpected log. Tx '%s'. Log index: %d. Log:\n%s",
				txIndex,
				i,
				mjwrite.LogToString(ae.convertLogToTestFormat(actualLog)),
			)
		}
	}

	return nil
}

func (ae *VMTestExecutor) checkTxLog(
	txIndex string,
	logIndex int,
	expectedLog *mj.LogEntry,
	actualLog *vmi.LogEntry) error {
	if !expectedLog.Address.Check(actualLog.Address) {
		return fmt.Errorf("bad log address. Tx '%s'. Log index: %d. Want:\n%s\nGot:\n%s",
			txIndex,
			logIndex,
			mjwrite.LogToString(expectedLog),
			mjwrite.LogToString(ae.convertLogToTestFormat(actualLog)))
	}
	if !expectedLog.Endpoint.Check(actualLog.Identifier) {
		return fmt.Errorf("bad log identifier. Tx '%s'. Log index: %d. Want:\n%s\nGot:\n%s",
			txIndex,
			logIndex,
			mjwrite.LogToString(expectedLog),
			mjwrite.LogToString(ae.convertLogToTestFormat(actualLog)))
	}
	if !expectedLog.Topics.CheckList(actualLog.Topics) {
		return fmt.Errorf("bad log topics. Tx '%s'. Log index: %d. Want: %s. Have: %s",
			txIndex,
			logIndex,
			checkBytesListPretty(expectedLog.Topics),
			ae.exprReconstructor.ReconstructList(actualLog.Topics, er.NoHint))
	}
	if !expectedLog.Data.Check(actualLog.Data) {
		return fmt.Errorf("bad log data. Tx '%s'. Log index: %d. Want:\n%s\nGot:\n%s",
			txIndex,
			logIndex,
			mjwrite.LogToString(expectedLog),
			mjwrite.LogToString(ae.convertLogToTestFormat(actualLog)))
	}
	return nil
}

// JSONCheckBytesString formats a list of JSONCheckBytes for printing to console.
// TODO: move somewhere else
func checkBytesListPretty(jcbl mj.JSONCheckValueList) string {
	str := "["
	for i, jcb := range jcbl.Values {
		if i > 0 {
			str += ", "
		}

		str += oj.JSONString(jcb.Original)
	}
	return str + "]"
}
