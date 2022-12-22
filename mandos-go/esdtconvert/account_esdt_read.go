package esdtconvert

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
)

// MockESDTData groups together all instances of a token (same token name, different nonces).
type MockESDTData struct {
	TokenIdentifier []byte
	Instances       []*esdt.ESDigitalToken
	LastNonce       uint64
	Roles           [][]byte
}

const (
	esdtIdentifierSeparator  = "-"
	esdtRandomSequenceLength = 6
)

// GetTokenBalance returns the ESDT balance of the account, specified by the
// token key.
func GetTokenBalance(tokenIdentifier []byte, nonce uint64, source map[string][]byte) (*big.Int, error) {
	tokenData, err := GetTokenData(tokenIdentifier, nonce, source, make(map[string][]byte))
	if err != nil {
		return nil, err
	}

	return tokenData.Value, nil
}

// GetTokenData gets the ESDT information related to a token from the storage of the account.
func GetTokenData(tokenIdentifier []byte, nonce uint64, source map[string][]byte, systemAccStorage map[string][]byte) (*esdt.ESDigitalToken, error) {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	return getTokenDataByKey(tokenKey, source, systemAccStorage)
}

func getTokenDataByKey(tokenKey []byte, source map[string][]byte, systemAccStorage map[string][]byte) (*esdt.ESDigitalToken, error) {
	// default value copied from the protocol
	esdtData := &esdt.ESDigitalToken{
		Value: big.NewInt(0),
	}

	marshaledData := source[string(tokenKey)]
	if len(marshaledData) == 0 {
		return esdtData, nil
	}

	err := esdtDataMarshalizer.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	marshaledData = systemAccStorage[string(tokenKey)]
	if len(marshaledData) == 0 {
		return esdtData, nil
	}
	esdtDataFromSystemAcc := &esdt.ESDigitalToken{}
	err = esdtDataMarshalizer.Unmarshal(esdtDataFromSystemAcc, marshaledData)
	if err != nil {
		return nil, err
	}

	esdtData.TokenMetaData = esdtDataFromSystemAcc.TokenMetaData

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

// GetFullMockESDTData returns the information about all the ESDT tokens held by the account.
func GetFullMockESDTData(source map[string][]byte, systemAccStorage map[string][]byte) (map[string]*MockESDTData, error) {
	resultMap := make(map[string]*MockESDTData)
	for key := range source {
		storageKeyBytes := []byte(key)
		if isTokenKey(storageKeyBytes) {
			tokenName, tokenInstance, err := loadMockESDTDataInstance(storageKeyBytes, source, systemAccStorage)
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

func extractTokenIdentifierAndNonceESDTWipe(args []byte) ([]byte, uint64) {
	argsSplit := bytes.Split(args, []byte(esdtIdentifierSeparator))
	if len(argsSplit) < 2 {
		return args, 0
	}

	if len(argsSplit[1]) <= esdtRandomSequenceLength {
		return args, 0
	}

	identifier := []byte(fmt.Sprintf("%s-%s", argsSplit[0], argsSplit[1][:esdtRandomSequenceLength]))
	nonce := big.NewInt(0).SetBytes(argsSplit[1][esdtRandomSequenceLength:])

	return identifier, nonce.Uint64()
}

// loads and prepared the ESDT instance
func loadMockESDTDataInstance(tokenKey []byte, source map[string][]byte, systemAccStorage map[string][]byte) (string, *esdt.ESDigitalToken, error) {
	tokenInstance, err := getTokenDataByKey(tokenKey, source, systemAccStorage)
	if err != nil {
		return "", nil, err
	}

	tokenNameFromKey := getTokenNameFromKey(tokenKey)
	tokenName, nonce := extractTokenIdentifierAndNonceESDTWipe(tokenNameFromKey)

	if tokenInstance.TokenMetaData == nil {
		tokenInstance.TokenMetaData = &esdt.MetaData{
			Name:  tokenName,
			Nonce: nonce,
		}
	}

	return string(tokenName), tokenInstance, nil
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
