package mandosjsonwrite

import (
	mj "github.com/multiversx/mx-chain-vm-go/mandos-go/model"
	oj "github.com/multiversx/mx-chain-vm-go/mandos-go/orderedjson"
)

func esdtTxDataToOJ(esdtItems []*mj.ESDTTxData) oj.OJsonObject {
	esdtItemList := oj.OJsonList{}
	for _, esdtItemRaw := range esdtItems {
		esdtItemOJ := esdtTxRawEntryToOJ(esdtItemRaw)
		esdtItemList = append(esdtItemList, esdtItemOJ)
	}

	return &esdtItemList

}

func esdtTxRawEntryToOJ(esdtItemRaw *mj.ESDTTxData) *oj.OJsonMap {
	esdtItemOJ := oj.NewMap()

	if len(esdtItemRaw.TokenIdentifier.Original) > 0 {
		esdtItemOJ.Put("tokenIdentifier", bytesFromStringToOJ(esdtItemRaw.TokenIdentifier))
	}
	if len(esdtItemRaw.Nonce.Original) > 0 {
		esdtItemOJ.Put("nonce", uint64ToOJ(esdtItemRaw.Nonce))
	}
	if len(esdtItemRaw.Value.Original) > 0 {
		esdtItemOJ.Put("value", bigIntToOJ(esdtItemRaw.Value))
	}

	return esdtItemOJ
}

func esdtDataToOJ(esdtItems []*mj.ESDTData) *oj.OJsonMap {
	esdtItemsOJ := oj.NewMap()
	for _, esdtItem := range esdtItems {
		esdtItemsOJ.Put(esdtItem.TokenIdentifier.Original, esdtItemToOJ(esdtItem))
	}
	return esdtItemsOJ
}

func esdtItemToOJ(esdtItem *mj.ESDTData) oj.OJsonObject {
	if isCompactESDT(esdtItem) {
		return bigIntToOJ(esdtItem.Instances[0].Balance)
	}

	esdtItemOJ := oj.NewMap()

	// instances
	if len(esdtItem.Instances) > 0 {
		var convertedList []oj.OJsonObject
		for _, esdtInstance := range esdtItem.Instances {
			esdtInstanceOJ := oj.NewMap()
			appendESDTInstanceToOJ(esdtInstance, esdtInstanceOJ)
			convertedList = append(convertedList, esdtInstanceOJ)
		}
		instancesOJList := oj.OJsonList(convertedList)
		esdtItemOJ.Put("instances", &instancesOJList)
	}

	if len(esdtItem.LastNonce.Original) > 0 {
		esdtItemOJ.Put("lastNonce", uint64ToOJ(esdtItem.LastNonce))
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
		esdtItemOJ.Put("frozen", uint64ToOJ(esdtItem.Frozen))
	}

	return esdtItemOJ
}

func appendESDTInstanceToOJ(esdtInstance *mj.ESDTInstance, targetOj *oj.OJsonMap) {
	targetOj.Put("nonce", uint64ToOJ(esdtInstance.Nonce))

	if len(esdtInstance.Balance.Original) > 0 {
		targetOj.Put("balance", bigIntToOJ(esdtInstance.Balance))
	}
	if len(esdtInstance.Creator.Original) > 0 {
		targetOj.Put("creator", bytesFromStringToOJ(esdtInstance.Creator))
	}
	if len(esdtInstance.Royalties.Original) > 0 {
		targetOj.Put("royalties", uint64ToOJ(esdtInstance.Royalties))
	}
	if len(esdtInstance.Hash.Original) > 0 {
		targetOj.Put("hash", bytesFromStringToOJ(esdtInstance.Hash))
	}
	if !esdtInstance.Uris.IsUnspecified() {
		targetOj.Put("uri", valueListToOJ(esdtInstance.Uris))
	}
	if len(esdtInstance.Attributes.Value) > 0 {
		targetOj.Put("attributes", bytesFromTreeToOJ(esdtInstance.Attributes))
	}
}

func isCompactESDT(esdtItem *mj.ESDTData) bool {
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
