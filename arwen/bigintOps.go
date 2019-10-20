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

func (host *vmContext) BigInsert(smallValue int64) bigIntHandle {
	handler := host.bigIntContainer.Insert(big.NewInt(smallValue))
	newIndex := bigIntHandle(len(host.bigIntHandles))
	host.bigIntHandles = append(host.bigIntHandles, handler)
	return newIndex
}

func (host *vmContext) BigAdd(destination, op1, op2 bigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Add(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) BigSub(destination, op1, op2 bigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Add(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) DebugPrintBig(value bigIntHandle) {
	bi := host.bigIntContainer.GetUnsafe(host.bigIntHandles[value])
	fmt.Printf(">>> Big Int: %d\n", bi)
}
