package wasmer

import (
	"errors"
	"fmt"
	"unsafe"
)

var ErrFailedInstantiation = errors.New("could not create wasmer instance")

var ErrFailedCacheImports = errors.New("could not cache imports")

var ErrInvalidBytecode = errors.New("invalid bytecode")

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

	fmt.Printf("%x\n", errorMessage)

	if -1 == errorResult {
		return "", errors.New("Cannot read last error")
	}

	return cGoString(errorMessagePointer), nil
}
