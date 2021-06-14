package fuzzForwarder

import (
	"fmt"
)

func (pfe *fuzzExecutor) initData() error {
	pfe.data = &fuzzData{
		mainCallerAddress: "address:main_caller",
		numForwarders:     5,
		programmedCalls:   make(map[int][]*programmedCall),
	}

	pfe.world.Clear()
	return nil
}

func (pfe *fuzzExecutor) setUp() error {
	err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0"
				}
			}
		}`,
		pfe.data.mainCallerAddress,
	))
	if err != nil {
		return err
	}
	for i := 1; i <= pfe.data.numForwarders; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "10000000",
					"storage": {},
					"code": "file:forwarder.wasm"
				}
			}
		}`,
			pfe.forwarderAddress(i),
		))
		if err != nil {
			return err
		}
	}
	return nil
}
