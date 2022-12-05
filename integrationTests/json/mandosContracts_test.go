package vmjsonintegrationtest

import (
	"testing"
)

func TestRustAdder(t *testing.T) {
	MandosTest(t).
		Folder("adder/mandos").
		Run().
		CheckNoError()
}

func TestRustFactorial(t *testing.T) {
	MandosTest(t).
		Folder("factorial/mandos").
		Run().
		CheckNoError()
}

func TestRustErc20(t *testing.T) {
	MandosTest(t).
		Folder("erc20-rust/mandos").
		Run().
		CheckNoError()
}

func TestCErc20(t *testing.T) {
	MandosTest(t).
		Folder("erc20-c").
		Run().
		CheckNoError()
}

func TestDigitalCash(t *testing.T) {
	MandosTest(t).
		Folder("digital-cash").
		Run().
		CheckNoError()
}

func TestESDTMultiTransferOnCallback(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		Run().
		CheckNoError()
}

func TestCreateAsyncCall(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos-promises").
		File("promises_single_transfer.scen.json").
		Run().
		CheckNoError()
}

func TestESDTMultiTransferOnCallAndCallback(t *testing.T) {
	MandosTest(t).
		Folder("features/composability/mandos").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		Run().
		CheckNoError()
}

func TestMultisig(t *testing.T) {
	MandosTest(t).
		Folder("multisig/mandos").
		Run().
		CheckNoError()
}

func TestDnsContract(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("dns").
		Run().
		CheckNoError()
}

func TestCrowdfundingEsdt(t *testing.T) {
	MandosTest(t).
		Folder("crowdfunding-esdt").
		Run().
		CheckNoError()
}

func TestEgldEsdtSwap(t *testing.T) {
	MandosTest(t).
		Folder("egld-esdt-swap").
		Run().
		CheckNoError()
}

func TestPingPongEgld(t *testing.T) {
	MandosTest(t).
		Folder("ping-pong-egld").
		Run().
		CheckNoError()
}

func TestRustAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("attestation-rust").
		Run().
		CheckNoError()
}

func TestCAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	MandosTest(t).
		Folder("attestation-c").
		Run().
		CheckNoError()
}
