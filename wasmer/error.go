package wasmer

import (
	"errors"
	"unsafe"
)

// ErrFailedInstantiation indicates that a Wasmer instance could not be created
var ErrFailedInstantiation = errors.New("could not create wasmer instance")

// ErrFailedCacheImports indicates that the imports could not be cached
var ErrFailedCacheImports = errors.New("could not cache imports")

// ErrInvalidBytecode indicates that the bytecode is invalid
var ErrInvalidBytecode = errors.New("invalid bytecode")

// ErrCachingFailed indicates that creating the precompilation cache of an instance has failed
var ErrCachingFailed = errors.New("instance caching failed")

// GetLastError returns the last error message if any, otherwise returns an error.
func GetLastError() (string, error) {
	var errorLength = cWasmerLastErrorLength()

	if errorLength == 0 {
		return "", nil
	}

	var errorMessage = make([]cChar, errorLength)
	var errorMessagePointer = (*cChar)(unsafe.Pointer(&errorMessage[0]))

	var errorResult = cWasmerLastErrorMessage(errorMessagePointer, errorLength)

	if -1 == errorResult {
		return "", errors.New("cannot read last error")
	}

	return cGoString(errorMessagePointer), nil
}
