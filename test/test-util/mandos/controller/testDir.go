package mandoscontroller

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func isExcluded(excludedFilePatterns []string, testPath string, generalTestPath string) bool {
	for _, et := range excludedFilePatterns {
		excludedFullPath := path.Join(generalTestPath, et)
		match, err := filepath.Match(excludedFullPath, testPath)
		if err != nil {
			panic(err)
		}
		if match {
			return true
		}
	}
	return false
}

// RunAllJSONTestsInDirectory walks directory, parses and prepares all json tests,
// then calls testExecutor for each of them.
func (r *TestRunner) RunAllJSONTestsInDirectory(
	generalTestPath string,
	specificTestPath string,
	allowedSuffix string,
	excludedFilePatterns []string) error {

	mainDirPath := path.Join(generalTestPath, specificTestPath)
	var nrPassed, nrFailed, nrSkipped int

	err := filepath.Walk(mainDirPath, func(testFilePath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(testFilePath, allowedSuffix) {
			fmt.Printf("Test: %s ... ", shortenTestPath(testFilePath, generalTestPath))
			if isExcluded(excludedFilePatterns, testFilePath, generalTestPath) {
				nrSkipped++
				fmt.Print("  skip\n")
			} else {
				testErr := r.RunSingleJSONTest(testFilePath)
				if testErr == nil {
					nrPassed++
					fmt.Print("  ok\n")
				} else {
					nrFailed++
					fmt.Print("  FAIL!!!\n")
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

func shortenTestPath(path string, generalTestPath string) string {
	if strings.HasPrefix(path, generalTestPath+"/") {
		return path[len(generalTestPath)+1:]
	}
	return path
}
