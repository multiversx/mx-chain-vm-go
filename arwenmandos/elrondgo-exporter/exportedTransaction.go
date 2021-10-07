package elrondgo_exporter

import (
	"math/big"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
)

type transaction struct {
	isDeploy   bool
	deployPath string

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

func NewTransaction() *transaction {
	return &transaction{
		args:      make([][]byte, 0),
		value:     big.NewInt(0),
		esdtValue: make([]*mj.ESDTTxData, 0),
		sndAddr:   make([]byte, 0),
		rcvAddr:   make([]byte, 0),
	}
}

func (tx *transaction) WithDeployWithPath(path string) *transaction {
	tx.isDeploy = true
	tx.deployPath = path
	return tx
}

func (tx *transaction) IsDeploy() bool {
	return tx.isDeploy
}

func (tx *transaction) GetDeployPath() string {
	return tx.deployPath
}

func (tx *transaction) WithNonce(nonce uint64) *transaction {
	tx.nonce = nonce
	return tx
}

func (tx *transaction) GetNonce() uint64 {
	return tx.nonce
}

func (tx *transaction) WithCallValue(value *big.Int) *transaction {
	tx.value.Set(value)
	return tx
}

func (tx *transaction) GetCallValue() *big.Int {
	return tx.value
}

func (tx *transaction) WithESDTTransfers(esdtTransfers []*mj.ESDTTxData) *transaction {
	copy(tx.esdtValue, esdtTransfers)
	return tx
}

func (tx *transaction) GetESDTTransfers() []*mj.ESDTTxData {
	return tx.esdtValue
}

func (tx *transaction) WithCallFunction(functionName string) *transaction {
	tx.function = functionName
	return tx
}

func (tx *transaction) GetCallFunction() string {
	return tx.function
}

func (tx *transaction) WithCallArguments(arguments [][]byte) *transaction {
	copy(tx.args, arguments)
	return tx
}

func (tx *transaction) GetCallArguments() [][]byte {
	return tx.args
}

func (tx *transaction) WithSenderAddress(address []byte) *transaction {
	copy(tx.sndAddr, address)
	return tx
}

func (tx *transaction) GetSenderAddress() []byte {
	return tx.sndAddr
}

func (tx *transaction) WithReceiverAddress(address []byte) *transaction {
	copy(tx.rcvAddr, address)
	return tx
}

func (tx *transaction) GetReceiverAddress() []byte {
	return tx.rcvAddr
}

func (tx *transaction) WithGasLimitAndPrice(gasLimit, gasPrice uint64) *transaction {
	tx.gasLimit = gasLimit
	tx.gasPrice = gasPrice
	return tx
}

func (tx *transaction) GetGasLimitAndPrice() (uint64, uint64) {
	return tx.gasLimit, tx.gasPrice
}

func CreateDeployTransaction(deployPath string, args [][]byte, nonce uint64, value *big.Int, sndAddr []byte, gasLimit uint64, gasPrice uint64) *transaction {
	return NewTransaction().WithDeployWithPath(deployPath[5:]).WithCallArguments(args).WithNonce(nonce).WithCallValue(value).WithSenderAddress(sndAddr).WithGasLimitAndPrice(gasLimit, gasPrice)
}

func CreateTransaction(function string, args [][]byte, nonce uint64, value *big.Int, esdtTransfers []*mj.ESDTTxData, sndAddr []byte, rcvAddr []byte, gasLimit uint64, gasPrice uint64) *transaction {
	return NewTransaction().WithCallFunction(function).WithCallArguments(args).WithNonce(nonce).WithCallValue(value).WithESDTTransfers(esdtTransfers).WithSenderAddress(sndAddr).WithReceiverAddress(rcvAddr).WithGasLimitAndPrice(gasLimit, gasPrice)
}
