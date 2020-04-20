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
		runner := controller.NewScenarioRunner(
			ajt.NewArwenScenarioExecutor(),
			ij.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONScenario(jsonFilePath)
	} else {
		runner := controller.NewTestRunner(
			ajt.NewArwenTestExecutor(),
			ij.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONTest(jsonFilePath)
	}

	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
