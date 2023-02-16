package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mc "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/controller"
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
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if shouldFormatFile(filePath) {
			fmt.Printf("Upgrade: %s\n ", filePath)
			upgradeScenariosFile(filePath)
		}
		return nil
	})
	return err
}

func upgradeScenariosFile(filePath string) {
	scenario, err := mc.ParseScenariosScenarioDefaultParser(filePath)
	if err == nil {
		_ = mc.WriteScenariosScenario(scenario, filePath)
	} else {
		fmt.Printf("Error upgrading: %s\n", err.Error())
	}
}
