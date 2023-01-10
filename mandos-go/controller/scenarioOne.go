package mandoscontroller

import (
	mj "github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/model"
)

type RunScenarioOptions struct {
	ForceTraceGas bool
}

func applyScenarioOptions(scenario *mj.Scenario, options *RunScenarioOptions) {
	if options.ForceTraceGas {
		scenario.TraceGas = true
	}
}

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
