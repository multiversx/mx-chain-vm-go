package contexts

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-common-go/builtInFunctions"
	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/multiversx/mx-chain-vm-go/vmhost/mock"
	"github.com/stretchr/testify/require"
)

func testImportNames() vmcommon.FunctionNames {
	importNames := make(vmcommon.FunctionNames)
	var empty struct{}
	importNames["getArgument"] = empty
	importNames["asyncCall"] = empty
	return importNames
}

func TestFunctionsGuard_isValidFunctionName(t *testing.T) {
	builtInFuncContainer := builtInFunctions.NewBuiltInFunctionContainer()
	_ = builtInFuncContainer.Add("protocolFunctionFoo", &mock.BuiltInFunctionStub{})
	_ = builtInFuncContainer.Add("protocolFunctionBar", &mock.BuiltInFunctionStub{})

	validator := newWASMValidator(testImportNames(), builtInFuncContainer, worldmock.EnableEpochsHandlerStubAllFlags())

	require.Nil(t, validator.verifyValidFunctionName("foo"))
	require.Nil(t, validator.verifyValidFunctionName("_"))
	require.Nil(t, validator.verifyValidFunctionName("a"))
	require.Nil(t, validator.verifyValidFunctionName("i"))

	require.NotNil(t, validator.verifyValidFunctionName(""))
	require.NotNil(t, validator.verifyValidFunctionName("3"))
	require.NotNil(t, validator.verifyValidFunctionName("π"))
	require.NotNil(t, validator.verifyValidFunctionName("2foo"))
	require.NotNil(t, validator.verifyValidFunctionName("-"))
	require.NotNil(t, validator.verifyValidFunctionName("â"))
	require.NotNil(t, validator.verifyValidFunctionName("ș"))
	require.NotNil(t, validator.verifyValidFunctionName("Ä"))

	require.NotNil(t, validator.verifyValidFunctionName("protocolFunctionFoo"))
	require.NotNil(t, validator.verifyValidFunctionName("protocolFunctionBar"))

	require.Nil(t, validator.verifyValidFunctionName(strings.Repeat("_", 255)))
	require.NotNil(t, validator.verifyValidFunctionName(strings.Repeat("_", 256)))

	require.NotNil(t, validator.verifyValidFunctionName("getArgument"))
	require.NotNil(t, validator.verifyValidFunctionName("asyncCall"))
	require.Nil(t, validator.verifyValidFunctionName("getArgument55"))
}

func TestFunctionsProtected(t *testing.T) {
	host := InitializeVMAndWasmer()

	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())

	world := worldmock.NewMockWorld()
	imb := contextmock.NewExecutorMock(world)
	instance := imb.CreateAndStoreInstanceMock(t, host, []byte{}, []byte{}, []byte{}, []byte{}, 0, 0, false)

	instance.AddMockMethod("transferValueOnly", func() *contextmock.InstanceMock {
		testHost := instance.Host
		testInstance := contextmock.GetMockInstance(testHost)
		return testInstance
	})

	err := validator.verifyProtectedFunctions(instance)
	require.NotNil(t, err)
}

type instanceWithCustomFunctionNames struct {
	*contextmock.InstanceMock
	customFunctionNames []string
}

func (inst *instanceWithCustomFunctionNames) GetFunctionNames() []string {
	return inst.customFunctionNames
}

func TestValidator_VerifyFunctions_NilGetFunctionNames_Fails(t *testing.T) {
	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())

	mockInst := &instanceWithCustomFunctionNames{
		InstanceMock:        contextmock.NewInstanceMock(nil),
		customFunctionNames: nil,
	}

	err := validator.verifyFunctions(mockInst)
	require.ErrorIs(t, err, vmhost.ErrContractInvalid)

	errProtected := validator.verifyProtectedFunctions(mockInst)
	require.ErrorIs(t, errProtected, vmhost.ErrContractInvalid)
}

func TestValidator_VerifyFunctions_NulByteInExport_Fails(t *testing.T) {
	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())

	mockInst := &instanceWithCustomFunctionNames{
		InstanceMock:        contextmock.NewInstanceMock(nil),
		customFunctionNames: []string{"a\x00", "transferValueOnly"},
	}

	err := validator.verifyFunctions(mockInst)
	require.ErrorIs(t, err, vmhost.ErrInvalidFunctionName)

	errProtected := validator.verifyProtectedFunctions(mockInst)
	require.ErrorIs(t, errProtected, vmhost.ErrContractInvalid)
}

func encodeUleb128(val uint32) []byte {
	var res []byte
	for {
		b := byte(val & 0x7f)
		val >>= 7
		if val != 0 {
			b |= 0x80
		}
		res = append(res, b)
		if val == 0 {
			break
		}
	}
	return res
}

func buildWasmModule(tableSec []byte, importSec []byte, codeSec []byte) []byte {
	header := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	var mod []byte
	mod = append(mod, header...)

	// Type section (type 0: func () -> ())
	typeSecPayload := []byte{0x01, 0x60, 0x00, 0x00}
	mod = append(mod, 0x01)
	mod = append(mod, encodeUleb128(uint32(len(typeSecPayload)))...)
	mod = append(mod, typeSecPayload...)

	if len(importSec) > 0 {
		mod = append(mod, 0x02)
		mod = append(mod, encodeUleb128(uint32(len(importSec)))...)
		mod = append(mod, importSec...)
	}

	// Function section: 1 func using type 0
	funcSecPayload := []byte{0x01, 0x00}
	mod = append(mod, 0x03)
	mod = append(mod, encodeUleb128(uint32(len(funcSecPayload)))...)
	mod = append(mod, funcSecPayload...)

	if len(tableSec) > 0 {
		mod = append(mod, 0x04)
		mod = append(mod, encodeUleb128(uint32(len(tableSec)))...)
		mod = append(mod, tableSec...)
	}

	if len(codeSec) > 0 {
		mod = append(mod, 0x0A)
		mod = append(mod, encodeUleb128(uint32(len(codeSec)))...)
		mod = append(mod, codeSec...)
	}

	return mod
}

func buildCodeSection(instructions []byte) []byte {
	var body []byte
	body = append(body, 0x00) // 0 locals
	body = append(body, instructions...)
	body = append(body, 0x0B) // end

	var codeSec []byte
	codeSec = append(codeSec, 0x01) // 1 function
	codeSec = append(codeSec, encodeUleb128(uint32(len(body)))...)
	codeSec = append(codeSec, body...)
	return codeSec
}

func buildTableSection(hasMax bool, min, max uint32) []byte {
	var payload []byte
	payload = append(payload, 0x01) // 1 table
	payload = append(payload, 0x70) // funcref
	if hasMax {
		payload = append(payload, 0x01)
		payload = append(payload, encodeUleb128(min)...)
		payload = append(payload, encodeUleb128(max)...)
	} else {
		payload = append(payload, 0x00)
		payload = append(payload, encodeUleb128(min)...)
	}
	return payload
}

func buildTableImportSection() []byte {
	var payload []byte
	payload = append(payload, 0x01) // 1 import
	modName := "env"
	payload = append(payload, encodeUleb128(uint32(len(modName)))...)
	payload = append(payload, []byte(modName)...)
	fieldName := "tbl"
	payload = append(payload, encodeUleb128(uint32(len(fieldName)))...)
	payload = append(payload, []byte(fieldName)...)
	payload = append(payload, 0x01) // kind = Table
	payload = append(payload, 0x70) // funcref
	payload = append(payload, 0x01) // flags
	payload = append(payload, encodeUleb128(1)...)
	payload = append(payload, encodeUleb128(10)...)
	return payload
}

func TestValidator_VerifyContractBytecode(t *testing.T) {
	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())

	t.Run("valid module with table", func(t *testing.T) {
		tbl := buildTableSection(true, 1, 10)
		code := buildCodeSection([]byte{0x01}) // nop
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.NoError(t, err)
	})

	t.Run("valid module with typed select opcode 0x1C", func(t *testing.T) {
		code := buildCodeSection([]byte{0x1C, 0x01, 0x7F})
		wasm := buildWasmModule(nil, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.NoError(t, err)
	})

	t.Run("typed select opcode 0x1C with truncated operand is rejected", func(t *testing.T) {
		code := buildCodeSection([]byte{0x1C, 0x05, 0x7F})
		wasm := buildWasmModule(nil, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("table without max declared is rejected", func(t *testing.T) {
		tbl := buildTableSection(false, 1, 0)
		code := buildCodeSection([]byte{0x01})
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("table with max exceeding limit is rejected", func(t *testing.T) {
		tbl := buildTableSection(true, 1, 10001)
		code := buildCodeSection([]byte{0x01})
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("table with min greater than max is rejected", func(t *testing.T) {
		tbl := buildTableSection(true, 20, 10)
		code := buildCodeSection([]byte{0x01})
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("imported table is rejected", func(t *testing.T) {
		imp := buildTableImportSection()
		code := buildCodeSection([]byte{0x01})
		wasm := buildWasmModule(nil, imp, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("table.grow opcode is rejected", func(t *testing.T) {
		tbl := buildTableSection(true, 1, 10)
		// i32.const 0, table.grow 0
		code := buildCodeSection([]byte{0x41, 0x00, 0xFC, 0x0F, 0x00})
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("table.set opcode is rejected", func(t *testing.T) {
		tbl := buildTableSection(true, 1, 10)
		// i32.const 0, ref.null 0x70, table.set 0
		code := buildCodeSection([]byte{0x41, 0x00, 0xD0, 0x70, 0x26, 0x00})
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})

	t.Run("i32.const 2044 does not trigger false positive on table.grow", func(t *testing.T) {
		tbl := buildTableSection(true, 1, 10)
		// i32.const 2044 (encoded in sleb128 as 0xFC 0x0F), drop
		code := buildCodeSection([]byte{0x41, 0xFC, 0x0F, 0x1A})
		wasm := buildWasmModule(tbl, nil, code)
		err := validator.verifyContractBytecode(wasm)
		require.NoError(t, err)
	})

	t.Run("invalid header is rejected", func(t *testing.T) {
		err := validator.verifyContractBytecode([]byte("not wasm"))
		require.ErrorIs(t, err, vmhost.ErrContractInvalid)
	})
}

func BenchmarkLockUnlockOSThread(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		runtime.LockOSThread()
		runtime.UnlockOSThread()
	}
}

func BenchmarkVerifyContractBytecode_Pair(b *testing.B) {
	code, err := os.ReadFile("../../test/contracts/pair/output/pair.wasm")
	if err != nil {
		b.Fatal(err)
	}
	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())
	b.SetBytes(int64(len(code)))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = validator.verifyContractBytecode(code)
	}
}

func BenchmarkVerifyContractBytecode_DexRouter(b *testing.B) {
	code, err := os.ReadFile("../../fuzz/dex/v0.1/wasms/dex_router.wasm")
	if err != nil {
		b.Fatal(err)
	}
	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())
	b.SetBytes(int64(len(code)))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = validator.verifyContractBytecode(code)
	}
}

func BenchmarkVerifyContractBytecode_Delegation(b *testing.B) {
	code, err := os.ReadFile("../../test/delegation/v0_3/output/delegation.wasm")
	if err != nil {
		b.Fatal(err)
	}
	validator := newWASMValidator(testImportNames(), builtInFunctions.NewBuiltInFunctionContainer(), worldmock.EnableEpochsHandlerStubAllFlags())
	b.SetBytes(int64(len(code)))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = validator.verifyContractBytecode(code)
	}
}
