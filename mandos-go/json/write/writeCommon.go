package mandosjsonwrite

import (
	"encoding/hex"
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func accountsToOJ(accounts []*mj.Account) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, account := range accounts {
		acctOJ := oj.NewMap()
		if len(account.Comment) > 0 {
			acctOJ.Put("comment", stringToOJ(account.Comment))
		}
		acctOJ.Put("nonce", uint64ToOJ(account.Nonce))
		acctOJ.Put("balance", bigIntToOJ(account.Balance))
		storageOJ := oj.NewMap()
		for _, st := range account.Storage {
			storageOJ.Put(bytesFromStringToString(st.Key), bytesFromTreeToOJ(st.Value))
		}
		acctOJ.Put("storage", storageOJ)
		esdtOJ := oj.NewMap()
		for _, esdtItem := range account.ESDTData {
			if len(esdtItem.Frozen.Original) == 0 {
				esdtOJ.Put(bytesFromStringToString(esdtItem.TokenName), bigIntToOJ(esdtItem.Balance))
			} else {
				esdtItemOJ := oj.NewMap()
				esdtItemOJ.Put("balance", bigIntToOJ(esdtItem.Balance))
				esdtItemOJ.Put("frozen", uint64ToOJ(esdtItem.Frozen))
				esdtOJ.Put(bytesFromStringToString(esdtItem.TokenName), esdtItemOJ)
			}
		}
		if esdtOJ.Size() > 0 {
			acctOJ.Put("esdt", esdtOJ)
		}
		acctOJ.Put("code", bytesFromStringToOJ(account.Code))
		if len(account.AsyncCallData) > 0 {
			acctOJ.Put("asyncCallData", stringToOJ(account.AsyncCallData))
		}

		acctsOJ.Put(bytesFromStringToString(account.Address), acctOJ)
	}

	return acctsOJ
}

func checkAccountsToOJ(checkAccounts *mj.CheckAccounts) oj.OJsonObject {
	acctsOJ := oj.NewMap()
	for _, checkAccount := range checkAccounts.Accounts {
		acctOJ := oj.NewMap()
		if len(checkAccount.Comment) > 0 {
			acctOJ.Put("comment", stringToOJ(checkAccount.Comment))
		}
		if !checkAccount.Nonce.IsUnspecified() {
			acctOJ.Put("nonce", checkUint64ToOJ(checkAccount.Nonce))
		}
		if !checkAccount.Balance.IsUnspecified() {
			acctOJ.Put("balance", checkBigIntToOJ(checkAccount.Balance))
		}
		storageOJ := oj.NewMap()
		for _, st := range checkAccount.CheckStorage {
			storageOJ.Put(bytesFromStringToString(st.Key), bytesFromTreeToOJ(st.Value))
		}
		if checkAccount.IgnoreStorage {
			acctOJ.Put("storage", stringToOJ("*"))
		} else {
			acctOJ.Put("storage", storageOJ)
		}
		if !checkAccount.Code.IsUnspecified() {
			acctOJ.Put("code", checkBytesToOJ(checkAccount.Code))
		}
		if !checkAccount.AsyncCallData.IsUnspecified() {
			acctOJ.Put("asyncCallData", checkBytesToOJ(checkAccount.AsyncCallData))
		}

		acctsOJ.Put(bytesFromStringToString(checkAccount.Address), acctOJ)
	}

	if checkAccounts.OtherAccountsAllowed {
		acctsOJ.Put("+", stringToOJ(""))
	}

	return acctsOJ
}

func blockHashesToOJ(blockHashes []mj.JSONBytesFromString) oj.OJsonObject {
	var blockhashesList []oj.OJsonObject
	for _, blh := range blockHashes {
		blockhashesList = append(blockhashesList, bytesFromStringToOJ(blh))
	}
	blockhashesOJ := oj.OJsonList(blockhashesList)
	return &blockhashesOJ
}

func resultToOJ(res *mj.TransactionResult) oj.OJsonObject {
	resultOJ := oj.NewMap()

	var outList []oj.OJsonObject
	for _, out := range res.Out {
		outList = append(outList, checkBytesToOJ(out))
	}
	outOJ := oj.OJsonList(outList)
	resultOJ.Put("out", &outOJ)

	if !res.Status.IsUnspecified() {
		resultOJ.Put("status", checkBigIntToOJ(res.Status))
	}
	if !res.Message.IsUnspecified() {
		resultOJ.Put("message", checkBytesToOJ(res.Message))
	}
	if !res.LogsUnspecified {
		if res.LogsStar {
			resultOJ.Put("logs", stringToOJ("*"))
		} else {
			if len(res.LogHash) > 0 {
				resultOJ.Put("logs", stringToOJ(res.LogHash))
			} else {
				resultOJ.Put("logs", logsToOJ(res.Logs))
			}
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
	logOJ.Put("identifier", checkBytesToOJ(logEntry.Identifier))

	var topicsList []oj.OJsonObject
	for _, topic := range logEntry.Topics {
		topicsList = append(topicsList, checkBytesToOJ(topic))
	}
	topicsOJ := oj.OJsonList(topicsList)
	logOJ.Put("topics", &topicsOJ)

	logOJ.Put("data", checkBytesToOJ(logEntry.Data))

	return logOJ
}

func logsToOJ(logEntries []*mj.LogEntry) oj.OJsonObject {
	var logList []oj.OJsonObject
	for _, logEntry := range logEntries {
		logOJ := logToOJ(logEntry)
		logList = append(logList, logOJ)
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

func uint64ToOJ(i mj.JSONUint64) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func checkUint64ToOJ(i mj.JSONCheckUint64) oj.OJsonObject {
	return &oj.OJsonString{Value: i.Original}
}

func stringToOJ(str string) oj.OJsonObject {
	return &oj.OJsonString{Value: str}
}
