package vmhost

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapError(t *testing.T) {
	t.Parallel()

	err := errors.New("base error")
	wrappedErr := WrapError(err, "info1")

	require.NotNil(t, wrappedErr)
	require.Equal(t, err, wrappedErr.GetBaseError())
	require.Equal(t, err, wrappedErr.GetLastError())

	allErrs, allInfo := wrappedErr.GetAllErrorsAndOtherInfo()
	require.Len(t, allErrs, 1)
	require.Equal(t, err, allErrs[0])
	require.Len(t, allInfo, 1)
	require.Equal(t, "info1", allInfo[0])
}

func TestWrappableError_WrapWithMessage(t *testing.T) {
	t.Parallel()

	err := errors.New("base error")
	wrappedErr := WrapError(err)
	wrappedErr = wrappedErr.WrapWithMessage("second layer")

	require.Equal(t, err, wrappedErr.GetBaseError())
	require.NotEqual(t, err, wrappedErr.GetLastError())
	require.Equal(t, "second layer", wrappedErr.GetLastError().Error())

	allErrs := wrappedErr.GetAllErrors()
	require.Len(t, allErrs, 2)
	require.Equal(t, "second layer", allErrs[1].Error())
}

func TestWrappableError_WrapWithStackTrace(t *testing.T) {
	t.Parallel()

	err := errors.New("base error")
	wrappedErr := WrapError(err)
	wrappedErr = wrappedErr.WrapWithStackTrace()

	require.Equal(t, err, wrappedErr.GetBaseError())
	require.NotEqual(t, err, wrappedErr.GetLastError())
	require.Equal(t, "", wrappedErr.GetLastError().Error())

	allErrs := wrappedErr.GetAllErrors()
	require.Len(t, allErrs, 2)
}

func TestWrappableError_WrapWithError(t *testing.T) {
	t.Parallel()

	err1 := errors.New("base error")
	err2 := errors.New("second error")
	wrappedErr := WrapError(err1)
	wrappedErr = wrappedErr.WrapWithError(err2, "info2")

	require.Equal(t, err1, wrappedErr.GetBaseError())
	require.Equal(t, err2, wrappedErr.GetLastError())

	allErrs, allInfo := wrappedErr.GetAllErrorsAndOtherInfo()
	require.Len(t, allErrs, 2)
	require.Len(t, allInfo, 1)
	require.Equal(t, "info2", allInfo[0])
}

func TestWrappableError_Error(t *testing.T) {
	t.Parallel()

	err1 := errors.New("base error")
	wrappedErr := WrapError(err1)
	wrappedErr = wrappedErr.WrapWithMessage("L2")
	wrappedErr = wrappedErr.WrapWithError(errors.New("L3"), "info3")

	errStr := wrappedErr.Error()
	require.True(t, strings.Contains(errStr, "errorWrapping_test.go"))
	require.True(t, strings.Contains(errStr, "[base error]"))
	require.True(t, strings.Contains(errStr, "[L2]"))
	require.True(t, strings.Contains(errStr, "[L3]"))
	require.True(t, strings.Contains(errStr, "[info3]"))
}

func TestWrappableError_Unwrap(t *testing.T) {
	t.Parallel()

	err1 := errors.New("base error")
	err2 := errors.New("L2")
	err3 := errors.New("L3")
	wrappedErr := WrapError(err1).WrapWithError(err2).WrapWithError(err3)

	unwrappedOnce := errors.Unwrap(wrappedErr)
	require.NotNil(t, unwrappedOnce)
	require.True(t, errors.Is(unwrappedOnce, err1))
	require.True(t, errors.Is(unwrappedOnce, err2))
	require.False(t, errors.Is(unwrappedOnce, err3))

	unwrappedTwice := errors.Unwrap(unwrappedOnce)
	require.NotNil(t, unwrappedTwice)
	require.True(t, errors.Is(unwrappedTwice, err1))
	require.False(t, errors.Is(unwrappedTwice, err2))

	unwrappedThrice := errors.Unwrap(unwrappedTwice)
	require.Equal(t, err1, unwrappedThrice)

	unwrappedFour := errors.Unwrap(unwrappedThrice)
	require.Nil(t, unwrappedFour)
}

func TestWrappableError_Is(t *testing.T) {
	t.Parallel()

	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")

	wrappedErr := WrapError(err1)
	wrappedErr = wrappedErr.WrapWithError(err2)

	require.True(t, wrappedErr.Is(err1))
	require.True(t, wrappedErr.Is(err2))
	require.False(t, wrappedErr.Is(err3))

	require.True(t, errors.Is(wrappedErr, err1))
	require.True(t, errors.Is(wrappedErr, err2))
	require.False(t, errors.Is(wrappedErr, err3))
}
