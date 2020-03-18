package marshaling

// MarshalizerKind is the kind of a message (that is passed between the Node and Arwen)
type MarshalizerKind uint32

const (
	// JSON is a marshalizer kind
	JSON MarshalizerKind = iota
	// Gob is a marshalizer kind
	Gob
)

// Marshalizer deals with messages serialization
type Marshalizer interface {
	MarshalItem(data interface{}) ([]byte, error)
	UnmarshalItem(dataBytes []byte, data interface{}) error
}
