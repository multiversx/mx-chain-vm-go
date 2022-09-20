package esdtconvert

import (
	"math/big"

	mj "github.com/ElrondNetwork/wasm-vm/mandos-go/model"
	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data/esdt"
	"github.com/ElrondNetwork/elrond-vm-common/builtInFunctions"
)

func MakeESDTUserMetadataBytes(frozen bool) []byte {
	metadata := &builtInFunctions.ESDTUserMetadata{
		Frozen: frozen,
	}

	return metadata.ToBytes()
}

func WriteMockESDTToStorage(esdtData map[string]*MockESDTData, destination map[string][]byte) error {
	for _, token := range esdtData {
		tokenIdentifier := token.TokenIdentifier
		for _, instance := range token.Instances {
			tokenNonce := instance.TokenMetaData.Nonce
			tokenKey := makeTokenKey(tokenIdentifier, tokenNonce)
			err := setTokenDataByKey(tokenKey, instance, destination)
			if err != nil {
				return err
			}
		}
		err := SetLastNonce(tokenIdentifier, token.LastNonce, destination)
		if err != nil {
			return err
		}

		rolesAsStrings := make([]string, len(token.Roles))
		for i, roleBytes := range token.Roles {
			rolesAsStrings[i] = string(roleBytes)
		}
		err = SetTokenRolesAsStrings(tokenIdentifier, rolesAsStrings, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func WriteMandosESDTToStorage(esdtData []*mj.ESDTData, destination map[string][]byte) error {
	for _, mandosESDTData := range esdtData {
		tokenIdentifier := mandosESDTData.TokenIdentifier.Value
		isFrozen := mandosESDTData.Frozen.Value > 0
		for _, instance := range mandosESDTData.Instances {
			tokenNonce := instance.Nonce.Value
			tokenKey := makeTokenKey(tokenIdentifier, tokenNonce)
			tokenBalance := instance.Balance.Value
			var uris [][]byte
			for _, jsonUri := range instance.Uris.Values {
				uris = append(uris, jsonUri.Value)
			}
			tokenData := &esdt.ESDigitalToken{
				Value:      tokenBalance,
				Type:       uint32(core.Fungible),
				Properties: MakeESDTUserMetadataBytes(isFrozen),
				TokenMetaData: &esdt.MetaData{
					Name:       []byte{},
					Nonce:      tokenNonce,
					Creator:    instance.Creator.Value,
					Royalties:  uint32(instance.Royalties.Value),
					Hash:       instance.Hash.Value,
					URIs:       uris,
					Attributes: instance.Attributes.Value,
				},
			}
			err := setTokenDataByKey(tokenKey, tokenData, destination)
			if err != nil {
				return err
			}
		}
		err := SetLastNonce(tokenIdentifier, mandosESDTData.LastNonce.Value, destination)
		if err != nil {
			return err
		}
		err = SetTokenRolesAsStrings(tokenIdentifier, mandosESDTData.Roles, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetTokenData sets the ESDT information related to a token into the storage of the account.
func setTokenDataByKey(tokenKey []byte, tokenData *esdt.ESDigitalToken, destination map[string][]byte) error {
	marshaledData, err := esdtDataMarshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}
	destination[string(tokenKey)] = marshaledData
	return nil
}

func SetTokenData(tokenIdentifier []byte, nonce uint64, tokenData *esdt.ESDigitalToken, destination map[string][]byte) error {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	return setTokenDataByKey(tokenKey, tokenData, destination)
}

// SetTokenRoles sets the specified roles to the account, corresponding to the given tokenIdentifier.
func SetTokenRoles(tokenIdentifier []byte, roles [][]byte, destination map[string][]byte) error {
	tokenRolesKey := makeTokenRolesKey(tokenIdentifier)
	tokenRolesData := &esdt.ESDTRoles{
		Roles: roles,
	}

	marshaledData, err := esdtDataMarshalizer.Marshal(tokenRolesData)
	if err != nil {
		return err
	}

	destination[string(tokenRolesKey)] = marshaledData
	return nil
}

// SetTokenRolesAsStrings sets the specified roles to the account, corresponding to the given tokenIdentifier.
func SetTokenRolesAsStrings(tokenIdentifier []byte, rolesAsStrings []string, destination map[string][]byte) error {
	roles := make([][]byte, len(rolesAsStrings))
	for i := 0; i < len(roles); i++ {
		roles[i] = []byte(rolesAsStrings[i])
	}

	return SetTokenRoles(tokenIdentifier, roles, destination)
}

// SetLastNonce writes the last nonce of a specified ESDT into the storage.
func SetLastNonce(tokenIdentifier []byte, lastNonce uint64, destination map[string][]byte) error {
	tokenNonceKey := makeLastNonceKey(tokenIdentifier)
	nonceBytes := big.NewInt(0).SetUint64(lastNonce).Bytes()
	destination[string(tokenNonceKey)] = nonceBytes
	return nil
}

// SetTokenBalance sets the ESDT balance of the account, specified by the token
// key.
func SetTokenBalance(tokenIdentifier []byte, nonce uint64, balance *big.Int, destination map[string][]byte) error {
	tokenKey := makeTokenKey(tokenIdentifier, nonce)
	tokenData, err := getTokenDataByKey(tokenKey, destination, make(map[string][]byte))
	if err != nil {
		return err
	}

	if balance.Sign() < 0 {
		return errNegativeValue
	}

	tokenData.Value = balance
	return setTokenDataByKey(tokenKey, tokenData, destination)
}
