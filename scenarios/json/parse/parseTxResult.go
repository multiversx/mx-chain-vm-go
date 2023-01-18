package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/multiversx/wasm-vm/scenarios/model"
	oj "github.com/multiversx/wasm-vm/scenarios/orderedjson"
)

func (p *Parser) processTxExpectedResult(blrRaw oj.OJsonObject) (*mj.TransactionResult, error) {
	blrMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block result is not a map")
	}

	blr := mj.TransactionResult{
		Status:  mj.JSONCheckBigIntUnspecified(),
		Message: mj.JSONCheckBytesUnspecified(),
		Gas:     mj.JSONCheckUint64Unspecified(),
		Refund:  mj.JSONCheckBigIntUnspecified(),
		Logs:    mj.LogList{IsUnspecified: true, IsStar: true},
	}
	var err error
	for _, kvp := range blrMap.OrderedKV {
		switch kvp.Key {
		case "out":
			blr.Out, err = p.parseCheckValueList(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block result out: %w", err)
			}
		case "status":
			blr.Status, err = p.processCheckBigInt(kvp.Value, bigIntSignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block result status: %w", err)
			}
		case "message":
			blr.Message, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block result message: %w", err)
			}
		case "logs":
			blr.Logs, err = p.processLogList(kvp.Value)
			if err != nil {
				return nil, err
			}
		case "gas":
			blr.Gas, err = p.processCheckUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid block result gas: %w", err)
			}
		case "refund":
			blr.Refund, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid block result refund: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown tx result field: %s", kvp.Key)
		}
	}

	return &blr, nil
}
