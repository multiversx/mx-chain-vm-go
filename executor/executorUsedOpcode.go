package executor

type OpcodeUsed uint32

const (
	OpcodeUsedMemoryCopy OpcodeUsed = iota
	OpcodeUsedMemoryFill
)
