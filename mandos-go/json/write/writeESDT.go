package mandosjsonwrite

import (
	"sort"

	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func appendESDTToOJ(esdtItems []*mj.ESDTData, parentOj *oj.OJsonMap) {
	esdtOJ := oj.NewMap()
	esdtMap := groupESDTIntoMap(esdtItems)
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
				if isCompactESDT(esdtItem) {
					esdtOJ.Put(tokenNameKey, bigIntToOJ(esdtItem.Balance))
				} else {
					esdtItemOJ := esdtToFullMapOJ(esdtItem)
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
				esdtItemOJ := esdtToFullMapOJ(esdtItem)
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

func isCompactESDT(esdtItem *mj.ESDTData) bool {
	if len(esdtItem.Frozen.Original) > 0 {
		return false
	}
	return false
}

func groupESDTIntoMap(esdtItems []*mj.ESDTData) map[string]map[int]*mj.ESDTData {
	m := make(map[string]map[int]*mj.ESDTData)
	for _, item := range esdtItems {
		var nonceMap map[int]*mj.ESDTData
		nonceMap, isPresent := m[string(item.TokenName.Value)]
		if !isPresent {
			nonceMap = make(map[int]*mj.ESDTData)
			m[string(item.TokenName.Value)] = nonceMap
		}
		nonceMap[int(item.Nonce.Value)] = item
	}
	return m
}

func esdtToFullMapOJ(esdtItem *mj.ESDTData) *oj.OJsonMap {
	esdtItemOJ := oj.NewMap()
	esdtItemOJ.Put("nonce", uint64ToOJ(esdtItem.Nonce))
	esdtItemOJ.Put("balance", bigIntToOJ(esdtItem.Balance))
	esdtItemOJ.Put("frozen", uint64ToOJ(esdtItem.Frozen))
	return esdtItemOJ
}
