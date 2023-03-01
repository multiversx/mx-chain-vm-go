package vmhooksgenerate

import (
	"fmt"
	"strings"
)

func rustVMHooksType(eiType EIType) string {
	switch eiType {
	case EITypeMemPtr:
		return "MemPtr"
	case EITypeMemLength:
		return "MemLength"
	case EITypeInt32:
		return "i32"
	case EITypeInt64:
		return "i64"
	default:
		panic("invalid type")
	}
}

func rustCapiType(eiType EIType) string {
	switch eiType {
	case EITypeMemPtr:
		fallthrough
	case EITypeMemLength:
		fallthrough
	case EITypeInt32:
		return "i32"
	case EITypeInt64:
		return "i64"
	default:
		panic("invalid type")
	}
}

func rustWasmerType(eiType EIType) string {
	return rustCapiType(eiType)
}

func rustWasmerConvertArg(arg *EIFunctionArg) string {
	argRustName := snakeCase(arg.Name)
	switch arg.Type {
	case EITypeMemPtr:
		return fmt.Sprintf("env.convert_mem_ptr(%s)", argRustName)
	case EITypeMemLength:
		return fmt.Sprintf("env.convert_mem_length(%s)", argRustName)
	default:
		return argRustName
	}
}

func rustCapiConvertArg(arg *EIFunctionArg) string {
	argRustName := snakeCase(arg.Name)
	switch arg.Type {
	case EITypeMemPtr:
		return fmt.Sprintf("self.convert_mem_ptr(%s)", argRustName)
	case EITypeMemLength:
		return fmt.Sprintf("self.convert_mem_length(%s)", argRustName)
	default:
		return argRustName
	}
}

func wasmerImportAdapterFunctionName(name string) string {
	return fmt.Sprintf("wasmer_import_%s", snakeCase(name))
}

func cgoFuncPointerFieldName(funcMetadata *EIFunction) string {
	return snakeCase(funcMetadata.Name) + "_func_ptr"
}

func writeRustFnDeclarationArguments(firstArgs string, funcMetadata *EIFunction, rustType func(EIType) string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("(%s", firstArgs))
	for _, arg := range funcMetadata.Arguments {
		sb.WriteString(fmt.Sprintf(", %s: %s", snakeCase(arg.Name), rustType(arg.Type)))
	}
	sb.WriteString(")")
	if funcMetadata.Result != nil {
		sb.WriteString(fmt.Sprintf(" -> %s", rustType(funcMetadata.Result.Type)))
	}
	return sb.String()
}
