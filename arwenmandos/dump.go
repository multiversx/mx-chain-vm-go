package arwenmandos

import (
	"fmt"
	"sort"

	worldhook "github.com/ElrondNetwork/arwen-wasm-vm/mock/world"
)

func dumpWorld(world *worldhook.MockWorld) error {
	fmt.Print("world state dump:\n")
	for addr, account := range world.AcctMap {
		fmt.Printf("\t%s\n", byteArrayPretty([]byte(addr)))
		fmt.Printf("\t\tnonce: %d\n", account.Nonce)
		fmt.Printf("\t\tbalance: %d\n", account.Balance)
		if len(account.Storage) > 0 {
			var keys []string
			for key := range account.Storage {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			fmt.Print("\t\tstorage:\n")
			for _, key := range keys {
				value := account.Storage[key]
				if len(value) > 0 {
					fmt.Printf("\t\t\t%s => %s\n", byteArrayPretty([]byte(key)), byteArrayPretty(value))
				}
			}

		}
	}
	return nil
}
