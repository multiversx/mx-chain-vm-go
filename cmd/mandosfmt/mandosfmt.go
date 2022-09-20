package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mc "github.com/ElrondNetwork/wasm-vm/mandos-go/controller"
)

func main() {
	if len(os.Args) != 2 {
		panic("One argument expected - the root path where to search.")
	}

	_ = convertAllInFolder(os.Args[1])
}

var suffixes = []string{".scen.json", ".step.json", ".steps.json"}

func shouldFormatFile(path string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}

func convertAllInFolder(path string) error {
	err := filepath.Walk(path, func(mandosFilePath string, info os.FileInfo, err error) error {
		if shouldFormatFile(mandosFilePath) {
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
		_ = mc.WriteMandosScenario(scenario, mandosFilePath)
	} else {
		fmt.Printf("Error upgrading: %s\n", err.Error())
	}
}
