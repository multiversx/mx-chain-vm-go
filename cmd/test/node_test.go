package main

import (
	"testing"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

var excludedTests = []string{}

func TestArwenTests(t *testing.T) {
	err := controller.RunAllIeleTestsInDirectory(
		getTestRoot(),
		"erc20",
		excludedTests,
		newArwenTestExecutor())

	if err != nil {
		t.Error(err)
	}
}
