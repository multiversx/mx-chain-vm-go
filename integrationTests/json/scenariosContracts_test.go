package vmjsonintegrationtest

import (
	"testing"

	"github.com/multiversx/mx-chain-scenario-go/worldmock"
)

func TestRustAdder(t *testing.T) {
	ScenariosTest(t).
		Folder("adder/scenarios").
		Run().
		CheckNoError()
}

func TestRustFactorial(t *testing.T) {
	ScenariosTest(t).
		Folder("factorial/scenarios").
		Run().
		CheckNoError()
}

func TestRustErc20(t *testing.T) {
	ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		Run().
		CheckNoError()
}

func TestCErc20(t *testing.T) {
	ScenariosTest(t).
		Folder("erc20-c").
		Run().
		CheckNoError()
}

func TestDigitalCash(t *testing.T) {
	ScenariosTest(t).
		Folder("digital-cash").
		Run().
		CheckNoError()
}

func TestESDTMultiTransferOnCallback(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		Run().
		CheckNoError()
}

func TestCreateAsyncCall(t *testing.T) {
	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("promises_single_transfer.scen.json").
		Run().
		CheckNoError()
}

func TestMultisig(t *testing.T) {
	ScenariosTest(t).
		Folder("multisig/scenarios").
		Exclude("multisig/scenarios/interactor*").
		Run().
		CheckNoError()
}

func TestDnsContract(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("dns").
		WithEnableEpochsHandler(worldmock.EnableEpochsHandlerStubNoFlags()).
		Run().
		CheckNoError()
}

func TestCrowdfundingEsdt(t *testing.T) {
	ScenariosTest(t).
		Folder("crowdfunding-esdt").
		Run().
		CheckNoError()
}

func TestWEgldSwap(t *testing.T) {
	ScenariosTest(t).
		Folder("wegld-swap").
		Run().
		CheckNoError()
}

func TestPingPongEgld(t *testing.T) {
	ScenariosTest(t).
		Folder("ping-pong-egld").
		Run().
		CheckNoError()
}

func TestRustAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("attestation-rust").
		Run().
		CheckNoError()
}

func TestCAttestation(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	ScenariosTest(t).
		Folder("attestation-c").
		Run().
		CheckNoError()
}

func TestRustPaymaster(t *testing.T) {
	ScenariosTest(t).
		Folder("paymaster/scenarios").
		Run().
		CheckNoError()
}
