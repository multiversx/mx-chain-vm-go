package contexts

// Cacher provides caching services
type Cacher interface {
	// Clear is used to completely clear the cache.
	Clear()
	// Put adds a value to the cache.  Returns true if an eviction occurred.
	Put(key []byte, value interface{}, sizeInBytes int) (evicted bool)
	// Get looks up a key's value from the cache.
	Get(key []byte) (value interface{}, ok bool)
	// Has checks if a key is in the cache, without updating the
	// recent-ness or deleting it for being stale.
	Has(key []byte) bool
	// Peek returns the key value (or undefined if not found) without updating
	// the "recently used"-ness of the key.
	Peek(key []byte) (value interface{}, ok bool)
	// HasOrAdd checks if a key is in the cache without updating the
	// recent-ness or deleting it for being stale, and if not adds the value.
	HasOrAdd(key []byte, value interface{}, sizeInBytes int) (has, added bool)
	// Remove removes the provided key from the cache.
	Remove(key []byte)
	// Keys returns a slice of the keys in the cache, from oldest to newest.
	Keys() [][]byte
	// Len returns the number of items in the cache.
	Len() int
	// SizeInBytesContained returns the size in bytes of all contained elements
	SizeInBytesContained() uint64
	// MaxSize returns the maximum number of items which can be stored in the cache.
	MaxSize() int
	// RegisterHandler registers a new handler to be called when a new data is added
	RegisterHandler(handler func(key []byte, value interface{}), id string)
	// UnRegisterHandler deletes the handler from the list
	UnRegisterHandler(id string)
	// Close closes the underlying temporary db if the cacher implementation has one,
	// otherwise it does nothing
	Close() error
	// IsInterfaceNil returns true if there is no value under the interface
	IsInterfaceNil() bool
}
