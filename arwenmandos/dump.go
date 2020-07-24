package arwenmandos

import (
	"fmt"

	worldhook "github.com/ElrondNetwork/elrond-vm-util/mock-hook-blockchain"
)

func dumpWorld(world *worldhook.BlockchainHookMock) error {
	fmt.Print("world state dump:\n")
	for addr, account := range world.AcctMap {
		fmt.Printf("\t%s\n", byteArrayPretty([]byte(addr)))
		fmt.Printf("\t\tnonce: %d\n", account.Nonce)
		fmt.Printf("\t\tbalance: %d\n", account.Balance)
		if len(account.Storage) > 0 {
			fmt.Print("\t\tstorage:\n")
			for key, value := range account.Storage {
				fmt.Printf("\t\t\t%s => %s\n", byteArrayPretty([]byte(key)), byteArrayPretty(value))
			}

		}
	}
	return nil
}
