package util

import (
	mj "github.com/ElrondNetwork/wasm-vm/mandos-go/model"
	"github.com/ElrondNetwork/elrond-go-core/core"
	txDataBuilder "github.com/ElrondNetwork/elrond-vm-common/txDataBuilder"
)

// CreateMultiTransferData builds data for a multiTransferESDT
func CreateMultiTransferData(to []byte, esdtData []*mj.ESDTTxData, endpointName string, arguments [][]byte) []byte {

	multiTransferData := make([]byte, 0)
	multiTransferData = append(multiTransferData, []byte(core.BuiltInFunctionMultiESDTNFTTransfer)...)
	tdb := txDataBuilder.NewBuilder()
	tdb.Bytes(to)
	tdb.Int(len(esdtData))

	for _, esdtDataTransfer := range esdtData {
		tdb.Bytes(esdtDataTransfer.TokenIdentifier.Value)
		tdb.Int64(int64(esdtDataTransfer.Nonce.Value))
		tdb.BigInt(esdtDataTransfer.Value.Value)
	}

	if len(endpointName) > 0 {
		tdb.Str(endpointName)

		for _, arg := range arguments {
			tdb.Bytes(arg)
		}
	}
	multiTransferData = append(multiTransferData, tdb.ToBytes()...)
	return multiTransferData
}
