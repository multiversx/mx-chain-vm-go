package arwenmandos

import (
	"encoding/hex"
	"fmt"
	"math/big"

	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
	mj "github.com/ElrondNetwork/elrond-vm-util/test-util/mandos/json/model"
)

func convertAccount(testAcct *mj.Account) *worldhook.Account {
	storage := make(map[string][]byte)
	for _, stkvp := range testAcct.Storage {
		key := string(stkvp.Key.Value)
		storage[key] = stkvp.Value.Value
	}

	if len(testAcct.Address.Value) != 32 {
		panic("bad test: account address should be 32 bytes long")
	}

	return &worldhook.Account{
		Address:       testAcct.Address.Value,
		Nonce:         testAcct.Nonce.Value,
		Balance:       big.NewInt(0).Set(testAcct.Balance.Value),
		Storage:       storage,
		Code:          []byte(testAcct.Code.Value),
		AsyncCallData: testAcct.AsyncCallData,
	}
}

func convertNewAddressMocks(testNAMs []*mj.NewAddressMock) []*worldhook.NewAddressMock {
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

func convertBlockInfo(testBlockInfo *mj.BlockInfo) *worldhook.BlockInfo {
	if testBlockInfo == nil {
		return nil
	}
	result := &worldhook.BlockInfo{
		BlockTimestamp: testBlockInfo.BlockTimestamp.Value,
		BlockNonce:     testBlockInfo.BlockNonce.Value,
		BlockRound:     testBlockInfo.BlockRound.Value,
		BlockEpoch:     uint32(testBlockInfo.BlockEpoch.Value),
	}

	return result
}

func convertLogToTestFormat(outputLog *vmi.LogEntry) *mj.LogEntry {
	testLog := mj.LogEntry{
		Address:    mj.JSONBytes{Value: outputLog.Address},
		Identifier: mj.JSONBytes{Value: outputLog.Identifier},
		Data:       mj.JSONBytes{Value: outputLog.Data},
		Topics:     make([]mj.JSONBytes, len(outputLog.Topics)),
	}
	for i, topic := range outputLog.Topics {
		testLog.Topics[i] = mj.JSONBytes{Value: topic}
	}

	return &testLog
}

func bigIntPretty(i *big.Int) string {
	return fmt.Sprintf("0x%x (%d)", i, i)
}

func byteArrayPretty(bytes []byte) string {
	if len(bytes) == 0 {
		return "[]"
	}

	if canInterpretAsString(bytes) {
		return fmt.Sprintf("0x%s (``%s)", hex.EncodeToString(bytes), string(bytes))
	}

	asInt := big.NewInt(0).SetBytes(bytes)
	return fmt.Sprintf("0x%s (%d)", hex.EncodeToString(bytes), asInt)
}

func canInterpretAsString(bytes []byte) bool {
	if len(bytes) == 0 {
		return false
	}
	for _, b := range bytes {
		if b < 32 || b > 126 {
			return false
		}
	}
	return true
}
