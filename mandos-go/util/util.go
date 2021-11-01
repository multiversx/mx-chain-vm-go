package util

import (
	"encoding/hex"
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

var separator = []byte("@")

// CreateMultiTransferData builds data for a multiTransferESDT
func CreateMultiTransferData(to []byte, esdtData []*mj.ESDTTxData, endpointName string, arguments [][]byte) []byte {
	multiTransferData := make([]byte, 0)
	multiTransferData = append(multiTransferData, []byte(core.BuiltInFunctionMultiESDTNFTTransfer)...)
	multiTransferData = append(multiTransferData, separator...)
	encodedReceiverAddress := hex.EncodeToString(to)
	multiTransferData = append(multiTransferData, []byte(encodedReceiverAddress)...)
	multiTransferData = append(multiTransferData, separator...)

	encodedNumberOfTransfers := hex.EncodeToString(big.NewInt(int64(len(esdtData))).Bytes())
	multiTransferData = append(multiTransferData, []byte(encodedNumberOfTransfers)...)
	multiTransferData = append(multiTransferData, separator...)

	for _, esdtDataTransfer := range esdtData {
		multiTransferData = append(multiTransferData, esdtDataTransfer.TokenIdentifier.Value...)
		multiTransferData = append(multiTransferData, separator...)

		if tokenIsESDTNFT(esdtDataTransfer.Nonce.Original) {
			encodedNonceValue := hex.EncodeToString(big.NewInt(int64(esdtDataTransfer.Nonce.Value)).Bytes())
			multiTransferData = append(multiTransferData, []byte(encodedNonceValue)...)
			multiTransferData = append(multiTransferData, separator...)
		}

		encodedAmountValue := hex.EncodeToString(esdtDataTransfer.Value.Value.Bytes())
		multiTransferData = append(multiTransferData, []byte(encodedAmountValue)...)
		multiTransferData = append(multiTransferData, separator...)
	}

	if len(endpointName) > 0 {
		encodedEndpointName := hex.EncodeToString([]byte(endpointName))
		multiTransferData = append(multiTransferData, []byte(encodedEndpointName)...)
		multiTransferData = append(multiTransferData, separator...)

		for _, arg := range arguments {
			encodedArg := hex.EncodeToString(arg)
			multiTransferData = append(multiTransferData, []byte(encodedArg)...)
			multiTransferData = append(multiTransferData, separator...)
		}
	}
	return multiTransferData[:len(multiTransferData)-1]
}

func tokenIsESDTNFT(nonceStringValue string) bool {
	return nonceStringValue == ""
}
