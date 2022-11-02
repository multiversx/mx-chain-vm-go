package mandoscontroller

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// RunAllJSONScenariosInDirectory walks directory, parses and prepares all json scenarios,
// then calls scenarioExecutor for each of them.
func (r *ScenarioRunner) RunAllJSONScenariosInDirectory(
	generalTestPath string,
	specificTestPath string,
	allowedSuffix string,
	excludedFilePatterns []string,
	options *RunScenarioOptions) error {

	mainDirPath := path.Join(generalTestPath, specificTestPath)
	var nrPassed, nrFailed, nrSkipped int

	err := filepath.Walk(mainDirPath, func(testFilePath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(testFilePath, allowedSuffix) {
			if isExcluded(excludedFilePatterns, testFilePath, generalTestPath) {
				nrSkipped++
				fmt.Printf("Scenario: %s ... ", shortenTestPath(testFilePath, generalTestPath))
				fmt.Print("  skip\n")
			} else {
				r.Executor.Reset()
				r.RunsNewTest = true
				fmt.Printf("Scenario: %s ... ", shortenTestPath(testFilePath, generalTestPath))
				testErr := r.RunSingleJSONScenario(testFilePath, options)
				if testErr == nil {
					nrPassed++
					fmt.Print("  ok\n")
				} else {
					nrFailed++
					fmt.Printf("  FAIL: %s\n", testErr.Error())
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Printf("Done. Passed: %d. Failed: %d. Skipped: %d.\n", nrPassed, nrFailed, nrSkipped)
	if nrFailed > 0 {
		return errors.New("some tests failed")
	}

	return nil
}
