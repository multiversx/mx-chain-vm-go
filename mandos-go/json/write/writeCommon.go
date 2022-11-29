package mandosjsonwrite

import (
	"encoding/hex"
	"math/big"

	mj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/model"
	oj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/orderedjson"
)

func resultToOJ(res *mj.TransactionResult) oj.OJsonObject {
	resultOJ := oj.NewMap()

	resultOJ.Put("out", checkValueListToOJ(res.Out))

	if !res.Status.IsUnspecified() {
		resultOJ.Put("status", checkBigIntToOJ(res.Status))
	}
	if !res.Message.IsUnspecified() {
		resultOJ.Put("message", checkBytesToOJ(res.Message))
	}
	if !res.Logs.IsUnspecified {
		if res.Logs.IsStar {
			resultOJ.Put("logs", stringToOJ("*"))
		} else {
			resultOJ.Put("logs", logsToOJ(res.Logs))

		}
	}
	if !res.Gas.IsUnspecified() {
		resultOJ.Put("gas", checkUint64ToOJ(res.Gas))
	}
	if !res.Refund.IsUnspecified() {
		resultOJ.Put("refund", checkBigIntToOJ(res.Refund))
	}

	return resultOJ
}

// LogToString returns a json representation of a log entry, we use it for debugging
func LogToString(logEntry *mj.LogEntry) string {
	logOJ := logToOJ(logEntry)
	return oj.JSONString(logOJ)
}

func logToOJ(logEntry *mj.LogEntry) oj.OJsonObject {
	logOJ := oj.NewMap()
	logOJ.Put("address", checkBytesToOJ(logEntry.Address))
	logOJ.Put("endpoint", checkBytesToOJ(logEntry.Endpoint))
	logOJ.Put("topics", checkValueListToOJ(logEntry.Topics))
	logOJ.Put("data", checkBytesToOJ(logEntry.Data))

	return logOJ
}

func logsToOJ(logEntries mj.LogList) oj.OJsonObject {
	var logList []oj.OJsonObject
	for _, logEntry := range logEntries.List {
		logOJ := logToOJ(logEntry)
		logList = append(logList, logOJ)
	}
	if logEntries.MoreAllowedAtEnd {
		logList = append(logList, stringToOJ("+"))
	}
	logOJList := oj.OJsonList(logList)
	return &logOJList
}

func intToString(i *big.Int) string {
	if i == nil {
		return ""
	}
	if i.Sign() == 0 {
		return "0x00"
	}

	isNegative := i.Sign() == -1
	str := i.Text(16)
	if isNegative {
		str = str[1:] // drop the minus in front
	}
	if len(str)%2 != 0 {
		str = "0" + str
	}
	str = "0x" + str
	if isNegative {
		str = "-" + str
	}
	return str
}

func bigIntToOJ(i mj.JSONBigInt) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func checkBigIntToOJ(i mj.JSONCheckBigInt) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func bytesFromStringToString(bytes mj.JSONBytesFromString) string {
	if len(bytes.Original) == 0 && len(bytes.Value) > 0 {
		bytes.Original = hex.EncodeToString(bytes.Value)
	}
	return bytes.Original
}

func bytesFromStringToOJ(bytes mj.JSONBytesFromString) oj.OJsonObject {
	return &oj.OJsonString{Value: bytesFromStringToString(bytes)}
}

func bytesFromTreeToOJ(bytes mj.JSONBytesFromTree) oj.OJsonObject {
	if bytes.OriginalEmpty() {
		bytes.Original = &oj.OJsonString{Value: hex.EncodeToString(bytes.Value)}
	}
	return bytes.Original
}

func checkBytesToOJ(checkBytes mj.JSONCheckBytes) oj.OJsonObject {
	if checkBytes.OriginalEmpty() && len(checkBytes.Value) > 0 {
		checkBytes.Original = &oj.OJsonString{Value: hex.EncodeToString(checkBytes.Value)}
	}
	return checkBytes.Original
}

func valueListToOJ(jsonBytesList mj.JSONValueList) oj.OJsonObject {
	var valuesList []oj.OJsonObject
	for _, blh := range jsonBytesList.Values {
		valuesList = append(valuesList, bytesFromStringToOJ(blh))
	}
	ojList := oj.OJsonList(valuesList)
	return &ojList
}

func checkValueListToOJ(jcbl mj.JSONCheckValueList) oj.OJsonObject {
	if jcbl.IsStar {
		return &oj.OJsonString{Value: "*"}
	}

	var valuesList []oj.OJsonObject
	for _, jcb := range jcbl.Values {
		valuesList = append(valuesList, checkBytesToOJ(jcb))
	}
	ojList := oj.OJsonList(valuesList)
	return &ojList
}

func uint64ToOJ(i mj.JSONUint64) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func checkUint64ToOJ(i mj.JSONCheckUint64) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func stringToOJ(str string) oj.OJsonObject {
	return &oj.OJsonString{Value: str}
}

func boolToOJ(val bool) oj.OJsonObject {
	obj := oj.OJsonBool(val)
	return &obj
}
