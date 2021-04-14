package worldmock

import (
	"fmt"

	mvi "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/valueinterpreter"
)

var numDNSAddresses = uint8(0xFF)

func makeDNSAddresses(numAddresses uint8) map[string]struct{} {
	vi := mvi.ValueInterpreter{}

	dnsMap := make(map[string]struct{}, numAddresses)
	for i := uint8(0); i < numAddresses; i++ {
		// using the value interpreter to generate the addresses
		// consistently to how they appear in the DNS mandos tests
		dnsAddress, _ := vi.InterpretString(fmt.Sprintf("sc:dns#%02x", i))
		dnsMap[string(dnsAddress)] = struct{}{}
	}

	return dnsMap
}
