package arwendebug

import (
	"math/big"

	worldmock "github.com/multiversx/mx-chain-vm-v1_4-go/mock/world"
)

// CreateAccountRequest is a CLI / REST request message
type CreateAccountRequest struct {
	RequestBase
	AddressHex      string
	Address         []byte
	Balance         string
	BalanceAsBigInt *big.Int
	Nonce           uint64
}

func (request *CreateAccountRequest) digest() error {
	err := request.RequestBase.digest()
	if err != nil {
		return err
	}

	if len(request.AddressHex) == 0 {
		return NewRequestErrorMessageInner("empty account address", err)
	}

	request.Address, err = fromHex(request.AddressHex)
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
	Account *worldmock.Account
}
