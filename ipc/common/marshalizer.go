package common

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func marshalMessage(data interface{}) ([]byte, error) {
	return marshalJSON(data)
}

func unmarshalMessage(dataBytes []byte, data interface{}) error {
	return unmarshalJSON(dataBytes, data)
}

func marshalJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func unmarshalJSON(dataBytes []byte, data interface{}) error {
	return json.Unmarshal(dataBytes, data)
}

func marshalGob(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func unmarshalGob(dataBytes []byte, data interface{}) error {
	buffer := bytes.NewBuffer(dataBytes)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	return nil
}
