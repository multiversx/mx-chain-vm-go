package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ajt "github.com/ElrondNetwork/arwen-wasm-vm/arwenjsontest"
	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func isDirectory(name string) (bool, error) {
	// use a switch to make it a bit cleaner
	fi, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the path to the json test.")
	}

	exeDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	jsonFilePath := filepath.Join(exeDir, os.Args[1])

	isDir, err := isDirectory(jsonFilePath)
	if err != nil {
		fmt.Printf("Path does not exist: %v\n", err)
		return
	}

	// init
	executor, err := ajt.NewArwenTestExecutor()
	if err != nil {
		panic("Could not instantiate Arwen VM")
	}

	// execute
	switch {
	case isDir:
		runner := controller.NewScenarioRunner(
			executor,
			ij.NewDefaultFileResolver(),
		)
		err = runner.RunAllJSONScenariosInDirectory(
			jsonFilePath,
			"",
			".scen.json",
			[]string{})
	case strings.HasSuffix(jsonFilePath, ".scen.json"):
		runner := controller.NewScenarioRunner(
			executor,
			ij.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONScenario(jsonFilePath)
	default:
		runner := controller.NewTestRunner(
			executor,
			ij.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONTest(jsonFilePath)
	}

	// print result
	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
