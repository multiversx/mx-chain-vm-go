package util

import (
	"math/big"
	"testing"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	"github.com/stretchr/testify/require"
)

func Test_CreateMultiTransferData_SingleTransfer(t *testing.T) {
	esdtTransfers := make([]*mj.ESDTTxData, 0)
	esdtTransfer := &mj.ESDTTxData{
		Nonce:           mj.JSONUint64{Value: 2},
		TokenIdentifier: mj.JSONBytesFromString{Value: []byte("TOK1-abcdef")},
		Value:           mj.JSONBigInt{Value: big.NewInt(100)},
	}
	esdtTransfers = append(esdtTransfers, esdtTransfer)
	data := CreateMultiTransferData(
		[]byte("receiverAddress"),
		esdtTransfers, "function1",
		[][]byte{
			[]byte("arg1"),
			[]byte("arg2")},
	)
	require.Equal(t, "MultiESDTNFTTransfer@receiverAddress@01@TOK1-abcdef@02@64@66756e6374696f6e31@61726731@61726732", string(data))
}

func Test_CreateMultiTransferData_MultipleTransfers(t *testing.T) {
	esdtTransfers := make([]*mj.ESDTTxData, 0)
	esdtTransfer1 := &mj.ESDTTxData{
		Nonce:           mj.JSONUint64{Value: 2},
		TokenIdentifier: mj.JSONBytesFromString{Value: []byte("TOK1-abcdef")},
		Value:           mj.JSONBigInt{Value: big.NewInt(100)},
	}
	esdtTransfer2 := &mj.ESDTTxData{
		Nonce:           mj.JSONUint64{Value: 14},
		TokenIdentifier: mj.JSONBytesFromString{Value: []byte("TOK2-abcdef")},
		Value:           mj.JSONBigInt{Value: big.NewInt(396)},
	}

	esdtTransfers = append(esdtTransfers, esdtTransfer1, esdtTransfer2)
	data := CreateMultiTransferData(
		[]byte("receiverAddress"),
		esdtTransfers, "function1",
		[][]byte{
			[]byte("arg1"),
			[]byte("arg2")},
	)
	require.Equal(t, "MultiESDTNFTTransfer@receiverAddress@02@TOK1-abcdef@02@64@TOK2-abcdef@0e@018c@66756e6374696f6e31@61726731@61726732", string(data))
}
