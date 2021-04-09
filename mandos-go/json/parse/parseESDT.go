package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processAppendESDTData(
	tokenName []byte,
	esdtDataRaw oj.OJsonObject,
	output []*mj.ESDTData) ([]*mj.ESDTData, error) {

	var err error

	switch data := esdtDataRaw.(type) {
	case *oj.OJsonString:
		// simple string representing balance "400,000,000,000"
		esdtData := mj.ESDTData{}
		esdtData.TokenIdentifier = mj.NewJSONBytesFromString(tokenName, string(tokenName))
		esdtData.Value, err = p.processBigInt(esdtDataRaw, bigIntUnsignedBytes)
		if err != nil {
			return output, fmt.Errorf("invalid ESDT balance: %w", err)
		}

		output = append(output, &esdtData)
		return output, nil
	case *oj.OJsonMap:
		esdtData, err := p.processESDTDataMap(tokenName, data)
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
			esdtData, err := p.processESDTDataMap(tokenName, itemAsMap)
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

func (p *Parser) processTxESDT(esdtRaw oj.OJsonObject) (*mj.ESDTData, error) {
	esdtDataMap, isMap := esdtRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}
	return p.processESDTDataMap([]byte{}, esdtDataMap)
}

// map containing other fields too, e.g.:
// {
// 	"balance": "400,000,000,000",
// 	"frozen": "true"
// }
func (p *Parser) processESDTDataMap(tokenNameKey []byte, esdtDataMap *oj.OJsonMap) (*mj.ESDTData, error) {
	esdtData := mj.ESDTData{
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
			esdtData.Nonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		case "value":
			esdtData.Value, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
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

func (p *Parser) processESDTRoles(esdtRolesRaw oj.OJsonObject) ([]*mj.ESDTRoles, error) {
	var rolesList []*mj.ESDTRoles
	esdtRolesMap, isMap := esdtRolesRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("ESDTRoles object is not a map")
	}
	for _, kvp := range esdtRolesMap.OrderedKV {
		tokenNameStr, err := p.ValueInterpreter.InterpretString(kvp.Key)
		if err != nil {
			return nil, fmt.Errorf("invalid esdt token identifer: %w", err)
		}
		tokenRoles, err := p.processStringList(kvp.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse esdtRoles: %w", err)
		}
		rolesList = append(rolesList, &mj.ESDTRoles{
			TokenIdentifier: mj.NewJSONBytesFromString(tokenNameStr, kvp.Key),
			Roles:           tokenRoles,
		})
	}

	return rolesList, nil
}

func (p *Parser) processESDTLastNonces(esdtLastNonces *oj.OJsonMap) (map[string]*mj.JSONUint64, error) {
	lastNonces := make(map[string]*mj.JSONUint64)
	for _, kvp := range esdtLastNonces.OrderedKV {
		tokenNameStr, err := p.ValueInterpreter.InterpretString(kvp.Key)
		if err != nil {
			return nil, fmt.Errorf("invalid esdt token identifer: %w", err)
		}
		nonce, err := p.processUint64(kvp.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid esdt last nonce: %w", err)
		}

		lastNonces[string(tokenNameStr)] = &nonce
	}

	return lastNonces, nil
}
