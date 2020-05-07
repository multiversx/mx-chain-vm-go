package arwenjsontest

import (
	"encoding/hex"
	"fmt"
	"math/big"

	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func convertAccount(testAcct *ij.Account) *worldhook.Account {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}

	if len(testAcct.Address.Value) != 32 {
		panic("bad test: account address should be 32 bytes long")
	}

	return &worldhook.Account{
		Exists:        true,
		Address:       testAcct.Address.Value,
		Nonce:         testAcct.Nonce.Value,
		Balance:       big.NewInt(0).Set(testAcct.Balance.Value),
		Storage:       storage,
		Code:          []byte(testAcct.Code.Value),
		AsyncCallData: testAcct.AsyncCallData,
	}
}

func convertNewAddressMocks(testNAMs []*ij.NewAddressMock) []*worldhook.NewAddressMock {
	var result []*worldhook.NewAddressMock
	for _, testNAM := range testNAMs {
		result = append(result, &worldhook.NewAddressMock{
			CreatorAddress: testNAM.CreatorAddress.Value,
			CreatorNonce:   testNAM.CreatorNonce.Value,
			NewAddress:     testNAM.NewAddress.Value,
		})
	}
	return result
}

func convertLogToTestFormat(outputLog *vmi.LogEntry) *ij.LogEntry {
	testLog := ij.LogEntry{
		Address:    ij.JSONBytes{Value: outputLog.Address},
		Identifier: ij.JSONBytes{Value: outputLog.Identifier},
		Data:       ij.JSONBytes{Value: outputLog.Data},
		Topics:     make([]ij.JSONBytes, len(outputLog.Topics)),
	}
	for i, topic := range outputLog.Topics {
		testLog.Topics[i] = ij.JSONBytes{Value: topic}
	}

	return &testLog
}

func bigIntPretty(i *big.Int) string {
	return fmt.Sprintf("0x%x (%d)", i, i)
}

func byteArrayPretty(b []byte) string {
	if len(b) == 0 {
		return "[]"
	}
	asInt := big.NewInt(0).SetBytes(b)
	return fmt.Sprintf("0x%s (%d)", hex.EncodeToString(b), asInt)
}
