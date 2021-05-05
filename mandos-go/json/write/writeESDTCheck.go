package mandosjsonwrite

import (
	mj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/model"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
)

func checkESDTDataToOJ(esdtItems []*mj.CheckESDTData) *oj.OJsonMap {
	esdtItemsOJ := oj.NewMap()
	for _, esdtItem := range esdtItems {
		esdtItemsOJ.Put(esdtItem.TokenIdentifier.Original, checkESDTItemToOJ(esdtItem))
	}
	return esdtItemsOJ
}

func checkESDTItemToOJ(esdtItem *mj.CheckESDTData) oj.OJsonObject {
	if isCompactCheckESDT(esdtItem) {
		return checkBigIntToOJ(esdtItem.Instances[0].Balance)
	}

	esdtItemOJ := oj.NewMap()

	// instances
	if len(esdtItem.Instances) == 1 {
		appendCheckESDTInstanceToOJ(esdtItem.Instances[0], esdtItemOJ)
	} else {
		var convertedList []oj.OJsonObject
		for _, esdtInstance := range esdtItem.Instances {
			esdtInstanceOJ := oj.NewMap()
			appendCheckESDTInstanceToOJ(esdtInstance, esdtInstanceOJ)
			convertedList = append(convertedList, esdtInstanceOJ)
		}
		instancesOJList := oj.OJsonList(convertedList)
		esdtItemOJ.Put("instances", &instancesOJList)
	}

	if len(esdtItem.LastNonce.Original) > 0 {
		esdtItemOJ.Put("lastNonce", checkUint64ToOJ(esdtItem.LastNonce))
	}

	// roles
	if len(esdtItem.Roles) > 0 {
		var convertedList []oj.OJsonObject
		for _, roleStr := range esdtItem.Roles {
			convertedList = append(convertedList, &oj.OJsonString{Value: roleStr})
		}
		rolesOJList := oj.OJsonList(convertedList)
		esdtItemOJ.Put("roles", &rolesOJList)
	}
	if len(esdtItem.Frozen.Original) > 0 {
		esdtItemOJ.Put("frozen", checkUint64ToOJ(esdtItem.Frozen))
	}

	return esdtItemOJ
}

func appendCheckESDTInstanceToOJ(esdtInstance *mj.CheckESDTInstance, targetOj *oj.OJsonMap) {
	if len(esdtInstance.Nonce.Original) > 0 {
		targetOj.Put("nonce", checkUint64ToOJ(esdtInstance.Nonce))
	}
	if len(esdtInstance.Balance.Original) > 0 {
		targetOj.Put("balance", checkBigIntToOJ(esdtInstance.Balance))
	}
	if !esdtInstance.Attributes.Unspecified && len(esdtInstance.Attributes.Value) > 0 {
		targetOj.Put("attributes", checkBytesToOJ(esdtInstance.Attributes))
	}
}

func isCompactCheckESDT(esdtItem *mj.CheckESDTData) bool {
	if len(esdtItem.Instances) != 1 {
		return false
	}
	if len(esdtItem.Instances[0].Nonce.Original) > 0 {
		return false
	}
	if len(esdtItem.Roles) > 0 {
		return false
	}
	if len(esdtItem.Frozen.Original) > 0 {
		return false
	}
	return true
}
