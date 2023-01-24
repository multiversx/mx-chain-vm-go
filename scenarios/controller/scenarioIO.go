package mandoscontroller

import (
	"io/ioutil"
	"os"
	"path/filepath"

	mjparse "github.com/multiversx/mx-chain-vm-go/scenarios/json/parse"
	mjwrite "github.com/multiversx/mx-chain-vm-go/scenarios/json/write"
	mj "github.com/multiversx/mx-chain-vm-go/scenarios/model"
)

// ParseMandosScenario reads and parses a Mandos scenario from a JSON file.
func ParseMandosScenario(parser mjparse.Parser, scenFilePath string) (*mj.Scenario, error) {
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

// ParseMandosScenarioDefaultParser reads and parses a Mandos scenario from a JSON file.
func ParseMandosScenarioDefaultParser(scenFilePath string) (*mj.Scenario, error) {
	parser := mjparse.NewParser(NewDefaultFileResolver())
	parser.ExprInterpreter.FileResolver.SetContext(scenFilePath)
	return ParseMandosScenario(parser, scenFilePath)
}

// WriteMandosScenario exports a Mandos scenario to a file, using the default formatting.
func WriteMandosScenario(scenario *mj.Scenario, toPath string) error {
	jsonString := mjwrite.ScenarioToJSONString(scenario)

	err := os.MkdirAll(filepath.Dir(toPath), os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(toPath, []byte(jsonString), 0644)
}
