package scenario_exporter

import (
	"math/big"

	txDataBuilder "github.com/multiversx/mx-chain-vm-common-go/txDataBuilder"
	mj "github.com/multiversx/mx-chain-vm-go/scenarios/model"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
)

const vmTypeHex = "0500"

const dummyCodeMetadataHex = "0102"

// length of "file:" in the mandos test
const contractCodePrefixLength = 5

// Transaction defines the test tranaction structure
type Transaction struct {
	function   string
	args       [][]byte
	deployData []byte
	nonce      uint64
	value      *big.Int
	esdtValue  []*mj.ESDTTxData
	sndAddr    []byte
	rcvAddr    []byte
	gasPrice   uint64
	gasLimit   uint64
}

// NewTransaction creates a new transaction instance
func NewTransaction() *Transaction {
	return &Transaction{
		args:       make([][]byte, 0),
		value:      big.NewInt(0),
		esdtValue:  make([]*mj.ESDTTxData, 0),
		sndAddr:    make([]byte, 0),
		rcvAddr:    make([]byte, 0),
		deployData: make([]byte, 0),
	}
}

// WithNonce sets the nonce
func (tx *Transaction) WithNonce(nonce uint64) *Transaction {
	tx.nonce = nonce
	return tx
}

// GetNonce gets the nonce
func (tx *Transaction) GetNonce() uint64 {
	return tx.nonce
}

// WithCallValue sets the call value
func (tx *Transaction) WithCallValue(value *big.Int) *Transaction {
	tx.value.Set(value)
	return tx
}

// GetCallValue gets the call value
func (tx *Transaction) GetCallValue() *big.Int {
	return tx.value
}

// WithESDTTransfers sets the ESDT transafers
func (tx *Transaction) WithESDTTransfers(esdtTransfers []*mj.ESDTTxData) *Transaction {
	tx.esdtValue = append(tx.esdtValue, esdtTransfers...)
	return tx
}

// GetESDTTransfers gets the ESDT transfers
func (tx *Transaction) GetESDTTransfers() []*mj.ESDTTxData {
	return tx.esdtValue
}

// WithCallFunction sets the call function
func (tx *Transaction) WithCallFunction(functionName string) *Transaction {
	tx.function = functionName
	return tx
}

// GetCallFunction gets the call function
func (tx *Transaction) GetCallFunction() string {
	return tx.function
}

// WithCallArguments sets the call arguments
func (tx *Transaction) WithCallArguments(arguments [][]byte) *Transaction {
	tx.args = append(tx.args, arguments...)
	return tx
}

// GetCallArguments gets the call arguments
func (tx *Transaction) GetCallArguments() [][]byte {
	return tx.args
}

// WithSenderAddress sets the sender address
func (tx *Transaction) WithSenderAddress(address []byte) *Transaction {
	tx.sndAddr = make([]byte, len(address))
	copy(tx.sndAddr, address)
	return tx
}

// GetSenderAddress gets the sender address
func (tx *Transaction) GetSenderAddress() []byte {
	return tx.sndAddr
}

// WithReceiverAddress sets the receiver address
func (tx *Transaction) WithReceiverAddress(address []byte) *Transaction {
	tx.rcvAddr = make([]byte, len(address))
	copy(tx.rcvAddr, address)
	return tx
}

// GetReceiverAddress gets the receiver address
func (tx *Transaction) GetReceiverAddress() []byte {
	return tx.rcvAddr
}

// WithGasLimitAndPrice sets the gas limit & gas price
func (tx *Transaction) WithGasLimitAndPrice(gasLimit, gasPrice uint64) *Transaction {
	tx.gasLimit = gasLimit
	tx.gasPrice = gasPrice
	return tx
}

// GetGasLimitAndPrice gets the gas limit & gas price
func (tx *Transaction) GetGasLimitAndPrice() (uint64, uint64) {
	return tx.gasLimit, tx.gasPrice
}

// WithDeployData sets the deploy data: sc code + arguments
func (tx *Transaction) WithDeployData(scCodePath string, args [][]byte) *Transaction {
	deployData := createDeployTxData(scCodePath, args)
	tx.deployData = append(tx.deployData, deployData...)
	return tx
}

func createDeployTxData(scCodePath string, args [][]byte) []byte {
	scCode := vmhost.GetSCCode(scCodePath[contractCodePrefixLength:])
	tdb := txDataBuilder.NewBuilder()
	tdb.Bytes(scCode)
	tdb.Bytes([]byte(vmTypeHex))
	tdb.Bytes([]byte(dummyCodeMetadataHex))
	if args != nil {
		for i := 0; i < len(args); i++ {
			tdb.Bytes(args[i])
		}
	}
	return tdb.ToBytes()
}

func (tx *Transaction) GetDeployData() []byte {
	return tx.deployData
}

// CreateTransaction will create a transaction based on the parameters provided
func CreateTransaction(
	function string,
	args [][]byte,
	nonce uint64,
	value *big.Int,
	esdtTransfers []*mj.ESDTTxData,
	sndAddr []byte,
	rcvAddr []byte,
	gasLimit uint64,
	gasPrice uint64,
) *Transaction {
	return NewTransaction().
		WithCallFunction(function).
		WithCallArguments(args).
		WithNonce(nonce).
		WithCallValue(value).
		WithESDTTransfers(esdtTransfers).
		WithSenderAddress(sndAddr).
		WithReceiverAddress(rcvAddr).
		WithGasLimitAndPrice(gasLimit, gasPrice)
}

// CreateDeployTransaction creates a deploy transaction
func CreateDeployTransaction(
	args [][]byte,
	scCodePath string,
	sndAddr []byte,
	gasLimit uint64,
	gasPrice uint64,
) *Transaction {
	return NewTransaction().
		WithDeployData(scCodePath, args).
		WithSenderAddress(sndAddr).
		WithGasLimitAndPrice(gasLimit, gasPrice)
}
