package arwenmandos

import (
	"fmt"
	"sort"

	vr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/valuereconstructor"
)

// DumpWorld prints the state of the MockWorld to stdout.
func (ae *ArwenTestExecutor) DumpWorld() error {
	fmt.Print("world state dump:\n")

	for addr, account := range ae.World.AcctMap {
		fmt.Printf("\t%s\n", ae.valueReconstructor.Reconstruct([]byte(addr), vr.AddressHint))
		fmt.Printf("\t\tnonce: %d\n", account.Nonce)
		fmt.Printf("\t\tbalance: %d\n", account.Balance)

		if len(account.Storage) > 0 {
			var keys []string
			for key := range account.Storage {
				keys = append(keys, key)
			}

			fmt.Print("\t\tstorage:\n")
			sort.Strings(keys)
			for _, key := range keys {
				value := account.Storage[key]
				if len(value) > 0 {
					fmt.Printf("\t\t\t%s => %s\n",
						ae.valueReconstructor.Reconstruct([]byte(key), vr.NoHint),
						ae.valueReconstructor.Reconstruct(value, vr.NoHint))
				}
			}
		}
	}

	return nil
}
