package main

import (
	"testing"

	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func Test1(t *testing.T) {
	var testCases []*pureFunctionIO

	testCases = append(testCases, &pureFunctionIO{
		functionName:    "add_big_int",
		arguments:       [][]byte{[]byte{1}, []byte{2}},
		expectedStatus:  vmi.Ok,
		expectedResults: [][]byte{[]byte{3}},
	})

	pureFunctionTest(t, testCases)
}
