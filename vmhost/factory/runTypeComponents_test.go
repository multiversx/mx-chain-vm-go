package factory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRunTypeComponentsFactory(t *testing.T) {
	t.Parallel()

	rcf := NewRunTypeComponentsFactory()
	rc := rcf.Create()
	require.NotNil(t, rc)

	require.NoError(t, rc.Close())
}
