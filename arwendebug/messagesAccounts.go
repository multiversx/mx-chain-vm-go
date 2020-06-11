package arwendebug

import (
	"encoding/hex"
	"math/big"
)

// CreateAccountRequest is a CLI / REST request message
type CreateAccountRequest struct {
	RequestBase
	AddressAsHex    string
	AddressAsBytes  []byte
	Balance         string
	BalanceAsBigInt *big.Int
	Nonce           uint64
}

func (request *CreateAccountRequest) digest() error {
	err := request.RequestBase.digest()
	if err != nil {
		return err
	}

	if len(request.AddressAsHex) == 0 {
		return NewRequestErrorMessageInner("empty account address", err)
	}

	request.AddressAsBytes, err = hex.DecodeString(request.AddressAsHex)
	if err != nil {
		return NewRequestErrorMessageInner("invalid account address", err)
	}

	request.BalanceAsBigInt, err = parseValue(request.Balance)
	if err != nil {
		return err
	}

	return nil
}

// CreateAccountResponse is a CLI / REST response message
type CreateAccountResponse struct {
	Account *Account
}
