package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests Scenarios consistency, no smart contracts.
func TestScenariosSelfTest(t *testing.T) {
	runTestsInFolder(t, "scenarios-self-test", []string{
		"scenarios-self-test/builtin-func-esdt-transfer.scen.json",
		"scenarios-self-test/esdt-zero-balance-check-err.scen.json",
		"scenarios-self-test/esdt-non-zero-balance-check-err.scen.json",
	})
}

func TestScenariosSetAccountAddressLengthErr1(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-account-addr-len.err1.json")
	require.EqualError(t, err,
		"error processing steps: cannot parse set state step: account address is not 32 bytes in length")
}

func TestScenariosSetAccountAddressLengthErr2(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-account-addr-len.err2.json")
	require.EqualError(t, err,
		"error processing steps: error parsing new addresses: account address is not 32 bytes in length")
}

func TestScenariosSetAccountSCAddressErr1(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-account-sc-addr.err1.json")
	require.EqualError(t, err,
		"\"setState\" step validation failed for account \"address:not-a-sc-address\": account has a smart contract address, but has no code: 0x6e6f742d612d73632d616464726573735f5f5f5f5f5f5f5f5f5f5f5f5f5f5f5f")
}

func TestScenariosSetAccountSCAddressErr2(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-account-sc-addr.err2.json")
	require.EqualError(t, err,
		"\"setState\" step validation failed for account \"sc:should-be-sc\": account has code but not a smart contract address: 0000000000000000000073686f756c642d62652d73635f5f5f5f5f5f5f5f5f5f")
}

func TestScenariosSetAccountSCAddressErr3(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-account-sc-addr.err3.json")
	require.EqualError(t, err,
		"address in \"setState\" \"newAddresses\" field should have SC format: address:not-a-sc-address")
}

func TestScenariosCheckNonceErr(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-nonce.err.json")
	require.EqualError(t, err,
		"Check state \"check-1\": bad account nonce. Account: address:the-address. Want: \"1002\". Have: \"1001\"")
}

func TestScenariosCheckOwnerErr1(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-owner.err1.json")
	require.EqualError(t, err,
		"Check state \"check-1\": bad account owner. Account: address:child. Want: \"address:other\". Have: \"address:parent\"")
}

func TestScenariosCheckOwnerErr2(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-owner.err2.json")
	require.EqualError(t, err,
		"Check state \"check-1\": bad account owner. Account: address:parent. Want: \"address:other\". Have: \"\"")
}

func TestScenariosCheckBalanceErr(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-balance.err.json")
	require.EqualError(t, err,
		"Check state \"check-1\": bad account balance. Account: address:the-address. Want: \"1,000,002\". Have: \"1000001\"")
}

func TestScenariosCheckUsernameErr(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-username.err.json")
	require.EqualError(t, err,
		"Check state \"check-1\": bad account username. Account: address:the-address. Want: \"str:wrong.domain\". Have: \"str:theusername.domain\"")
}

func TestScenariosCheckCodeErr(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-code.err.json")
	require.EqualError(t, err,
		"Check state \"check-1\": bad account code. Account: sc:contract-address. Want: \"file:set-check-code.scen.json\". Have: \"0x7b0a2020202022636f6d...\"")
}

func TestScenariosCheckStorageErr1(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-storage.err1.json")
	require.EqualError(t, err,
		"Check state \"check-1\": wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"str:another-value\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestScenariosCheckStorageErr2(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-storage.err2.json")
	require.EqualError(t, err,
		"Check state \"check-1\": wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestScenariosCheckStorageErr3(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-storage.err3.json")
	require.EqualError(t, err,
		"Check state \"check-1\": wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d64 (str:key-d): Want: \"str:value-d\". Have: \"\"")
}

func TestScenariosCheckStorageErr4(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-storage.err4.json")
	require.EqualError(t, err,
		"Check state \"check-1\": wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d63 (str:key-c): Want: \"\". Have: \"0x76616c75652d63 (str:value-c)\"")
}

func TestScenariosCheckStorageErr5(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-storage.err5.json")
	require.EqualError(t, err,
		"Check state \"check-1\": wrong account storage for account \"address:the-address\":\n"+
			"  for key 0x6b65792d62 (str:key-b): Want: \"str:another-b\". Have: \"0x76616c75652d62 (str:value-b)\"")
}

func TestScenariosCheckESDTErr1(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test/set-check", "set-check-esdt.err1.json")
	require.EqualError(t, err,
		`Check state "check-1": mismatch for account "address:the-address":
  for token: NFT-123456, nonce: 1: Bad balance. Want: "4". Have: "1"
  for token: NFT-123456, nonce: 1: Bad creator. Want: "address:another-address". Have: "address:the-address"
  for token: NFT-123456, nonce: 1: Bad royalties. Want: "2001". Have: "2000"
  for token: NFT-123456, nonce: 1: Bad hash. Want: "keccak256:str:another_hash". Have: 0x54e3ea4bdef3b22154767a2cae081fca2bec2eae1ec62ee71308cb2a300d675d (str:"T\xe3\xeaK\xde\xf3\xb2!Tvz,\xae\b\x1f\xca+\xec.\xae\x1e\xc6.\xe7\x13\b\xcb*0\rg]")
  for token: NFT-123456, nonce: 1: Bad URI. Want: ["str:www.cool_nft.com/another_nft.jpg", "*"]. Have: ["str:www.cool_nft.com/my_nft.jpg", "str:www.cool_nft.com/my_nft.json"]
  for token: NFT-123456, nonce: 1: Bad attributes. Want: "str:other_attributes". Have: "str:serialized_attributes"`)
}

func TestScenariosEsdtZeroBalance(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test", "esdt-zero-balance-check-err.scen.json")
	require.EqualError(t, err,
		`Check state "check-1": mismatch for account "address:A":
  for token: TOK-123456, nonce: 0: Bad balance. Want: "". Have: "150"`)
}

func TestScenariosEsdtNonZeroBalance(t *testing.T) {
	err := runSingleTestReturnError("scenarios-self-test", "esdt-non-zero-balance-check-err.scen.json")
	require.EqualError(t, err,
		`Check state "check-1": mismatch for account "address:B":
  for token: TOK-123456, nonce: 0: Bad balance. Want: "100". Have: "0"`)
}
