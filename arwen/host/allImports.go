package host

import (
	"github.com/ElrondNetwork/wasm-vm/arwen/cryptoapi"
	"github.com/ElrondNetwork/wasm-vm/arwen/elrondapi"
	"github.com/ElrondNetwork/wasm-vm/arwen/elrondapimeta"
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
