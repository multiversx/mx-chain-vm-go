package tests

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var bytecodeCounter []byte

func init() {
	bytecodeCounter = getSCCode("./../../test/contracts/counter.wasm")
}

func getSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("Cannot read file [%s].", fileName))
	}

	return code
}
