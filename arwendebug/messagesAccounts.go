package arwendebug

import (
	"encoding/hex"
	"math/big"
)

// CreateAccountRequest is a CLI / REST request message
type CreateAccountRequest struct {
	RequestBase
	Address string
	Balance string
	Nonce   uint64
}

func (request *CreateAccountRequest) getAddress() ([]byte, error) {
	addressAsHex := request.Address
	addressAsBytes, err := hex.DecodeString(addressAsHex)
	if err != nil {
		return nil, NewRequestErrorMessageInner("invalid account address", err)
	}

	return addressAsBytes, nil
}

func (request *CreateAccountRequest) getBalance() (*big.Int, error) {
	balance, ok := big.NewInt(0).SetString(request.Balance, 10)
	if !ok {
		return nil, NewRequestError("invalid balance")
	}

	return balance, nil
}

// CreateAccountResponse is a CLI / REST response message
type CreateAccountResponse struct {
	Account *Account
}
