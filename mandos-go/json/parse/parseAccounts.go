package mandosjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func (p *Parser) parseAccountAddress(addrRaw string) (mj.JSONBytesFromString, error) {
	if len(addrRaw) == 0 {
		return mj.JSONBytesFromString{}, errors.New("missing account address")
	}
	addrBytes, err := p.ValueInterpreter.InterpretString(addrRaw)
	if err == nil && len(addrBytes) != 32 {
		return mj.JSONBytesFromString{}, errors.New("account addressis not 32 bytes in length")
	}
	return mj.NewJSONBytesFromString(addrBytes, addrRaw), err
}

func (p *Parser) processAccount(acctRaw oj.OJsonObject) (*mj.Account, error) {
	acctMap, isMap := acctRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}

	acct := mj.Account{}
	var err error

	for _, kvp := range acctMap.OrderedKV {
		switch kvp.Key {
		case "comment":
			acct.Comment, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account comment: %w", err)
			}
		case "shard":
			acct.Shard, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid shard number: %w", err)
			}
		case "isSmartContract":
			jsonInt, err := p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid boolean isContract: %w", err)
			}
			acct.IsSmartContract = jsonInt.Value != 0
		case "nonce":
			acct.Nonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		case "balance":
			acct.Balance, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, errors.New("invalid account balance")
			}
		case "storage":
			storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
			if !storageOk {
				return nil, errors.New("invalid account storage")
			}
			for _, storageKvp := range storageMap.OrderedKV {
				byteKey, err := p.ValueInterpreter.InterpretString(storageKvp.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid account storage key: %w", err)
				}
				byteVal, err := p.processSubTreeAsByteArray(storageKvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid account storage value: %w", err)
				}
				stElem := mj.StorageKeyValuePair{
					Key:   mj.NewJSONBytesFromString(byteKey, storageKvp.Key),
					Value: byteVal,
				}
				acct.Storage = append(acct.Storage, &stElem)
			}
		case "code":
			acct.Code, err = p.processStringAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account code: %w", err)
			}
		case "asyncCallData":
			acct.AsyncCallData, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid asyncCallData string: %w", err)
			}
		case "esdt":
			esdtMap, esdtOk := kvp.Value.(*oj.OJsonMap)
			if !esdtOk {
				return nil, errors.New("invalid ESDT map")
			}
			for _, esdtKvp := range esdtMap.OrderedKV {
				tickerStr, err := p.ValueInterpreter.InterpretString(esdtKvp.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid esdt ticker: %w", err)
				}
				esdtData, err := p.processESDTData(esdtKvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid esdt value: %w", err)
				}
				esdtData.TokenName = mj.NewJSONBytesFromString(tickerStr, esdtKvp.Key)
				acct.ESDTData = append(acct.ESDTData, esdtData)
			}
		default:
			return nil, fmt.Errorf("unknown account field: %s", kvp.Key)
		}
	}

	return &acct, nil
}

func (p *Parser) processAccountMap(acctMapRaw oj.OJsonObject) ([]*mj.Account, error) {
	var accounts []*mj.Account
	preMap, isPreMap := acctMapRaw.(*oj.OJsonMap)
	if !isPreMap {
		return nil, errors.New("unmarshalled account map object is not a map")
	}
	for _, acctKVP := range preMap.OrderedKV {
		acct, acctErr := p.processAccount(acctKVP.Value)
		if acctErr != nil {
			return nil, acctErr
		}
		acctAddr, hexErr := p.parseAccountAddress(acctKVP.Key)
		if hexErr != nil {
			return nil, hexErr
		}
		acct.Address = acctAddr
		accounts = append(accounts, acct)

	}
	return accounts, nil
}

func (p *Parser) processCheckAccount(acctRaw oj.OJsonObject) (*mj.CheckAccount, error) {
	acctMap, isMap := acctRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled account object is not a map")
	}

	acct := mj.CheckAccount{
		Nonce:         mj.JSONCheckUint64Default(),
		Balance:       mj.JSONCheckBigIntDefault(),
		IgnoreStorage: true,
		Code:          mj.JSONCheckBytesDefault(),
		AsyncCallData: mj.JSONCheckBytesDefault(),
	}
	var err error

	for _, kvp := range acctMap.OrderedKV {
		switch kvp.Key {
		case "comment":
			acct.Comment, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid check account comment: %w", err)
			}
		case "nonce":
			acct.Nonce, err = p.processCheckUint64(kvp.Value)
			if err != nil {
				return nil, errors.New("invalid account nonce")
			}
		case "balance":
			acct.Balance, err = p.processCheckBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, errors.New("invalid account balance")
			}
		case "storage":
			acct.IgnoreStorage = IsStar(kvp.Value)
			if !acct.IgnoreStorage {
				// TODO: convert to a more permissive format
				storageMap, storageOk := kvp.Value.(*oj.OJsonMap)
				if !storageOk {
					return nil, errors.New("invalid account storage")
				}
				for _, storageKvp := range storageMap.OrderedKV {
					byteKey, err := p.ValueInterpreter.InterpretString(storageKvp.Key)
					if err != nil {
						return nil, fmt.Errorf("invalid account storage key: %w", err)
					}
					byteVal, err := p.processSubTreeAsByteArray(storageKvp.Value)
					if err != nil {
						return nil, fmt.Errorf("invalid account storage value: %w", err)
					}
					stElem := mj.StorageKeyValuePair{
						Key:   mj.NewJSONBytesFromString(byteKey, storageKvp.Key),
						Value: byteVal,
					}
					acct.CheckStorage = append(acct.CheckStorage, &stElem)
				}
			}
		case "code":
			acct.Code, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid account code: %w", err)
			}
		case "asyncCallData":
			acct.AsyncCallData, err = p.parseCheckBytes(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid asyncCallData: %w", err)
			}
		case "esdt":
			esdtMap, esdtOk := kvp.Value.(*oj.OJsonMap)
			if !esdtOk {
				return nil, errors.New("invalid ESDT map")
			}
			for _, esdtKvp := range esdtMap.OrderedKV {
				tickerStr, err := p.ValueInterpreter.InterpretString(esdtKvp.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid esdt ticker: %w", err)
				}
				esdtData, err := p.processCheckESDTData(esdtKvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid esdt value: %w", err)
				}
				esdtData.TokenName = mj.NewJSONBytesFromString(tickerStr, esdtKvp.Key)
				acct.ESDTData = append(acct.ESDTData, esdtData)
			}
		default:
			return nil, fmt.Errorf("unknown account field: %s", kvp.Key)
		}
	}

	return &acct, nil
}

func (p *Parser) processCheckAccountMap(acctMapRaw oj.OJsonObject) (*mj.CheckAccounts, error) {
	var checkAccounts = &mj.CheckAccounts{
		OtherAccountsAllowed: false,
		Accounts:             nil,
	}

	preMap, isPreMap := acctMapRaw.(*oj.OJsonMap)
	if !isPreMap {
		return nil, errors.New("unmarshalled check account map object is not a map")
	}
	for _, acctKVP := range preMap.OrderedKV {
		if acctKVP.Key == "+" {
			checkAccounts.OtherAccountsAllowed = true
		} else {
			acct, acctErr := p.processCheckAccount(acctKVP.Value)
			if acctErr != nil {
				return nil, acctErr
			}
			acctAddr, hexErr := p.parseAccountAddress(acctKVP.Key)
			if hexErr != nil {
				return nil, hexErr
			}
			acct.Address = acctAddr
			checkAccounts.Accounts = append(checkAccounts.Accounts, acct)
		}
	}
	return checkAccounts, nil
}
