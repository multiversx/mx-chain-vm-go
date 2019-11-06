package main

import (
	"testing"

	controller "github.com/ElrondNetwork/elrond-vm-util/test-util/testcontroller"
)

var excludedTests = []string{}

func TestArwenTests(t *testing.T) {
	err := controller.RunAllJSONTestsInDirectory(
		getTestRoot(),
		"erc20",
		".json",
		excludedTests,
		newArwenTestExecutor())

	if err != nil {
		t.Error(err)
	}
}
