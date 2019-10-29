package arwen

import (
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
	host.bigIntHandles[destination] = host.bigIntContainer.Update(
		host.bigIntHandles[destination],
		newValue)
}

func (host *vmContext) BigGet(reference BigIntHandle) *big.Int {
	return host.bigIntContainer.Get(host.bigIntHandles[reference])
}

func (host *vmContext) BigByteLength(reference BigIntHandle) int32 {
	return int32(host.bigIntContainer.ByteLen(host.bigIntHandles[reference]))
}

func (host *vmContext) BigGetBytes(reference BigIntHandle) []byte {
	return host.bigIntContainer.GetBytes(host.bigIntHandles[reference])
}

func (host *vmContext) BigSetBytes(destination BigIntHandle, bytes []byte) {
	host.bigIntHandles[destination] = host.bigIntContainer.SetBytes(
		host.bigIntHandles[destination],
		bytes)
}

func (host *vmContext) BigIsInt64(destination BigIntHandle) bool {
	return host.bigIntContainer.IsInt64(host.bigIntHandles[destination])
}

func (host *vmContext) BigGetInt64(destination BigIntHandle) int64 {
	return host.bigIntContainer.ToInt64(host.bigIntHandles[destination])
}

func (host *vmContext) BigSetInt64(destination BigIntHandle, value int64) {
	host.bigIntHandles[destination] = host.bigIntContainer.SetInt64(
		host.bigIntHandles[destination],
		value)
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

func (host *vmContext) BigMul(destination, op1, op2 BigIntHandle) {
	host.bigIntHandles[destination] = host.bigIntContainer.Mul(
		host.bigIntHandles[destination],
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) BigCmp(op1, op2 BigIntHandle) int {
	return host.bigIntContainer.Cmp(
		host.bigIntHandles[op1],
		host.bigIntHandles[op2],
	)
}

func (host *vmContext) ReturnBigInt(reference BigIntHandle) {
	host.returnData = append(host.returnData, host.bigIntContainer.Get(host.bigIntHandles[reference]))
}
