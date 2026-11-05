package wasmer2

import (
	"reflect"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	"github.com/stretchr/testify/require"
)

// setAllFields sets every field of the struct pointed to by ptr to a distinct non-zero value.
func setAllFields(t *testing.T, ptr interface{}) {
	value := reflect.ValueOf(ptr).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.Uint32, reflect.Uint64:
			field.SetUint(uint64(i) + 1)
		case reflect.Bool:
			field.SetBool(true)
		default:
			t.Fatalf("unsupported field kind %s for %s", field.Kind(), value.Type().Field(i).Name)
		}
	}
}

// namedFields returns the fields of a struct type, skipping the blank padding
// fields that cgo may add.
func namedFields(structType reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		if structType.Field(i).Name != "_" {
			fields = append(fields, structType.Field(i))
		}
	}
	return fields
}

// requireAllFieldsSet checks that no (non-padding) field of the struct is left at its zero value.
func requireAllFieldsSet(t *testing.T, structValue reflect.Value, what string) {
	for _, field := range namedFields(structValue.Type()) {
		require.False(t, structValue.FieldByName(field.Name).IsZero(), "%s: field %s not set", what, field.Name)
	}
}

// Guards toCWasmerCompilationOptions against a field being added to either
// executor.CompilationOptions or vm_exec_compilation_options_t without the
// conversion being updated: every field is set to a non-zero value, and the
// converted struct must have the same number of fields, all non-zero.
func TestToCWasmerCompilationOptions_CoversAllFields(t *testing.T) {
	goOptions := executor.CompilationOptions{}
	setAllFields(t, &goOptions)

	cValue := reflect.ValueOf(toCWasmerCompilationOptions(goOptions))
	requireAllFieldsSet(t, cValue, "compilation options")
	require.Equal(t, reflect.TypeOf(goOptions).NumField(), len(namedFields(cValue.Type())),
		"field count mismatch between Go and C compilation options")
}

// OpcodeCost is reinterpreted as vm_exec_opcode_cost_t via unsafe.Pointer, so
// the two must have identical layout: same fields, in the same order, at the same offsets.
func TestOpcodeCost_LayoutMatchesC(t *testing.T) {
	goType := reflect.TypeOf(OpcodeCost{})
	cType := reflect.TypeOf(cWasmerOpcodeCostT{})
	cFields := namedFields(cType)

	require.Equal(t, goType.Size(), cType.Size(), "size mismatch between Go and C opcode costs")
	require.Equal(t, goType.NumField(), len(cFields), "field count mismatch between Go and C opcode costs")
	for i, cField := range cFields {
		goField := goType.Field(i)
		cName := strings.TrimPrefix(cField.Name, "opcode_")
		require.Equal(t, strings.ToLower(goField.Name), strings.ReplaceAll(cName, "_", ""),
			"field %d name mismatch", i)
		require.Equal(t, goField.Type.Kind(), cField.Type.Kind(), "field %s kind mismatch", goField.Name)
		require.Equal(t, goField.Offset, cField.Offset, "field %s offset mismatch", goField.Name)
	}
}

// Every OpcodeCost field must be populated from executor.WASMOpcodeCost.
func TestExtractOpcodeCost_CoversAllFields(t *testing.T) {
	wasmOps := executor.WASMOpcodeCost{}
	setAllFields(t, &wasmOps)

	opcodeCost := (&Wasmer2Executor{}).extractOpcodeCost(&wasmOps)
	requireAllFieldsSet(t, reflect.ValueOf(*opcodeCost), "opcode cost")
}

// Every hook in vm_exec_vm_hook_c_func_pointers must be given a function pointer.
func TestPopulateCgoFunctionPointers_CoversAllFields(t *testing.T) {
	pointers := reflect.ValueOf(*populateCgoFunctionPointers())
	for _, field := range namedFields(pointers.Type()) {
		require.False(t, pointers.FieldByName(field.Name).IsNil(), "hook %s not populated", field.Name)
	}
}
