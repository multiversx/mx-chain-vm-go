package arwenmandos

import (
	"fmt"
	"sort"
	"strings"

	er "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/expression/reconstructor"
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	mjwrite "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/write"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
	"github.com/ElrondNetwork/elrond-go/core"
)

const includeElrondProtectedStorage = false

func (ae *ArwenTestExecutor) convertMockAccountToMandosFormat(account *worldmock.Account) (*mj.Account, error) {
	var storageKeys []string
	for storageKey := range account.Storage {
		storageKeys = append(storageKeys, storageKey)
	}

	sort.Strings(storageKeys)
	var storageKvps []*mj.StorageKeyValuePair
	for _, storageKey := range storageKeys {
		storageValue := account.Storage[storageKey]
		includeKey := includeElrondProtectedStorage || !strings.HasPrefix(storageKey, core.ElrondProtectedKeyPrefix)
		if includeKey && len(storageValue) > 0 {
			storageKvps = append(storageKvps, &mj.StorageKeyValuePair{
				Key: mj.JSONBytesFromString{
					Value:    []byte(storageKey),
					Original: ae.exprReconstructor.Reconstruct([]byte(storageKey), er.NoHint),
				},
				Value: mj.JSONBytesFromTree{
					Value:    storageValue,
					Original: &oj.OJsonString{Value: ae.exprReconstructor.Reconstruct(storageValue, er.NoHint)},
				},
			})
		}
	}

	tokenData, err := account.GetFullMockESDTData()
	if err != nil {
		return nil, err
	}
	var esdtNames []string
	for esdtName := range tokenData {
		esdtNames = append(esdtNames, esdtName)
	}
	sort.Strings(esdtNames)
	var mandosESDT []*mj.ESDTData
	for _, esdtName := range esdtNames {
		esdtObj := tokenData[esdtName]

		var mandosRoles []string
		for _, mockRoles := range esdtObj.Roles {
			mandosRoles = append(mandosRoles, string(mockRoles))
		}

		var mandosInstances []*mj.ESDTInstance
		for _, mockInstance := range esdtObj.Instances {
			var uri mj.JSONBytesFromTree
			if len(mockInstance.TokenMetaData.URIs) > 0 {
				uri = mj.JSONBytesFromTree{
					Value:    mockInstance.TokenMetaData.URIs[0],
					Original: &oj.OJsonString{Value: ae.exprReconstructor.Reconstruct(mockInstance.TokenMetaData.URIs[0], er.NoHint)},
				}
			}

			mandosInstances = append(mandosInstances, &mj.ESDTInstance{
				Nonce: mj.JSONUint64{
					Value:    mockInstance.TokenMetaData.Nonce,
					Original: ae.exprReconstructor.ReconstructFromUint64(mockInstance.TokenMetaData.Nonce),
				},
				Balance: mj.JSONBigInt{
					Value:    mockInstance.Value,
					Original: ae.exprReconstructor.ReconstructFromBigInt(mockInstance.Value),
				},
				Uri: uri,
			})
		}

		mandosESDT = append(mandosESDT, &mj.ESDTData{
			TokenIdentifier: mj.JSONBytesFromString{
				Value:    esdtObj.TokenIdentifier,
				Original: ae.exprReconstructor.Reconstruct(esdtObj.TokenIdentifier, er.StrHint),
			},
			Instances: mandosInstances,
			LastNonce: mj.JSONUint64{
				Value:    esdtObj.LastNonce,
				Original: ae.exprReconstructor.ReconstructFromUint64(esdtObj.LastNonce),
			},
			Roles: mandosRoles,
		})
	}

	return &mj.Account{
		Address: mj.JSONBytesFromString{
			Value:    account.Address,
			Original: ae.exprReconstructor.Reconstruct([]byte(account.Address), er.AddressHint),
		},
		Nonce: mj.JSONUint64{
			Value:    account.Nonce,
			Original: ae.exprReconstructor.ReconstructFromUint64(account.Nonce),
		},
		Balance: mj.JSONBigInt{
			Value:    account.Balance,
			Original: ae.exprReconstructor.ReconstructFromBigInt(account.Balance),
		},
		Storage:  storageKvps,
		ESDTData: mandosESDT,
	}, nil
}

// DumpWorld prints the state of the MockWorld to stdout.
func (ae *ArwenTestExecutor) DumpWorld() error {
	fmt.Print("world state dump:\n")
	var mandosAccounts []*mj.Account

	for _, account := range ae.World.AcctMap {
		mandosAccount, err := ae.convertMockAccountToMandosFormat(account)
		if err != nil {
			return err
		}
		mandosAccounts = append(mandosAccounts, mandosAccount)
	}

	ojAccount := mjwrite.AccountsToOJ(mandosAccounts)
	s := oj.JSONString(ojAccount)
	fmt.Println(s)

	return nil
}
