package scencontroller

import (
	"io/ioutil"
	"os"
	"path/filepath"

	mjparse "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/json/parse"
	mjwrite "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/json/write"
	mj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/model"
)

// ParseScenariosScenario reads and parses a Scenarios scenario from a JSON file.
func ParseScenariosScenario(parser mjparse.Parser, scenFilePath string) (*mj.Scenario, error) {
	var err error
	scenFilePath, err = filepath.Abs(scenFilePath)
	if err != nil {
		return nil, err
	}

	// Open our jsonFile
	var jsonFile *os.File
	jsonFile, err = os.Open(scenFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		_ = jsonFile.Close()
	}()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	parser.ExprInterpreter.FileResolver.SetContext(scenFilePath)
	return parser.ParseScenarioFile(byteValue)
}

// ParseScenariosScenarioDefaultParser reads and parses a Scenarios scenario from a JSON file.
func ParseScenariosScenarioDefaultParser(scenFilePath string) (*mj.Scenario, error) {
	parser := mjparse.NewParser(NewDefaultFileResolver())
	parser.ExprInterpreter.FileResolver.SetContext(scenFilePath)
	return ParseScenariosScenario(parser, scenFilePath)
}

// WriteScenariosScenario exports a Scenarios scenario to a file, using the default formatting.
func WriteScenariosScenario(scenario *mj.Scenario, toPath string) error {
	jsonString := mjwrite.ScenarioToJSONString(scenario)

	err := os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(toPath, []byte(jsonString), 0644)
}
