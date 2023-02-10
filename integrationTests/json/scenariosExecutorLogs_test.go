package vmjsonintegrationtest

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

func TestRustCompareAdderLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustFactorialLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("factorial/scenarios").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("factorial/scenarios").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustErc20Log(t *testing.T) {
	t.Skip("not a working test")

	expected := ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCErc20Log(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDigitalCashLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("digital-cash").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("digital-cash").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestESDTMultiTransferOnCallbackLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCreateAsyncCallLog(t *testing.T) {
	t.Skip("not a working test")

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios-promises").
		File("promises_single_transfer.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios-promises").
		File("promises_single_transfer.scen.json").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestESDTMultiTransferOnCallAndCallbackLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestMultisigLog(t *testing.T) {
	t.Skip("not a working test")

	expected := ScenariosTest(t).
		Folder("multisig/scenarios").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("multisig/scenarios").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDnsContractLog(t *testing.T) {
	t.Skip("not a working test")

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("dns").
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("dns").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCrowdfundingEsdtLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("crowdfunding-esdt").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("crowdfunding-esdt").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestEgldEsdtSwapLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("egld-esdt-swap").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("egld-esdt-swap").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestPingPongEgldLog(t *testing.T) {
	expected := ScenariosTest(t).
		Folder("ping-pong-egld").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("ping-pong-egld").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustAttestationLog(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("attestation-rust").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("attestation-rust").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCAttestationLog(t *testing.T) {
	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("attestation-c").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("attestation-c").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}
