package arwen

import (
	"fmt"
	"math/big"

	mbig "github.com/ElrondNetwork/managed-big-int"
)

func (host *vmContext) initBigIntContainer() {
	host.bigIntHandles = nil
	host.bigIntContainer = mbig.NewBigIntContainer()
}

func (host *vmContext) bigInsert(bi *big.Int) BigIntHandle {
	handler := host.bigIntContainer.Insert(bi)
	newIndex := BigIntHandle(len(host.bigIntHandles))
	host.bigIntHandles = append(host.bigIntHandles, handler)
	return newIndex
}

func (host *vmContext) BigInsertInt64(smallValue int64) BigIntHandle {
	return host.bigInsert(big.NewInt(smallValue))
}

func (host *vmContext) BigUpdate(destination BigIntHandle, newValue *big.Int) {
	host.bigIntContainer.Update(host.bigIntHandles[destination], newValue)
}

func (host *vmContext) BigAdd(destination, op1, op2 BigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Add(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) BigSub(destination, op1, op2 BigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Add(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) DebugPrintBig(value BigIntHandle) {
	bi := host.bigIntContainer.GetUnsafe(host.bigIntHandles[value])
	fmt.Printf(">>> Big Int: %d\n", bi)
}
