package main

import (
	"fmt"
	"os"
	"path/filepath"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

func main() {

	//debugArwenTest("erc20/create.iele.json")
	//debugArwenTest("erc20/balanceOf_Caller.iele.json")
	//debugArwenTest("erc20/transfer_Caller-EntireBalance.iele.json")
	debugArwenTest("dummy/dummy1.wasm.json")

}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../test")
	return arwenTestRoot
}

func debugArwenTest(testFile string) {
	arwenTestRoot := getTestRoot()
	err := controller.RunSingleIeleTest(
		filepath.Join(arwenTestRoot, testFile),
		newArwenTestExecutor())

	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
