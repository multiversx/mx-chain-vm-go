package worldmock

import (
	"fmt"

	mei "github.com/ElrondNetwork/wasm-vm/mandos-go/expression/interpreter"
)

var numDNSAddresses = uint8(0xFF)

func makeDNSAddresses(numAddresses uint8) map[string]struct{} {
	ei := mei.ExprInterpreter{}

	dnsMap := make(map[string]struct{}, numAddresses)
	for i := uint8(0); i < numAddresses; i++ {
		// using the value interpreter to generate the addresses
		// consistently to how they appear in the DNS mandos tests
		dnsAddress, _ := ei.InterpretString(fmt.Sprintf("sc:dns#%02x", i))
		dnsMap[string(dnsAddress)] = struct{}{}
	}

	return dnsMap
}
