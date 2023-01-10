package host

import (
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/cryptoapi"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooks"
	"github.com/multiversx/mx-chain-vm-v1_4-go/vmhost/vmhooksmeta"
)

// PopulateAllImports fills a function container with all existing EI functions.
func PopulateAllImports(imports vmhooksmeta.EIFunctionReceiver) error {
	err := vmhooks.ElrondEIImports(imports)
	if err != nil {
		return err
	}

	err = vmhooks.BigIntImports(imports)
	if err != nil {
		return err
	}

	err = vmhooks.BigFloatImports(imports)
	if err != nil {
		return err
	}

	err = vmhooks.SmallIntImports(imports)
	if err != nil {
		return err
	}

	err = vmhooks.ManagedEIImports(imports)
	if err != nil {
		return err
	}

	err = vmhooks.ManagedBufferImports(imports)
	if err != nil {
		return err
	}

	err = cryptoapi.CryptoImports(imports)
	if err != nil {
		return err
	}

	return nil
}
