package arwen

import (
	"errors"
	"fmt"
	"runtime"
)

const skipLevels = 3

// WrappableError - TODO
type WrappableError interface {
	error

	WrapWithMessage(errMessage string) WrappableError
	WrapWithStackTrace() WrappableError
	WrapWithError(err error) WrappableError
	GetBaseError() error
	GetLastError() error

	Unwrap() error
	Is(target error) bool
}

type errorWithLocation struct {
	err      error
	location string
}

type wrappableError struct {
	errsWithLocation []errorWithLocation
}

func (werr *wrappableError) getTopErrorWithLocation() errorWithLocation {
	return werr.errsWithLocation[len(werr.errsWithLocation)-1]
}

func WrapError(err error) WrappableError {
	errAsWrappable, ok := err.(WrappableError)
	if ok {
		return errAsWrappable
	}
	return &wrappableError{
		errsWithLocation: []errorWithLocation{createErrorWithLocation(err, 2)},
	}
}

func (werr *wrappableError) WrapWithMessage(errMessage string) WrappableError {
	return werr.wrapWithErrorWithSkipLevels(errors.New(errMessage), skipLevels)
}

func (werr *wrappableError) WrapWithStackTrace() WrappableError {
	return werr.wrapWithErrorWithSkipLevels(errors.New(""), skipLevels)
}

func (werr *wrappableError) WrapWithError(err error) WrappableError {
	return werr.wrapWithErrorWithSkipLevels(err, skipLevels)
}

func (werr *wrappableError) GetBaseError() error {
	errors := werr.errsWithLocation
	return errors[0].err
}

func (werr *wrappableError) GetLastError() error {
	errors := werr.errsWithLocation
	return errors[len(errors)-1].err
}

func (werr *wrappableError) wrapWithErrorWithSkipLevels(err error, skipStackLevels int) *wrappableError {
	newErrs := make([]errorWithLocation, len(werr.errsWithLocation))
	copy(newErrs, werr.errsWithLocation)
	if err == nil {
		return &wrappableError{
			errsWithLocation: newErrs,
		}
	}
	return &wrappableError{
		errsWithLocation: append(newErrs, createErrorWithLocation(err, skipStackLevels)),
	}
}

func createErrorWithLocation(err error, skipStackLevels int) errorWithLocation {
	_, file, line, _ := runtime.Caller(skipStackLevels)
	locationLine := fmt.Sprintf("%s:%d", file, line)
	errWithLocation := errorWithLocation{err: err, location: locationLine}
	return errWithLocation
}

func (werr *wrappableError) Error() string {
	strErr := ""
	errors := werr.errsWithLocation
	for idxErr := range errors {
		errWithLocation := errors[len(errors)-1-idxErr]
		errMsg := errWithLocation.err.Error()
		suffix := ""
		if errMsg != "" {
			suffix = " [" + errMsg + "]"
		}
		strErr += "\n\t" + errWithLocation.location + suffix
	}
	return strErr
}

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

func (werr *wrappableError) Is(target error) bool {
	for _, err := range werr.errsWithLocation {
		if errors.Is(err.err, target) {
			return true
		}
	}
	return false
}
