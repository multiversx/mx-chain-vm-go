package orderedjson

import (
	"strings"
)

// OJsonObject is an ordered JSON tree object interface.
type OJsonObject interface {
	writeJSON(sb *strings.Builder, indent int)
}

// OJsonKeyValuePair is a key-value pair in a JSON map.
// Since this is ordered JSON, maps are really ordered lists of key value pairs.
type OJsonKeyValuePair struct {
	Key   string
	Value OJsonObject
}

// OJsonMap is an ordered map, actually a list of key value pairs.
type OJsonMap struct {
	KeySet    map[string]bool
	OrderedKV []*OJsonKeyValuePair
}

// OJsonList is a JSON list.
type OJsonList []OJsonObject

// OJsonString is a JSON string value.
type OJsonString struct {
	Value string
}

// OJsonBool is a JSON bool value.
type OJsonBool bool

// NewMap is a create new ordered "map" instance.
func NewMap() *OJsonMap {
	KeySet := make(map[string]bool)
	return &OJsonMap{KeySet: KeySet, OrderedKV: nil}
}

// Put puts into map. Does nothing if key exists in map.
func (j *OJsonMap) Put(key string, value OJsonObject) {
	_, alreadyInserted := j.KeySet[key]
	if !alreadyInserted {
		j.KeySet[key] = true
		keyValuePair := &OJsonKeyValuePair{Key: key, Value: value}
		j.OrderedKV = append(j.OrderedKV, keyValuePair)
	}
}

// Size yields the size of ordered map.
func (j *OJsonMap) Size() int {
	return len(j.OrderedKV)
}

// RefreshKeySet recreates the key set from the key value pairs.
func (j *OJsonMap) RefreshKeySet() {
	j.KeySet = make(map[string]bool)
	for _, kv := range j.OrderedKV {
		j.KeySet[kv.Key] = true
	}
}

// AsList converts a JSON list to a slice of objects.
func (j *OJsonList) AsList() []OJsonObject {
	return []OJsonObject(*j)
}
