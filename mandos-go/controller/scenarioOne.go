package mandoscontroller

// RunSingleJSONScenario parses and prepares test, then calls testCallback.
func (r *ScenarioRunner) RunSingleJSONScenario(contextPath string) error {
	scenario, parseErr := ParseMandosScenario(r.Parser, contextPath)

	if r.RunsNewTest {
		scenario.IsNewTest = true
		r.RunsNewTest = false
	}

	if parseErr != nil {
		return parseErr
	}

	return r.Executor.ExecuteScenario(scenario, r.Parser.ExprInterpreter.FileResolver)
}
