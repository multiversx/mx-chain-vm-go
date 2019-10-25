package arwen

import (
	"fmt"
	"math/big"

	mbig "github.com/ElrondNetwork/managed-big-int"
)

func (host *vmContext) initBigIntContainer() {
	host.nextAllocMemIndex = -1
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

func (host *vmContext) BigByteLength(reference BigIntHandle) int32 {
	return int32(host.bigIntContainer.ByteLen(host.bigIntHandles[reference]))
}

func (host *vmContext) BigGetBytes(reference BigIntHandle) []byte {
	return host.bigIntContainer.GetBytes(host.bigIntHandles[reference])
}

func (host *vmContext) GetNextAllocMemIndex(allocSize int32, totalMemSize int32) (newIndex int32) {
	if host.nextAllocMemIndex == -1 {
		// first allocation
		// to be sure that we don't overwrite anything in memory,
		// we start allocating immediately beyond the current memory size
		// and force it to grow
		host.nextAllocMemIndex = totalMemSize
	}
	newIndex = host.nextAllocMemIndex
	host.nextAllocMemIndex += allocSize
	return
}

func (host *vmContext) BigSetBytes(destination BigIntHandle, bytes []byte) {
	host.bigIntContainer.SetBytes(host.bigIntHandles[destination], bytes)
}

func (host *vmContext) BigAdd(destination, op1, op2 BigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Add(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) BigSub(destination, op1, op2 BigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Sub(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) ReturnBigInt(reference BigIntHandle) {
	host.returnData = host.bigIntContainer.Get(host.bigIntHandles[reference])
}

func (host *vmContext) ReturnInt32(value int32) {
	host.returnData = big.NewInt(int64(value))
}

func (host *vmContext) DebugPrintBig(value BigIntHandle) {
	bi := host.bigIntContainer.GetUnsafe(host.bigIntHandles[value])
	fmt.Printf(">>> Big Int: %d\n", bi)
}
