package hosttest

import (
	"testing"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

func TestAsyncBuiltin(t *testing.T) {
	t.Skip()
	log := logger.GetOrCreate("test")

	TestGasUsed_LegacyAsyncCall_InShard_BuiltinCall(t)
	log.Info("TestGasUsed_AsyncCall_BuiltinCall ok")

	TestGasUsed_LegacyAsyncCall_CrossShard_BuiltinCall(t)
	log.Info("TestGasUsed_AsyncCall_CrossShard_BuiltinCall ok")

	TestGasUsed_ESDTTransferInCallback(t)
	log.Info("TestGasUsed_ESDTTransferInCallback ok")

	TestGasUsed_ESDTTransfer_CallbackFail(t)
	log.Info("TestGasUsed_ESDTTransfer_CallbackFail ok")

	TestESDT_GettersAPI_ExecuteAfterBuiltinCall(t)
	log.Info("TestESDT_GettersAPI_ExecuteAfterBuiltinCall ok")
}
