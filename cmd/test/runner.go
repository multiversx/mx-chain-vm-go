package main

import (
	"fmt"
	"os"
	"strings"

	ajt "github.com/ElrondNetwork/arwen-wasm-vm/arwenjsontest"
	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
	ij "github.com/ElrondNetwork/elrond-vm-util/test-util/vmtestjson"
)

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the path to the json test.")
	}

	jsonFilePath := os.Args[1]
	var err error
	if strings.HasSuffix(jsonFilePath, ".scen.json") {
		executor, err := ajt.NewArwenScenarioExecutor()
		if err != nil {
			panic("Could not instantiate Arwen VM")
		}
		runner := controller.NewScenarioRunner(
			executor,
			ij.NewDefaultFileResolver(),
		)
		if err == nil {
			err = runner.RunSingleJSONScenario(jsonFilePath)
		}
	} else {
		executor, err := ajt.NewArwenTestExecutor()
		if err != nil {
			panic("Could not instantiate Arwen VM")
		}
		runner := controller.NewTestRunner(
			executor,
			ij.NewDefaultFileResolver(),
		)
		if err == nil {
			err = runner.RunSingleJSONTest(jsonFilePath)
		}
	}

	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
