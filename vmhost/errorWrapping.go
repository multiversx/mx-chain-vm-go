package vmhost

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const skipStackLevels = 2

// WrappableError - an interface that extends error and represents a multi-layer error
type WrappableError interface {
	error

	WrapWithMessage(errMessage string) WrappableError
	WrapWithStackTrace() WrappableError
	WrapWithError(err error, otherInfo ...string) WrappableError
	GetBaseError() error
	GetLastError() error
	GetAllErrors() []error
	GetAllErrorsAndOtherInfo() ([]error, []string)

	Unwrap() error
	Is(target error) bool
}

type errorWithLocation struct {
	err       error
	location  string
	otherInfo []string
}

type wrappableError struct {
	errsWithLocation []errorWithLocation
}

// WrapError constructs a WrappableError from an error
func WrapError(err error, otherInfo ...string) WrappableError {
	errAsWrappable, ok := err.(WrappableError)
	if ok {
		return errAsWrappable
	}
	return &wrappableError{
		errsWithLocation: []errorWithLocation{createErrorWithLocation(err, skipStackLevels, otherInfo...)},
	}
}

// WrapWithMessage wraps the target error with a new one, created using the input message
func (werr *wrappableError) WrapWithMessage(errMessage string) WrappableError {
	return werr.wrapWithErrorWithSkipLevels(errors.New(errMessage), skipStackLevels+1)
}

// WrapWithStackTrace wraps the target error with a new one, without any message only a stack frame trace
func (werr *wrappableError) WrapWithStackTrace() WrappableError {
	return werr.wrapWithErrorWithSkipLevels(errors.New(""), skipStackLevels+1)
}

// WrapWithError wraps the target error with the provided one
func (werr *wrappableError) WrapWithError(err error, otherInfo ...string) WrappableError {
	return werr.wrapWithErrorWithSkipLevels(err, skipStackLevels+1, otherInfo...)
}

// GetBaseError gets the core error
func (werr *wrappableError) GetBaseError() error {
	errs := werr.errsWithLocation
	return errs[0].err
}

// GetLastError gets the last wrapped error
func (werr *wrappableError) GetLastError() error {
	errs := werr.errsWithLocation
	return errs[len(errs)-1].err
}

// GetAllErrors gets all the wrapped errors
func (werr *wrappableError) GetAllErrors() []error {
	errs := werr.errsWithLocation
	allErrors := make([]error, 0)
	for _, err := range errs {
		allErrors = append(allErrors, err.err)
	}
	return allErrors
}

// GetAllErrorsAndOtherInfo gets all the wrapped errors + otherInfos
func (werr *wrappableError) GetAllErrorsAndOtherInfo() ([]error, []string) {
	errs := werr.errsWithLocation
	allErrors := make([]error, 0)
	allOtherInfo := make([]string, 0)
	for _, err := range errs {
		allErrors = append(allErrors, err.err)
		allOtherInfo = append(allOtherInfo, err.otherInfo...)
	}
	return allErrors, allOtherInfo
}

func (werr *wrappableError) wrapWithErrorWithSkipLevels(err error, skipStackLevels int, otherInfo ...string) *wrappableError {
	newErrs := make([]errorWithLocation, len(werr.errsWithLocation))
	copy(newErrs, werr.errsWithLocation)
	if err == nil {
		return &wrappableError{
			errsWithLocation: newErrs,
		}
	}

	var errsWithLocation []errorWithLocation
	inputWrappableError, ok := err.(*wrappableError)
	if !ok {
		errsWithLocation = append(newErrs, createErrorWithLocation(err, skipStackLevels, otherInfo...))
	} else {
		errsWithLocation = append(newErrs, inputWrappableError.errsWithLocation...)
	}
	return &wrappableError{
		errsWithLocation: errsWithLocation,
	}
}

func createErrorWithLocation(err error, skipStackLevels int, otherInfo ...string) errorWithLocation {
	_, file, line, _ := runtime.Caller(skipStackLevels)

	splitString := strings.Split(file, "/")
	fileName := splitString[len(splitString)-1]
	locationLine := fmt.Sprintf("%s:%d", fileName, line)

	errWithLocation := errorWithLocation{err: err, location: locationLine, otherInfo: otherInfo}
	return errWithLocation
}

// Error - standard error function implementation for wrappable errors
func (werr *wrappableError) Error() string {
	strErr := ""
	errs := werr.errsWithLocation
	for idxErr := range errs {
		errWithLocation := errs[len(errs)-1-idxErr]
		errMsg := errWithLocation.err.Error()
		suffix := ""
		if errMsg != "" {
			suffix = " [" + errMsg + "]"
		}
		if errWithLocation.otherInfo != nil {
			suffix += " [" + strings.Join(errWithLocation.otherInfo, ",") + "]"
		}
		strErr += "\n\t" + errWithLocation.location + suffix
	}
	return strErr
}

// Unwrap - standard error function implementation for wrappable errors
func (werr *wrappableError) Unwrap() error {
	wrappingErr := werr.unwrapWrapping()
	if len(wrappingErr.errsWithLocation) == 1 {
		return wrappingErr.errsWithLocation[0].err
	} else {
		return wrappingErr
	}
}

func (werr *wrappableError) unwrapWrapping() *wrappableError {
	if len(werr.errsWithLocation) == 0 {
		return nil
	}
	return &wrappableError{
		errsWithLocation: werr.errsWithLocation[:len(werr.errsWithLocation)-1],
	}
}

// Is - standard error function implementation for wrappable errors
func (werr *wrappableError) Is(target error) bool {
	for _, err := range werr.errsWithLocation {
		if errors.Is(err.err, target) {
			return true
		}
	}
	return false
}
