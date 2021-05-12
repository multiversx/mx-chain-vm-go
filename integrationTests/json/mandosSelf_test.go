package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests Mandos consistency, no smart contracts.
func TestMandosSelfTest(t *testing.T) {
	runTestsInFolder(t, "mandos-self-test", []string{
		"mandos-self-test/builtin-func-esdt-transfer.scen.json",
	})
}

func TestMandosCheckNonceErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-nonce.err.json")
	require.EqualError(t, err,
		"bad account nonce. Account: address:the-address. Want: \"1002\". Have: \"1001\"")
}

func TestMandosCheckBalanceErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-balance.err.json")
	require.EqualError(t, err,
		"bad account balance. Account: address:the-address. Want: \"1,000,002\". Have: \"1000001\"")
}

func TestMandosCheckUsernameErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-username.err.json")
	require.EqualError(t, err,
		"bad account username. Account: address:the-address. Want: \"str:wrong.elrond\". Have: \"str:theusername.elrond\"")
}

func TestMandosCheckCodeErr(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-code.err.json")
	require.EqualError(t, err,
		"bad account code. Account: address:the-address. Want: \"file:set-check-code.scen.json\". Have: \"0x7b0a2020202022636f6d...\"")
}

func TestMandosCheckStorageErr1(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err1.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"str:another-value\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestMandosCheckStorageErr2(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err2.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestMandosCheckStorageErr3(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err3.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d64 (str:key-d): Want: \"str:value-d\". Have: \"\"")
}

func TestMandosCheckStorageErr4(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err4.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestMandosCheckStorageErr5(t *testing.T) {
	err := runSingleTest(t, "mandos-self-test/set-check", "set-check-storage.err5.json")
	require.EqualError(t, err,
		"wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d62 (str:key-b): Want: \"str:another-b\". Have: \"0x76616c75652d62 (str:value-b)\"")
}
