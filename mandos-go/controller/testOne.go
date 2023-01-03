package mandoscontroller

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// RunSingleJSONTest parses and prepares test, then calls testCallback.
func (r *TestRunner) RunSingleJSONTest(contextPath string) error {
	var err error
	contextPath, err = filepath.Abs(contextPath)
	if err != nil {
		return err
	}

	// Open our jsonFile
	var jsonFile *os.File
	jsonFile, err = os.Open(contextPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		_ = jsonFile.Close()
	}()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	r.Parser.ExprInterpreter.FileResolver.SetContext(contextPath)
	top, parseErr := r.Parser.ParseTestFile(byteValue)
	if parseErr != nil {
		return parseErr
	}

	for _, test := range top {
		testErr := r.Executor.ExecuteTest(test)
		if testErr != nil {
			return testErr
		}
	}

	return nil
}
