package vmjsonintegrationtest

import (
	"testing"

	"github.com/multiversx/mx-chain-vm-go/testcommon/testexecutor"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
)

func TestRustCompareAdderLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}
	expected := ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("adder/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustFactorialLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("factorial/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("factorial/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustErc20Log(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	t.Skip("not a working test")

	expected := ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("erc20-rust/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCErc20Log(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("erc20-c").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDigitalCashLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("digital-cash").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("digital-cash").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestESDTMultiTransferOnCallbackLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_call_async_retrieve_multi_transfer.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCreateAsyncCallLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	t.Skip("not a working test")

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios-promises").
		File("promises_single_transfer.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios-promises").
		File("promises_single_transfer.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestESDTMultiTransferOnCallAndCallbackLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("features/composability/scenarios").
		File("forw_raw_async_send_and_retrieve_multi_transfer_funds.scen.json").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestMultisigLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	t.Skip("not a working test")

	expected := ScenariosTest(t).
		Folder("multisig/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("multisig/scenarios").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestDnsContractLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	t.Skip("not a working test")

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("dns").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("dns").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCrowdfundingEsdtLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("crowdfunding-esdt").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("crowdfunding-esdt").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestEgldEsdtSwapLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("egld-esdt-swap").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("egld-esdt-swap").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestPingPongEgldLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	expected := ScenariosTest(t).
		Folder("ping-pong-egld").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("ping-pong-egld").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestRustAttestationLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("attestation-rust").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("attestation-rust").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}

func TestCAttestationLog(t *testing.T) {
	if !testexecutor.IsWasmer1Allowed() {
		t.Skip("run exclusively with wasmer1")
	}

	if testing.Short() {
		t.Skip("not a short test")
	}

	expected := ScenariosTest(t).
		Folder("attestation-c").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		ExtractLog()

	ScenariosTest(t).
		Folder("attestation-c").
		WithExecutorFactory(wasmer2.ExecutorFactory()).
		WithExecutorLogs().
		Run().
		CheckNoError().
		CheckLog(expected)
}
