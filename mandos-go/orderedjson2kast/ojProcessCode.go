package orderedjson2kast

import (
	oj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/orderedjson"
)

func processTestCode(jobj oj.OJsonObject, testPath string, processCodeCallback ProcessCodeFunc) {
	switch j := jobj.(type) {
	case *oj.OJsonMap:
		isCreateTx := false
		for _, keyValuePair := range j.OrderedKV {
			if keyValuePair.Key == "to" {
				if strVal, isStr := keyValuePair.Value.(*oj.OJsonString); isStr {
					if strVal.Value == "" {
						isCreateTx = true
						break
					}
				}
			}
		}

		for _, keyValuePair := range j.OrderedKV {
			if keyValuePair.Key == "code" ||
				(keyValuePair.Key == "contractCode" && isCreateTx) {
				if strVal, isStr := keyValuePair.Value.(*oj.OJsonString); isStr {
					strVal.Value = processCodeCallback(testPath, strVal.Value)
				}
			} else {
				processTestCode(keyValuePair.Value, testPath, processCodeCallback)
			}
		}
	case *oj.OJsonList:
		collection := []oj.OJsonObject(*j)
		for _, elem := range collection {
			processTestCode(elem, testPath, processCodeCallback)
		}
	default:
	}
}
