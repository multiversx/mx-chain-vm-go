package main

import (
	"fmt"
	"os"
	"path/filepath"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the path to the json test.")
	}
	jsonTestPath := os.Args[1]

	testExec := newArwenTestExecutor().replaceCode(
		"erc20.wasm",
		filepath.Join(getTestRoot(), "contracts/erc20-c.wasm"))
	err := controller.RunSingleJSONTest(
		jsonTestPath,
		testExec)

	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}
