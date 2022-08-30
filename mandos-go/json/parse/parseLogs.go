package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/wasm-vm/mandos-go/model"
	oj "github.com/ElrondNetwork/wasm-vm/mandos-go/orderedjson"
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
	result := mj.LogList{
		IsUnspecified:    false,
		IsStar:           false,
		MoreAllowedAtEnd: false,
		List:             nil,
	}
	var err error
	for _, logRaw := range logList.AsList() {
		switch logItem := logRaw.(type) {
		case *oj.OJsonString:
			if logItem.Value == "+" {
				result.MoreAllowedAtEnd = true
			} else {
				return mj.LogList{}, errors.New("unmarshalled log entry is an invalid string")
			}
		case *oj.OJsonMap:
			if result.MoreAllowedAtEnd {
				return mj.LogList{}, errors.New("log entry ")
			}

			logEntry := mj.LogEntry{}
			for _, kvp := range logItem.OrderedKV {
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
			result.List = append(result.List, &logEntry)
		default:
			return mj.LogList{}, errors.New("log entry should be either string or object")
		}
	}

	return result, nil
}
