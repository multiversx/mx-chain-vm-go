package worldmock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/esdt"
)

func (bf *BuiltinFunctionsWrapper) MakeTokenKey(tokenName []byte) []byte {
	keyPrefix := []byte(core.ElrondProtectedKeyPrefix + core.ESDTKeyIdentifier)

	tokenKey := append(keyPrefix, tokenName...)
	return tokenKey
}

func (bf *BuiltinFunctionsWrapper) GetTokenData(address []byte, tokenKey []byte) (*esdt.ESDigitalToken, error) {
	account := bf.World.AcctMap.GetAccount(address)
	esdtData := &esdt.ESDigitalToken{Value: big.NewInt(0), Type: uint32(core.Fungible)}

	marshaledData, err := account.DataTrieTracker().RetrieveValue(tokenKey)
	if err != nil || len(marshaledData) == 0 {
		return esdtData, nil
	}

	err = bf.Marshalizer.Unmarshal(esdtData, marshaledData)
	if err != nil {
		return nil, err
	}

	return esdtData, nil
}

func (bf *BuiltinFunctionsWrapper) SetTokenData(address []byte, tokenKey []byte, tokenData *esdt.ESDigitalToken) error {
	account := bf.World.AcctMap.GetAccount(address)

	marshaledData, err := bf.Marshalizer.Marshal(tokenData)
	if err != nil {
		return err
	}

	return account.DataTrieTracker().SaveKeyValue(tokenKey, marshaledData)
}
