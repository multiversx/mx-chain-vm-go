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
	excludedFilePatterns []string) error {

	mainDirPath := path.Join(generalTestPath, specificTestPath)
	var nrPassed, nrFailed, nrSkipped int

	err := filepath.Walk(mainDirPath, func(testFilePath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(testFilePath, allowedSuffix) {
			fmt.Printf("Scenario: %s ... ", shortenTestPath(testFilePath, generalTestPath))
			if isExcluded(excludedFilePatterns, testFilePath, generalTestPath) {
				nrSkipped++
				fmt.Print("  skip\n")
			} else {
				r.Executor.Reset()
				testErr := r.RunSingleJSONScenario(testFilePath)
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
		return errors.New("Some tests failed")
	}

	return nil
}
