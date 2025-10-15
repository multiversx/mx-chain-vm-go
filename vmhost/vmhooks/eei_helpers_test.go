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

		data := "\n" + string(rune((len("ESDTRoleLocalMint")))) + "ESDTRoleLocalMint"

		roles := getESDTRoles([]byte(data), false)
		require.Equal(t, int64(RoleMint), roles)
	})

	t.Run("two roles", func(t *testing.T) {
		t.Parallel()

		data := "\n" + string(rune((len("ESDTRoleLocalMint")))) + "ESDTRoleLocalMint"
		data += "\n" + string(rune((len("ESDTRoleLocalBurn")))) + "ESDTRoleLocalBurn"

		roles := getESDTRoles([]byte(data), false)
		require.Equal(t, int64(RoleMint|RoleBurn), roles)
	})

	t.Run("two roles v2", func(t *testing.T) {
		t.Parallel()
		data := "\n" + string(rune((len("ESDTRoleLocalMint")))) + "ESDTRoleLocalMint"
		data += "\n" + string(rune((len("ESDTRoleNFTUpdateAttributes")))) + "ESDTRoleNFTUpdateAttributes"

		roles := getESDTRoles([]byte(data), true)
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
