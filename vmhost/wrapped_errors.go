package vmhost

import "errors"

// WrappableError is an error type that can be wrapped with additional context.
type WrappableError interface {
	error
	Unwrap() error
	Wrap(err error) WrappableError
	WrapWithMessage(msg string) WrappableError
	WrapWithError(err error, otherInfo ...string) WrappableError
	GetAllErrors() []error
	GetAllErrorsAndOtherInfo() ([]error, [][]string)
}

type wrappableError struct {
	err      error
	info     []string
	innerErr error
}

// NewWrappableError creates a new WrappableError.
func NewWrappableError(err error, otherInfo ...string) WrappableError {
	return &wrappableError{
		err:  err,
		info: otherInfo,
	}
}

func (we *wrappableError) Error() string {
	return we.err.Error()
}

func (we *wrappableError) Unwrap() error {
	return we.innerErr
}

func (we *wrappableError) Wrap(err error) WrappableError {
	we.innerErr = err
	return we
}

func (we *wrappableError) WrapWithMessage(msg string) WrappableError {
	we.innerErr = errors.New(msg)
	return we
}

func (we *wrappableError) WrapWithError(err error, otherInfo ...string) WrappableError {
	newErr := &wrappableError{
		err:      err,
		info:     otherInfo,
		innerErr: we,
	}
	return newErr
}

func (we *wrappableError) GetAllErrors() []error {
	var errs []error
	var currErr error = we
	for currErr != nil {
		errs = append(errs, currErr)
		currErr = errors.Unwrap(currErr)
	}
	return errs
}

func (we *wrappableError) GetAllErrorsAndOtherInfo() ([]error, [][]string) {
	var errs []error
	var infos [][]string
	var currErr WrappableError = we
	for currErr != nil {
		errs = append(errs, currErr)
		if werr, ok := currErr.(*wrappableError); ok {
			infos = append(infos, werr.info)
		}
		unwrapped := currErr.Unwrap()
		var ok bool
		currErr, ok = unwrapped.(WrappableError)
		if !ok {
			if unwrapped != nil {
				errs = append(errs, unwrapped)
			}
			break
		}
	}
	return errs, infos
}
