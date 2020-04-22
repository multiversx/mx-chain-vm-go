package arwendebug

import (
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
)

func fixTestAddress(address string) []byte {
	if len(address) > arwen.AddressLen {
		address = address[0:arwen.AddressLen]
	}

	address = fmt.Sprintf("%032s", address)
	return []byte(address)
}
