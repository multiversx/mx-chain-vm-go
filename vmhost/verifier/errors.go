package verifier

import "errors"

// Invariant violation errors returned by VMOutputVerifier
var (
	// ErrNilVMOutput signals that the VM output is nil
	ErrNilVMOutput = errors.New("vm output is nil")

	// ErrNilVerificationInput signals that the verification input is nil
	ErrNilVerificationInput = errors.New("verification input is nil")

	// ReturnCode state isolation errors
	ErrStorageUpdateOnFailedExecution  = errors.New("storage update found on failed execution")
	ErrBalanceDeltaOnFailedExecution   = errors.New("balance delta found on failed execution")
	ErrCodeDeploymentOnFailedExecution = errors.New("code deployment found on failed execution")
	ErrOutputTransferOnFailedExecution = errors.New("output transfer found on failed execution")
	ErrDeletedAccountOnFailedExecution = errors.New("deleted account found on failed execution")
	ErrGasRefundOnFailedExecution      = errors.New("gas refund found on failed execution")
	ErrTokenTransferOnFailedExecution  = errors.New("token transfer or mint/burn log found on failed execution")

	// Identity and address integrity errors (UF-208)
	ErrInvalidAccountAddress          = errors.New("invalid account address: must be exactly 32 bytes")
	ErrAddressKeyMismatch             = errors.New("output account address does not match map key")
	ErrDeletedAccountInOutputAccounts = errors.New("deleted account exists in output accounts")
	ErrDuplicateDeletedAccount        = errors.New("duplicate address in deleted accounts")

	// Native EGLD conservation errors
	ErrNativeTokenConservationViolated   = errors.New("native token conservation violated: balance deltas net sum does not match deposit")
	ErrNegativeDeposit                   = errors.New("negative call value deposit")
	ErrNegativeOutputTransferValue       = errors.New("negative output transfer value")
	ErrOutputTransferExceedsBalanceDelta = errors.New("output transfer value exceeds balance delta")

	// Account balance errors
	ErrAccountOverdrawn        = errors.New("account balance overdrawn below zero")
	ErrNegativeBalanceInOutput = errors.New("account balance is negative in output")

	// ESDT token conservation errors
	ErrTokenConservationViolated        = errors.New("ESDT token conservation violated: outflow plus burn exceeds inflow plus net mint")
	ErrNegativeTokenAmount              = errors.New("negative ESDT token amount")
	ErrInvalidTokenNonce                = errors.New("invalid token nonce: exceeds 8 bytes")
	ErrMalformedOutputTransferData      = errors.New("malformed ESDT transfer data in output transfer")
	ErrOutputTransferWithoutTransferLog = errors.New("output transfer has no matching transfer log")
	ErrTransferLogWithoutOutputTransfer = errors.New("transfer log has no matching output transfer")

	// Gas bounds errors
	ErrGasRemainingExceedsProvided = errors.New("gas remaining exceeds gas provided")
	ErrGasRefundNegative           = errors.New("gas refund is negative")
	ErrGasRefundExceedsCap         = errors.New("gas refund exceeds 50% cap of spent gas")

	// Read-only purity errors
	ErrReadOnlyStateChange = errors.New("read-only execution violation: attempted state mutation")
)
