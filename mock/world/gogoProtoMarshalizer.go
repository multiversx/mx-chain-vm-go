package worldmock

import (
	"errors"
	"fmt"

	gproto "github.com/gogo/protobuf/proto"
	proto "github.com/golang/protobuf/proto" //nolint TODO:deprecated
)

// GogoProtoObj groups the necessary of a gogo protobuf marshalizeble object
type GogoProtoObj interface {
	gproto.Marshaler
	gproto.Unmarshaler
	proto.Message
}

// ErrMarshallingProto is raised when the object does not implement proto.Message
var ErrMarshallingProto = errors.New("can not serialize the object")

// ErrUnmarshallingProto is raised when the object that needs to be unmarshaled does not implement proto.Message
var ErrUnmarshallingProto = errors.New("obj does not implement proto.Message")

// GogoProtoMarshalizer implements marshaling with protobuf
type GogoProtoMarshalizer struct {
}

// Marshal does the actual serialization of an object
// The object to be serialized must implement the gogoProtoObj interface
func (x *GogoProtoMarshalizer) Marshal(obj interface{}) ([]byte, error) {
	if msg, ok := obj.(GogoProtoObj); ok {
		return msg.Marshal()
	}
	return nil, fmt.Errorf("%T, %w", obj, ErrMarshallingProto)
}

// Unmarshal does the actual deserialization of an object
// The object to be deserialized must implement the gogoProtoObj interface
func (x *GogoProtoMarshalizer) Unmarshal(obj interface{}, buff []byte) error {
	if msg, ok := obj.(GogoProtoObj); ok {
		msg.Reset()
		return msg.Unmarshal(buff)
	}

	return fmt.Errorf("%T, %w", obj, ErrUnmarshallingProto)
}

// IsInterfaceNil returns true if there is no value under the interface
func (x *GogoProtoMarshalizer) IsInterfaceNil() bool {
	return x == nil
}
