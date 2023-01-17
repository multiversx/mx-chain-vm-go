package vmjsonintegrationtest

import (
	"testing"
)

func TestRustAdder(t *testing.T) {
	runAllTestsInFolder(t, "adder/scenarios")
}

func TestRustErc20(t *testing.T) {
	runAllTestsInFolder(t, "erc20-rust/scenarios")
}

func TestCErc20(t *testing.T) {
	runAllTestsInFolder(t, "erc20-c")
}

func TestDigitalCash(t *testing.T) {
	runAllTestsInFolder(t, "digital-cash")
}

func TestMultisig(t *testing.T) {
	runAllTestsInFolder(t, "multisig/scenarios")
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
