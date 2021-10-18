package esdtconvert

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
)

// MockESDTData groups together all instances of a token (same token name, different nonces).
type MockESDTData struct {
	TokenIdentifier []byte
	Instances       []*esdt.ESDigitalToken
	LastNonce       uint64
	Roles           [][]byte
}

// GetTokenBalance returns the ESDT balance of the account, specified by the
// token key.
func GetTokenBalance(tokenIdentifier []byte, nonce uint64, source map[string][]byte) (*big.Int, error) {
	tokenData, err := GetTokenData(tokenIdentifier, nonce, source)
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

// GetTokenData gets the ESDT information related to a token from the storage of the account.
func GetTokenData(tokenIdentifier []byte, nonce uint64, source map[string][]byte) (*esdt.ESDigitalToken, error) {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	return getTokenDataByKey(tokenKey, source)
}

func getTokenDataByKey(tokenKey []byte, source map[string][]byte) (*esdt.ESDigitalToken, error) {
	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
		Type:  uint32(core.Fungible),
		TokenMetaData: &esdt.MetaData{
			Name:  getTokenNameFromKey(tokenKey),
			Nonce: 0,
		},
	}

	marshaledData := source[string(tokenKey)]
	if len(marshaledData) == 0 {
		return esdtData, nil
	}

	err := esdtDataMarshalizer.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	return esdtData, nil
}

// GetTokenRoles returns the roles of the account for the specified tokenName.
func GetTokenRoles(tokenName []byte, source map[string][]byte) ([][]byte, error) {
	tokenRolesKey := makeTokenRolesKey(tokenName)
	tokenRolesData := &esdt.ESDTRoles{
		Roles: make([][]byte, 0),
	}

	marshaledData := source[string(tokenRolesKey)]
	if len(marshaledData) == 0 {
		return tokenRolesData.Roles, nil
	}

	err := esdtDataMarshalizer.Unmarshal(tokenRolesData, marshaledData)
	if err != nil {
		return nil, err
	}

	return tokenRolesData.Roles, nil

}

// GetTokenKeys returns the storage keys of all the ESDT tokens owned by the account.
func GetTokenKeys(source map[string][]byte) [][]byte {
	tokenKeys := make([][]byte, 0)
	for key := range source {
		if isTokenKey([]byte(key)) {
			tokenKeys = append(tokenKeys, []byte(key))
		}
	}

	return tokenKeys
}

// GetFullMockESDTData returns the information about all the ESDT tokens held by the account.
func GetFullMockESDTData(source map[string][]byte) (map[string]*MockESDTData, error) {
	resultMap := make(map[string]*MockESDTData)
	for key := range source {
		storageKeyBytes := []byte(key)
		if isTokenKey(storageKeyBytes) {
			tokenName, tokenInstance, err := loadMockESDTDataInstance(storageKeyBytes, source)
			if err != nil {
				return nil, err
			}
			if tokenInstance.Value.Sign() > 0 {
				resultObj := getOrCreateMockESDTData(tokenName, resultMap)
				resultObj.Instances = append(resultObj.Instances, tokenInstance)
			}
		} else if isNonceKey(storageKeyBytes) {
			tokenName := key[len(esdtNonceKeyPrefix):]
			resultObj := getOrCreateMockESDTData(tokenName, resultMap)
			resultObj.LastNonce = big.NewInt(0).SetBytes(source[key]).Uint64()
		} else if isRoleKey(storageKeyBytes) {
			tokenName := key[len(esdtRoleKeyPrefix):]
			roles, err := GetTokenRoles([]byte(tokenName), source)
			if err != nil {
				return nil, err
			}
			resultObj := getOrCreateMockESDTData(tokenName, resultMap)
			resultObj.Roles = roles
		}
	}

	return resultMap, nil
}

// loads and prepared the ESDT instance
func loadMockESDTDataInstance(tokenKey []byte, source map[string][]byte) (string, *esdt.ESDigitalToken, error) {
	tokenInstance, err := getTokenDataByKey(tokenKey, source)
	if err != nil {
		return "", nil, err
	}

	tokenNameFromKey := getTokenNameFromKey(tokenKey)

	var tokenName string
	if tokenInstance.TokenMetaData == nil || tokenInstance.TokenMetaData.Nonce == 0 {
		// ESDT, no nonce in the key
		tokenInstance.TokenMetaData = &esdt.MetaData{
			Name:  tokenNameFromKey,
			Nonce: 0,
		}
		tokenName = string(tokenNameFromKey)
	} else {
		nonceAsBytes := big.NewInt(0).SetUint64(tokenInstance.TokenMetaData.Nonce).Bytes()
		tokenNameLen := len(tokenNameFromKey) - len(nonceAsBytes)

		if !bytes.Equal(nonceAsBytes, tokenNameFromKey[tokenNameLen:]) {
			return "", nil, errors.New("invalid key for NFT (key does not end in nonce)")
		}

		tokenName = string(tokenNameFromKey[:tokenNameLen])
	}

	return tokenName, tokenInstance, nil
}

func getOrCreateMockESDTData(tokenName string, resultMap map[string]*MockESDTData) *MockESDTData {
	resultObj := resultMap[tokenName]
	if resultObj == nil {
		resultObj = &MockESDTData{
			TokenIdentifier: []byte(tokenName),
			Instances:       nil,
			LastNonce:       0,
			Roles:           nil,
		}
		resultMap[tokenName] = resultObj
	}
	return resultObj
}
