package mock

// MemoryMock mocks the linear memory of a Wasmer instance and is used by the
// InstanceMock.
type MemoryMock struct {
	Pages    uint32
	PageSize uint32
	Contents []byte
}

// NewMemoryMock creates a new MemoryMock instance
func NewMemoryMock() *MemoryMock {
	memory := &MemoryMock{}
	memory.Pages = 2
	memory.PageSize = 65536
	memory.initMemory()
	return memory
}

// Length mocked method
func (memory *MemoryMock) Length() uint32 {
	return uint32(len(memory.Contents))
}

// Data mocked method
func (memory *MemoryMock) Data() []byte {
	return memory.Contents
}

// Grow mocked method
func (memory *MemoryMock) Grow(pages uint32) error {
	newPages := make([]byte, pages*memory.PageSize)
	memory.Contents = append(memory.Contents, newPages...)
	return nil
}

// Destroy mocked method
func (memory *MemoryMock) Destroy() {
	memory.Contents = nil
}

// initMemory mocked method
func (memory *MemoryMock) initMemory() {
	memory.Contents = make([]byte, memory.Pages*memory.PageSize)
}
