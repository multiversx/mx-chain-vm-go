package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/model"
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/orderedjson"
)

func (p *Parser) processBlockInfo(blockInfoRaw oj.OJsonObject) (*mj.BlockInfo, error) {
	blockMap, isMap := blockInfoRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled block info object is not a map")
	}
	blockInfo := &mj.BlockInfo{}
	var err error

	for _, kvp := range blockMap.OrderedKV {
		switch kvp.Key {
		case "blockTimestamp":
			blockInfo.BlockTimestamp, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockTimestamp: %w", err)
			}
		case "blockNonce":
			blockInfo.BlockNonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockNonce: %w", err)
			}
		case "blockRound":
			blockInfo.BlockRound, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockRound: %w", err)
			}
		case "blockEpoch":
			blockInfo.BlockEpoch, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockEpoch: %w", err)
			}
		case "blockRandomSeed":
			blockRandomSeed, err := p.processSubTreeAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("error parsing blockEpoch: %w", err)
			}
			if len(blockRandomSeed.Value) != 48 {
				return nil, fmt.Errorf("blockRandomSeed must be 48 bytes long. Actual length: %d", len(blockRandomSeed.Value))
			}
			blockInfo.BlockRandomSeed = &blockRandomSeed
		default:
			return nil, fmt.Errorf("unknown block info field: %s", kvp.Key)
		}
	}

	return blockInfo, nil
}
