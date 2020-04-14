package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"testing"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	arwenHost "github.com/ElrondNetwork/arwen-wasm-vm/arwen/host"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	cryptohook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-crypto"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

type pureFunctionIO struct {
	functionName    string
	arguments       [][]byte
	expectedStatus  vmi.ReturnCode
	expectedMessage string
	expectedResults [][]byte
}

func pureFunctionTest(t *testing.T, testCases []*pureFunctionIO) {

	contractPathFilePath := filepath.Join(getTestRoot(), "features/features.wasm")
	scCode, err := ioutil.ReadFile(contractPathFilePath)
	if err != nil {
		panic(err)
	}

	world := worldhook.NewMock()
	world.EnableMockAddressGeneration()

	contractAddrHex := "c0879ac700000000000000000000000000000000000000000000000000000000"
	account1AddrHex := "acc1000000000000000000000000000000000000000000000000000000000000"

	contractAddr, _ := hex.DecodeString(contractAddrHex)
	account1Addr, _ := hex.DecodeString(account1AddrHex)

	world.AcctMap.PutAccount(&worldhook.Account{
		Exists:  true,
		Address: contractAddr,
		Nonce:   0,
		Balance: big.NewInt(0),
		Storage: make(map[string][]byte),
		Code:    scCode,
	})

	world.AcctMap.PutAccount(&worldhook.Account{
		Exists:  true,
		Address: account1Addr,
		Nonce:   0,
		Balance: big.NewInt(0x100000000),
		Storage: make(map[string][]byte),
		Code:    []byte{},
	})

	// create VM
	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)
	vm, err := arwenHost.NewArwenVM(world, cryptohook.KryptoHookMockInstance, &arwen.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasSchedule,
		ProtocolBuiltinFunctions: make(vmcommon.FunctionNames),
	})
	if err != nil {
		panic(err)
	}

	// RUN!
	for _, testCase := range testCases {

		input := &vmi.ContractCallInput{
			RecipientAddr: contractAddr,
			Function:      testCase.functionName,
			VMInput: vmi.VMInput{
				CallerAddr:  account1Addr,
				Arguments:   testCase.arguments,
				CallValue:   big.NewInt(0),
				GasPrice:    1,
				GasProvided: 100000000,
			},
		}

		output, err := vm.RunSmartContractCall(input)
		if err != nil {
			panic(err)
		}

		if output.ReturnCode != testCase.expectedStatus {
			t.Error(fmt.Errorf("result code mismatch. Want: %d. Have: %d (%s). Message: %s",
				int(testCase.expectedStatus), int(output.ReturnCode), output.ReturnCode.String(), output.ReturnMessage))
		}

		if output.ReturnMessage != testCase.expectedMessage {
			t.Error(fmt.Errorf("result message mismatch. Want: %s. Have: %s",
				testCase.expectedMessage, output.ReturnMessage))
		}

		// check result
		if len(output.ReturnData) != len(testCase.expectedResults) {
			t.Error(fmt.Errorf("result length mismatch. Want: %s. Have: %s",
				ij.ResultAsString(testCase.expectedResults), ij.ResultAsString(output.ReturnData)))
		}
		for i, expected := range testCase.expectedResults {
			if !ij.ResultEqual(expected, output.ReturnData[i]) {
				t.Error(fmt.Errorf("result mismatch. Want: %s. Have: %s",
					ij.ResultAsString(testCase.expectedResults), ij.ResultAsString(output.ReturnData)))
			}
		}

	}

	//fmt.Println("Returned: " + string(output.ReturnData[0].Bytes()))

	//lastReturnCode = output.ReturnCode

}
