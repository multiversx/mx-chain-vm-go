package logger

import (
	"encoding/json"
)

func marshalLog(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func unmarshalLog(dataBytes []byte, data interface{}) error {
	return json.Unmarshal(dataBytes, data)
}
