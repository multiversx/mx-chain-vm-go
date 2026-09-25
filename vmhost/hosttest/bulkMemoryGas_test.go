package hostCoretest

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	"github.com/multiversx/mx-chain-vm-go/config"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	gasschedules "github.com/multiversx/mx-chain-vm-go/scenario/gasSchedules"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/require"
)

// Deliberately not round numbers and not multiples of one another, so that a wrong
// per-byte cost cannot accidentally produce the expected gas delta.
const (
	memoryCopyCost        = uint64(300)
	memoryCopyPerByteCost = uint64(7)
	memoryFillCost        = uint64(400)
	memoryFillPerByteCost = uint64(11)
)

const bulkMemoryGasProvided = uint64(100_000_000)

// bulkMemoryGasSchedule takes the default test gas schedule and gives the bulk memory
// opcodes distinctive costs, so that their contribution can be isolated.
func bulkMemoryGasSchedule() config.GasScheduleMap {
	gasSchedule := config.MakeGasMapForTests()
	wasmOpcodeCost := gasSchedule["WASMOpcodeCost"]
	wasmOpcodeCost["MemoryCopy"] = memoryCopyCost
	wasmOpcodeCost["MemoryCopyPerByte"] = memoryCopyPerByteCost
	wasmOpcodeCost["MemoryFill"] = memoryFillCost
	wasmOpcodeCost["MemoryFillPerByte"] = memoryFillPerByteCost
	return gasSchedule
}

// bulkMemoryGasUsed calls the given endpoint of the bulk-memory contract with the number
// of bytes as argument 0 and returns the gas consumed by the call.
func bulkMemoryGasUsed(tb testing.TB, gasSchedule config.GasScheduleMap, function string, numBytes uint64) uint64 {
	var gasUsed uint64

	testcommon.BuildInstanceCallTest(tb).
		WithContracts(
			testcommon.CreateInstanceContract(testcommon.ParentAddress).
				WithCode(testcommon.GetTestSCCode("bulk-memory", "../../"))).
		WithGasSchedule(gasSchedule).
		WithInput(testcommon.CreateTestContractCallInputBuilder().
			WithGasProvided(bulkMemoryGasProvided).
			WithFunction(function).
			WithArguments(big.NewInt(0).SetUint64(numBytes).Bytes()).
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
			verify.Ok()
			require.LessOrEqual(tb, verify.VmOutput.GasRemaining, bulkMemoryGasProvided)
			gasUsed = bulkMemoryGasProvided - verify.VmOutput.GasRemaining
		})

	return gasUsed
}

// TestBulkMemoryGas_PerByteCostsSurviveGasConfigDecoding checks that the per-byte costs make
// it from the gas schedule map into the WASMOpcodeCost struct, which is what the host hands
// to the executor via SetOpcodeConfig.
func TestBulkMemoryGas_PerByteCostsSurviveGasConfigDecoding(t *testing.T) {
	gasSchedule := bulkMemoryGasSchedule()

	wasmOpcodeCost := gasSchedule["WASMOpcodeCost"]
	require.Contains(t, wasmOpcodeCost, "MemoryCopyPerByte")
	require.Contains(t, wasmOpcodeCost, "MemoryFillPerByte")

	gasCostConfig, err := config.CreateGasConfig(gasSchedule)
	require.Nil(t, err)

	require.Equal(t, uint32(memoryCopyCost), gasCostConfig.WASMOpcodeCost.MemoryCopy)
	require.Equal(t, uint32(memoryCopyPerByteCost), gasCostConfig.WASMOpcodeCost.MemoryCopyPerByte)
	require.Equal(t, uint32(memoryFillCost), gasCostConfig.WASMOpcodeCost.MemoryFill)
	require.Equal(t, uint32(memoryFillPerByteCost), gasCostConfig.WASMOpcodeCost.MemoryFillPerByte)
}

// TestBulkMemoryGas_ShippedSchedulesDefinePerByteCosts guards against a gas schedule that
// the node would accept while leaving the new per-byte costs unset.
func TestBulkMemoryGas_ShippedSchedulesDefinePerByteCosts(t *testing.T) {
	shippedSchedules := map[string]string{
		"V3": gasschedules.GetV3(),
		"V4": gasschedules.GetV4(),
	}

	for name, contents := range shippedSchedules {
		gasSchedule, err := gasschedules.LoadGasScheduleConfig(contents)
		require.Nil(t, err, name)

		require.Contains(t, gasSchedule["WASMOpcodeCost"], "MemoryCopyPerByte", name)
		require.Contains(t, gasSchedule["WASMOpcodeCost"], "MemoryFillPerByte", name)

		// CreateGasConfig rejects any WASMOpcodeCost field left at zero
		gasCostConfig, err := config.CreateGasConfig(gasSchedule)
		require.Nil(t, err, name)

		require.NotZero(t, gasCostConfig.WASMOpcodeCost.MemoryCopyPerByte, name)
		require.NotZero(t, gasCostConfig.WASMOpcodeCost.MemoryFillPerByte, name)
	}
}

// TestBulkMemoryGas_PerByteCostOfMemoryCopy checks that the gas difference between copying
// more and fewer bytes is exactly the number of extra bytes times MemoryCopyPerByte.
func TestBulkMemoryGas_PerByteCostOfMemoryCopy(t *testing.T) {
	gasSchedule := bulkMemoryGasSchedule()

	gasForZero := bulkMemoryGasUsed(t, gasSchedule, "memoryCopy", 0)
	gasForTen := bulkMemoryGasUsed(t, gasSchedule, "memoryCopy", 10)
	gasForThousand := bulkMemoryGasUsed(t, gasSchedule, "memoryCopy", 1000)

	require.Equal(t, 10*memoryCopyPerByteCost, gasForTen-gasForZero)
	require.Equal(t, 990*memoryCopyPerByteCost, gasForThousand-gasForTen)
	require.Equal(t, 1000*memoryCopyPerByteCost, gasForThousand-gasForZero)
}

// TestBulkMemoryGas_PerByteCostOfMemoryFill is TestBulkMemoryGas_PerByteCostOfMemoryCopy
// for memory.fill, which has its own, different per-byte cost.
func TestBulkMemoryGas_PerByteCostOfMemoryFill(t *testing.T) {
	gasSchedule := bulkMemoryGasSchedule()

	gasForZero := bulkMemoryGasUsed(t, gasSchedule, "memoryFill", 0)
	gasForTen := bulkMemoryGasUsed(t, gasSchedule, "memoryFill", 10)
	gasForThousand := bulkMemoryGasUsed(t, gasSchedule, "memoryFill", 1000)

	require.Equal(t, 10*memoryFillPerByteCost, gasForTen-gasForZero)
	require.Equal(t, 990*memoryFillPerByteCost, gasForThousand-gasForTen)
	require.Equal(t, 1000*memoryFillPerByteCost, gasForThousand-gasForZero)
}

// TestBulkMemoryGas_PerByteCostFollowsTheSchedule doubles the per-byte costs and checks that
// the gas delta doubles with them, which the base costs alone could not explain.
func TestBulkMemoryGas_PerByteCostFollowsTheSchedule(t *testing.T) {
	doubledSchedule := bulkMemoryGasSchedule()
	wasmOpcodeCost := doubledSchedule["WASMOpcodeCost"]
	wasmOpcodeCost["MemoryCopyPerByte"] = 2 * memoryCopyPerByteCost
	wasmOpcodeCost["MemoryFillPerByte"] = 2 * memoryFillPerByteCost

	copyDelta := bulkMemoryGasUsed(t, doubledSchedule, "memoryCopy", 1000) -
		bulkMemoryGasUsed(t, doubledSchedule, "memoryCopy", 0)
	require.Equal(t, 1000*2*memoryCopyPerByteCost, copyDelta)

	fillDelta := bulkMemoryGasUsed(t, doubledSchedule, "memoryFill", 1000) -
		bulkMemoryGasUsed(t, doubledSchedule, "memoryFill", 0)
	require.Equal(t, 1000*2*memoryFillPerByteCost, fillDelta)
}

// TestBulkMemoryGas_BaseCostIsPaidForZeroBytes isolates the base cost of each opcode, by
// comparing against an endpoint that does everything except the bulk memory opcode itself.
func TestBulkMemoryGas_BaseCostIsPaidForZeroBytes(t *testing.T) {
	gasSchedule := bulkMemoryGasSchedule()

	gasWithoutBulkMemory := bulkMemoryGasUsed(t, gasSchedule, "noBulkMemory", 0)

	// the *ViaLocal endpoints only add the opcode and the three operands it pops, which
	// the default test schedule prices at config.GasValueForTests each
	const operandsCost = 3 * uint64(config.GasValueForTests)

	copyGas := bulkMemoryGasUsed(t, gasSchedule, "memoryCopyViaLocal", 0)
	require.Equal(t, memoryCopyCost+operandsCost, copyGas-gasWithoutBulkMemory)

	fillGas := bulkMemoryGasUsed(t, gasSchedule, "memoryFillViaLocal", 0)
	require.Equal(t, memoryFillCost+operandsCost, fillGas-gasWithoutBulkMemory)
}

// TestBulkMemoryGas_SizeOperandSurvivesTheInjection checks that the metering code injected
// before the opcode gives the size operand back unchanged, whether it comes straight off the
// stack or through a local.
func TestBulkMemoryGas_SizeOperandSurvivesTheInjection(t *testing.T) {
	gasSchedule := bulkMemoryGasSchedule()

	// the *ViaLocal endpoints only add a local.set/local.get pair on top of their counterparts
	const localCost = 2 * uint64(config.GasValueForTests)

	for _, numBytes := range []uint64{0, 10, 1000} {
		copyDirect := bulkMemoryGasUsed(t, gasSchedule, "memoryCopy", numBytes)
		copyViaLocal := bulkMemoryGasUsed(t, gasSchedule, "memoryCopyViaLocal", numBytes)
		require.Equal(t, copyDirect+localCost, copyViaLocal, "memory.copy, numBytes = %d", numBytes)

		fillDirect := bulkMemoryGasUsed(t, gasSchedule, "memoryFill", numBytes)
		fillViaLocal := bulkMemoryGasUsed(t, gasSchedule, "memoryFillViaLocal", numBytes)
		require.Equal(t, fillDirect+localCost, fillViaLocal, "memory.fill, numBytes = %d", numBytes)
	}
}

// TestBulkMemoryGas_StillForbiddenOnOpcodeVersionV1 makes sure that metering the bulk memory
// opcodes did not accidentally allow them on the old opcode version, which is the one the
// host picks while the OpcodeV2 flag is off.
func TestBulkMemoryGas_StillForbiddenOnOpcodeVersionV1(t *testing.T) {
	testcommon.BuildInstanceCallTest(t).
		WithContracts(
			testcommon.CreateInstanceContract(testcommon.ParentAddress).
				WithCode(testcommon.GetTestSCCode("bulk-memory", "../../"))).
		WithGasSchedule(bulkMemoryGasSchedule()).
		WithEnableEpochsHandler(&worldmock.EnableEpochsHandlerStub{
			IsFlagEnabledCalled: func(flag core.EnableEpochFlag) bool {
				return flag != vmhost.OpcodeV2Flag
			},
		}).
		WithInput(testcommon.CreateTestContractCallInputBuilder().
			WithGasProvided(bulkMemoryGasProvided).
			WithFunction("memoryCopy").
			WithArguments(big.NewInt(10).Bytes()).
			Build()).
		AndAssertResults(func(_ vmhost.VMHost, _ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
			verify.ContractInvalid()
		})
}
