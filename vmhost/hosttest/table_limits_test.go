package hostCoretest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	contextmock "github.com/multiversx/mx-chain-vm-go/mock/context"
	"github.com/multiversx/mx-chain-vm-go/testcommon"
	"github.com/stretchr/testify/require"
)

// must match maxAllowedTableElements in vmhost/contexts/validator.go
const maxAllowedTableElements = 10000

const tableLimitWatTemplate = `(module
  (type $void (func))
  (func $init (type $void))
  (table $tab %d %d funcref)
  (elem (i32.const 0) $init)
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "init" (func $init))
)
`

const noTableWat = `(module
  (type $void (func))
  (func $init (type $void))
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "init" (func $init))
)
`

func compileWat(t *testing.T, wat string) []byte {
	wat2wasm, err := exec.LookPath("wat2wasm")
	if err != nil {
		t.Skip("wat2wasm not found in PATH")
	}

	dir := t.TempDir()
	watPath := filepath.Join(dir, "contract.wat")
	wasmPath := filepath.Join(dir, "contract.wasm")
	require.NoError(t, os.WriteFile(watPath, []byte(wat), 0644))

	output, err := exec.Command(wat2wasm, watPath, "-o", wasmPath).CombinedOutput()
	require.NoError(t, err, string(output))

	wasm, err := os.ReadFile(wasmPath)
	require.NoError(t, err)
	return wasm
}

func deployWatContract(t *testing.T, wat string, assertResults func(*contextmock.BlockchainHookStub, *testcommon.VMOutputVerifier)) {
	code := compileWat(t, wat)

	testcommon.BuildInstanceCreatorTest(t).
		WithInput(testcommon.CreateTestContractCreateInputBuilder().
			WithGasProvided(1000000).
			WithContractCode(code).
			Build()).
		WithAddress(newAddress).
		AndAssertResults(assertResults)
}

func TestTableLimits_MaxAllowed(t *testing.T) {
	deployWatContract(t, fmt.Sprintf(tableLimitWatTemplate, 1, maxAllowedTableElements), func(_ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
		verify.Ok()
	})
}

func TestTableLimits_MaxAllowedPlusOne(t *testing.T) {
	deployWatContract(t, fmt.Sprintf(tableLimitWatTemplate, 1, maxAllowedTableElements+1), func(_ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
		verify.ContractInvalid()
	})
}

func TestTableLimits_NoTable(t *testing.T) {
	deployWatContract(t, noTableWat, func(_ *contextmock.BlockchainHookStub, verify *testcommon.VMOutputVerifier) {
		verify.Ok()
	})
}
