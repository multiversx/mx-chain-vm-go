package host

import (
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen/cryptoapi"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen/elrondapi"
	"github.com/multiversx/mx-chain-vm-v1_4-go/arwen/elrondapimeta"
)

// PopulateAllImports fills a function container with all existing EI functions.
func PopulateAllImports(imports elrondapimeta.EIFunctionReceiver) error {
	err := elrondapi.ElrondEIImports(imports)
	if err != nil {
		return err
	}

	err = elrondapi.BigIntImports(imports)
	if err != nil {
		return err
	}

	err = elrondapi.BigFloatImports(imports)
	if err != nil {
		return err
	}

	err = elrondapi.SmallIntImports(imports)
	if err != nil {
		return err
	}

	err = elrondapi.ManagedEIImports(imports)
	if err != nil {
		return err
	}

	err = elrondapi.ManagedBufferImports(imports)
	if err != nil {
		return err
	}

	err = cryptoapi.CryptoImports(imports)
	if err != nil {
		return err
	}

	return nil
}
