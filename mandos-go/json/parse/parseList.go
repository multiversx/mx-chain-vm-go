package mandosjsonparse

import (
	"errors"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) processStringList(obj interface{}) ([]string, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []string
	for _, elemRaw := range listRaw.AsList() {
		strVal, err := p.parseString(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, strVal)
	}
	return result, nil
}

func (p *Parser) parseByteArrayList(obj interface{}) ([]mj.JSONBytesFromString, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []mj.JSONBytesFromString
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processStringAsByteArray(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, ba)
	}
	return result, nil
}

func (p *Parser) parseSubTreeList(obj interface{}) ([]mj.JSONBytesFromTree, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []mj.JSONBytesFromTree
	for _, elemRaw := range listRaw.AsList() {
		ba, err := p.processSubTreeAsByteArray(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, ba)
	}
	return result, nil
}

func (p *Parser) parseCheckBytesList(obj interface{}) ([]mj.JSONCheckBytes, error) {
	listRaw, listOk := obj.(*oj.OJsonList)
	if !listOk {
		return nil, errors.New("not a JSON list")
	}
	var result []mj.JSONCheckBytes
	for _, elemRaw := range listRaw.AsList() {
		checkBytes, err := p.parseCheckBytes(elemRaw)
		if err != nil {
			return nil, err
		}
		result = append(result, checkBytes)
	}
	return result, nil
}
