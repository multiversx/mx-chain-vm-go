package marshaling

import (
	"encoding/json"
)

var _ Marshalizer = (*jsonMarshalizer)(nil)

type jsonMarshalizer struct {
}

func (marshalizer *jsonMarshalizer) MarshalItem(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (marshalizer *jsonMarshalizer) UnmarshalItem(dataBytes []byte, data interface{}) error {
	return json.Unmarshal(dataBytes, data)
}
