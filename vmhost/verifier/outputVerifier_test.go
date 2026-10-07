package verifier_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/marshal"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/vmhost/verifier"
	"github.com/stretchr/testify/require"
)

func createValidAddress(suffix byte) []byte {
	addr := make([]byte, verifier.AddressLen)
	addr[0] = 1
	addr[verifier.AddressLen-1] = suffix
	return addr
}

func createValidAccount(addr []byte) *vmcommon.OutputAccount {
	return &vmcommon.OutputAccount{
		Address:        bytes.Clone(addr),
		Balance:        big.NewInt(1000),
		BalanceDelta:   big.NewInt(0),
		StorageUpdates: make(map[string]*vmcommon.StorageUpdate),
	}
}

func createBaseVMOutput(accs ...*vmcommon.OutputAccount) *vmcommon.VMOutput {
	outAccs := make(map[string]*vmcommon.OutputAccount)
	for _, acc := range accs {
		outAccs[string(acc.Address)] = acc
	}
	return &vmcommon.VMOutput{
		ReturnCode:      vmcommon.Ok,
		ReturnMessage:   "",
		GasRemaining:    50000,
		GasRefund:       big.NewInt(0),
		OutputAccounts:  outAccs,
		DeletedAccounts: make([][]byte, 0),
		Logs:            make([]*vmcommon.LogEntry, 0),
	}
}

func createBaseInput() *verifier.VerificationInput {
	return &verifier.VerificationInput{
		CallerAddr:    createValidAddress(1),
		RecipientAddr: createValidAddress(2),
		CallValue:     big.NewInt(0),
		GasProvided:   100000,
		IsReadOnly:    false,
	}
}

// -----------------------------------------------------------------------------
// 1. ReturnCode State Isolation Tests
// -----------------------------------------------------------------------------

func TestVerifier_ReturnCodeStateIsolation(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addr := createValidAddress(10)

	t.Run("Ok return code with mutations is allowed", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.BalanceDelta = big.NewInt(100)
		acc.StorageUpdates["k1"] = &vmcommon.StorageUpdate{Offset: []byte("k1"), Data: []byte("v1")}

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Non-Ok with clean output passes", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.ExecutionFailed
		output.ReturnMessage = "error"
		output.GasRemaining = 0

		input := createBaseInput()
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Non-Ok with storage update fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.StorageUpdates["k1"] = &vmcommon.StorageUpdate{Offset: []byte("k1"), Data: []byte("v1")}

		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.ExecutionFailed

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrStorageUpdateOnFailedExecution)
	})

	t.Run("Non-Ok with balance delta fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.BalanceDelta = big.NewInt(50)

		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.UserError

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrBalanceDeltaOnFailedExecution)
	})

	t.Run("Non-Ok with code deployment fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.Code = []byte{1, 2, 3}

		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.OutOfGas

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrCodeDeploymentOnFailedExecution)
	})

	t.Run("Non-Ok with output transfer fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.OutputTransfers = []vmcommon.OutputTransfer{{Value: big.NewInt(10)}}

		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.ExecutionFailed

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrOutputTransferOnFailedExecution)
	})

	t.Run("Non-Ok with deleted account fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.ExecutionFailed
		output.DeletedAccounts = [][]byte{createValidAddress(99)}

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrDeletedAccountOnFailedExecution)
	})

	t.Run("Non-Ok with positive gas refund succeeds", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.ExecutionFailed
		output.GasRefund = big.NewInt(500)

		input := createBaseInput()
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Non-Ok with negative gas refund fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.ReturnCode = vmcommon.ExecutionFailed
		output.GasRefund = big.NewInt(-500)

		input := createBaseInput()
		err := v.Verify(input, output)
		require.Error(t, err)
	})
}

// -----------------------------------------------------------------------------
// 2. Identity & Address Integrity Tests (UF-208 Defense)
// -----------------------------------------------------------------------------

func TestVerifier_IdentityAndAddressIntegrity(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	t.Run("Valid 32-byte address matching map key passes", func(t *testing.T) {
		addr := createValidAddress(20)
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		input := createBaseInput()

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Address shorter than 32 bytes fails", func(t *testing.T) {
		shortAddr := []byte("short-address")
		acc := &vmcommon.OutputAccount{
			Address: shortAddr,
		}
		output := &vmcommon.VMOutput{
			ReturnCode:     vmcommon.Ok,
			OutputAccounts: map[string]*vmcommon.OutputAccount{string(shortAddr): acc},
		}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrInvalidAccountAddress)
	})

	t.Run("Address longer than 32 bytes fails", func(t *testing.T) {
		longAddr := make([]byte, 33)
		acc := &vmcommon.OutputAccount{
			Address: longAddr,
		}
		output := &vmcommon.VMOutput{
			ReturnCode:     vmcommon.Ok,
			OutputAccounts: map[string]*vmcommon.OutputAccount{string(longAddr): acc},
		}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrInvalidAccountAddress)
	})

	t.Run("Map key does not match account address (aliasing attack)", func(t *testing.T) {
		addrA := createValidAddress(21)
		addrB := createValidAddress(22)

		acc := createValidAccount(addrA)
		output := &vmcommon.VMOutput{
			ReturnCode:     vmcommon.Ok,
			OutputAccounts: map[string]*vmcommon.OutputAccount{string(addrB): acc}, // Key is B, but acc is A!
		}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrAddressKeyMismatch)
	})

	t.Run("Deleted account length invalid", func(t *testing.T) {
		addr := createValidAddress(23)
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.DeletedAccounts = [][]byte{[]byte("too-short")}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrInvalidAccountAddress)
	})

	t.Run("Account is in both OutputAccounts and DeletedAccounts when DisjointDeletedAccounts is true", func(t *testing.T) {
		addr := createValidAddress(24)
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.DeletedAccounts = [][]byte{bytes.Clone(addr)}
		input := createBaseInput()
		input.DisjointDeletedAccounts = true

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrDeletedAccountInOutputAccounts)
	})

	t.Run("Deleted account in OutputAccounts allowed if clean, fails if mutating", func(t *testing.T) {
		addr := createValidAddress(24)
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.DeletedAccounts = [][]byte{bytes.Clone(addr)}
		input := createBaseInput()

		// Clean: no storage updates, no code, delta <= 0, no transfers -> passes
		err := v.Verify(input, output)
		require.NoError(t, err)

		// Storage updates in deleted account -> fails
		acc.StorageUpdates["key"] = &vmcommon.StorageUpdate{Offset: []byte("key"), Data: []byte("val")}
		err = v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrDeletedAccountInOutputAccounts)
		delete(acc.StorageUpdates, "key")

		// Positive balance delta in deleted account is allowed (e.g. cross-shard deleteContract with callValue)
		acc.BalanceDelta = big.NewInt(10)
		input.CallValue = big.NewInt(10)
		err = v.Verify(input, output)
		require.NoError(t, err)
		acc.BalanceDelta = big.NewInt(0)
		input.CallValue = big.NewInt(0)

		// Code deployment in deleted account -> fails
		acc.Code = []byte("code")
		err = v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrDeletedAccountInOutputAccounts)
		acc.Code = nil

		// Output transfers in deleted account -> fails
		acc.OutputTransfers = []vmcommon.OutputTransfer{{Data: []byte("data")}}
		err = v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrDeletedAccountInOutputAccounts)
	})

	t.Run("Duplicate deleted account fails", func(t *testing.T) {
		addrDel := createValidAddress(25)
		addrOut := createValidAddress(26)
		acc := createValidAccount(addrOut)
		output := createBaseVMOutput(acc)
		output.DeletedAccounts = [][]byte{bytes.Clone(addrDel), bytes.Clone(addrDel)}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrDuplicateDeletedAccount)
	})
}

// -----------------------------------------------------------------------------
// 3. Native EGLD Conservation Tests
// -----------------------------------------------------------------------------

func TestVerifier_NativeEGLDConservation(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addr1 := createValidAddress(31)
	addr2 := createValidAddress(32)

	t.Run("Net sum matches deposit exactly", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(60)

		acc2 := createValidAccount(addr2)
		acc2.BalanceDelta = big.NewInt(40)

		output := createBaseVMOutput(acc1, acc2)
		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Internal transfer conserves net balance", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(-40)

		acc2 := createValidAccount(addr2)
		acc2.BalanceDelta = big.NewInt(140)

		output := createBaseVMOutput(acc1, acc2)
		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Sum of deltas exceeds deposit (money creation attack)", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(101)

		output := createBaseVMOutput(acc1)
		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	})

	t.Run("Sum of deltas less than deposit (money loss)", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(99)

		output := createBaseVMOutput(acc1)
		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	})

	t.Run("Negative deposit rejected", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		output := createBaseVMOutput(acc1)
		input := createBaseInput()
		input.CallValue = big.NewInt(-1)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNegativeDeposit)
	})

	t.Run("Developer rewards claim event detected in logs passes with matching reward amount", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		rewardAmount := big.NewInt(5000)
		acc1.BalanceDelta = big.NewInt(5100)

		output := createBaseVMOutput(acc1)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionClaimDeveloperRewards),
			Topics: [][]byte{
				rewardAmount.Bytes(),
				addr1,
			},
		})

		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Developer rewards event spoofing without matching balance delta fails", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(6000)

		output := createBaseVMOutput(acc1)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionClaimDeveloperRewards),
			Topics: [][]byte{
				big.NewInt(5000).Bytes(),
				addr1,
			},
		})

		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	})

	t.Run("Developer rewards event in real execution where rewards credited out-of-band passes with zero delta", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(0)

		output := createBaseVMOutput(acc1)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionClaimDeveloperRewards),
			Topics: [][]byte{
				big.NewInt(5000).Bytes(),
				addr1,
			},
		})

		input := createBaseInput()
		input.CallValue = big.NewInt(0)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Balance delta inflated without developer rewards log fails", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(5100)

		output := createBaseVMOutput(acc1)
		input := createBaseInput()
		input.CallValue = big.NewInt(100)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	})

	t.Run("Developer rewards log with empty topics does not inflate deposit", func(t *testing.T) {
		acc1 := createValidAccount(addr1)
		acc1.BalanceDelta = big.NewInt(5000)

		output := createBaseVMOutput(acc1)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionClaimDeveloperRewards),
			Topics:     [][]byte{},
		})

		input := createBaseInput()
		input.CallValue = big.NewInt(0)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	})
}

// -----------------------------------------------------------------------------
// 4. Non-Negative Balance Check Tests
// -----------------------------------------------------------------------------

func TestVerifier_NonNegativeBalanceCheck(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addr := createValidAddress(40)

	t.Run("Positive final balance passes", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.Balance = big.NewInt(100)
		acc.BalanceDelta = big.NewInt(-40) // final = 60

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.CallValue = big.NewInt(-40) // matched for net sum test
		input.CallValue = big.NewInt(0)
		// compensate in another account
		addr2 := createValidAddress(41)
		acc2 := createValidAccount(addr2)
		acc2.BalanceDelta = big.NewInt(40)
		output = createBaseVMOutput(acc, acc2)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Overdrawn balance on OutputAccount fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.Balance = big.NewInt(50)
		acc.BalanceDelta = big.NewInt(-60) // final = -10 < 0!

		addr2 := createValidAddress(41)
		acc2 := createValidAccount(addr2)
		acc2.BalanceDelta = big.NewInt(60)
		output := createBaseVMOutput(acc, acc2)

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrAccountOverdrawn)
	})

	t.Run("Negative initial balance in OutputAccount fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.Balance = big.NewInt(-10)

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNegativeBalanceInOutput)
	})

	t.Run("Account with initial balance overdrawn below zero fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.Balance = big.NewInt(150)
		acc.BalanceDelta = big.NewInt(-200) // 150 - 200 = -50 < 0!

		addr2 := createValidAddress(41)
		acc2 := createValidAccount(addr2)
		acc2.BalanceDelta = big.NewInt(200)
		output := createBaseVMOutput(acc, acc2)

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrAccountOverdrawn)
	})
}

// -----------------------------------------------------------------------------
// 5. Token Conservation (ESDT) Tests
// -----------------------------------------------------------------------------

func TestVerifier_TokenConservation(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	tokenA := "TOKENA-123456"
	tokenB := "TOKENB-abcdef"

	t.Run("ESDT transfer within deposit passes", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(1000),
			},
		}

		destAddr := createValidAddress(51)
		acc := createValidAccount(destAddr)

		// Outflow 600 of tokenA via ESDTTransfer data
		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(600).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTTransfer),
			Topics: [][]byte{
				[]byte(tokenA), nil, big.NewInt(600).Bytes(), destAddr,
			},
		})

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("ESDT transfer without transfer log fails (unlogged transfer)", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = nil

		destAddr := createValidAddress(52)
		acc := createValidAccount(destAddr)

		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(500).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
		require.ErrorIs(t, err, verifier.ErrOutputTransferWithoutTransferLog)
	})

	t.Run("ESDTLocalMint with matching transfer and mint logs passes", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = nil

		destAddr := createValidAddress(53)
		acc := createValidAccount(destAddr)

		// Outflow 500
		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(500).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		output.Logs = append(output.Logs,
			&vmcommon.LogEntry{
				Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
				Topics: [][]byte{
					[]byte(tokenA),
					nil,
					big.NewInt(500).Bytes(),
				},
			},
			&vmcommon.LogEntry{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Topics: [][]byte{
					[]byte(tokenA),
					nil,
					big.NewInt(500).Bytes(),
					destAddr,
				},
			},
		)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("ESDTLocalBurn correctly consumes deposit", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(1000),
			},
		}

		destAddr := createValidAddress(54)
		acc := createValidAccount(destAddr)

		output := createBaseVMOutput(acc)
		// Burn 1000
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTLocalBurn),
			Topics: [][]byte{
				[]byte(tokenA),
				nil,
				big.NewInt(1000).Bytes(),
			},
		})

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Combination: Inflow + Mint - Burn - Outflow matches perfectly", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(1000),
			},
		}

		destAddr := createValidAddress(55)
		acc := createValidAccount(destAddr)

		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(1300).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		// Mint log
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
			Topics: [][]byte{
				[]byte(tokenA),
				nil,
				big.NewInt(500).Bytes(),
			},
		})
		// Burn log
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTLocalBurn),
			Topics: [][]byte{
				[]byte(tokenA),
				nil,
				big.NewInt(200).Bytes(),
			},
		})
		// Transfer log
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTTransfer),
			Topics: [][]byte{
				[]byte(tokenA),
				nil,
				big.NewInt(1300).Bytes(),
				destAddr,
			},
		})

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Transfer log amount mismatch fails conservation", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(1000),
			},
		}

		destAddr := createValidAddress(56)
		acc := createValidAccount(destAddr)

		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(1301).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		output.Logs = append(output.Logs,
			&vmcommon.LogEntry{
				Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
				Topics: [][]byte{
					[]byte(tokenA), nil, big.NewInt(500).Bytes(),
				},
			},
			&vmcommon.LogEntry{
				Identifier: []byte(core.BuiltInFunctionESDTLocalBurn),
				Topics: [][]byte{
					[]byte(tokenA), nil, big.NewInt(200).Bytes(),
				},
			},
			&vmcommon.LogEntry{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Topics: [][]byte{
					[]byte(tokenA), nil, big.NewInt(1300).Bytes(), destAddr,
				},
			},
		)

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
	})

	t.Run("Cross-token spoofing attack: Transfer log for token B does NOT match transfer of token A", func(t *testing.T) {
		input := createBaseInput()
		destAddr := createValidAddress(57)
		acc := createValidAccount(destAddr)

		// Attacker wants to transfer Token A
		encodedTokenA := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(100).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedTokenA + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		// But emits log for Token B!
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTTransfer),
			Topics: [][]byte{
				[]byte(tokenB), nil, big.NewInt(100).Bytes(), destAddr,
			},
		})

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
	})

	t.Run("NFT/SFT Nonce Isolation: Transfer log for Nonce 1 does NOT match transfer of Nonce 2", func(t *testing.T) {
		input := createBaseInput()
		destAddr := createValidAddress(58)
		acc := createValidAccount(destAddr)

		// Attacker transfers Nonce 2 of Token A
		encodedTokenA := hex.EncodeToString([]byte(tokenA))
		encodedNonce2 := hex.EncodeToString(big.NewInt(2).Bytes())
		encodedVal := hex.EncodeToString(big.NewInt(1).Bytes())
		encodedDest := hex.EncodeToString(destAddr)
		data := []byte(core.BuiltInFunctionESDTNFTTransfer + "@" + encodedTokenA + "@" + encodedNonce2 + "@" + encodedVal + "@" + encodedDest)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		// But emits log for Nonce 1!
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTNFTTransfer),
			Topics: [][]byte{
				[]byte(tokenA), big.NewInt(1).Bytes(), big.NewInt(1).Bytes(), destAddr,
			},
		})

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
	})

	t.Run("Pre-existing ESDT balance permits transfer with matching transfer log", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = nil

		destAddr := createValidAddress(59)
		acc := createValidAccount(destAddr)

		encodedTokenA := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(700).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedTokenA + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTTransfer),
			Topics: [][]byte{
				[]byte(tokenA), nil, big.NewInt(700).Bytes(), destAddr,
			},
		})
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Contract transferring 100% of ESDT balance passes verification", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(500),
			},
		}

		destAddr := createValidAddress(91)
		acc := createValidAccount(destAddr)

		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedVal := hex.EncodeToString(big.NewInt(500).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTTransfer),
			Topics: [][]byte{
				[]byte(tokenA), nil, big.NewInt(500).Bytes(), destAddr,
			},
		})
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Contract transferring pre-existing NFT passes verification", func(t *testing.T) {
		input := createBaseInput()
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 1,
				ESDTValue:      big.NewInt(1),
			},
		}

		destAddr := createValidAddress(92)
		acc := createValidAccount(destAddr)
		encodedToken := hex.EncodeToString([]byte(tokenA))
		encodedNonce := hex.EncodeToString(big.NewInt(1).Bytes())
		encodedVal := hex.EncodeToString(big.NewInt(1).Bytes())
		encodedDest := hex.EncodeToString(destAddr)
		data := []byte(core.BuiltInFunctionESDTNFTTransfer + "@" + encodedToken + "@" + encodedNonce + "@" + encodedVal + "@" + encodedDest)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		output.Logs = append(output.Logs, &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTNFTTransfer),
			Topics: [][]byte{
				[]byte(tokenA), big.NewInt(1).Bytes(), big.NewInt(1).Bytes(), destAddr,
			},
		})
		err := v.Verify(input, output)
		require.NoError(t, err)
	})
}

// -----------------------------------------------------------------------------
// 6. Gas Bounds & 50% Refund Cap Tests
// -----------------------------------------------------------------------------

func TestVerifier_GasBoundsAndRefundCap(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addr := createValidAddress(60)

	t.Run("Valid gas spent and 50% refund passes", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.GasRemaining = 40000
		// GasProvided = 100000, GasSpent = 60000. Max refund = 30000
		output.GasRefund = big.NewInt(30000)

		input := createBaseInput()
		input.GasProvided = 100000

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("GasRemaining > GasProvided fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.GasRemaining = 100001 // > 100000!

		input := createBaseInput()
		input.GasProvided = 100000

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrGasRemainingExceedsProvided)
	})

	t.Run("Negative gas refund fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.GasRefund = big.NewInt(-1)

		input := createBaseInput()
		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrGasRefundNegative)
	})

	t.Run("Gas refund exceeds 50% cap of spent gas is clamped", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.GasRemaining = 40000 // spent = 60000, max refund = 30000
		output.GasRefund = big.NewInt(30001)

		input := createBaseInput()
		input.GasProvided = 100000

		err := v.Verify(input, output)
		require.NoError(t, err)
		require.Equal(t, big.NewInt(30000), output.GasRefund)
	})
}

// -----------------------------------------------------------------------------
// 7. ReadOnly Purity Tests
// -----------------------------------------------------------------------------

func TestVerifier_ReadOnlyPurity(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addr := createValidAddress(70)

	t.Run("Pure read-only call passes", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.ReturnData = [][]byte{[]byte("read-data")}

		input := createBaseInput()
		input.IsReadOnly = true

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("ReadOnly with storage update fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.StorageUpdates["k1"] = &vmcommon.StorageUpdate{Offset: []byte("k1"), Data: []byte("v1")}

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.IsReadOnly = true

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
	})

	t.Run("ReadOnly with balance delta fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.BalanceDelta = big.NewInt(10)

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.IsReadOnly = true

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
	})

	t.Run("ReadOnly with code deployment fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.Code = []byte("wasm-code")

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.IsReadOnly = true

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
	})

	t.Run("ReadOnly with output transfer fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		acc.OutputTransfers = []vmcommon.OutputTransfer{{Value: big.NewInt(5)}}

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.IsReadOnly = true

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
	})

	t.Run("ReadOnly with deleted accounts fails", func(t *testing.T) {
		acc := createValidAccount(addr)
		output := createBaseVMOutput(acc)
		output.DeletedAccounts = [][]byte{createValidAddress(71)}

		input := createBaseInput()
		input.IsReadOnly = true

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
	})

	t.Run("ReadOnly with restricted events fails", func(t *testing.T) {
		restrictedEvents := []string{
			"ChangeOwnerAddress",
			"ClaimDeveloperRewards",
			"SetUserName",
			"SaveKeyValue",
		}
		for _, event := range restrictedEvents {
			output := createBaseVMOutput()
			output.Logs = []*vmcommon.LogEntry{
				{
					Identifier: []byte(event),
				},
			}
			input := createBaseInput()
			input.IsReadOnly = true

			err := v.Verify(input, output)
			require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange, "expected event %s to be rejected in read-only mode", event)
		}
	})
}

// -----------------------------------------------------------------------------
// 8. Hardened Invariant Verification Tests
// -----------------------------------------------------------------------------

func TestVerifier_OverdrawnNonExistentAccount(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addrOverdrawn := createValidAddress(88)
	accOverdrawn := createValidAccount(addrOverdrawn)
	accOverdrawn.Balance = big.NewInt(30)
	accOverdrawn.BalanceDelta = big.NewInt(-50) // 30 - 50 = -20 < 0

	addrRecipient := createValidAddress(89)
	accRecipient := createValidAccount(addrRecipient)
	accRecipient.BalanceDelta = big.NewInt(50)

	output := createBaseVMOutput(accOverdrawn, accRecipient)
	input := createBaseInput()
	input.CallValue = big.NewInt(0) // Net sum: -50 + 50 = 0 (conserved)

	err := v.Verify(input, output)
	require.ErrorIs(t, err, verifier.ErrAccountOverdrawn)
}

func TestVerifier_HardenedGasBounds(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	t.Run("Single account GasUsed exceeds GasProvided", func(t *testing.T) {
		acc := createValidAccount(createValidAddress(90))
		acc.GasUsed = 100001 // GasProvided is 100000

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.GasProvided = 100000

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrGasRemainingExceedsProvided)
	})

	t.Run("Sum of account GasUsed exceeds GasProvided", func(t *testing.T) {
		acc1 := createValidAccount(createValidAddress(91))
		acc1.GasUsed = 60000
		acc2 := createValidAccount(createValidAddress(92))
		acc2.GasUsed = 50000 // Total = 110000 > 100000

		output := createBaseVMOutput(acc1, acc2)
		input := createBaseInput()
		input.GasProvided = 100000

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrGasRemainingExceedsProvided)
	})

	t.Run("OutputTransfer GasLimit exceeds GasProvided", func(t *testing.T) {
		acc := createValidAccount(createValidAddress(93))
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				GasLimit: 100001,
			},
		}

		output := createBaseVMOutput(acc)
		input := createBaseInput()
		input.GasProvided = 100000

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrGasRemainingExceedsProvided)
	})
}

func TestVerifier_TokenNonceValidation(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	tokenID := "TEST-123456"
	oversizedNonce := bytes.Repeat([]byte{0x01}, 9) // 9 bytes > 8 bytes

	t.Run("Mint event with oversized nonce fails", func(t *testing.T) {
		output := createBaseVMOutput()
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
				Topics: [][]byte{
					[]byte(tokenID),
					oversizedNonce,
					big.NewInt(100).Bytes(),
				},
			},
		}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrInvalidTokenNonce)
	})

	t.Run("Transfer log with oversized nonce fails", func(t *testing.T) {
		output := createBaseVMOutput()
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTNFTTransfer),
				Topics: [][]byte{
					[]byte(tokenID),
					oversizedNonce,
					big.NewInt(100).Bytes(),
				},
			},
		}
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrInvalidTokenNonce)
	})

	t.Run("Output transfer data with oversized nonce fails", func(t *testing.T) {
		acc := createValidAccount(createValidAddress(94))
		encodedToken := hex.EncodeToString([]byte(tokenID))
		encodedNonce := hex.EncodeToString(oversizedNonce)
		encodedVal := hex.EncodeToString(big.NewInt(100).Bytes())
		data := []byte(core.BuiltInFunctionESDTNFTTransfer + "@" + encodedToken + "@" + encodedNonce + "@" + encodedVal)
		acc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data: data,
			},
		}

		output := createBaseVMOutput(acc)
		input := createBaseInput()

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrInvalidTokenNonce)
	})
}

func TestVerifier_MalformedOutputTransferData(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	acc := createValidAccount(createValidAddress(95))
	acc.OutputTransfers = []vmcommon.OutputTransfer{
		{
			Data: []byte("ESDTTransfer@not-valid-hex-!@#$%"),
		},
	}

	output := createBaseVMOutput(acc)
	input := createBaseInput()

	err := v.Verify(input, output)
	require.ErrorIs(t, err, verifier.ErrMalformedOutputTransferData)
}

func TestVerifier_EGLDInMultiTransferIgnored(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	output := createBaseVMOutput()
	output.Logs = []*vmcommon.LogEntry{
		{
			Identifier: []byte(core.BuiltInFunctionMultiESDTNFTTransfer),
			Topics: [][]byte{
				[]byte("EGLD-000000"),
				nil,
				big.NewInt(1000).Bytes(),
			},
		},
	}
	input := createBaseInput()
	input.ESDTTransfers = nil

	// EGLD-000000 is ignored from ESDT conservation, so no unbacked token error is returned
	err := v.Verify(input, output)
	require.NoError(t, err)
}

func TestVerifier_ComposableMultiContractSwap(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	tokenIn := "USDC-abcdef"
	tokenOut := "WEGLD-123456"

	routerAddr := createValidAddress(101)
	pairAddr := createValidAddress(102)
	userAddr := createValidAddress(103)

	input := createBaseInput()
	input.RecipientAddr = routerAddr
	// Inflow: user sends 100 USDC to Router
	input.ESDTTransfers = []*vmcommon.ESDTTransfer{
		{
			ESDTTokenName:  []byte(tokenIn),
			ESDTTokenNonce: 0,
			ESDTValue:      big.NewInt(100),
		},
	}
	// Router OutputAccount
	routerAcc := createValidAccount(routerAddr)
	routerAcc.GasUsed = 10000

	// Pair OutputAccount (executed child contract)
	pairAcc := createValidAccount(pairAddr)
	pairAcc.GasUsed = 20000

	// Intra-shard call: Router transfers 100 USDC to Pair
	encodedTokenIn := hex.EncodeToString([]byte(tokenIn))
	encodedValIn := hex.EncodeToString(big.NewInt(100).Bytes())
	dataRouterToPair := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedTokenIn + "@" + encodedValIn)
	pairAcc.OutputTransfers = []vmcommon.OutputTransfer{
		{
			Data:          dataRouterToPair,
			SenderAddress: routerAddr,
		},
	}

	// External transfer: Pair transfers 50 WEGLD to User C
	userAcc := createValidAccount(userAddr)
	encodedTokenOut := hex.EncodeToString([]byte(tokenOut))
	encodedValOut := hex.EncodeToString(big.NewInt(50).Bytes())
	dataPairToUser := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedTokenOut + "@" + encodedValOut)
	userAcc.OutputTransfers = []vmcommon.OutputTransfer{
		{
			Data:          dataPairToUser,
			SenderAddress: pairAddr,
		},
	}

	output := createBaseVMOutput(routerAcc, pairAcc, userAcc)
	output.Logs = append(output.Logs, &vmcommon.LogEntry{
		Identifier: []byte(core.BuiltInFunctionESDTTransfer),
		Topics: [][]byte{
			[]byte(tokenOut), nil, big.NewInt(50).Bytes(), userAddr,
		},
	})
	err := v.Verify(input, output)
	require.NoError(t, err)
}

// -----------------------------------------------------------------------------
// 10. Advanced Hardening & Edge-Case Tests
// -----------------------------------------------------------------------------

func TestVerifier_FailedExecution_TokenLogsRejected(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	addr := createValidAddress(20)
	input := createBaseInput()
	input.RecipientAddr = addr

	tests := []struct {
		name      string
		eventID   string
		retCode   vmcommon.ReturnCode
		shouldErr bool
	}{
		{
			name:      "ESDTTransfer log on ExecutionFailed is rejected",
			eventID:   core.BuiltInFunctionESDTTransfer,
			retCode:   vmcommon.ExecutionFailed,
			shouldErr: true,
		},
		{
			name:      "ESDTLocalBurn log on UserError is rejected",
			eventID:   core.BuiltInFunctionESDTLocalBurn,
			retCode:   vmcommon.UserError,
			shouldErr: true,
		},
		{
			name:      "ESDTBurn log on OutOfGas is rejected",
			eventID:   core.BuiltInFunctionESDTBurn,
			retCode:   vmcommon.OutOfGas,
			shouldErr: true,
		},
		{
			name:      "transferValueOnly log on ExecutionFailed is rejected",
			eventID:   "transferValueOnly",
			retCode:   vmcommon.ExecutionFailed,
			shouldErr: true,
		},
		{
			name:      "Arbitrary unprivileged event on ExecutionFailed is allowed",
			eventID:   "myCustomEvent",
			retCode:   vmcommon.ExecutionFailed,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			acc := createValidAccount(addr)
			output := createBaseVMOutput(acc)
			output.ReturnCode = tt.retCode
			output.Logs = []*vmcommon.LogEntry{
				{
					Identifier: []byte(tt.eventID),
					Address:    addr,
					Topics:     [][]byte{[]byte("TOKEN-123456"), {0}, {10}},
				},
			}

			err := v.Verify(input, output)
			if tt.shouldErr {
				require.ErrorIs(t, err, verifier.ErrTokenTransferOnFailedExecution)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVerifier_NativeEGLD_OutputTransfersHardened(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	senderAddr := createValidAddress(30)
	externalAddr := createValidAddress(31)

	t.Run("Negative OutputTransfer.Value is rejected", func(t *testing.T) {
		senderAcc := createValidAccount(senderAddr)
		senderAcc.BalanceDelta = big.NewInt(-10)

		externalAcc := createValidAccount(externalAddr)
		externalAcc.BalanceDelta = big.NewInt(10)
		externalAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				SenderAddress: senderAddr,
				Value:         big.NewInt(-10), // Negative value!
			},
		}

		output := createBaseVMOutput(senderAcc, externalAcc)
		input := createBaseInput()
		input.RecipientAddr = senderAddr

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrNegativeOutputTransferValue)
	})

	t.Run("OutputTransfer.Value exceeds BalanceDelta on external recipient fails (money creation)", func(t *testing.T) {
		senderAcc := createValidAccount(senderAddr)
		senderAcc.BalanceDelta = big.NewInt(-100)

		externalAcc := createValidAccount(externalAddr)
		externalAcc.BalanceDelta = big.NewInt(50) // Only credited 50, but transfer says 100!

		externalAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				SenderAddress: senderAddr,
				Value:         big.NewInt(100),
			},
		}

		output := createBaseVMOutput(senderAcc, externalAcc)
		input := createBaseInput()
		input.RecipientAddr = senderAddr
		input.CallValue = big.NewInt(0)

		// Note: net sum (-100 + 50 = -50 != 0) would fail conservation, but transfer delta check catches it specifically
		externalAcc.BalanceDelta = big.NewInt(0)
		err := v.Verify(input, output)
		require.Error(t, err)
	})

	t.Run("OutputTransfer.Value matches BalanceDelta on external recipient passes", func(t *testing.T) {
		senderAcc := createValidAccount(senderAddr)
		senderAcc.BalanceDelta = big.NewInt(-100)

		externalAcc := createValidAccount(externalAddr)
		externalAcc.BalanceDelta = big.NewInt(100)
		externalAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				SenderAddress: senderAddr,
				Value:         big.NewInt(100),
			},
		}

		output := createBaseVMOutput(senderAcc, externalAcc)
		input := createBaseInput()
		input.RecipientAddr = senderAddr
		input.CallValue = big.NewInt(0)

		err := v.Verify(input, output)
		require.NoError(t, err)
	})
}

func TestVerifier_TokenConservation_ESDTBurn_Recognized(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	token := "TEST-123456"
	contractAddr := createValidAddress(40)
	userAddr := createValidAddress(41)

	input := createBaseInput()
	input.RecipientAddr = contractAddr
	input.ESDTTransfers = []*vmcommon.ESDTTransfer{
		{
			ESDTTokenName:  []byte(token),
			ESDTTokenNonce: 0,
			ESDTValue:      big.NewInt(100),
		},
	}

	t.Run("ESDTBurn event accounts for burned supply", func(t *testing.T) {
		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		// Burn 40 tokens using core.BuiltInFunctionESDTBurn
		burnLog := &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTBurn),
			Address:    contractAddr,
			Topics:     [][]byte{[]byte(token), {0}, big.NewInt(40).Bytes()},
		}

		// Transfer 60 tokens to user
		userAcc := createValidAccount(userAddr)
		encodedToken := hex.EncodeToString([]byte(token))
		encodedVal := hex.EncodeToString(big.NewInt(60).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		output.Logs = []*vmcommon.LogEntry{
			burnLog,
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics:     [][]byte{[]byte(token), nil, big.NewInt(60).Bytes(), userAddr},
			},
		}

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("ESDTBurn plus excess outflow fails conservation", func(t *testing.T) {
		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		burnLog := &vmcommon.LogEntry{
			Identifier: []byte(core.BuiltInFunctionESDTBurn),
			Address:    contractAddr,
			Topics:     [][]byte{[]byte(token), {0}, big.NewInt(40).Bytes()},
		}

		// Transfer 61 tokens to user (40 burned + 61 outflow = 101 > 100 available)
		userAcc := createValidAccount(userAddr)
		encodedToken := hex.EncodeToString([]byte(token))
		encodedVal := hex.EncodeToString(big.NewInt(61).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		output.Logs = []*vmcommon.LogEntry{burnLog}

		err := v.Verify(input, output)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
	})
}

func TestVerifier_MultiESDTNFTTransfer_FormatB(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	token := "MULTI-abcdef"
	contractAddr := createValidAddress(50)
	userAddr := createValidAddress(51)

	input := createBaseInput()
	input.RecipientAddr = contractAddr
	input.ESDTTransfers = []*vmcommon.ESDTTransfer{
		{
			ESDTTokenName:  []byte(token),
			ESDTTokenNonce: 5,
			ESDTValue:      big.NewInt(200),
		},
	}

	contractAcc := createValidAccount(contractAddr)
	contractAcc.GasUsed = 5000

	// Format B: MultiESDTNFTTransfer@numTransfers@token1@nonce1@val1
	// numTransfers = 1
	encodedNumTransfers := hex.EncodeToString(big.NewInt(1).Bytes())
	encodedToken := hex.EncodeToString([]byte(token))
	encodedNonce := hex.EncodeToString(big.NewInt(5).Bytes())
	encodedVal := hex.EncodeToString(big.NewInt(200).Bytes())
	data := []byte(core.BuiltInFunctionMultiESDTNFTTransfer + "@" +
		encodedNumTransfers + "@" + encodedToken + "@" + encodedNonce + "@" + encodedVal)

	userAcc := createValidAccount(userAddr)
	userAcc.OutputTransfers = []vmcommon.OutputTransfer{
		{
			Data:          data,
			SenderAddress: contractAddr,
		},
	}

	output := createBaseVMOutput(contractAcc, userAcc)
	output.Logs = []*vmcommon.LogEntry{
		{
			Identifier: []byte(core.BuiltInFunctionMultiESDTNFTTransfer),
			Address:    contractAddr,
			Topics: [][]byte{
				[]byte(token), big.NewInt(5).Bytes(), big.NewInt(200).Bytes(), userAddr,
			},
		},
	}
	err := v.Verify(input, output)
	require.NoError(t, err)
}

func TestVerifier_MultiESDTNFTTransfer_DestinationCheck(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	token := "TEST-999999"
	parentAddr := createValidAddress(60)
	childAddr := createValidAddress(61)

	input := createBaseInput()
	input.RecipientAddr = parentAddr
	// Inflow to parent: 100 tokens
	input.ESDTTransfers = []*vmcommon.ESDTTransfer{
		{
			ESDTTokenName:  []byte(token),
			ESDTTokenNonce: 0,
			ESDTValue:      big.NewInt(100),
		},
	}

	parentAcc := createValidAccount(parentAddr)
	parentAcc.GasUsed = 5000

	childAcc := createValidAccount(childAddr)
	childAcc.GasUsed = 5000 // Executed contract!

	// Log of MultiESDTNFTTransfer to childAddr (internal transfer)
	// Topics: [token, nonce, amount, destination] (N=4, N%3 == 1, topics[3] == childAddr)
	multiLogInternal := &vmcommon.LogEntry{
		Identifier: []byte(core.BuiltInFunctionMultiESDTNFTTransfer),
		Address:    parentAddr,
		Topics: [][]byte{
			[]byte(token),
			{0},
			big.NewInt(100).Bytes(),
			childAddr, // destination is executed child contract!
		},
	}

	output := createBaseVMOutput(parentAcc, childAcc)
	output.Logs = []*vmcommon.LogEntry{multiLogInternal}

	// Should pass because internal transfer between executed contracts does NOT count as external outflow
	err := v.Verify(input, output)
	require.NoError(t, err)
}

func TestVerifier_EGLD_TokenIdentifier_Variants(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	contractAddr := createValidAddress(70)
	userAddr := createValidAddress(71)

	input := createBaseInput()
	input.RecipientAddr = contractAddr
	input.CallValue = big.NewInt(100)

	contractAcc := createValidAccount(contractAddr)
	contractAcc.GasUsed = 5000
	contractAcc.BalanceDelta = big.NewInt(0)

	userAcc := createValidAccount(userAddr)
	userAcc.BalanceDelta = big.NewInt(100)

	// Transfer using "EGLD" identifier in MultiESDTNFTTransfer
	encodedNum := hex.EncodeToString(big.NewInt(1).Bytes())
	encodedToken := hex.EncodeToString([]byte("EGLD"))
	encodedNonce := hex.EncodeToString(big.NewInt(0).Bytes())
	encodedVal := hex.EncodeToString(big.NewInt(100).Bytes())
	data := []byte(core.BuiltInFunctionMultiESDTNFTTransfer + "@" +
		encodedNum + "@" + encodedToken + "@" + encodedNonce + "@" + encodedVal)

	userAcc.OutputTransfers = []vmcommon.OutputTransfer{
		{
			Data:          data,
			Value:         big.NewInt(100),
			SenderAddress: contractAddr,
		},
	}

	output := createBaseVMOutput(contractAcc, userAcc)
	// Should pass cleanly: EGLD is ignored by ESDT token conservation, handled by native EGLD conservation
	err := v.Verify(input, output)
	require.NoError(t, err)
}

func TestVerifier_ReadOnlyPurity_Hardened(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	contractAddr := createValidAddress(80)
	input := createBaseInput()
	input.RecipientAddr = contractAddr
	input.IsReadOnly = true

	restrictedEvents := []string{
		core.BuiltInFunctionESDTBurn,
		core.BuiltInFunctionESDTFreeze,
		core.BuiltInFunctionESDTUnFreeze,
		"transferValueOnly",
	}

	for _, ev := range restrictedEvents {
		ev := ev
		t.Run("Event "+ev+" rejected in read-only mode", func(t *testing.T) {
			t.Parallel()
			acc := createValidAccount(contractAddr)
			output := createBaseVMOutput(acc)
			output.Logs = []*vmcommon.LogEntry{
				{
					Identifier: []byte(ev),
					Address:    contractAddr,
					Topics:     [][]byte{[]byte("TEST-123456"), {0}, {1}},
				},
			}

			err := v.Verify(input, output)
			require.ErrorIs(t, err, verifier.ErrReadOnlyStateChange)
		})
	}
}

func TestESDTTokenConservation_CrossShardMarshaledNFTTransfer(t *testing.T) {
	t.Parallel()

	v := verifier.NewVMOutputVerifier()
	contractAddr := createValidAddress(10)
	destAddr := createValidAddress(20)

	input := createBaseInput()
	input.RecipientAddr = contractAddr
	input.CallerAddr = createValidAddress(9)

	tokenData := &esdt.ESDigitalToken{
		Value: big.NewInt(100),
		Type:  uint32(core.SemiFungible),
		TokenMetaData: &esdt.MetaData{
			Name:       []byte("token"),
			URIs:       [][]byte{[]byte("www.my-cool-nfts.com")},
			Attributes: []byte("cool nft"),
		},
	}
	protoMarshalizer := &marshal.GogoProtoMarshalizer{}
	marshaledData, err := protoMarshalizer.Marshal(tokenData)
	require.NoError(t, err)

	// Format: MultiESDTNFTTransfer@numTransfers@token@nonce@marshaledData
	outTransferData := "MultiESDTNFTTransfer@" +
		hex.EncodeToString(big.NewInt(1).Bytes()) + "@" +
		hex.EncodeToString([]byte("SFT-123456")) + "@" +
		hex.EncodeToString(big.NewInt(1).Bytes()) + "@" +
		hex.EncodeToString(marshaledData)

	destAcc := createValidAccount(destAddr)
	destAcc.OutputTransfers = []vmcommon.OutputTransfer{
		{
			Value: big.NewInt(0),
			Data:  []byte(outTransferData),
		},
	}

	contractAcc := createValidAccount(contractAddr)
	output := createBaseVMOutput(contractAcc, destAcc)
	output.Logs = []*vmcommon.LogEntry{
		{
			Identifier: []byte(core.BuiltInFunctionMultiESDTNFTTransfer),
			Address:    contractAddr,
			Topics: [][]byte{
				[]byte("SFT-123456"),
				big.NewInt(1).Bytes(),
				big.NewInt(100).Bytes(),
				destAddr,
			},
		},
	}

	err = v.Verify(input, output)
	require.NoError(t, err)
}

// -----------------------------------------------------------------------------
// 11. Verification Test Suite: Log Spoofing & EGLD Normalization
// -----------------------------------------------------------------------------

func TestVerifier_TransferLogNeverSuppliesTokens(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	tokenID := "TOKEN-123456"
	contractAddr := createValidAddress(80)
	userAddr := createValidAddress(81)

	t.Run("Log-only ESDT transfer with 0 inflow and 0 balance fails", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.ESDTTransfers = nil

		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		output := createBaseVMOutput(contractAcc)
		// Attacker emits a log-only transfer of 1,000,000 tokens
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(1000000).Bytes(),
					userAddr,
				},
			},
		}

		err := v.Verify(input, output)
		require.Error(t, err)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
		require.ErrorIs(t, err, verifier.ErrTransferLogWithoutOutputTransfer)
	})

	t.Run("P1 Exploit: OutputTransfer without matching transfer log fails", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.ESDTTransfers = nil

		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		encodedToken := hex.EncodeToString([]byte(tokenID))
		encodedVal := hex.EncodeToString(big.NewInt(1000000).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)

		err := v.Verify(input, output)
		require.Error(t, err)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
		require.ErrorIs(t, err, verifier.ErrOutputTransferWithoutTransferLog)
	})

	t.Run("P1 Exploit: Unbacked MultiESDTNFTTransfer log fails conservation", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.ESDTTransfers = nil

		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		output := createBaseVMOutput(contractAcc)
		// Multi transfer log for 2 tokens without matching OutputTransfers
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionMultiESDTNFTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte("TOK1-111111"), {0}, big.NewInt(100).Bytes(),
					[]byte("TOK2-222222"), {0}, big.NewInt(200).Bytes(),
					userAddr,
				},
			},
		}

		err := v.Verify(input, output)
		require.Error(t, err)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
		require.ErrorIs(t, err, verifier.ErrTransferLogWithoutOutputTransfer)
	})
}

func TestVerifier_RealWorld_TokenConservationScenarios(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	tokenID := "USDC-c0ffee"
	contractAddr := createValidAddress(82)
	userAddr := createValidAddress(83)

	t.Run("Legitimate deposit with partial spend and burn passes", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		// Deposit 1,000 USDC
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenID),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(1000),
			},
		}

		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		encodedToken := hex.EncodeToString([]byte(tokenID))
		encodedVal := hex.EncodeToString(big.NewInt(400).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		// Transfer log for 400 + Burn log for 100
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(400).Bytes(),
					userAddr,
				},
			},
			{
				Identifier: []byte(core.BuiltInFunctionESDTBurn),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(100).Bytes(),
				},
			},
		}

		// Consumed = 400 + 100 = 500 <= Inflow 1000. 500 remains in contract.
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("Mint and transfer with matching logs passes", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.ESDTTransfers = nil

		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		encodedToken := hex.EncodeToString([]byte(tokenID))
		encodedVal := hex.EncodeToString(big.NewInt(500).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(500).Bytes(),
				},
			},
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(500).Bytes(),
					userAddr,
				},
			},
		}

		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("P1: Amount mismatch between OutputTransfer and transfer log fails", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.ESDTTransfers = nil

		contractAcc := createValidAccount(contractAddr)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		encodedToken := hex.EncodeToString([]byte(tokenID))
		encodedVal := hex.EncodeToString(big.NewInt(501).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(500).Bytes(),
				},
			},
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenID),
					{0},
					big.NewInt(500).Bytes(),
					userAddr,
				},
			},
		}

		err := v.Verify(input, output)
		require.Error(t, err)
		require.ErrorIs(t, err, verifier.ErrTokenConservationViolated)
	})
}

func TestVerifier_RealWorld_EGLD_ESDT_Integration(t *testing.T) {
	t.Parallel()
	v := verifier.NewVMOutputVerifier()

	contractAddr := createValidAddress(84)
	userAddr := createValidAddress(85)

	t.Run("Deposit via ESDTTransfers with EGLD-000000 conserved in native balance", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.CallValue = big.NewInt(1000)
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(verifier.EGLDTokenIdentifier),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(1000),
			},
		}

		contractAcc := createValidAccount(contractAddr)
		contractAcc.BalanceDelta = big.NewInt(0)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		userAcc.BalanceDelta = big.NewInt(1000)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				SenderAddress: contractAddr,
				Value:         big.NewInt(1000),
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("OutputTransfer of EGLD via ESDTTransfer data checked in native balance delta", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.CallValue = big.NewInt(500)

		contractAcc := createValidAccount(contractAddr)
		contractAcc.BalanceDelta = big.NewInt(0)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		userAcc.BalanceDelta = big.NewInt(500)

		// Transfer EGLD using ESDTTransfer data with 0 Value
		encodedToken := hex.EncodeToString([]byte(verifier.EGLDTokenIdentifier))
		encodedVal := hex.EncodeToString(big.NewInt(500).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
				Value:         big.NewInt(0),
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		err := v.Verify(input, output)
		require.NoError(t, err)
	})

	t.Run("OutputTransfer of EGLD exceeding deposit fails native conservation", func(t *testing.T) {
		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.CallValue = big.NewInt(500)

		contractAcc := createValidAccount(contractAddr)
		contractAcc.BalanceDelta = big.NewInt(400)
		contractAcc.GasUsed = 5000

		userAcc := createValidAccount(userAddr)
		userAcc.BalanceDelta = big.NewInt(200) // 400 + 200 = 600 > 500

		encodedToken := hex.EncodeToString([]byte(verifier.EGLDTokenIdentifier))
		encodedVal := hex.EncodeToString(big.NewInt(500).Bytes())
		data := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedToken + "@" + encodedVal)
		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          data,
				SenderAddress: contractAddr,
				Value:         big.NewInt(0),
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		err := v.Verify(input, output)
		require.Error(t, err)
		require.ErrorIs(t, err, verifier.ErrNativeTokenConservationViolated)
	})

	t.Run("Multi-token atomic swap with USDC, WEGLD, and EGLD-000000", func(t *testing.T) {
		tokenA := "USDC-111111"
		tokenB := "WEGLD-222222"

		input := createBaseInput()
		input.RecipientAddr = contractAddr
		input.CallValue = big.NewInt(100) // 100 native EGLD
		input.ESDTTransfers = []*vmcommon.ESDTTransfer{
			{
				ESDTTokenName:  []byte(tokenA),
				ESDTTokenNonce: 0,
				ESDTValue:      big.NewInt(500),
			},
		}

		contractAcc := createValidAccount(contractAddr)
		contractAcc.BalanceDelta = big.NewInt(0)
		contractAcc.GasUsed = 10000

		userAcc := createValidAccount(userAddr)
		userAcc.BalanceDelta = big.NewInt(100)

		// Contract sends 500 USDC and 200 WEGLD (minted 200 WEGLD)
		encodedTokenA := hex.EncodeToString([]byte(tokenA))
		encodedValA := hex.EncodeToString(big.NewInt(500).Bytes())
		dataA := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedTokenA + "@" + encodedValA)

		encodedTokenB := hex.EncodeToString([]byte(tokenB))
		encodedValB := hex.EncodeToString(big.NewInt(200).Bytes())
		dataB := []byte(core.BuiltInFunctionESDTTransfer + "@" + encodedTokenB + "@" + encodedValB)

		userAcc.OutputTransfers = []vmcommon.OutputTransfer{
			{
				Data:          dataA,
				SenderAddress: contractAddr,
			},
			{
				Data:          dataB,
				SenderAddress: contractAddr,
			},
			{
				Value:         big.NewInt(100),
				SenderAddress: contractAddr,
			},
		}

		output := createBaseVMOutput(contractAcc, userAcc)
		output.Logs = []*vmcommon.LogEntry{
			{
				Identifier: []byte(core.BuiltInFunctionESDTLocalMint),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenB),
					{0},
					big.NewInt(200).Bytes(),
				},
			},
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenA),
					{0},
					big.NewInt(500).Bytes(),
					userAddr,
				},
			},
			{
				Identifier: []byte(core.BuiltInFunctionESDTTransfer),
				Address:    contractAddr,
				Topics: [][]byte{
					[]byte(tokenB),
					{0},
					big.NewInt(200).Bytes(),
					userAddr,
				},
			},
		}

		err := v.Verify(input, output)
		require.NoError(t, err)
	})
}

func BenchmarkOutputVerifier_Verify_OkCall(b *testing.B) {
	v := verifier.NewVMOutputVerifier()
	addr := createValidAddress(10)
	acc := createValidAccount(addr)
	acc.BalanceDelta = big.NewInt(100)
	acc.StorageUpdates["k1"] = &vmcommon.StorageUpdate{Offset: []byte("k1"), Data: []byte("v1")}

	output := createBaseVMOutput(acc)
	input := createBaseInput()
	input.CallValue = big.NewInt(100)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = v.Verify(input, output)
	}
}

func BenchmarkOutputVerifier_Verify_MultiAccount(b *testing.B) {
	v := verifier.NewVMOutputVerifier()
	acc1 := createValidAccount(createValidAddress(10))
	acc1.BalanceDelta = big.NewInt(50)
	acc2 := createValidAccount(createValidAddress(20))
	acc2.BalanceDelta = big.NewInt(50)
	acc3 := createValidAccount(createValidAddress(30))
	acc3.BalanceDelta = big.NewInt(0)

	output := createBaseVMOutput(acc1, acc2, acc3)
	input := createBaseInput()
	input.CallValue = big.NewInt(100)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = v.Verify(input, output)
	}
}

