package vmhooks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetESDTRoles(t *testing.T) {
	t.Parallel()

	t.Run("no roles", func(t *testing.T) {
		t.Parallel()
		roles := getESDTRoles([]byte{}, false)
		require.Equal(t, int64(0), roles)
	})

	t.Run("one role", func(t *testing.T) {
		t.Parallel()
		data := []byte("\n\x0eESDTRoleLocalMint")
		roles := getESDTRoles(data, false)
		require.Equal(t, int64(RoleMint), roles)
	})

	t.Run("two roles", func(t *testing.T) {
		t.Parallel()
		data := []byte("\n\x0eESDTRoleLocalMint\n\x0fESDTRoleLocalBurn")
		roles := getESDTRoles(data, false)
		require.Equal(t, int64(RoleMint|RoleBurn), roles)
	})

	t.Run("two roles v2", func(t *testing.T) {
		t.Parallel()
		data := []byte("\n\x0eESDTRoleLocalMint\n\x15ESDTRoleNFTUpdateAttributes")
		roles := getESDTRoles(data, true)
		require.Equal(t, int64(RoleMint|RoleNFTUpdateAttributes), roles)
	})
}

func TestValidateToken(t *testing.T) {
	t.Parallel()

	require.True(t, ValidateToken([]byte("TEST-123456")))
	require.False(t, ValidateToken([]byte("TEST-12345")))
	require.False(t, ValidateToken([]byte("TEST-1234567")))
	require.False(t, ValidateToken([]byte("test-123456")))
	require.False(t, ValidateToken([]byte("TEST-12345G")))
}
