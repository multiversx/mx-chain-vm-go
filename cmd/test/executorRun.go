package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"path/filepath"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen/context"
	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
	cryptohook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-crypto"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

const ignoreGas = true
const ignoreAllLogs = false

type arwenTestExecutor struct {
	world                    *worldhook.BlockchainHookMock
	vm                       vmi.VMExecutionHandler
	contractPathReplacements map[string]string
	checkGas                 bool
}

func newArwenTestExecutor() *arwenTestExecutor {
	world := worldhook.NewMock()
	world.EnableMockAddressGeneration()

	blockGasLimit := uint64(10000000)
	gasSchedule := config.MakeGasMap(1)
	vm, err := arwen.NewArwenVM(world, cryptohook.KryptoHookMockInstance, TestVMType, blockGasLimit, gasSchedule)
	if err != nil {
		panic(err)
	}
	return &arwenTestExecutor{
		world:                    world,
		vm:                       vm,
		contractPathReplacements: make(map[string]string),
		checkGas:                 false,
	}
}

func (te *arwenTestExecutor) replaceCode(pathInTest, actualPath string) *arwenTestExecutor {
	te.contractPathReplacements[pathInTest] = actualPath
	return te
}

// ProcessCode takes the contract file path, assembles it and yields the bytecode.
func (te *arwenTestExecutor) ProcessCode(testPath string, value string) (string, error) {
	if len(value) == 0 {
		return "", nil
	}
	var fullPath string
	if replacement, shouldReplace := te.contractPathReplacements[value]; shouldReplace {
		fullPath = replacement
	} else {
		fullPath = filepath.Join(testPath, value)
	}
	scCode, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return string(scCode), nil
}

// Run executes an individual test.
func (te *arwenTestExecutor) Run(test *ij.Test) error {
	world := te.world
	vm := te.vm

	// reset world
	world.Clear()
	world.Blockhashes = test.BlockHashes

	for _, acct := range test.Pre {
		world.AcctMap.PutAccount(convertAccount(acct))
	}

	//spew.Dump(world.AcctMap)

	for _, block := range test.Blocks {
		for txIndex, tx := range block.Transactions {
			//fmt.Printf("%d\n", txIndex)
			beforeErr := world.UpdateWorldStateBefore(tx.From, tx.GasLimit, tx.GasPrice)
			if beforeErr != nil {
				return beforeErr
			}

			arguments := make([][]byte, len(tx.Arguments))
			for i, arg := range tx.Arguments {
				arguments[i] = append(arguments[i], arg.ToBytes()...)
			}
			var output *vmi.VMOutput

			if tx.IsCreate {
				input := &vmi.ContractCreateInput{
					ContractCode: []byte(tx.AssembledCode),
					VMInput: vmi.VMInput{
						CallerAddr:  tx.From,
						Arguments:   arguments,
						CallValue:   tx.Value,
						GasPrice:    tx.GasPrice,
						GasProvided: tx.GasLimit,
					},
				}

				var err error
				output, err = vm.RunSmartContractCreate(input)
				if err != nil {
					return err
				}
			} else {
				input := &vmi.ContractCallInput{
					RecipientAddr: tx.To,
					Function:      tx.Function,
					VMInput: vmi.VMInput{
						CallerAddr:  tx.From,
						Arguments:   arguments,
						CallValue:   tx.Value,
						GasPrice:    tx.GasPrice,
						GasProvided: tx.GasLimit,
					},
				}

				var err error
				output, err = vm.RunSmartContractCall(input)
				if err != nil {
					return err
				}
			}

			// subtract call value from sender (this is not reflected in the delta)
			world.UpdateBalanceWithDelta(tx.From, big.NewInt(0).Neg(tx.Value))

			// update accounts based on deltas
			updErr := world.UpdateAccounts(output.OutputAccounts, output.DeletedAccounts)
			if updErr != nil {
				return updErr
			}

			blResult := block.Results[txIndex]

			// sum of all balance deltas should equal call value
			sumOfBalanceDeltas := big.NewInt(0)
			for _, oa := range output.OutputAccounts {
				sumOfBalanceDeltas = sumOfBalanceDeltas.Add(sumOfBalanceDeltas, oa.BalanceDelta)
			}
			if sumOfBalanceDeltas.Cmp(tx.Value) != 0 {
				return fmt.Errorf("sum of balance deltas should equal call value. Sum of balance deltas: %d (0x%x). Call value: %d (0x%x)",
					sumOfBalanceDeltas, sumOfBalanceDeltas, tx.Value, tx.Value)
			}

			// check return code
			expectedStatus := 0
			if blResult.Status != nil {
				expectedStatus = int(blResult.Status.Int64())
			}
			if expectedStatus != int(output.ReturnCode) {
				return fmt.Errorf("result code mismatch. Tx #%d. Want: %d. Have: %d", txIndex, expectedStatus, int(output.ReturnCode))
			}

			// check result
			if len(output.ReturnData) != len(blResult.Out) {
				return fmt.Errorf("result length mismatch. Tx #%d. Want: %s. Have: %s",
					txIndex, ij.ResultAsString(blResult.Out), ij.ResultAsString(output.ReturnData))
			}
			for i, expected := range blResult.Out {
				if !ij.ResultEqual(expected, output.ReturnData[i]) {
					return fmt.Errorf("result mismatch. Tx #%d. Want: %s. Have: %s",
						txIndex, ij.ResultAsString(blResult.Out), ij.ResultAsString(output.ReturnData))
				}
			}

			// check refund
			if !ignoreGas && blResult.Refund != nil {
				if blResult.Refund.Cmp(output.GasRefund) != 0 {
					return fmt.Errorf("result gas refund mismatch. Want: 0x%x. Have: 0x%x",
						blResult.Refund, output.GasRefund)
				}
			}

			// check gas
			if te.checkGas && test.CheckGas && blResult.CheckGas {
				if blResult.Gas != output.GasRemaining {
					return fmt.Errorf("result gas mismatch. Want: %d (0x%x). Got: %d (0x%x)",
						blResult.Gas, blResult.Gas, output.GasRemaining, output.GasRemaining)
				}
			}
			// burned := big.NewInt(0).Sub(tx.GasLimit, output.GasRemaining)
			// fmt.Printf("all: 0x%x  remaining: 0x%x  consumed: 0x%x   refund: 0x%x\n", tx.GasLimit, output.GasRemaining, burned, output.GasRefund)

			// check empty logs, this seems to be the value
			if blResult.IgnoreLogs || ignoreAllLogs {
				// nothing, ignore
			} else {
				// this is the real log check
				if len(blResult.Logs) != len(output.Logs) {
					return fmt.Errorf("wrong number of logs. Want:%d. Got:%d",
						len(blResult.Logs), len(output.Logs))
				}
				for i, outLog := range output.Logs {
					testLog := blResult.Logs[i]
					if !bytes.Equal(outLog.Address, testLog.Address) {
						return fmt.Errorf("bad log address. Want:\n%s\nGot:\n%s",
							ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
					}
					if len(outLog.Topics) != len(testLog.Topics) {
						return fmt.Errorf("wrong number of log topics. Want:\n%s\nGot:\n%s",
							ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
					}
					for ti := range outLog.Topics {
						if !bytes.Equal(outLog.Topics[ti], testLog.Topics[ti]) {
							return fmt.Errorf("bad log topic. Want:\n%s\nGot:\n%s",
								ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
						}
					}
					if big.NewInt(0).SetBytes(outLog.Data).Cmp(big.NewInt(0).SetBytes(testLog.Data)) != 0 {
						return fmt.Errorf("bad log data. Want:\n%s\nGot:\n%s",
							ij.LogToString(testLog), ij.LogToString(convertLogToTestFormat(outLog)))
					}
				}
			}
		}
	}

	for worldAcctAddr := range world.AcctMap {
		postAcctMatch := ij.FindAccount(test.PostState, []byte(worldAcctAddr))
		if postAcctMatch == nil {
			return fmt.Errorf("unexpected account address: %s", hex.EncodeToString([]byte(worldAcctAddr)))
		}
	}

	for _, postAcctFromTest := range test.PostState {
		postAcct := convertAccount(postAcctFromTest)
		matchingAcct, isMatch := world.AcctMap[string(postAcct.Address)]
		if !isMatch {
			return fmt.Errorf("account %s expected but not found after running test",
				hex.EncodeToString(postAcct.Address))
		}

		if !bytes.Equal(matchingAcct.Address, postAcct.Address) {
			return fmt.Errorf("bad account address %s", hex.EncodeToString(matchingAcct.Address))
		}

		if matchingAcct.Nonce != postAcct.Nonce {
			return fmt.Errorf("bad account nonce. Account: %s. Want: %d. Have: %d",
				hex.EncodeToString(matchingAcct.Address), postAcct.Nonce, matchingAcct.Nonce)
		}

		if matchingAcct.Balance.Cmp(postAcct.Balance) != 0 {
			return fmt.Errorf("bad account balance. Account: %s. Want: 0x%x. Have: 0x%x",
				hex.EncodeToString(matchingAcct.Address), postAcct.Balance, matchingAcct.Balance)
		}

		if !bytes.Equal(matchingAcct.Code, postAcct.Code) {
			return fmt.Errorf("bad account code. Account: %s. Want: [%s]. Have: [%s]",
				hex.EncodeToString(matchingAcct.Address), postAcct.Code, matchingAcct.Code)
		}

		// compare storages
		allKeys := make(map[string]bool)
		for k := range postAcct.Storage {
			allKeys[k] = true
		}
		for k := range matchingAcct.Storage {
			allKeys[k] = true
		}
		storageError := ""
		for k := range allKeys {
			want := postAcct.StorageValue(k)
			have := matchingAcct.StorageValue(k)
			if !bytes.Equal(want, have) {
				storageError += fmt.Sprintf(
					"\n  for key %s: Want: 0x%s. Have: 0x%s",
					hex.EncodeToString([]byte(k)), hex.EncodeToString(want), hex.EncodeToString(have))
			}
		}
		if len(storageError) > 0 {
			return fmt.Errorf("wrong account storage for account 0x%s:%s",
				hex.EncodeToString(postAcct.Address), storageError)
		}
	}

	return nil
}
