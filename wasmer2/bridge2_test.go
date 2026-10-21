package wasmer2

import (
	"reflect"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

// Guards toCWasmerCompilationOptions against a field being added to either
// executor.CompilationOptions or vm_exec_compilation_options_t without the
// conversion being updated: every field is set to a non-zero value, and the
// converted struct must have the same number of fields, all non-zero.
func TestToCWasmerCompilationOptions_CoversAllFields(t *testing.T) {
	goOptions := executor.CompilationOptions{}
	goValue := reflect.ValueOf(&goOptions).Elem()
	for i := 0; i < goValue.NumField(); i++ {
		field := goValue.Field(i)
		switch field.Kind() {
		case reflect.Uint64:
			field.SetUint(uint64(i) + 1)
		case reflect.Bool:
			field.SetBool(true)
		default:
			t.Fatalf("unsupported field kind %s for %s", field.Kind(), goValue.Type().Field(i).Name)
		}
	}

	// cgo may add blank padding fields, which the conversion cannot set
	cValue := reflect.ValueOf(toCWasmerCompilationOptions(goOptions))
	numCFields := 0
	for i := 0; i < cValue.NumField(); i++ {
		name := cValue.Type().Field(i).Name
		if name == "_" {
			continue
		}
		numCFields++
		require.False(t, cValue.Field(i).IsZero(), "C field %s not set by conversion", name)
	}
	require.Equal(t, goValue.NumField(), numCFields, "field count mismatch between Go and C compilation options")
}
