package elrondapigenerate

import "fmt"

func rustType(goType string) string {
	if goType == "int32" {
		return "i32"
	}
	if goType == "int64" {
		return "i64"
	}
	return goType
}

func adapter_function_name(name string) string {
	return fmt.Sprintf("wasmer_import_%s", snakeCase(name))
}
