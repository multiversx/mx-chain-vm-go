package verifier

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-core-go/marshal"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/parsers"
)

var bigZero = big.NewInt(0)

func getSortedAccountKeys(m map[string]*vmcommon.OutputAccount) []string {
	n := len(m)
	if n == 0 {
		return nil
	}
	keys := make([]string, 0, n)
	for k := range m {
		keys = append(keys, k)
	}
	if n > 1 {
		sort.Strings(keys)
	}
	return keys
}

const AddressLen = 32
const EGLDTokenIdentifier = "EGLD-000000"

func isEGLDToken(tokenID string) bool {
	return tokenID == EGLDTokenIdentifier || tokenID == "EGLD"
}

func normalizeTokenID(tokenID string) string {
	if isEGLDToken(tokenID) {
		return EGLDTokenIdentifier
	}
	return tokenID
}

// isExecutedContract returns true if addr corresponds to a contract that executed code or had its state modified
func isExecutedContract(input *VerificationInput, output *vmcommon.VMOutput, addr []byte) bool {
	if len(addr) == 0 {
		return false
	}
	if input != nil && bytes.Equal(addr, input.CallerAddr) {
		return false
	}
	if input != nil && bytes.Equal(addr, input.RecipientAddr) {
		return true
	}
	if output == nil || output.OutputAccounts == nil {
		return false
	}
	acc, ok := output.OutputAccounts[string(addr)]
	if !ok || acc == nil {
		return false
	}
	return acc.GasUsed > 0 || len(acc.StorageUpdates) > 0 || len(acc.Code) > 0 || len(acc.OutputTransfers) == 0
}

// OutputVerifier defines the exit-boundary verification contract for VMOutput
type OutputVerifier interface {
	Verify(input *VerificationInput, output *vmcommon.VMOutput) error
	IsInterfaceNil() bool
}

// VerificationInput provides the transaction-level context needed to verify a VMOutput
type VerificationInput struct {
	CallerAddr               []byte
	RecipientAddr            []byte
	CallValue                *big.Int
	ESDTTransfers            []*vmcommon.ESDTTransfer
	GasProvided              uint64
	GasLocked                uint64
	IsReadOnly               bool
	AllowCustomAddressLength bool
	DisjointDeletedAccounts  bool
}

// TokenKey uniquely identifies a token (fungible or NFT/SFT)
type TokenKey struct {
	Identifier string
	Nonce      uint64
}

type transferRecord struct {
	key    TokenKey
	amount *big.Int
}

type vmOutputVerifier struct {
	callArgsParser vmcommon.CallArgsParser
	marshaller     vmcommon.Marshalizer
}

// NewVMOutputVerifier creates a new exit-boundary verifier for VMOutput
func NewVMOutputVerifier() *vmOutputVerifier {
	return &vmOutputVerifier{
		callArgsParser: parsers.NewCallArgsParser(),
		marshaller:     &marshal.GogoProtoMarshalizer{},
	}
}

func (v *vmOutputVerifier) getMarshaller() vmcommon.Marshalizer {
	return v.marshaller
}

// IsInterfaceNil returns true if the verifier interface is nil
func (v *vmOutputVerifier) IsInterfaceNil() bool {
	return v == nil
}

// Verify executes all 7 invariant checks on the provided VMOutput
func (v *vmOutputVerifier) Verify(input *VerificationInput, output *vmcommon.VMOutput) error {
	if output == nil {
		return ErrNilVMOutput
	}
	if input == nil {
		return ErrNilVerificationInput
	}

	sortedAccountKeys := getSortedAccountKeys(output.OutputAccounts)

	err := v.verifyReturnCodeStateIsolation(output, sortedAccountKeys)
	if err != nil {
		return err
	}

	err = v.verifyIdentityAndAddressIntegrity(input, output, sortedAccountKeys)
	if err != nil {
		return err
	}

	err = v.verifyReadOnlyPurity(input, output, sortedAccountKeys)
	if err != nil {
		return err
	}

	err = v.verifyNativeEGLDConservation(input, output, sortedAccountKeys)
	if err != nil {
		return err
	}

	err = v.verifyNonNegativeBalances(input, output, sortedAccountKeys)
	if err != nil {
		return err
	}

	err = v.verifyTokenConservation(input, output, sortedAccountKeys)
	if err != nil {
		return err
	}

	err = v.verifyGasBounds(input, output, sortedAccountKeys)
	if err != nil {
		return err
	}

	return nil
}

// 1. ReturnCode State Isolation: If ReturnCode != Ok, no storage updates,
// balance deltas, code deployments, or output transfers are allowed.
func (v *vmOutputVerifier) verifyReturnCodeStateIsolation(output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	if output.ReturnCode == vmcommon.Ok {
		return nil
	}

	if len(output.DeletedAccounts) > 0 {
		return fmt.Errorf("%w: returnCode=%s deletedAccounts=%d",
			ErrDeletedAccountOnFailedExecution, output.ReturnCode, len(output.DeletedAccounts))
	}

	if output.GasRefund != nil && output.GasRefund.Sign() < 0 {
		return fmt.Errorf("%w: returnCode=%s gasRefund=%s",
			ErrGasRefundNegative, output.ReturnCode, output.GasRefund.String())
	}

	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil {
			continue
		}
		if len(acc.StorageUpdates) > 0 {
			return fmt.Errorf("%w: returnCode=%s account=%x storageUpdates=%d",
				ErrStorageUpdateOnFailedExecution, output.ReturnCode, acc.Address, len(acc.StorageUpdates))
		}
		if acc.BalanceDelta != nil && acc.BalanceDelta.Sign() != 0 {
			return fmt.Errorf("%w: returnCode=%s account=%x balanceDelta=%s",
				ErrBalanceDeltaOnFailedExecution, output.ReturnCode, acc.Address, acc.BalanceDelta.String())
		}
		if len(acc.Code) > 0 || len(acc.CodeMetadata) > 0 || acc.CodeDeployerAddress != nil {
			return fmt.Errorf("%w: returnCode=%s account=%x",
				ErrCodeDeploymentOnFailedExecution, output.ReturnCode, acc.Address)
		}
		if len(acc.OutputTransfers) > 0 {
			return fmt.Errorf("%w: returnCode=%s account=%x outputTransfers=%d",
				ErrOutputTransferOnFailedExecution, output.ReturnCode, acc.Address, len(acc.OutputTransfers))
		}
	}

	for _, logEntry := range output.Logs {
		if logEntry == nil {
			continue
		}
		eventID := string(logEntry.Identifier)
		if isStateMutatingEvent(eventID) {
			return fmt.Errorf("%w: returnCode=%s event=%s",
				ErrTokenTransferOnFailedExecution, output.ReturnCode, eventID)
		}
	}

	return nil
}

// 2. Identity & Address Integrity: UF-208 defense (output account map key == acc.Address,
// length == 32 bytes, disjoint or non-mutating deleted accounts).
func (v *vmOutputVerifier) verifyIdentityAndAddressIntegrity(input *VerificationInput, output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	var deletedSet map[string]struct{}
	if len(output.DeletedAccounts) > 0 {
		deletedSet = make(map[string]struct{}, len(output.DeletedAccounts))
		for _, delAddr := range output.DeletedAccounts {
			if len(delAddr) == 0 {
				return fmt.Errorf("%w: deleted account address is empty", ErrInvalidAccountAddress)
			}
			if !input.AllowCustomAddressLength && len(delAddr) != AddressLen {
				return fmt.Errorf("%w: deleted account address length %d != %d",
					ErrInvalidAccountAddress, len(delAddr), AddressLen)
			}
			delKey := string(delAddr)
			if _, duplicate := deletedSet[delKey]; duplicate {
				return fmt.Errorf("%w: address %x duplicated in deleted accounts",
					ErrDuplicateDeletedAccount, delAddr)
			}
			deletedSet[delKey] = struct{}{}
		}
	}

	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil {
			return fmt.Errorf("%w: nil output account for key %x", ErrInvalidAccountAddress, []byte(addrKey))
		}
		if len(acc.Address) == 0 {
			return fmt.Errorf("%w: account address is empty", ErrInvalidAccountAddress)
		}
		if !input.AllowCustomAddressLength && len(acc.Address) != AddressLen {
			return fmt.Errorf("%w: account address length %d != %d",
				ErrInvalidAccountAddress, len(acc.Address), AddressLen)
		}
		if addrKey != string(acc.Address) {
			return fmt.Errorf("%w: map key %x != acc.Address %x",
				ErrAddressKeyMismatch, []byte(addrKey), acc.Address)
		}
		if len(deletedSet) > 0 {
			if _, inDeleted := deletedSet[addrKey]; inDeleted {
				if input.DisjointDeletedAccounts {
					return fmt.Errorf("%w: account %x is both in OutputAccounts and DeletedAccounts",
						ErrDeletedAccountInOutputAccounts, acc.Address)
				}
				if len(acc.StorageUpdates) > 0 {
					return fmt.Errorf("%w: deleted account %x has storage updates",
						ErrDeletedAccountInOutputAccounts, acc.Address)
				}
				if len(acc.Code) > 0 || len(acc.CodeMetadata) > 0 || acc.CodeDeployerAddress != nil {
					return fmt.Errorf("%w: deleted account %x has code deployment",
						ErrDeletedAccountInOutputAccounts, acc.Address)
				}
				if len(acc.OutputTransfers) > 0 {
					return fmt.Errorf("%w: deleted account %x has output transfers",
						ErrDeletedAccountInOutputAccounts, acc.Address)
				}
			}
		}

		for _, outTransfer := range acc.OutputTransfers {
			if len(outTransfer.SenderAddress) > 0 && !input.AllowCustomAddressLength && len(outTransfer.SenderAddress) != AddressLen {
				return fmt.Errorf("%w: output transfer sender address length %d != %d",
					ErrInvalidAccountAddress, len(outTransfer.SenderAddress), AddressLen)
			}
		}
	}

	return nil
}

// 3. Native EGLD Conservation: Net sum of BalanceDelta across touched accounts
// matches deposit (input.CallValue), and no negative or unbounded output transfer values.
func (v *vmOutputVerifier) verifyNativeEGLDConservation(input *VerificationInput, output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	deposit := input.CallValue
	if deposit == nil {
		deposit = bigZero
	}
	if deposit.Sign() < 0 {
		return fmt.Errorf("%w: %s", ErrNegativeDeposit, deposit.String())
	}

	var claimedRewards *big.Int
	for _, logEntry := range output.Logs {
		if logEntry == nil {
			continue
		}
		if string(logEntry.Identifier) == core.BuiltInFunctionClaimDeveloperRewards {
			if len(logEntry.Topics) > 0 && len(logEntry.Topics[0]) > 0 {
				claimedAmount := big.NewInt(0).SetBytes(logEntry.Topics[0])
				if claimedRewards == nil {
					claimedRewards = claimedAmount
				} else {
					claimedRewards.Add(claimedRewards, claimedAmount)
				}
			}
		}
	}

	var netSum *big.Int
	multiDeltas := false
	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil {
			continue
		}
		if acc.BalanceDelta != nil && acc.BalanceDelta.Sign() != 0 {
			if netSum == nil {
				netSum = acc.BalanceDelta
			} else {
				if !multiDeltas {
					netSum = new(big.Int).Add(netSum, acc.BalanceDelta)
					multiDeltas = true
				} else {
					netSum.Add(netSum, acc.BalanceDelta)
				}
			}
		}

		for _, outTransfer := range acc.OutputTransfers {
			if outTransfer.Value != nil && outTransfer.Value.Sign() < 0 {
				return fmt.Errorf("%w: account %x transfer has negative value %s",
					ErrNegativeOutputTransferValue, acc.Address, outTransfer.Value.String())
			}
		}
	}
	if netSum == nil {
		netSum = bigZero
	}

	if output.ReturnCode != vmcommon.Ok {
		// On non-Ok returnCode, check 1 guarantees balance deltas are all zero/nil.
		if netSum.Sign() != 0 {
			return fmt.Errorf("%w: non-Ok returnCode has non-zero net delta %s",
				ErrNativeTokenConservationViolated, netSum.String())
		}
		return nil
	}

	if claimedRewards == nil {
		if netSum.Cmp(deposit) != 0 {
			return fmt.Errorf("%w: deposit=%s netSumDelta=%s",
				ErrNativeTokenConservationViolated, deposit.String(), netSum.String())
		}
	} else {
		depositWithRewards := new(big.Int).Add(deposit, claimedRewards)
		if netSum.Cmp(deposit) != 0 && (claimedRewards.Sign() == 0 || netSum.Cmp(depositWithRewards) != 0) {
			return fmt.Errorf("%w: deposit=%s netSumDelta=%s",
				ErrNativeTokenConservationViolated, deposit.String(), netSum.String())
		}
	}

	return nil
}

// 4. Non-Negative Balance Check: Ensures no account balance is overdrawn below zero.
func (v *vmOutputVerifier) verifyNonNegativeBalances(_ *VerificationInput, output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil {
			continue
		}

		if acc.Balance != nil {
			if acc.Balance.Sign() < 0 {
				return fmt.Errorf("%w: initial balance %s < 0 for account %x",
					ErrNegativeBalanceInOutput, acc.Balance.String(), acc.Address)
			}
			if acc.BalanceDelta != nil && acc.BalanceDelta.Sign() < 0 {
				finalBalance := new(big.Int).Add(acc.Balance, acc.BalanceDelta)
				if finalBalance.Sign() < 0 {
					return fmt.Errorf("%w: account %x final balance %s < 0 (initial %s + delta %s)",
						ErrAccountOverdrawn, acc.Address, finalBalance.String(), acc.Balance.String(), acc.BalanceDelta.String())
				}
			}
		}
	}

	return nil
}

// 5. Token Conservation (ESDT): Validates token integrity and ensures all external output transfers
// and transfer logs match 1-to-1 with identical tokens, nonces, and amounts.
// Note: Balance solvency for ESDT transfers is enforced during VM execution by the blockchain
// built-in functions (e.g. ESDTTransfer, MultiESDTNFTTransfer). The OutputVerifier acts as an
// exit-boundary structural validator: it guarantees that no unlogged output transfers are created,
// no transfer logs are spoofed, all amounts are non-negative, and intra-shard execution boundaries
// are respected.
func (v *vmOutputVerifier) verifyTokenConservation(input *VerificationInput, output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	// 5.1 Validate Inflows from input.ESDTTransfers
	for _, transfer := range input.ESDTTransfers {
		if transfer == nil || transfer.ESDTValue == nil {
			continue
		}
		if transfer.ESDTValue.Sign() < 0 {
			return fmt.Errorf("%w: inflow transfer has negative value %s",
				ErrNegativeTokenAmount, transfer.ESDTValue.String())
		}
	}

	// 5.2 Validate Mint and Burn operations from output.Logs
	for _, logEntry := range output.Logs {
		if logEntry == nil {
			continue
		}
		eventID := string(logEntry.Identifier)
		switch eventID {
		case core.BuiltInFunctionESDTLocalMint,
			core.BuiltInFunctionESDTNFTAddQuantity,
			core.BuiltInFunctionESDTNFTCreate:
			_, _, _, err := parseTokenEventTopics(logEntry.Topics, "mint")
			if err != nil {
				return err
			}
		case core.BuiltInFunctionESDTBurn,
			core.BuiltInFunctionESDTLocalBurn,
			core.BuiltInFunctionESDTNFTBurn,
			core.BuiltInFunctionESDTWipe:
			_, _, _, err := parseTokenEventTopics(logEntry.Topics, "burn")
			if err != nil {
				return err
			}
		}
	}

	// 5.3 Extract Outflows from OutputTransfers
	var outflowTransfers []transferRecord
	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil || len(acc.OutputTransfers) == 0 {
			continue
		}
		if isExecutedContract(input, output, acc.Address) {
			continue
		}
		for _, outTransfer := range acc.OutputTransfers {
			if len(outTransfer.Data) == 0 {
				continue
			}
			extracted, err := v.extractESDTTransfersFromData(outTransfer.Data)
			if err != nil {
				return err
			}
			for _, rec := range extracted {
				if rec.amount.Sign() < 0 {
					return fmt.Errorf("%w: output transfer has negative amount %s",
						ErrNegativeTokenAmount, rec.amount.String())
				}
				rec.key.Identifier = normalizeTokenID(rec.key.Identifier)
				outflowTransfers = append(outflowTransfers, rec)
			}
		}
	}

	// 5.4 Extract Transfer Logs from output.Logs
	var transferLogs []transferRecord
	for _, logEntry := range output.Logs {
		if logEntry == nil {
			continue
		}
		eventID := string(logEntry.Identifier)
		if eventID == core.BuiltInFunctionESDTTransfer || eventID == core.BuiltInFunctionESDTNFTTransfer {
			if len(logEntry.Topics) >= 3 {
				tokenID := string(logEntry.Topics[0])
				if len(tokenID) == 0 || isEGLDToken(tokenID) {
					continue
				}
				if len(logEntry.Topics) >= 4 && isExecutedContract(input, output, logEntry.Topics[3]) {
					continue
				}
				nonce, err := parseNonce(logEntry.Topics[1])
				if err != nil {
					return err
				}
				amount := parseValue(logEntry.Topics[2])
				if amount.Sign() < 0 {
					return fmt.Errorf("%w: transfer event has negative amount %s",
						ErrNegativeTokenAmount, amount.String())
				}
				transferLogs = append(transferLogs, transferRecord{
					key:    TokenKey{Identifier: normalizeTokenID(tokenID), Nonce: nonce},
					amount: amount,
				})
			}
		} else if eventID == core.BuiltInFunctionMultiESDTNFTTransfer {
			numOfTopics := len(logEntry.Topics)
			if numOfTopics >= 4 && numOfTopics%3 == 1 {
				destAddr := logEntry.Topics[numOfTopics-1]
				if isExecutedContract(input, output, destAddr) {
					continue
				}
			}
			numTriplets := numOfTopics / 3
			for i := 0; i < numTriplets; i++ {
				tokenID := string(logEntry.Topics[3*i])
				if len(tokenID) == 0 || isEGLDToken(tokenID) {
					continue
				}
				nonce, err := parseNonce(logEntry.Topics[3*i+1])
				if err != nil {
					return err
				}
				amount := parseValue(logEntry.Topics[3*i+2])
				if amount.Sign() < 0 {
					return fmt.Errorf("%w: multi transfer event has negative amount %s",
						ErrNegativeTokenAmount, amount.String())
				}
				transferLogs = append(transferLogs, transferRecord{
					key:    TokenKey{Identifier: normalizeTokenID(tokenID), Nonce: nonce},
					amount: amount,
				})
			}
		}
	}

	// Short-circuit if neither output transfers nor transfer logs exist
	if len(outflowTransfers) == 0 && len(transferLogs) == 0 {
		return nil
	}

	// 5.5 Verify 1-to-1 Matching Between OutputTransfers and Transfer Logs
	// Uses O(1) swap-delete to avoid slice shifting allocations.
	unmatchedLogs := make([]transferRecord, len(transferLogs))
	copy(unmatchedLogs, transferLogs)

	for _, tr := range outflowTransfers {
		matchedIdx := -1
		for i, logRec := range unmatchedLogs {
			if logRec.key == tr.key && logRec.amount.Cmp(tr.amount) == 0 {
				matchedIdx = i
				break
			}
		}
		if matchedIdx < 0 {
			return fmt.Errorf("%w: %w: token %s (nonce %d, amount %s)",
				ErrTokenConservationViolated, ErrOutputTransferWithoutTransferLog,
				tr.key.Identifier, tr.key.Nonce, tr.amount.String())
		}
		lastIdx := len(unmatchedLogs) - 1
		unmatchedLogs[matchedIdx] = unmatchedLogs[lastIdx]
		unmatchedLogs = unmatchedLogs[:lastIdx]
	}

	if len(unmatchedLogs) > 0 {
		first := unmatchedLogs[0]
		return fmt.Errorf("%w: %w: token %s (nonce %d, amount %s)",
			ErrTokenConservationViolated, ErrTransferLogWithoutOutputTransfer,
			first.key.Identifier, first.key.Nonce, first.amount.String())
	}

	return nil
}

func (v *vmOutputVerifier) extractESDTTransfersFromData(data []byte) ([]transferRecord, error) {
	results := make([]transferRecord, 0)

	fn, args, err := v.callArgsParser.ParseData(string(data))
	if err != nil {
		strData := string(data)
		if strings.HasPrefix(strData, core.BuiltInFunctionESDTTransfer) ||
			strings.HasPrefix(strData, core.BuiltInFunctionESDTNFTTransfer) ||
			strings.HasPrefix(strData, core.BuiltInFunctionMultiESDTNFTTransfer) {
			return nil, fmt.Errorf("%w: %s", ErrMalformedOutputTransferData, err.Error())
		}
		return results, nil
	}

	switch fn {
	case core.BuiltInFunctionESDTTransfer:
		if len(args) >= 2 {
			tokenID := string(args[0])
			if len(tokenID) == 0 || isEGLDToken(tokenID) {
				return results, nil
			}
			amount := big.NewInt(0).SetBytes(args[1])
			if amount.Sign() < 0 {
				return nil, fmt.Errorf("%w: negative amount %s in ESDTTransfer",
					ErrNegativeTokenAmount, amount.String())
			}
			results = append(results, transferRecord{
				key:    TokenKey{Identifier: tokenID, Nonce: 0},
				amount: amount,
			})
		}
	case core.BuiltInFunctionESDTNFTTransfer:
		if len(args) >= 3 {
			tokenID := string(args[0])
			if len(tokenID) == 0 || isEGLDToken(tokenID) {
				return results, nil
			}
			nonce, err := parseNonce(args[1])
			if err != nil {
				return nil, err
			}
			amount := v.extractTokenAmount(args[2], nonce)
			if amount.Sign() < 0 {
				return nil, fmt.Errorf("%w: negative amount %s in ESDTNFTTransfer",
					ErrNegativeTokenAmount, amount.String())
			}
			results = append(results, transferRecord{
				key:    TokenKey{Identifier: tokenID, Nonce: nonce},
				amount: amount,
			})
		}
	case core.BuiltInFunctionMultiESDTNFTTransfer:
		// Format A: MultiESDTNFTTransfer@destination@numTransfers@token1@nonce1@val1...
		// Format B: MultiESDTNFTTransfer@numTransfers@token1@nonce1@val1...
		if len(args) >= 2 {
			firstArg := args[0]
			tripletStart := 1
			numTransfers := big.NewInt(0).SetBytes(firstArg).Int64()
			if len(firstArg) == AddressLen && len(args) >= 3 {
				tripletStart = 2
				numTransfers = big.NewInt(0).SetBytes(args[1]).Int64()
			}
			if numTransfers <= 0 || numTransfers > 1000 {
				return results, nil
			}
			for i := int64(0); i < numTransfers; i++ {
				idx := tripletStart + int(i*3)
				if idx+2 >= len(args) {
					break
				}
				tokenID := string(args[idx])
				if len(tokenID) == 0 || isEGLDToken(tokenID) {
					continue
				}
				nonce, err := parseNonce(args[idx+1])
				if err != nil {
					return nil, err
				}
				amount := v.extractTokenAmount(args[idx+2], nonce)
				if amount.Sign() < 0 {
					return nil, fmt.Errorf("%w: negative amount %s in MultiESDTNFTTransfer",
						ErrNegativeTokenAmount, amount.String())
				}
				results = append(results, transferRecord{
					key:    TokenKey{Identifier: tokenID, Nonce: nonce},
					amount: amount,
				})
			}
		}
	}

	return results, nil
}

func parseNonce(topic []byte) (uint64, error) {
	if len(topic) == 0 {
		return 0, nil
	}
	if len(topic) > 8 {
		return 0, fmt.Errorf("%w: topic length %d > 8", ErrInvalidTokenNonce, len(topic))
	}
	var nonce uint64
	for _, b := range topic {
		nonce = (nonce << 8) | uint64(b)
	}
	return nonce, nil
}

func parseValue(topic []byte) *big.Int {
	if len(topic) == 0 {
		return big.NewInt(0)
	}
	return big.NewInt(0).SetBytes(topic)
}

// 6. Gas Bounds & 50% Refund Cap: GasRemaining <= GasProvided,
// GasRefund >= 0 && GasRefund <= (GasProvided - GasRemaining)/2.
func (v *vmOutputVerifier) verifyGasBounds(input *VerificationInput, output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	if output.GasRemaining > input.GasProvided {
		return fmt.Errorf("%w: gasRemaining %d > gasProvided %d",
			ErrGasRemainingExceedsProvided, output.GasRemaining, input.GasProvided)
	}

	sumGasUsed := uint64(0)
	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil {
			continue
		}
		if acc.GasUsed > input.GasProvided {
			return fmt.Errorf("%w: account %x gasUsed %d > gasProvided %d",
				ErrGasRemainingExceedsProvided, acc.Address, acc.GasUsed, input.GasProvided)
		}
		sumGasUsed += acc.GasUsed
		for _, tr := range acc.OutputTransfers {
			if tr.GasLimit > input.GasProvided {
				return fmt.Errorf("%w: transfer gasLimit %d > gasProvided %d",
					ErrGasRemainingExceedsProvided, tr.GasLimit, input.GasProvided)
			}
		}
	}
	if sumGasUsed > input.GasProvided {
		return fmt.Errorf("%w: sum of accounts gasUsed %d > gasProvided %d",
			ErrGasRemainingExceedsProvided, sumGasUsed, input.GasProvided)
	}

	if output.GasRefund == nil {
		return nil
	}

	if output.GasRefund.Sign() < 0 {
		return fmt.Errorf("%w: gasRefund %s < 0", ErrGasRefundNegative, output.GasRefund.String())
	}

	gasSpent := input.GasProvided - output.GasRemaining
	maxRefund := gasSpent / 2

	if !output.GasRefund.IsUint64() || output.GasRefund.Uint64() > maxRefund {
		output.GasRefund.SetUint64(maxRefund)
	}

	return nil
}

// 7. ReadOnly Purity: ReadOnly executions cannot delete accounts, deploy code,
// transfer value, or write storage.
func (v *vmOutputVerifier) verifyReadOnlyPurity(input *VerificationInput, output *vmcommon.VMOutput, sortedAccountKeys []string) error {
	if !input.IsReadOnly {
		return nil
	}

	if len(output.DeletedAccounts) > 0 {
		return fmt.Errorf("%w: deleted accounts count %d > 0",
			ErrReadOnlyStateChange, len(output.DeletedAccounts))
	}

	if output.GasRefund != nil && output.GasRefund.Sign() > 0 {
		return fmt.Errorf("%w: gas refund > 0 in read-only mode", ErrReadOnlyStateChange)
	}

	for _, addrKey := range sortedAccountKeys {
		acc := output.OutputAccounts[addrKey]
		if acc == nil {
			continue
		}
		if len(acc.StorageUpdates) > 0 {
			return fmt.Errorf("%w: account %x has storage updates in read-only mode",
				ErrReadOnlyStateChange, acc.Address)
		}
		if acc.BalanceDelta != nil && acc.BalanceDelta.Sign() != 0 {
			return fmt.Errorf("%w: account %x has balance delta %s in read-only mode",
				ErrReadOnlyStateChange, acc.Address, acc.BalanceDelta.String())
		}
		if len(acc.Code) > 0 || len(acc.CodeMetadata) > 0 || acc.CodeDeployerAddress != nil {
			return fmt.Errorf("%w: account %x has code deployment in read-only mode",
				ErrReadOnlyStateChange, acc.Address)
		}
		if len(acc.OutputTransfers) > 0 {
			return fmt.Errorf("%w: account %x has output transfers in read-only mode",
				ErrReadOnlyStateChange, acc.Address)
		}
	}

	for _, logEntry := range output.Logs {
		if logEntry == nil {
			continue
		}
		eventID := string(logEntry.Identifier)
		if isStateMutatingEvent(eventID) {
			return fmt.Errorf("%w: event %s in read-only mode", ErrReadOnlyStateChange, eventID)
		}
	}

	return nil
}

func isStateMutatingEvent(eventID string) bool {
	switch eventID {
	case core.BuiltInFunctionESDTTransfer,
		core.BuiltInFunctionESDTNFTTransfer,
		core.BuiltInFunctionMultiESDTNFTTransfer,
		core.BuiltInFunctionESDTBurn,
		core.BuiltInFunctionESDTLocalBurn,
		core.BuiltInFunctionESDTNFTBurn,
		core.BuiltInFunctionESDTWipe,
		core.BuiltInFunctionESDTLocalMint,
		core.BuiltInFunctionESDTNFTAddQuantity,
		core.BuiltInFunctionESDTNFTCreate,
		core.BuiltInFunctionESDTFreeze,
		core.BuiltInFunctionESDTUnFreeze,
		core.BuiltInFunctionChangeOwnerAddress,
		core.BuiltInFunctionClaimDeveloperRewards,
		core.BuiltInFunctionSetUserName,
		core.BuiltInFunctionSaveKeyValue,
		"transferValueOnly":
		return true
	default:
		return false
	}
}

// Helper constructor to build VerificationInput from ContractCallInput
func NewVerificationInputFromCall(
	input *vmcommon.ContractCallInput,
	isReadOnly bool,
) *VerificationInput {
	if input == nil {
		return &VerificationInput{}
	}

	callValue := input.CallValue
	if callValue == nil {
		callValue = big.NewInt(0)
	}

	readOnly := isReadOnly

	return &VerificationInput{
		CallerAddr:    bytes.Clone(input.CallerAddr),
		RecipientAddr: bytes.Clone(input.RecipientAddr),
		CallValue:     big.NewInt(0).Set(callValue),
		ESDTTransfers: input.ESDTTransfers,
		GasProvided:   input.GasProvided,
		GasLocked:     input.GasLocked,
		IsReadOnly:    readOnly,
	}
}

// Helper constructor to build VerificationInput from ContractCreateInput
func NewVerificationInputFromCreate(
	input *vmcommon.ContractCreateInput,
) *VerificationInput {
	if input == nil {
		return &VerificationInput{}
	}

	callValue := input.CallValue
	if callValue == nil {
		callValue = big.NewInt(0)
	}

	return &VerificationInput{
		CallerAddr:    bytes.Clone(input.CallerAddr),
		RecipientAddr: nil,
		CallValue:     big.NewInt(0).Set(callValue),
		ESDTTransfers: nil,
		GasProvided:   input.GasProvided,
		GasLocked:     input.GasLocked,
		IsReadOnly:    false,
	}
}

func parseTokenEventTopics(topics [][]byte, eventType string) (TokenKey, *big.Int, bool, error) {
	if len(topics) < 3 {
		return TokenKey{}, nil, false, nil
	}
	tokenID := string(topics[0])
	if len(tokenID) == 0 || isEGLDToken(tokenID) {
		return TokenKey{}, nil, false, nil
	}
	nonce, err := parseNonce(topics[1])
	if err != nil {
		return TokenKey{}, nil, false, err
	}
	amount := parseValue(topics[2])
	if amount.Sign() < 0 {
		return TokenKey{}, nil, false, fmt.Errorf("%w: %s event has negative amount %s",
			ErrNegativeTokenAmount, eventType, amount.String())
	}
	return TokenKey{Identifier: tokenID, Nonce: nonce}, amount, true, nil
}

func (v *vmOutputVerifier) extractTokenAmount(valBytes []byte, nonce uint64) *big.Int {
	amount := big.NewInt(0)
	if nonce > 0 && len(valBytes) > vmcommon.MaxLengthForValueToOptTransfer {
		transferESDTData := &esdt.ESDigitalToken{}
		err := v.getMarshaller().Unmarshal(transferESDTData, valBytes)
		if err == nil && transferESDTData.Value != nil {
			amount.Set(transferESDTData.Value)
			return amount
		}
	}
	amount.SetBytes(valBytes)
	return amount
}
