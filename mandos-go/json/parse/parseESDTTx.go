package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/model"
	oj "github.com/ElrondNetwork/wasm-vm-v1_4/mandos-go/orderedjson"
)

func (p *Parser) processTxESDT(txEsdtRaw oj.OJsonObject) ([]*mj.ESDTTxData, error) {
	allEsdtData := []*mj.ESDTTxData{}

	switch txEsdt := txEsdtRaw.(type) {
	case *oj.OJsonMap:
		if !p.AllowEsdtTxLegacySyntax {
			return nil, fmt.Errorf("wrong ESDT Multi-Transfer format, list expected")
		}
		entry, err := p.parseSingleTxEsdtEntry(txEsdt)
		if err != nil {
			return nil, err
		}

		allEsdtData = append(allEsdtData, entry)
	case *oj.OJsonList:
		for _, txEsdtListItem := range txEsdt.AsList() {
			txEsdtMap, isMap := txEsdtListItem.(*oj.OJsonMap)
			if !isMap {
				return nil, fmt.Errorf("wrong ESDT Multi-Transfer format")
			}

			entry, err := p.parseSingleTxEsdtEntry(txEsdtMap)
			if err != nil {
				return nil, err
			}

			allEsdtData = append(allEsdtData, entry)
		}
	default:
		return nil, fmt.Errorf("wrong ESDT transfer format, expected list")
	}

	return allEsdtData, nil
}

func (p *Parser) parseSingleTxEsdtEntry(esdtTxEntry *oj.OJsonMap) (*mj.ESDTTxData, error) {
	esdtData := mj.ESDTTxData{}
	var err error

	for _, kvp := range esdtTxEntry.OrderedKV {
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
		default:
			return nil, fmt.Errorf("unknown transaction ESDT data field: %s", kvp.Key)
		}
	}

	return &esdtData, nil
}
