package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processESDTData(esdtDataRaw oj.OJsonObject) (*mj.ESDTData, error) {
	esdtData := mj.ESDTData{}
	var err error

	if _, isStr := esdtDataRaw.(*oj.OJsonString); isStr {
		// simple string representing balance "400,000,000,000"
		esdtData.Balance, err = p.processBigInt(esdtDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid ESDT balance: %w", err)
		}
		return &esdtData, nil
	}

	// map containing other fields too, e.g.:
	// {
	// 	"balance": "400,000,000,000",
	// 	"frozen": "true"
	// }
	esdtDataMap, isMap := esdtDataRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("account ESDT data should be either JSON string or map")
	}

	for _, kvp := range esdtDataMap.OrderedKV {
		switch kvp.Key {
		case "balance":
			esdtData.Balance, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid ESDT balance: %w", err)
			}
		case "frozen":
			esdtData.Frozen, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid ESDT frozen flag: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown ESDT data field: %s", kvp.Key)
		}
	}

	return &esdtData, nil
}

func (p *Parser) processCheckESDTData(esdtDataRaw oj.OJsonObject) (*mj.CheckESDTData, error) {
	esdtData := mj.CheckESDTData{}
	var err error

	if _, isStr := esdtDataRaw.(*oj.OJsonString); isStr {
		// simple string representing balance "400,000,000,000"
		esdtData.Balance, err = p.processCheckBigInt(esdtDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return nil, fmt.Errorf("invalid ESDT balance: %w", err)
		}
		return &esdtData, nil
	}

	// map containing other fields too, e.g.:
	// {
	// 	"balance": "400,000,000,000",
	// 	"frozen": "true"
	// }
	esdtDataMap, isMap := esdtDataRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("account ESDT data should be either JSON string or map")
	}

	for _, kvp := range esdtDataMap.OrderedKV {
		switch kvp.Key {
		case "balance":
			esdtData.Balance, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid ESDT balance: %w", err)
			}
		case "frozen":
			esdtData.Frozen, err = p.processCheckUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid ESDT frozen flag: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown ESDT data field: %s", kvp.Key)
		}
	}

	return &esdtData, nil
}
