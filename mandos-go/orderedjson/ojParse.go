package orderedjson

import (
	"bytes"
	"errors"
	"strings"
)

type jsonParserState interface {
}

type jsonParserStateAnyObjPlaceholder struct {
}

type jsonParserStateSingleValue struct {
	buffer       bytes.Buffer
	stringEscape bool
}

type jsonParserStateMap struct {
	currentMap *OJsonMap
}

type jsonStateMapKeyValue struct {
	keyBuffer bytes.Buffer
	state     int // 0=key, 1=':', 2=value
	currentKV OJsonKeyValuePair
}

type jsonParserStateList struct {
	list OJsonList
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\r' || c == '\t'
}

// ParseOrderedJSON parses JSON preserving order in maps
func ParseOrderedJSON(input []byte) (OJsonObject, error) {
	stateStack := &jsonParserStateStack{}
	stateStack.push(&jsonParserStateAnyObjPlaceholder{})
	var pendingResult OJsonObject

	for i, c := range input {
		done := false
		for !done {
			done = true

			if stateStack.size() == 0 {
				if isWhitespace(c) {
					continue
				} else {
					return nil, errors.New("unexpected characters at the end")
				}
			}

			state := stateStack.peek()
			switch specificState := state.(type) {
			case *jsonParserStateAnyObjPlaceholder:
				if pendingResult != nil {
					return nil, errors.New("invalid state")
				}
				if isWhitespace(c) {
					// leading whitespace, ignore
				} else if c == '{' {
					// replace with map state
					stateStack.replaceTop(&jsonParserStateMap{currentMap: NewMap()})
				} else if c == '[' {
					// replace with list state
					stateStack.replaceTop(&jsonParserStateList{})
				} else if c == ']' || c == '}' || c == ',' {
					return nil, errors.New("misplaced character")
				} else {
					// replace with single value
					stateStack.replaceTop(&jsonParserStateSingleValue{})
					done = false
				}
			case *jsonParserStateSingleValue:
				if specificState.buffer.Len() == 0 {
					specificState.stringEscape = (c == '"')
					specificState.buffer.WriteByte(c)
				} else {
					prevChar := input[i-1]
					if specificState.stringEscape {
						specificState.buffer.WriteByte(c)
						if c == '"' && prevChar != '\\' {
							stateStack.pop()
							var err error
							pendingResult, err = specificState.finalize()
							if err != nil {
								return nil, err
							}
						}
					} else {
						if c == ']' || c == '}' || c == ',' || isWhitespace(c) {
							stateStack.pop()
							var err error
							pendingResult, err = specificState.finalize()
							if err != nil {
								return nil, err
							}
							done = false
						} else {
							specificState.buffer.WriteByte(c)
						}
					}
				}
			case *jsonParserStateList:
				if pendingResult != nil {
					specificState.list = append(specificState.list, pendingResult)
					pendingResult = nil
				}
				if isWhitespace(c) {
					// ignore
				} else {
					if c == ']' {
						pendingResult = &specificState.list
						stateStack.pop()
					} else if len(specificState.list) == 0 {
						// new empty list
						stateStack.push(&jsonParserStateAnyObjPlaceholder{})
						done = false
					} else if c == ',' {
						stateStack.push(&jsonParserStateAnyObjPlaceholder{})
					}
				}
			case *jsonParserStateMap:
				if isWhitespace(c) {
					// ignore
				} else if c == '}' {
					pendingResult = specificState.currentMap
					stateStack.pop()
				} else if c == ',' {
					stateStack.push(&jsonStateMapKeyValue{})
				} else if specificState.currentMap.Size() == 0 {
					stateStack.push(&jsonStateMapKeyValue{})
					done = false
				} else {
					return nil, errors.New("invalid map state")
				}
			case *jsonStateMapKeyValue:
				switch specificState.state {
				case 0: // key
					if specificState.keyBuffer.Len() == 0 {
						if isWhitespace(c) {
							// ignore
						} else {
							if c != '"' {
								return nil, errors.New("map key must start with a quote")
							}
							specificState.keyBuffer.WriteByte(c)
						}
					} else {
						specificState.keyBuffer.WriteByte(c)
						prevChar := input[i-1]
						if c == '"' && prevChar != '\\' {
							specificState.state = 1
						}
					}
				case 1: // ':'
					if isWhitespace(c) {
						// ignore
					} else if c == ':' {
						specificState.state = 2
						stateStack.push(&jsonParserStateAnyObjPlaceholder{})
					} else {
						return nil, errors.New("invalid character in map definition, colon expected")
					}
				case 2: // value
					if pendingResult == nil {
						return nil, errors.New("missing value in map")
					}
					key := specificState.keyBuffer.String()
					if !strings.HasPrefix(key, "\"") || !strings.HasSuffix(key, "\"") {
						return nil, errors.New("map key should be a string enclosed in quotes")
					}
					key = key[1 : len(key)-1]
					stateStack.pop()
					mapState, isMap := stateStack.peek().(*jsonParserStateMap)
					if !isMap {
						return nil, errors.New("map key value state, but no map state underneath")
					}
					mapState.currentMap.Put(key, pendingResult)
					pendingResult = nil
					done = false
				default:
					return nil, errors.New("unknown jsonStateMapKeyValue state")
				}
			default:
				return nil, errors.New("invalid parser state")
			}
		}
	}

	if stateStack.size() != 0 {
		return nil, errors.New("state stack should be empty at the end")
	}

	return pendingResult, nil
}

func (s *jsonParserStateSingleValue) finalize() (OJsonObject, error) {
	str := s.buffer.String()
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		str = str[1 : len(str)-1]
		return &OJsonString{Value: str}, nil
	}
	if str == "true" {
		result := OJsonBool(true)
		return &result, nil
	}
	if str == "false" {
		result := OJsonBool(false)
		return &result, nil
	}
	return nil, errors.New("Invalid value: " + str)
}

type jsonParserStateStack struct {
	stack []jsonParserState
}

func (s *jsonParserStateStack) push(state jsonParserState) {
	s.stack = append(s.stack, state)
}

func (s *jsonParserStateStack) replaceTop(state jsonParserState) {
	s.stack[len(s.stack)-1] = state
}

func (s *jsonParserStateStack) peek() jsonParserState {
	return s.stack[len(s.stack)-1]
}

func (s *jsonParserStateStack) pop() jsonParserState {
	top := s.stack[len(s.stack)-1]
	s.stack = s.stack[0 : len(s.stack)-1]
	return top
}

func (s *jsonParserStateStack) size() int {
	return len(s.stack)
}
