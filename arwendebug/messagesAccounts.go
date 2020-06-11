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
	var err error
	var ok bool

	err = request.RequestBase.digest()
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

	// todo move to common
	request.BalanceAsBigInt = big.NewInt(0)
	if len(request.Balance) > 0 {
		_, ok = request.BalanceAsBigInt.SetString(request.Balance, 10)
		if !ok {
			return NewRequestError("invalid value (erd)")
		}
	}

	return nil
}

// CreateAccountResponse is a CLI / REST response message
type CreateAccountResponse struct {
	Account *Account
}
