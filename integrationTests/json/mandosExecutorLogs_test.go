package vmjsonintegrationtest

import (
	"testing"

	"github.com/ElrondNetwork/wasm-vm/wasmer2"
)

func TestRustAdderLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("adder/mandos").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("adder/mandos").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustFactorialLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("factorial/mandos").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("factorial/mandos").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustErc20Log(t *testing.T) {
	t.Skip("not a working test")

	expected := MandosTest(t).
		Folder("erc20-rust/mandos").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("erc20-rust/mandos").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCErc20Log(t *testing.T) {
	expected := MandosTest(t).
		Folder("erc20-c").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("erc20-c").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDigitalCashLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("digital-cash").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("digital-cash").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestESDTMultiTransferOnCallbackLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("features/composability/mandos").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("features/composability/mandos").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCreateAsyncCallLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("features/composability/mandos-promises").
		File("promises_single_transfer.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("features/composability/mandos-promises").
		File("promises_single_transfer.scen.json").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestESDTMultiTransferOnCallAndCallbackLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("features/composability/mandos").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("features/composability/mandos").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestMultisigLog(t *testing.T) {
	t.Skip("not a working test")

	expected := MandosTest(t).
		Folder("multisig/mandos").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("multisig/mandos").
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

	expected := MandosTest(t).
		Folder("dns").
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("dns").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCrowdfundingEsdtLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("crowdfunding-esdt").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("crowdfunding-esdt").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestEgldEsdtSwapLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("egld-esdt-swap").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("egld-esdt-swap").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestPingPongEgldLog(t *testing.T) {
	expected := MandosTest(t).
		Folder("ping-pong-egld").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
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

	expected := MandosTest(t).
		Folder("attestation-rust").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
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

	expected := MandosTest(t).
		Folder("attestation-c").
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	MandosTest(t).
		Folder("attestation-c").
		WithExecutorLogs().
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		Run().
		CheckNoError().
		CheckLog(expected)
}
