package vmhooks

import (
	"encoding/binary"
	"math/big"
	"strconv"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/config"
	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/multiversx/mx-chain-vm-go/math"
	mockcontext "github.com/multiversx/mx-chain-vm-go/mock/context"
	twos "github.com/multiversx/mx-components-big-int/twos-complement"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInt32Array_ShortOrTruncatedDataShouldNotPanic(t *testing.T) {
	t.Parallel()

	// Empty
	res := createInt32Array([]byte{}, 0)
	assert.Equal(t, []int32{}, res)

	// 1 byte (truncated, less than 4 bytes)
	require.NotPanics(t, func() {
		res = createInt32Array([]byte{0x01}, 1)
		assert.Equal(t, []int32{0}, res)
	})

	// 3 bytes (truncated)
	require.NotPanics(t, func() {
		res = createInt32Array([]byte{0x01, 0x02, 0x03}, 1)
		assert.Equal(t, []int32{0}, res)
	})

	// 5 bytes with numIntegers = 1
	require.NotPanics(t, func() {
		res = createInt32Array([]byte{0x01, 0x00, 0x00, 0x00, 0x02}, 1)
		assert.Equal(t, []int32{1}, res)
	})

	// 8 bytes with numIntegers = 2
	require.NotPanics(t, func() {
		res = createInt32Array([]byte{0x05, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00}, 2)
		assert.Equal(t, []int32{5, 10}, res)
	})
}

func TestMakeCrossShardCallFromInput(t *testing.T) {
	t.Parallel()

	// No arguments
	res := makeCrossShardCallFromInput("testFunc", nil)
	assert.Equal(t, "testFunc", res)

	// Single argument
	res = makeCrossShardCallFromInput("testFunc", [][]byte{[]byte("hello")})
	assert.Equal(t, "testFunc@68656c6c6f", res)

	// Multiple arguments
	res = makeCrossShardCallFromInput("testFunc", [][]byte{
		[]byte("hello"),
		{0x00, 0x01, 0x02},
		{},
	})
	assert.Equal(t, "testFunc@68656c6c6f@000102@", res)
}

// the two byte length hooks must report exactly what the corresponding
// bytes hooks write out, without building those bytes
func TestBigIntByteLength_MatchesBytesHooks(t *testing.T) {
	t.Parallel()

	testCases := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		big.NewInt(127),
		big.NewInt(128),
		big.NewInt(255),
		big.NewInt(256),
		big.NewInt(65535),
		big.NewInt(65536),
		big.NewInt(16777215),
		big.NewInt(16777216),
		big.NewInt(-1),
		big.NewInt(-127),
		big.NewInt(-128),
		big.NewInt(-129),
		big.NewInt(-255),
		big.NewInt(-256),
		big.NewInt(-65536),
		new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)), // 2^256 - 1
		new(big.Int).Lsh(big.NewInt(1), 256),                                  // 2^256
		new(big.Int).Neg(new(big.Int).Lsh(big.NewInt(1), 256)),                // -2^256
		new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 512), big.NewInt(1)), // 2^512 - 1
	}

	for _, tc := range testCases {
		// BigIntGetUnsignedBytes / BigIntFinishUnsigned write value.Bytes()
		require.Equal(t, len(tc.Bytes()), math.UnsignedByteLength(tc),
			"unsigned length mismatch for value %s", tc.String())

		// BigIntGetSignedBytes / BigIntFinishSigned write twos.ToBytes(value)
		require.Equal(t, len(twos.ToBytes(tc)), math.TwosComplementByteLength(tc),
			"signed length mismatch for value %s", tc.String())
	}
}

func TestGetESDTRoles_TruncatedOrCorruptedDoesNotPanic(t *testing.T) {
	t.Parallel()

	// Empty buffer
	assert.Equal(t, int64(0), getESDTRoles([]byte{}, false))

	// Truncated at \n
	require.NotPanics(t, func() {
		assert.Equal(t, int64(0), getESDTRoles([]byte("\n"), false))
	})

	// Truncated role length byte
	require.NotPanics(t, func() {
		assert.Equal(t, int64(0), getESDTRoles([]byte("\n\x05"), false))
	})

	// Role length claims 10 bytes, but only 3 provided
	require.NotPanics(t, func() {
		assert.Equal(t, int64(0), getESDTRoles([]byte("\n\x0aabc"), false))
	})

	// Valid role followed by truncated role
	validRole := "\n" + string(byte(len(esdtRoleLocalMint))) + esdtRoleLocalMint
	res := getESDTRoles([]byte(validRole), false)
	assert.Equal(t, int64(RoleMint), res)

	// Valid role followed by corrupt trailer
	resCorrupt := getESDTRoles([]byte(validRole+"\n\x0fcorrupt"), false)
	assert.Equal(t, int64(RoleMint), resCorrupt)
}

func TestBaseOps_GetArgumentsFromMemory_TruncationGated(t *testing.T) {
	t.Parallel()

	for _, flagEnabled := range []bool{false, true} {
		flagEnabled := flagEnabled
		t.Run("FixCorrectSmallIntAndGasLimits="+strconv.FormatBool(flagEnabled), func(t *testing.T) {
			t.Parallel()

			mem := &mockcontext.MemoryMock{
				Pages:    1,
				PageSize: 100,
				Contents: make([]byte, 100),
			}
			// Write length = 3 into the 4 bytes at offset 96
			binary.LittleEndian.PutUint32(mem.Contents[96:100], 3)
			// Write 3 bytes of argument data at offset 10
			copy(mem.Contents[10:13], []byte{0xaa, 0xbb, 0xcc})

			instanceMock := mockcontext.NewInstanceMock(nil)
			instanceMock.Memory = mem

			runtimeMock := &mockcontext.RuntimeContextMock{
				Instance: instanceMock,
				FixCorrectSmallIntAndGasLimitsFunc: func() bool {
					return flagEnabled
				},
			}

			meteringMock := &mockcontext.MeteringContextMock{
				GasLeftMock: 1000000,
				GasCost: &config.GasCost{
					BaseOperationCost: config.BaseOperationCost{
						DataCopyPerByte: 1,
					},
				},
			}

			hostMock := &mockcontext.VMHostMock{
				RuntimeContext:  runtimeMock,
				MeteringContext: meteringMock,
			}

			hooks := NewVMHooksImpl(hostMock)

			// Memory has 100 bytes (valid offsets 0..99).
			// If argumentsLengthOffset is 96, and numArguments is 2 (requires 2 * 4 = 8 bytes):
			// MemLoad(96, 8) will return only 4 bytes (96..99) because offset+8 (104) exceeds memory length (100).
			numArguments := int32(2)
			argumentsLengthOffset := executor.MemPtr(96)
			dataOffset := executor.MemPtr(10)

			args, totalLen, err := hooks.getArgumentsFromMemory(hostMock, numArguments, argumentsLengthOffset, dataOffset)

			if flagEnabled {
				require.Error(t, err)
				assert.ErrorIs(t, err, executor.ErrMemoryBadBounds)
				assert.Contains(t, err.Error(), "arguments length data truncated")
				assert.Nil(t, args)
				assert.Equal(t, int32(0), totalLen)
			} else {
				require.NoError(t, err)
				assert.Equal(t, int32(3), totalLen)
				assert.Equal(t, [][]byte{{0xaa, 0xbb, 0xcc}, {}}, args)
			}
		})
	}
}
