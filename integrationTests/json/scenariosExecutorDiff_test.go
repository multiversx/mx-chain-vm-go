package vmjsonintegrationtest

import (
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	executorwrapper "github.com/multiversx/mx-chain-vm-go/executor/wrapper"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/multiversx/mx-chain-vm-go/wasmer3"
)

// hookSequenceUnder runs a scenario folder under one executor and returns the recorded sequence
// of VM hook calls.
//
// The wrapper logs each hook with its arguments and its result, so the string is a trace of
// what the contract asked of the VM and what it was told, in order. Two executors that agree
// on it agree on everything the contract can observe through the hooks, which is a far stronger
// statement than matching final VMOutputs.
func hookSequenceUnder(t *testing.T, folder string, factory executor.ExecutorAbstractFactory) string {
	logger := executorwrapper.NewStringLogger()

	ScenariosTest(t).
		Folder(folder).
		WithExecutorFactory(factory).
		WithCustomExecutorLogs(logger).
		Run().
		CheckNoError()

	return logger.String()
}

// TestExecutorsAgreeOnHookSequence is the differential check between the two backends.
//
// This is the harness that was used to migrate from Wasmer 1 to Wasmer 2, pointed at the
// Wasmer 2 to Wasmer 6 migration. It is the acceptance criterion of PLAN-go-ffi.md Phase 5,
// at the granularity of individual hook calls rather than of a VMOutput.
//
// Skipped by default while the Go side is still being converged: the sequences differ in the
// instance-lifecycle lines, because the new executor has no reset or module cache yet
// (PLAN-go-ffi.md, Phase 6). Run it explicitly to see the current difference:
//
//	VMEXECUTOR_DIFF=1 go test ./integrationTests/json/ -run TestExecutorsAgreeOnHookSequence
func TestExecutorsAgreeOnHookSequence(t *testing.T) {
	if os.Getenv("VMEXECUTOR_DIFF") == "" {
		t.Skip("set VMEXECUTOR_DIFF=1 to compare the two backends' hook sequences")
	}

	folder := os.Getenv("VMEXECUTOR_DIFF_FOLDER")
	if folder == "" {
		folder = "adder/scenarios/adder.scen.json"
	}

	withWasmer2 := hookSequenceUnder(t, folder, wasmer2.ExecutorFactory())
	withWasmer3 := hookSequenceUnder(t, folder, wasmer3.ExecutorFactory())

	if withWasmer2 == withWasmer3 {
		return
	}

	// Written out rather than diffed inline, because the sequences run to thousands of lines
	// for anything bigger than adder and a testing.T failure message is the wrong place for
	// them.
	writeDiffArtifact(t, "hookSequence_wasmer2.txt", withWasmer2)
	writeDiffArtifact(t, "hookSequence_wasmer3.txt", withWasmer3)

	t.Errorf("the two executors produced different hook sequences; see hookSequence_wasmer*.txt")
}

func writeDiffArtifact(t *testing.T, name string, content string) {
	err := os.WriteFile(name, []byte(content), 0644)
	if err != nil {
		t.Errorf("could not write %s: %v", name, err)
	}
}
