package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mc "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/controller"
)

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the root path where to search.")
	}

	convertAllInFolder(os.Args[1])
}

func convertAllInFolder(path string) error {
	err := filepath.Walk(path, func(mandosFilePath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(mandosFilePath, ".scen.json") {
			fmt.Printf("Upgrade: %s\n ", mandosFilePath)
			upgradeMandosFile(mandosFilePath)
		}
		return nil
	})
	return err
}

func upgradeMandosFile(mandosFilePath string) {
	scenario, err := mc.ParseMandosScenarioDefaultParser(mandosFilePath)
	if err == nil {
		mc.WriteMandosScenario(scenario, mandosFilePath)
	}
}
