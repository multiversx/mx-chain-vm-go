package executorinterface

import "fmt"

// ImportFunctionError represents any kind of errors related to a
// WebAssembly imported function. It is returned by `Import` or `Imports`
// functions only.
type ImportFunctionError struct {
	functionName string
	message      string
}

// NewImportFunctionError constructs a new `ImportedFunctionError`,
// where `functionName` is the name of the imported function, and
// `message` is the error message. If the error message contains `%s`,
// then this parameter will be replaced by `functionName`.
func NewImportFunctionError(functionName string, message string) *ImportFunctionError {
	return &ImportFunctionError{functionName, message}
}

// ImportedFunctionError is an actual error. The `Error` function
// returns the error message.
func (error *ImportFunctionError) Error() string {
	return fmt.Sprintf(error.message, error.functionName)
}
