package elrondgo_exporter

import (
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
)

type Transaction struct {
	function  string
	args      [][]byte
	nonce     uint64
	value     *big.Int
	esdtValue []*mj.ESDTTxData
	sndAddr   []byte
	rcvAddr   []byte
	gasPrice  uint64
	gasLimit  uint64
}

func NewTransaction() *Transaction {
	return &Transaction{
		args:      make([][]byte, 0),
		value:     big.NewInt(0),
		esdtValue: make([]*mj.ESDTTxData, 0),
		sndAddr:   make([]byte, 0),
		rcvAddr:   make([]byte, 0),
	}
}

func (tx *Transaction) WithNonce(nonce uint64) *Transaction {
	tx.nonce = nonce
	return tx
}

func (tx *Transaction) GetNonce() uint64 {
	return tx.nonce
}

func (tx *Transaction) WithCallValue(value *big.Int) *Transaction {
	tx.value.Set(value)
	return tx
}

func (tx *Transaction) GetCallValue() *big.Int {
	return tx.value
}

func (tx *Transaction) WithESDTTransfers(esdtTransfers []*mj.ESDTTxData) *Transaction {
	tx.esdtValue = append(tx.esdtValue, esdtTransfers...)
	return tx
}

func (tx *Transaction) GetESDTTransfers() []*mj.ESDTTxData {
	return tx.esdtValue
}

func (tx *Transaction) WithCallFunction(functionName string) *Transaction {
	tx.function = functionName
	return tx
}

func (tx *Transaction) GetCallFunction() string {
	return tx.function
}

func (tx *Transaction) WithCallArguments(arguments [][]byte) *Transaction {
	tx.args = append(tx.args, arguments...)
	return tx
}

func (tx *Transaction) GetCallArguments() [][]byte {
	return tx.args
}

func (tx *Transaction) WithSenderAddress(address []byte) *Transaction {
	tx.sndAddr = append(tx.sndAddr, address...)
	return tx
}

func (tx *Transaction) GetSenderAddress() []byte {
	return tx.sndAddr
}

func (tx *Transaction) WithReceiverAddress(address []byte) *Transaction {
	tx.rcvAddr = append(tx.rcvAddr, address...)
	return tx
}

func (tx *Transaction) GetReceiverAddress() []byte {
	return tx.rcvAddr
}

func (tx *Transaction) WithGasLimitAndPrice(gasLimit, gasPrice uint64) *Transaction {
	tx.gasLimit = gasLimit
	tx.gasPrice = gasPrice
	return tx
}

func (tx *Transaction) GetGasLimitAndPrice() (uint64, uint64) {
	return tx.gasLimit, tx.gasPrice
}

func CreateTransaction(function string, args [][]byte, nonce uint64, value *big.Int, esdtTransfers []*mj.ESDTTxData, sndAddr []byte, rcvAddr []byte, gasLimit uint64, gasPrice uint64) *Transaction {
	return NewTransaction().WithCallFunction(function).WithCallArguments(args).WithNonce(nonce).WithCallValue(value).WithESDTTransfers(esdtTransfers).WithSenderAddress(sndAddr).WithReceiverAddress(rcvAddr).WithGasLimitAndPrice(gasLimit, gasPrice)
}