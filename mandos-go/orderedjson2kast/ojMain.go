package orderedjson2kast

import (
	"path/filepath"

	oj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/orderedjson"
)

// ProcessCodeFunc represents a callback to assemble the code in the test
type ProcessCodeFunc func(testPath string, value string) string

// ConvertOrderedJSONToKast parses data as an ordered JSON,
// assembles code if necessary
// and converts to KAST format, readable by K
func ConvertOrderedJSONToKast(data []byte, testFilePath string, processCodeCallback ProcessCodeFunc) (string, error) {
	jsonObj, err := oj.ParseOrderedJSON(data)
	if err != nil {
		return "", err
	}
	testDirPath := filepath.Dir(testFilePath)
	processTestCode(jsonObj, testDirPath, processCodeCallback)
	kast := jsonToKastOrdered(jsonObj)

	return kast, nil
}
