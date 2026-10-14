package vmjsonintegrationtest

import (
	"os"
	"testing"

	"github.com/multiversx/mx-chain-vm-go/executor"
	executorwrapper "github.com/multiversx/mx-chain-vm-go/executor/wrapper"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	"github.com/multiversx/mx-chain-vm-go/wasmer2legacy"
	"github.com/multiversx/mx-chain-vm-go/wasmer7"
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

// TestExecutorsAgreeOnHookSequence is the end-to-end differential: production today against
// where the migration is going.
//
// This is the harness that was used to migrate from Wasmer 1 to Wasmer 2, pointed at the
// Wasmer 2 to Wasmer 7 migration. It is the acceptance criterion of PLAN-go-ffi.md Phase 5,
// at the granularity of individual hook calls rather than of a VMOutput.
//
// It moves two things at once — the boundary and the wasm runtime — so a difference it reports
// does not say which one caused it. TestMigrationPreservesHookSequence below is the one that
// does.
//
// Skipped by default while the Go side is still being converged: the sequences differ in the
// instance-lifecycle lines, because the new executor has no reset or module cache yet
// (PLAN-go-ffi.md, Phase 6). Run it explicitly to see the current difference:
//
//	VMEXECUTOR_DIFF=1 go test ./integrationTests/json/ -run TestExecutorsAgreeOnHookSequence
func TestExecutorsAgreeOnHookSequence(t *testing.T) {
	requireDiffEnabled(t)

	compareHookSequences(t,
		"wasmer2legacy", wasmer2legacy.ExecutorFactory(),
		"wasmer7", wasmer7.ExecutorFactory())
}

// TestMigrationPreservesHookSequence holds the wasm runtime constant and varies only the
// boundary.
//
// Both sides run Wasmer 2.2; the first reaches it through the legacy C API and executor
// interfaces, the second through the new ones. Anything that differs here is the migration's
// doing, with no runtime change available to blame — which is the claim the test above cannot
// make. It is the stricter of the two, and the one to get green first.
//
//	VMEXECUTOR_DIFF=1 go test ./integrationTests/json/ -run TestMigrationPreservesHookSequence
func TestMigrationPreservesHookSequence(t *testing.T) {
	requireDiffEnabled(t)

	compareHookSequences(t,
		"wasmer2legacy", wasmer2legacy.ExecutorFactory(),
		"wasmer2", wasmer2.ExecutorFactory())
}

func requireDiffEnabled(t *testing.T) {
	if os.Getenv("VMEXECUTOR_DIFF") == "" {
		t.Skip("set VMEXECUTOR_DIFF=1 to compare two executors' hook sequences")
	}
}

func diffFolder() string {
	folder := os.Getenv("VMEXECUTOR_DIFF_FOLDER")
	if folder == "" {
		return "adder/scenarios/adder.scen.json"
	}

	return folder
}

func compareHookSequences(
	t *testing.T,
	nameA string, factoryA executor.ExecutorAbstractFactory,
	nameB string, factoryB executor.ExecutorAbstractFactory,
) {
	folder := diffFolder()

	sequenceA := hookSequenceUnder(t, folder, factoryA)
	sequenceB := hookSequenceUnder(t, folder, factoryB)

	if sequenceA == sequenceB {
		return
	}

	// Written out rather than diffed inline, because the sequences run to thousands of lines
	// for anything bigger than adder and a testing.T failure message is the wrong place for
	// them.
	writeDiffArtifact(t, "hookSequence_"+nameA+".txt", sequenceA)
	writeDiffArtifact(t, "hookSequence_"+nameB+".txt", sequenceB)

	t.Errorf("%s and %s produced different hook sequences; see hookSequence_*.txt", nameA, nameB)
}

func writeDiffArtifact(t *testing.T, name string, content string) {
	err := os.WriteFile(name, []byte(content), 0644)
	if err != nil {
		t.Errorf("could not write %s: %v", name, err)
	}
}
