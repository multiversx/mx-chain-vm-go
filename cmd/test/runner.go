package main

import (
	"fmt"
	"os"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the path to the json test.")
	}

	jsonTestPath := os.Args[1]

	testExecutor := newArwenTestExecutor()
	err := controller.RunSingleJSONTest(jsonTestPath, testExecutor)
	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
