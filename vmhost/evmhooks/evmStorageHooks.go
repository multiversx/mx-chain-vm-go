package evmhooks

import (
	"github.com/ethereum/go-ethereum/common"
)

const getState = "getState"

func (context *EVMHooksImpl) GetState(key common.Hash) common.Hash {
	value, trieDepth, usedCache, err := context.GetStorageContext().GetStorage(key.Bytes())
	if err != nil {
		context.FailExecution(err)
		return common.Hash{}
	}

	loadCost := context.GetMeteringContext().GasSchedule().EVMOpcodeCost.Sload
	err = context.GetStorageContext().UseGasForStorageLoad(getState, int64(trieDepth), loadCost, usedCache)
	if err != nil {
		context.FailExecution(err)
		return common.Hash{}
	}

	return common.BytesToHash(value)
}

func (context *EVMHooksImpl) SetState(key common.Hash, value common.Hash) {
	_, err := context.GetStorageContext().SetStorage(key.Bytes(), trimValue(value))
	if err != nil {
		context.FailExecution(err)
	}
}

func trimValue(hash common.Hash) []byte {
	for currentPosition, currentByte := range hash {
		if currentByte != 0 {
			return hash.Bytes()[currentPosition:]
		}
	}
	return make([]byte, 0)
}
