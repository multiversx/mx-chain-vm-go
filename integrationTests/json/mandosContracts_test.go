package vmjsonintegrationtest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRustAdder(t *testing.T) {
	runAllTestsInFolder(t, "adder/mandos")
}

func TestRustFactorial(t *testing.T) {
	runAllTestsInFolder(t, "factorial/mandos")
}

func TestRustErc20(t *testing.T) {
	runAllTestsInFolder(t, "erc20-rust/mandos")
}

func TestCErc20(t *testing.T) {
	runAllTestsInFolder(t, "erc20-c")
}

func TestDigitalCash(t *testing.T) {
	runAllTestsInFolder(t, "digital-cash")
}

func TestESDTMultiTransferOnCallback(t *testing.T) {
	err := runSingleTestReturnError(
		"features/composability/mandos",
		"forw_raw_call_async_retrieve_multi_transfer.scen.json")
	require.Nil(t, err)
}

func TestCreateAsyncCall(t *testing.T) {
	err := runSingleTestReturnError(
		"features/composability/mandos-promises",
		"promises_single_transfer.scen.json")
	require.Nil(t, err)
}

func TestESDTMultiTransferOnCallAndCallback(t *testing.T) {
	err := runSingleTestReturnError(
		"features/composability/mandos",
		"forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json")
	require.Nil(t, err)
}

func TestMultisig(t *testing.T) {
	runAllTestsInFolder(t, "multisig/mandos")
}

func TestDnsContract(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "dns")
}

func TestCrowdfundingEsdt(t *testing.T) {
	runAllTestsInFolder(t, "crowdfunding-esdt")
}

func TestEgldEsdtSwap(t *testing.T) {
	runAllTestsInFolder(t, "egld-esdt-swap")
}

func TestPingPongEgld(t *testing.T) {
	runAllTestsInFolder(t, "ping-pong-egld")
}

func TestRustAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	runAllTestsInFolder(t, "attestation-rust")
}

func TestCAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}
	runAllTestsInFolder(t, "attestation-c")
}
