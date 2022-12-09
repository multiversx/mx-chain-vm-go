package hosttest

import (
	"bytes"
	"fmt"
	"testing"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/wasm-vm-v1_4/arwen"
	worldmock "github.com/ElrondNetwork/wasm-vm-v1_4/mock/world"
	"github.com/ElrondNetwork/wasm-vm-v1_4/testcommon"
	"github.com/stretchr/testify/require"
)

var codeUpgrader []byte
var codeInitSimple []byte
var codeUpgraderMarkOffset int
var codeInitSimpleMarkOffset int

var logUpgTest = logger.GetOrCreate("arwen/test")

var defaultDeployInput *vmcommon.ContractCreateInput
var defaultCallInput *vmcommon.ContractCallInput

func TestUpgrade_WithWorldMock(t *testing.T) {
	codeUpgrader = testcommon.GetTestSCCode("upgrader", "../../")
	codeUpgraderMarkOffset = bytes.Index(codeUpgrader, []byte("finish0000"))

	codeInitSimple = testcommon.GetTestSCCode("init-simple", "../../")
	codeInitSimpleMarkOffset = bytes.Index(codeInitSimple, []byte("finish0000"))

	setupStructs()

	usc := newUpgradeScenario(t)

	nUpgradeIterations := 4000
	nPairs := 1000
	contractPairs := make([]*upgradeSCPair, nPairs)
	for i := 0; i < nPairs; i++ {
		pair := newUpgradeSCPair(usc.ownerAccount.Address, i)
		pair.initialize(usc)
		contractPairs[i] = pair
	}

	validateTest(t, usc, contractPairs)

	logger.SetLogLevel("*:NONE,arwen/test:TRACE")
	for u := 0; u < nUpgradeIterations; u++ {
		logUpgTest.Trace("beginning upgrade iteration", "u", u)
		for _, pair := range contractPairs {
			pair.upgradeChild(usc, pair.index+u)
		}
	}
}

func validateTest(
	t *testing.T,
	usc *upgradeScenario,
	contractPairs []*upgradeSCPair,
) {
	for i, pair := range contractPairs {
		vmOutput := usc.callTestSC(
			pair.ownerAddress,
			pair.upgraderAddress,
			"dummy",
			nil,
		)
		require.Equal(t, fmt.Sprintf("finish%04d", i), string(vmOutput.ReturnData[0]))

		vmOutput = usc.callTestSC(
			pair.ownerAddress,
			pair.childAddress,
			"dummy",
			nil,
		)
		require.Equal(t, fmt.Sprintf("finish%04d", i), string(vmOutput.ReturnData[0]))
	}
}

type upgradeSCPair struct {
	index           int
	ownerAddress    Address
	upgraderAddress Address
	childAddress    Address
}

func newUpgradeSCPair(ownerAddress Address, index int) *upgradeSCPair {
	return &upgradeSCPair{
		ownerAddress:    ownerAddress,
		index:           index,
		upgraderAddress: testcommon.MakeTestSCAddress(fmt.Sprintf("upgrader%04d", index)),
		childAddress:    testcommon.MakeTestSCAddress(fmt.Sprintf("child%04d", index)),
	}
}

func (pair *upgradeSCPair) initialize(usc *upgradeScenario) {
	pair.deployUpgrader(usc)
	pair.deployChild(usc)
}

func (pair *upgradeSCPair) deployUpgrader(usc *upgradeScenario) {
	usc.deployTestSC(
		pair.ownerAddress,
		pair.upgraderAddress,
		makeModifiedBytecodeUpgrader(pair.index),
	)
}

func (pair *upgradeSCPair) deployChild(usc *upgradeScenario) {
	usc.setupNewSCAddress(pair.upgraderAddress, pair.childAddress, false)
	usc.callTestSC(
		pair.ownerAddress,
		pair.upgraderAddress,
		"deployChildContract",
		[][]byte{makeModifiedBytecodeInitSimple(pair.index)},
	)
}

func (pair *upgradeSCPair) upgradeChild(usc *upgradeScenario, index int) {
	usc.callTestSC(
		pair.ownerAddress,
		pair.upgraderAddress,
		"upgradeChildContract",
		[][]byte{
			pair.childAddress,
			makeModifiedBytecodeInitSimple(pair.index),
		},
	)
}

type upgradeScenario struct {
	tb           testing.TB
	host         arwen.VMHost
	world        *worldmock.MockWorld
	ownerAccount *worldmock.Account
}

func newUpgradeScenario(tb testing.TB) *upgradeScenario {
	world, ownerAccount, host, err := prepare(tb)
	require.Nil(tb, err)
	require.NotNil(tb, world)
	require.NotNil(tb, ownerAccount)
	require.NotNil(tb, host)

	return &upgradeScenario{
		tb:           tb,
		host:         host,
		world:        world,
		ownerAccount: ownerAccount,
	}
}

func (usc *upgradeScenario) deployTestSC(
	parentAddress Address,
	scAddress Address,
	code []byte,
) {
	deployInput := makeDeployInput(usc.ownerAccount.Address, code, nil)
	usc.setupNewSCAddress(usc.ownerAccount.Address, scAddress, true)
	usc.deployTestSCOnWorld(deployInput)
}

func (usc *upgradeScenario) callTestSC(
	callerAddress Address,
	scAddress Address,
	function string,
	args [][]byte,
) *vmcommon.VMOutput {
	callInput := makeCallInput(callerAddress, scAddress, function, args)
	return usc.callTestSCOnWorld(callInput)
}

func (usc *upgradeScenario) deployTestSCOnWorld(
	deployInput *vmcommon.ContractCreateInput,
) {
	vmOutput, err := usc.host.RunSmartContractCreate(deployInput)
	require.Nil(usc.tb, err)
	require.NotNil(usc.tb, vmOutput)
	require.Equal(usc.tb, "", vmOutput.ReturnMessage)
	require.Equal(usc.tb, vmcommon.Ok, vmOutput.ReturnCode)
	_ = usc.world.UpdateAccounts(vmOutput.OutputAccounts, nil)
}

func (usc *upgradeScenario) callTestSCOnWorld(
	callInput *vmcommon.ContractCallInput,
) *vmcommon.VMOutput {
	vmOutput, err := usc.host.RunSmartContractCall(callInput)
	require.Nil(usc.tb, err)
	require.NotNil(usc.tb, vmOutput)
	require.Equal(usc.tb, "", vmOutput.ReturnMessage)
	require.Equal(usc.tb, vmcommon.Ok, vmOutput.ReturnCode)
	_ = usc.world.UpdateAccounts(vmOutput.OutputAccounts, nil)

	return vmOutput
}

func (usc *upgradeScenario) setupNewSCAddress(
	ownerAddress Address,
	newAddress Address,
	incrementOwnerNonce bool,
) {
	ownerAccount := usc.world.AcctMap.GetAccount(ownerAddress)
	require.NotNil(usc.tb, ownerAccount)

	usc.world.NewAddressMocks = append(usc.world.NewAddressMocks, &worldmock.NewAddressMock{
		CreatorAddress: ownerAccount.Address,
		CreatorNonce:   ownerAccount.Nonce,
		NewAddress:     newAddress,
	})

	logUpgTest.Trace("setupNewSCAddress",
		"creatorAddress", string(ownerAccount.Address),
		"creatorNonce", ownerAccount.Nonce,
		"newAddress", string(newAddress),
	)

	if incrementOwnerNonce {
		ownerAccount.Nonce++
	}
}

func makeDeployInput(ownerAddress Address, code []byte, args [][]byte) *vmcommon.ContractCreateInput {
	defaultDeployInput.CallerAddr = ownerAddress
	defaultDeployInput.ContractCode = code
	defaultDeployInput.Arguments = args

	return defaultDeployInput
}

func makeCallInput(callerAddress Address, scAddress Address, function string, args [][]byte) *vmcommon.ContractCallInput {
	defaultCallInput.CallerAddr = callerAddress
	defaultCallInput.RecipientAddr = scAddress
	defaultCallInput.Function = function
	defaultCallInput.Arguments = args

	return defaultCallInput
}

func makeModifiedBytecodeUpgrader(index int) []byte {
	return makeModifiedBytecode(
		codeUpgrader,
		codeUpgraderMarkOffset,
		index,
	)
}

func makeModifiedBytecodeInitSimple(index int) []byte {
	return makeModifiedBytecode(
		codeInitSimple,
		codeInitSimpleMarkOffset,
		index,
	)
}

func makeModifiedBytecode(originalCode []byte, offset int, index int) []byte {
	replacement := []byte(fmt.Sprintf("finish%04d", index))

	newCode := make([]byte, len(originalCode))
	copy(newCode, originalCode)

	for i, b := range replacement {
		newCode[offset+i] = b
	}
	return newCode
}

func setupStructs() {
	deployBuilder := testcommon.CreateTestContractCreateInputBuilder().
		WithGasProvided(0xFFFFFFFFFFFFFFFF)

	defaultDeployInput = deployBuilder.Build()

	callBuilder := testcommon.CreateTestContractCallInputBuilder().
		WithCallValue(0).
		WithGasProvided(0xFFFFFFFFFFFFFFFF)

	defaultCallInput = callBuilder.Build()
}
