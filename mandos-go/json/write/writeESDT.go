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
					esdtOJ.Put(tokenNameKey, bigIntToOJ(esdtItem.Value))
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
	if len(esdtItem.Nonce.Original) > 0 {
		return false
	}
	if len(esdtItem.Frozen.Original) > 0 {
		return false
	}
	return true
}

func groupESDTIntoMap(esdtItems []*mj.ESDTData) map[string]map[int]*mj.ESDTData {
	m := make(map[string]map[int]*mj.ESDTData)
	for _, item := range esdtItems {
		var nonceMap map[int]*mj.ESDTData
		nonceMap, isPresent := m[string(item.TokenIdentifier.Value)]
		if !isPresent {
			nonceMap = make(map[int]*mj.ESDTData)
			m[string(item.TokenIdentifier.Value)] = nonceMap
		}
		nonceMap[int(item.Nonce.Value)] = item
	}
	return m
}

func esdtToFullMapOJ(esdtItem *mj.ESDTData) *oj.OJsonMap {
	esdtItemOJ := oj.NewMap()
	if len(esdtItem.TokenIdentifier.Original) > 0 {
		esdtItemOJ.Put("tokenIdentifier", bytesFromStringToOJ(esdtItem.TokenIdentifier))
	}
	if len(esdtItem.Nonce.Original) > 0 {
		esdtItemOJ.Put("nonce", uint64ToOJ(esdtItem.Nonce))
	}
	if len(esdtItem.Value.Original) > 0 {
		esdtItemOJ.Put("value", bigIntToOJ(esdtItem.Value))
	}
	if len(esdtItem.Frozen.Original) > 0 {
		esdtItemOJ.Put("frozen", uint64ToOJ(esdtItem.Frozen))
	}
	return esdtItemOJ
}

func esdtRolesToMapOJ(esdtRoles []*mj.ESDTRoles) *oj.OJsonMap {
	esdtRolesOJ := oj.NewMap()
	for _, rolesItem := range esdtRoles {
		var convertedList []oj.OJsonObject
		for _, roleStr := range rolesItem.Roles {
			convertedList = append(convertedList, &oj.OJsonString{Value: roleStr})
		}
		rolesOJList := oj.OJsonList(convertedList)
		esdtRolesOJ.Put(rolesItem.TokenIdentifier.Original, &rolesOJList)
	}
	return esdtRolesOJ
}

func esdtLastNoncesToMapOJ(lastNonces map[string]*mj.JSONUint64) *oj.OJsonMap {
	esdtLastNonces := oj.NewMap()
	for tokenName, nonce := range lastNonces {
		esdtLastNonces.Put(tokenName, uint64ToOJ(*nonce))
	}

	return esdtLastNonces
}
