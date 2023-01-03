package contexts

// Cacher provides caching services
type Cacher interface {
	Clear()
	Put(key []byte, value interface{}, sizeInBytes int) (evicted bool)
	Get(key []byte) (value interface{}, ok bool)
	Has(key []byte) bool
	Peek(key []byte) (value interface{}, ok bool)
	HasOrAdd(key []byte, value interface{}, sizeInBytes int) (has, added bool)
	Remove(key []byte)
	Keys() [][]byte
	Len() int
	SizeInBytesContained() uint64
	MaxSize() int
	RegisterHandler(handler func(key []byte, value interface{}), id string)
	UnRegisterHandler(id string)
	Close() error
	IsInterfaceNil() bool
}
