package mandosjsonwrite

import (
	"sort"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func appendCheckESDTToOJ(esdtItems []*mj.CheckESDTData, parentOj *oj.OJsonMap) {
	esdtOJ := oj.NewMap()
	esdtMap := groupCheckESDTIntoMap(esdtItems)
	var tokenNameKeys []string
	for k := range esdtMap {
		tokenNameKeys = append(tokenNameKeys, k)
	}
	sort.Strings(tokenNameKeys)
	for _, k := range tokenNameKeys {
		tokenNameKey := "str:" + string(k)
		nonceMap := esdtMap[k]
		switch len(nonceMap) {
		case 0:
		case 1:
			for _, esdtItem := range nonceMap {
				if isCompactCheckESDT(esdtItem) {
					esdtOJ.Put(tokenNameKey, checkBigIntToOJ(esdtItem.Value))
				} else {
					esdtItemOJ := checkESDTToFullMapOJ(esdtItem)
					esdtOJ.Put(tokenNameKey, esdtItemOJ)
				}
			}
		default:
			var objList []oj.OJsonObject
			var nonceKeys []int
			for nk := range nonceMap {
				nonceKeys = append(nonceKeys, nk)
			}
			sort.Ints(nonceKeys)
			for _, nk := range nonceKeys {
				esdtItem := nonceMap[nk]
				esdtItemOJ := checkESDTToFullMapOJ(esdtItem)
				objList = append(objList, esdtItemOJ)
			}
			jsonList := oj.OJsonList(objList)
			esdtOJ.Put(tokenNameKey, &jsonList)
		}
	}
	if esdtOJ.Size() > 0 {
		parentOj.Put("esdt", esdtOJ)
	}
}

func isCompactCheckESDT(esdtItem *mj.CheckESDTData) bool {
	if len(esdtItem.Nonce.Original) > 0 {
		return false
	}
	if len(esdtItem.Frozen.Original) > 0 {
		return false
	}
	return true
}

func groupCheckESDTIntoMap(esdtItems []*mj.CheckESDTData) map[string]map[int]*mj.CheckESDTData {
	m := make(map[string]map[int]*mj.CheckESDTData)
	for _, item := range esdtItems {
		var nonceMap map[int]*mj.CheckESDTData
		nonceMap, isPresent := m[string(item.TokenIdentifier.Value)]
		if !isPresent {
			nonceMap = make(map[int]*mj.CheckESDTData)
			m[string(item.TokenIdentifier.Value)] = nonceMap
		}
		nonceMap[int(item.Nonce.Value)] = item
	}
	return m
}

func checkESDTToFullMapOJ(esdtItem *mj.CheckESDTData) *oj.OJsonMap {
	esdtItemOJ := oj.NewMap()
	if len(esdtItem.TokenIdentifier.Original) > 0 {
		esdtItemOJ.Put("tokenIdentifier", bytesFromStringToOJ(esdtItem.TokenIdentifier))
	}
	if len(esdtItem.Nonce.Original) > 0 {
		esdtItemOJ.Put("nonce", checkUint64ToOJ(esdtItem.Nonce))
	}
	if len(esdtItem.Value.Original) > 0 {
		esdtItemOJ.Put("value", checkBigIntToOJ(esdtItem.Value))
	}
	if len(esdtItem.Frozen.Original) > 0 {
		esdtItemOJ.Put("frozen", checkUint64ToOJ(esdtItem.Frozen))
	}
	return esdtItemOJ
}
