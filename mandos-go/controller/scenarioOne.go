package mandoscontroller

import (
	mj "github.com/ElrondNetwork/wasm-vm/mandos-go/model"
)

// RunScenarioOptions defines the scenario options component
type RunScenarioOptions struct {
	ForceTraceGas bool
	UseWasmer2    bool
}

func applyScenarioOptions(scenario *mj.Scenario, options *RunScenarioOptions) {
	if options.ForceTraceGas {
		scenario.TraceGas = true
	}
}

// DefaultRunScenarioOptions creates a new RunScenarioOptions instance
func DefaultRunScenarioOptions() *RunScenarioOptions {
	return &RunScenarioOptions{
		ForceTraceGas: false,
	}
}

// RunSingleJSONScenario parses and prepares test, then calls testCallback.
func (r *ScenarioRunner) RunSingleJSONScenario(contextPath string, options *RunScenarioOptions) error {
	scenario, parseErr := ParseMandosScenario(r.Parser, contextPath)

	if parseErr != nil {
		return parseErr
	}

	if r.RunsNewTest {
		scenario.IsNewTest = true
		r.RunsNewTest = false
	}

	applyScenarioOptions(scenario, options)

	return r.Executor.ExecuteScenario(scenario, r.Parser.ExprInterpreter.FileResolver)
}
