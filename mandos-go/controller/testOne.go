package mandoscontroller

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	mjwrite "github.com/ElrondNetwork/wasm-vm/mandos-go/json/write"
	mj "github.com/ElrondNetwork/wasm-vm/mandos-go/model"
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

// tool to convert .test.json -> .scen.json
// use with extreme caution
func convertTestToScenario(contextPath string, top []*mj.Test) {
	if strings.HasSuffix(contextPath, ".test.json") {
		scenario, err := mj.ConvertTestToScenario(top)
		if err != nil {
			panic(err)
		}
		scenarioSerialized := mjwrite.ScenarioToJSONString(scenario)

		newPath := contextPath[:len(contextPath)-len(".test.json")] + ".scen.json"
		err = os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(newPath, []byte(scenarioSerialized), 0644)
		if err != nil {
			panic(err)
		}
	}
}

// tool to modify tests
// use with extreme caution
func saveModifiedTest(toPath string, top []*mj.Test) {
	resultJSON := mjwrite.TestToJSONString(top)

	err := os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(toPath, []byte(resultJSON), 0644)
	if err != nil {
		panic(err)
	}
}
