package orderedjson2kast

import (
	"fmt"
	"strings"

	oj "github.com/ElrondNetwork/wasm-vm/mandos-go/orderedjson"
)

func jsonToKastOrdered(j oj.OJsonObject) string {
	var sb strings.Builder
	writeKast(j, &sb)
	return sb.String()
}

func writeStringKast(sb *strings.Builder, value string) {
	sb.WriteString(fmt.Sprintf("#token(\"\\\"%s\\\"\",\"String\")", value))
}

func writeKast(jobj oj.OJsonObject, sb *strings.Builder) {
	switch j := jobj.(type) {
	case *oj.OJsonMap:
		sb.WriteString("`{_}_IELE-DATA`(")
		for _, keyValuePair := range j.OrderedKV {
			sb.WriteString("`_,__IELE-DATA`(`_:__IELE-DATA`(")
			writeStringKast(sb, keyValuePair.Key)
			sb.WriteString(",")
			writeKast(keyValuePair.Value, sb)
			sb.WriteString("),")
		}
		sb.WriteString("`.List{\"_,__IELE-DATA\"}`(.KList)")
		for i := 0; i < j.Size(); i++ {
			sb.WriteString(")")
		}
		sb.WriteString(")")
	case *oj.OJsonList:
		collection := []oj.OJsonObject(*j)

		sb.WriteString("`[_]_IELE-DATA`(")
		for _, elem := range collection {
			sb.WriteString("`_,__IELE-DATA`(")
			writeKast(elem, sb)
			sb.WriteString(",")
		}
		sb.WriteString("`.List{\"_,__IELE-DATA\"}`(.KList)")
		for i := 0; i < len(collection); i++ {
			sb.WriteString(")")
		}
		sb.WriteString(")")
	case *oj.OJsonString:
		writeStringKast(sb, j.Value)
	case *oj.OJsonBool:
		value := bool(*j)
		sb.WriteString(fmt.Sprintf("#token(\"%t\",\"Bool\")", value))
	default:
		panic("unknown OJsonObject type")
	}
}
