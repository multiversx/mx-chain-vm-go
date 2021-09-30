package esdtconvert

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
)

// GetTokenBalanceByName returns the ESDT balance of the account, specified by
// the token name.
func GetTokenBalanceByName(tokenName string, source map[string][]byte) (*big.Int, error) {
	tokenKey := MakeTokenKey([]byte(tokenName), 0)
	return GetTokenBalance(tokenKey, source)
}

// GetTokenBalance returns the ESDT balance of the account, specified by the
// token key.
func GetTokenBalance(tokenKey []byte, source map[string][]byte) (*big.Int, error) {
	tokenData, err := GetTokenData(tokenKey, source)
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

// SetTokenBalance sets the ESDT balance of the account, specified by the token
// key.
func SetTokenBalance(tokenKey []byte, balance *big.Int, source map[string][]byte) error {
	tokenData, err := GetTokenData(tokenKey, source)
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return ErrNegativeValue
	}

	tokenData.Value = balance
	return SetTokenData(tokenKey, tokenData, source)
}

// GetTokenBalanceUint64 returns the ESDT balance of the account, specified by
// the token key.
func GetTokenBalanceUint64(tokenKey []byte, source map[string][]byte) (uint64, error) {
	tokenData, err := GetTokenData(tokenKey, source)
	if err != nil {
		return 0, err
	}

	return tokenData.Value.Uint64(), nil
}

// GetTokenData gets the ESDT information related to a token from the storage of the account.
func GetTokenData(tokenKey []byte, source map[string][]byte) (*esdt.ESDigitalToken, error) {
	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
		Type:  uint32(core.Fungible),
		TokenMetaData: &esdt.MetaData{
			Name:  GetTokenNameFromKey(tokenKey),
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
	tokenRolesKey := MakeTokenRolesKey(tokenName)
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
		if IsTokenKey([]byte(key)) {
			tokenKeys = append(tokenKeys, []byte(key))
		}
	}

	return tokenKeys
}

// MockESDTData groups together all instances of a token (same token name, different nonces).
type MockESDTData struct {
	TokenIdentifier []byte
	Instances       []*esdt.ESDigitalToken
	LastNonce       uint64
	Roles           [][]byte
}

// GetFullMockESDTData returns the information about all the ESDT tokens held by the account.
func GetFullMockESDTData(source map[string][]byte) (map[string]*MockESDTData, error) {
	resultMap := make(map[string]*MockESDTData)
	for key := range source {
		storageKeyBytes := []byte(key)
		if IsTokenKey(storageKeyBytes) {
			tokenName, tokenInstance, err := loadMockESDTDataInstance(storageKeyBytes, source)
			if err != nil {
				return nil, err
			}
			if tokenInstance.Value.Sign() > 0 {
				resultObj := getOrCreateMockESDTData(tokenName, resultMap)
				resultObj.Instances = append(resultObj.Instances, tokenInstance)
			}
		} else if IsNonceKey(storageKeyBytes) {
			tokenName := key[len(ESDTNonceKeyPrefix):]
			resultObj := getOrCreateMockESDTData(tokenName, resultMap)
			resultObj.LastNonce = big.NewInt(0).SetBytes(source[key]).Uint64()
		} else if IsRoleKey(storageKeyBytes) {
			tokenName := key[len(ESDTRoleKeyPrefix):]
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
	tokenInstance, err := GetTokenData(tokenKey, source)
	if err != nil {
		return "", nil, err
	}

	tokenNameFromKey := GetTokenNameFromKey(tokenKey)

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
