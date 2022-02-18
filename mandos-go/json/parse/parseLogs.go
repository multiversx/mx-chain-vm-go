package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mandos-go/orderedjson"
)

func (p *Parser) processLogList(logsRaw oj.OJsonObject) (mj.LogList, error) {
	if IsStar(logsRaw) {
		return mj.LogList{
			IsUnspecified: false,
			IsStar:        true,
		}, nil
	}

	logList, isList := logsRaw.(*oj.OJsonList)
	if !isList {
		return mj.LogList{}, errors.New("unmarshalled logs list is not a list")
	}
	var logEntries []*mj.LogEntry
	var err error
	for _, logRaw := range logList.AsList() {
		logMap, isMap := logRaw.(*oj.OJsonMap)
		if !isMap {
			return mj.LogList{}, errors.New("unmarshalled log entry is not a map")
		}
		logEntry := mj.LogEntry{}
		for _, kvp := range logMap.OrderedKV {
			switch kvp.Key {
			case "address":
				logEntry.Address, err = p.parseCheckBytes(kvp.Value)
				if err != nil {
					return mj.LogList{}, fmt.Errorf("invalid log address: %w", err)
				}
			case "endpoint":
				logEntry.Endpoint, err = p.parseCheckBytes(kvp.Value)
				if err != nil {
					return mj.LogList{}, fmt.Errorf("invalid log identifier: %w", err)
				}
			case "topics":
				logEntry.Topics, err = p.parseCheckValueList(kvp.Value)
				if err != nil {
					return mj.LogList{}, fmt.Errorf("invalid log entry topics: %w", err)
				}
			case "data":
				logEntry.Data, err = p.parseCheckBytes(kvp.Value)
				if err != nil {
					return mj.LogList{}, fmt.Errorf("invalid log data: %w", err)
				}
			default:
				return mj.LogList{}, fmt.Errorf("unknown log field: %s", kvp.Key)
			}
		}
		logEntries = append(logEntries, &logEntry)
	}

	return mj.LogList{
		IsUnspecified: false,
		IsStar:        false,
		List:          logEntries,
	}, nil
}
