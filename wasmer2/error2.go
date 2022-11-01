package wasmer2

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

	if errorResult == -1 {
		return "", errors.New("cannot read last error")
	}

	return cGoString(errorMessagePointer), nil
}

func newWrappedError(target error) error {
	var lastError string
	var err error
	lastError, err = GetLastError()

	if err != nil {
		lastError = "unknown details"
	}

	return fmt.Errorf("%w: %s", target, lastError)
}
