package worldmock

import (
	"bytes"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
)

func MakeTokenKey(tokenName []byte) []byte {
	tokenKey := append(ESDTKeyPrefix, tokenName...)
	return tokenKey
}

func IsTokenKey(key []byte) bool {
	if len(key) <= len(ESDTKeyPrefix) {
		return false
	}

	if !bytes.HasPrefix(key, ESDTKeyPrefix) {
		return false
	}

	return true
}

func GetTokenNameFromKey(key []byte) []byte {
	return key[len(ESDTKeyPrefix):]
}

func (a *Account) GetTokenBalance(tokenKey []byte) (*big.Int, error) {
	tokenData, err := a.GetTokenData(tokenKey)
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

func (a *Account) SetTokenBalance(tokenKey []byte, balance *big.Int) error {
	tokenData, err := a.GetTokenData(tokenKey)
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return data.ErrNegativeValue
	}

	tokenData.Value = balance
	return a.SetTokenData(tokenKey, tokenData)
}

func (a *Account) GetTokenData(tokenKey []byte) (*esdt.ESDigitalToken, error) {
	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
		Type:  uint32(core.Fungible),
		TokenMetaData: &esdt.MetaData{
			Name:  GetTokenNameFromKey(tokenKey),
			Nonce: 0,
		},
	}

	marshaledData, err := a.DataTrieTracker().RetrieveValue(tokenKey)
	if err != nil || len(marshaledData) == 0 {
		return esdtData, nil
	}

	err = WorldMarshalizer.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	return esdtData, nil
}

func (a *Account) SetTokenData(tokenKey []byte, tokenData *esdt.ESDigitalToken) error {
	marshaledData, err := WorldMarshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}

	return a.DataTrieTracker().SaveKeyValue(tokenKey, marshaledData)
}

func (a *Account) GetAllTokenData() (map[string]*esdt.ESDigitalToken, error) {
	tokenDataMap := make(map[string]*esdt.ESDigitalToken)
	for _, tokenKey := range a.GetTokenKeys() {
		tokenData, err := a.GetTokenData(tokenKey)
		if err != nil {
			return nil, err
		}

		tokenDataMap[string(tokenData.TokenMetaData.Name)] = tokenData
	}

	return tokenDataMap, nil
}

func (a *Account) GetTokenKeys() [][]byte {
	tokenKeys := make([][]byte, 0)
	for key := range a.Storage {
		if IsTokenKey([]byte(key)) {
			tokenKeys = append(tokenKeys, []byte(key))
		}
	}

	return tokenKeys
}

func (a *Account) GetTokenNames() ([][]byte, error) {
	tokenKeys := a.GetTokenKeys()
	tokenNames := make([][]byte, len(tokenKeys))

	for i := 0; i < len(tokenKeys); i++ {
		tokenKey := tokenKeys[i]
		tokenData, err := a.GetTokenData(tokenKey)
		if err != nil {
			return nil, err
		}

		tokenNames = append(tokenNames, tokenData.TokenMetaData.Name)
	}

	return tokenNames, nil
}
