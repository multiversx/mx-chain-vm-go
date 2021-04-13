package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processAppendCheckESDTData(
	tokenName []byte,
	esdtDataRaw oj.OJsonObject,
	output []*mj.CheckESDTData) ([]*mj.CheckESDTData, error) {

	var err error

	switch data := esdtDataRaw.(type) {
	case *oj.OJsonString:
		// simple string representing balance "400,000,000,000"
		esdtData := mj.CheckESDTData{}
		esdtData.TokenIdentifier = mj.NewJSONBytesFromString(tokenName, string(tokenName))
		esdtData.Value, err = p.processCheckBigInt(esdtDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return output, fmt.Errorf("invalid ESDT balance: %w", err)
		}

		output = append(output, &esdtData)
		return output, nil
	case *oj.OJsonMap:
		esdtData, err := p.processCheckESDTDataMap(tokenName, data)
		if err != nil {
			return output, err
		}
		output = append(output, esdtData)
		return output, nil
	case *oj.OJsonList:
		for _, item := range data.AsList() {
			itemAsMap, isMap := item.(*oj.OJsonMap)
			if !isMap {
				return nil, errors.New("JSON map expected in ESDT list")
			}
			esdtData, err := p.processCheckESDTDataMap(tokenName, itemAsMap)
			if err != nil {
				return output, err
			}
			output = append(output, esdtData)
		}
		return output, nil
	default:
		return output, errors.New("invalid JSON object for ESDT")
	}
}

// map containing other fields too, e.g.:
// {
// 	"balance": "400,000,000,000",
// 	"frozen": "true"
// }
func (p *Parser) processCheckESDTDataMap(tokenNameKey []byte, esdtDataMap *oj.OJsonMap) (*mj.CheckESDTData, error) {
	esdtData := mj.CheckESDTData{
		TokenIdentifier: mj.NewJSONBytesFromString(tokenNameKey, ""),
	}
	var err error

	for _, kvp := range esdtDataMap.OrderedKV {
		switch kvp.Key {
		case "tokenIdentifier":
			esdtData.TokenIdentifier, err = p.processStringAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid ESDT token name: %w", err)
			}
		case "nonce":
			esdtData.Nonce, err = p.processCheckUint64(kvp.Value)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		case "value":
			esdtData.Value, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
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
