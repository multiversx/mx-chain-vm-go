package fuzzForwarder

import (
	"fmt"
)

func (pfe *fuzzExecutor) nextTxIndex() int {
	pfe.txIndex++
	return pfe.txIndex
}

func (pfe *fuzzExecutor) forwarderAddress(index int) string {
	return fmt.Sprintf("sc:forwarder-%02d", index)
}
