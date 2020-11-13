package mandosjsonmodel

// Scenario is a json object representing a test scenario with steps.
type Scenario struct {
	Name     string
	Comment  string
	CheckGas bool
	Steps    []Step
}

// Step is the basic block of a scenario.
type Step interface {
	StepTypeName() string
}

// NewAddressMock allows tests to specify what new addresses to generate
type NewAddressMock struct {
	CreatorAddress JSONBytesFromString
	CreatorNonce   JSONUint64
	NewAddress     JSONBytesFromString
}

// BlockInfo contains data for the block info hooks
type BlockInfo struct {
	BlockTimestamp JSONUint64
	BlockNonce     JSONUint64
	BlockRound     JSONUint64
	BlockEpoch     JSONUint64
}

// ExternalStepsStep allows including steps from another file
type ExternalStepsStep struct {
	Path string
}

// SetStateStep is a step where data is saved to the blockchain mock.
type SetStateStep struct {
	Comment           string
	Accounts          []*Account
	PreviousBlockInfo *BlockInfo
	CurrentBlockInfo  *BlockInfo
	BlockHashes       []JSONBytesFromString
	NewAddressMocks   []*NewAddressMock
}

// CheckStateStep is a step where the state of the blockchain mock is verified.
type CheckStateStep struct {
	Comment       string
	CheckAccounts *CheckAccounts
}

// DumpStateStep is a step that simply prints the entire state to console. Useful for debugging.
type DumpStateStep struct {
	Comment string
}

// TxStep is a step where a transaction is executed.
type TxStep struct {
	TxIdent        string
	Comment        string
	Tx             *Transaction
	ExpectedResult *TransactionResult
}

var _ Step = (*ExternalStepsStep)(nil)
var _ Step = (*SetStateStep)(nil)
var _ Step = (*CheckStateStep)(nil)
var _ Step = (*DumpStateStep)(nil)
var _ Step = (*TxStep)(nil)

// StepNameExternalSteps is a json step type name.
const StepNameExternalSteps = "externalSteps"

// StepTypeName type as string
func (*ExternalStepsStep) StepTypeName() string {
	return StepNameExternalSteps
}

// StepNameSetState is a json step type name.
const StepNameSetState = "setState"

// StepTypeName type as string
func (*SetStateStep) StepTypeName() string {
	return StepNameSetState
}

// StepNameCheckState is a json step type name.
const StepNameCheckState = "checkState"

// StepTypeName type as string
func (*CheckStateStep) StepTypeName() string {
	return StepNameCheckState
}

// StepNameDumpState is a json step type name.
const StepNameDumpState = "dumpState"

// StepTypeName type as string
func (*DumpStateStep) StepTypeName() string {
	return StepNameDumpState
}

// StepNameScCall is a json step type name.
const StepNameScCall = "scCall"

// StepNameScDeploy is a json step type name.
const StepNameScDeploy = "scDeploy"

// StepNameTransfer is a json step type name.
const StepNameTransfer = "transfer"

// StepNameValidatorReward is a json step type name.
const StepNameValidatorReward = "validatorReward"

// StepTypeName type as string
func (t *TxStep) StepTypeName() string {
	switch t.Tx.Type {
	case ScCall:
		return StepNameScCall
	case ScDeploy:
		return StepNameScDeploy
	case Transfer:
		return StepNameTransfer
	case ValidatorReward:
		return StepNameValidatorReward
	default:
		panic("unknown TransactionType")
	}
}
