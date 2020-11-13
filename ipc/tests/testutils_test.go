package tests

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"

	vmcommon "github.com/ElrondNetwork/elrond-go/core/vm-common"
)

var bytecodeCounter []byte

func init() {
	bytecodeCounter = getSCCode("./../../test/contracts/counter/output/counter.wasm")
}

func getSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("Cannot read file [%s].", fileName))
	}

	return code
}

func createDeployInput(contractCode []byte) *vmcommon.ContractCreateInput {
	return &vmcommon.ContractCreateInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  []byte("me"),
			Arguments:   [][]byte{},
			CallValue:   big.NewInt(0),
			GasPrice:    100000000,
			GasProvided: 2000000,
		},
		ContractCode: contractCode,
	}
}

func createCallInput(function string) *vmcommon.ContractCallInput {
	return &vmcommon.ContractCallInput{
		VMInput: vmcommon.VMInput{
			CallerAddr:  []byte("me"),
			Arguments:   [][]byte{},
			CallValue:   big.NewInt(0),
			GasPrice:    100000000,
			GasProvided: 2000000,
		},
		RecipientAddr: []byte("mycontract"),
		Function:      function,
	}
}
