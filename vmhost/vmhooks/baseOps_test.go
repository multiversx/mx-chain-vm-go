package vmhooks

import (
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
	"github.com/multiversx/mx-chain-scenario-go/worldmock"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-go/mock/mockery"
	"github.com/multiversx/mx-chain-vm-go/vmhost"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVMHooksImpl_GetGasLeft(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _ := createTestVMHooks()

	gasLeft := hooks.GetGasLeft()
	require.Equal(t, int64(100), gasLeft)
}

func TestVMHooksImpl_GetSCAddress(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, _, _, _ := createTestVMHooks()

	scAddress := []byte("sc-address")
	runtime.On("GetContextAddress").Return(scAddress)

	hooks.GetSCAddress(0)
	runtime.AssertCalled(t, "GetContextAddress")
}

func TestVMHooksImpl_GetOwnerAddress(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	ownerAddress := []byte("owner-address")
	blockchain.On("GetOwnerAddress").Return(ownerAddress, nil)

	hooks.GetOwnerAddress(0)
	blockchain.AssertCalled(t, "GetOwnerAddress")
}

func TestVMHooksImpl_GetShardOfAddress(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockchain.On("GetShardOfAddress", mock.Anything).Return(uint32(1))

	shard := hooks.GetShardOfAddress(0)
	require.Equal(t, int32(1), shard)
}

func TestVMHooksImpl_IsSmartContract(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockchain.On("IsSmartContract", mock.Anything).Return(true)

	isSC := hooks.IsSmartContract(0)
	require.Equal(t, int32(1), isSC)
}

func TestVMHooksImpl_SignalError(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	errorMessage := "error message"
	instance := runtime.GetInstance().(*mockery.MockInstance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return([]byte(errorMessage), nil)
	runtime.On("SignalUserError", errorMessage).Return()

	hooks.SignalError(0, 0)
	runtime.AssertCalled(t, "SignalUserError", errorMessage)
}

func TestVMHooksImpl_GetExternalBalance(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	balance := big.NewInt(100)
	blockchain.On("GetBalance", mock.Anything).Return(balance)

	hooks.GetExternalBalance(0, 0)
	blockchain.AssertCalled(t, "GetBalance", mock.Anything)
}

func TestVMHooksImpl_GetBlockHash(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockHash := []byte("block-hash")
	blockchain.On("BlockHash", mock.Anything).Return(blockHash)

	ret := hooks.GetBlockHash(0, 0)
	require.Equal(t, int32(0), ret)
	blockchain.AssertCalled(t, "BlockHash", mock.Anything)
}

func TestVMHooksImpl_GetESDTBalance(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	esdtToken := &esdt.ESDigitalToken{
		Value: big.NewInt(100),
	}
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(esdtToken, nil)

	ret := hooks.GetESDTBalance(0, 0, 0, 0, 0)
	require.NotEqual(t, int32(-1), ret)
}

func TestVMHooksImpl_GetESDTNFTNameLength(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	esdtToken := &esdt.ESDigitalToken{
		TokenMetaData: &esdt.MetaData{
			Name: []byte("test-token"),
		},
	}
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(esdtToken, nil)

	ret := hooks.GetESDTNFTNameLength(0, 0, 0, 0)
	require.Equal(t, int32(len(esdtToken.TokenMetaData.Name)), ret)
}

func TestVMHooksImpl_GetESDTNFTAttributeLength(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	esdtToken := &esdt.ESDigitalToken{
		TokenMetaData: &esdt.MetaData{
			Attributes: []byte("test-attributes"),
		},
	}
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(esdtToken, nil)

	ret := hooks.GetESDTNFTAttributeLength(0, 0, 0, 0)
	require.Equal(t, int32(len(esdtToken.TokenMetaData.Attributes)), ret)
}

func TestVMHooksImpl_GetESDTNFTURILength(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	esdtToken := &esdt.ESDigitalToken{
		TokenMetaData: &esdt.MetaData{
			URIs: [][]byte{[]byte("test-uri")},
		},
	}
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(esdtToken, nil)

	ret := hooks.GetESDTNFTURILength(0, 0, 0, 0)
	require.Equal(t, int32(len(esdtToken.TokenMetaData.URIs[0])), ret)
}

func TestVMHooksImpl_GetESDTTokenData(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, managedType := createTestVMHooksFull()

	esdtToken := &esdt.ESDigitalToken{
		Value:      big.NewInt(100),
		Properties: []byte("properties"),
		TokenMetaData: &esdt.MetaData{
			Hash:       []byte("hash"),
			Name:       []byte("name"),
			Attributes: []byte("attributes"),
			Creator:    []byte("creator"),
			Royalties:  10,
			URIs:       [][]byte{[]byte("uri")},
		},
	}
	blockchain.On("GetESDTToken", mock.Anything, mock.Anything, mock.Anything).Return(esdtToken, nil)
	managedType.On("GetBigIntOrCreate", mock.Anything).Return(big.NewInt(0))

	ret := hooks.GetESDTTokenData(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	require.NotEqual(t, int32(-1), ret)
}

func TestVMHooksImpl_GetESDTLocalRoles(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)
	enableEpochs := &worldmock.EnableEpochsHandlerStub{}
	host.On("EnableEpochsHandler").Return(enableEpochs)

	managedType.On("GetBytes", mock.Anything).Return([]byte("token-id"), nil)
	storage.On("GetStorage", mock.Anything).Return(nil, uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.GetESDTLocalRoles(0)
	require.NotEqual(t, int64(-1), ret)
}

func TestVMHooksImpl_ValidateTokenIdentifier(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	managedType.On("GetBytes", mock.Anything).Return([]byte("TEST-123456"), nil)

	ret := hooks.ValidateTokenIdentifier(0)
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_TransferValue(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	output.On("Transfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.TransferValue(0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_TransferValueExecute(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	host.On("AreInSameShard", mock.Anything, mock.Anything).Return(false)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	output.On("Transfer", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.TransferValueExecute(0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_TransferESDTExecute(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	host.On("AreInSameShard", mock.Anything, mock.Anything).Return(false)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	output.On("TransferESDT", mock.Anything, mock.Anything).Return(uint64(0), nil)

	ret := hooks.TransferESDTExecute(0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_MultiTransferESDTNFTExecute(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	host.On("AreInSameShard", mock.Anything, mock.Anything).Return(false)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	output.On("TransferESDT", mock.Anything, mock.Anything).Return(uint64(0), nil)

	instance := runtime.GetInstance().(*mockery.MockInstance)
	instance.On("MemLoad", mock.Anything, mock.Anything).Return([]byte{0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 3}, nil)
	instance.On("MemLoadMultiple", mock.Anything, mock.Anything).Return([][]byte{[]byte("token"), big.NewInt(1).Bytes(), big.NewInt(100).Bytes()}, nil)

	ret := hooks.MultiTransferESDTNFTExecute(0, 1, 0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_CreateAsyncCall(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, _, _, async, _ := createTestVMHooksWithSetMetering()

	async.On("RegisterAsyncCall", mock.Anything, mock.Anything).Return(nil)

	ret := hooks.CreateAsyncCall(0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_SetAsyncContextCallback(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, _, _, async, _ := createTestVMHooksWithSetMetering()

	async.On("SetContextCallback", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.SetAsyncContextCallback(0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_UpgradeContract(t *testing.T) {
	t.Parallel()

	hooks, _, runtime, _, _, _, _, _, async, _ := createTestVMHooksWithSetMetering()

	runtime.On("SetRuntimeBreakpointValue", mock.Anything).Return()
	async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.UpgradeContract(0, 1000000, 0, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_UpgradeFromSourceContract(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _, async, _ := createTestVMHooksWithSetMetering()

	blockchain.On("GetCode", mock.Anything).Return([]byte("code"), nil)
	async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	hooks.UpgradeFromSourceContract(0, 1000000, 0, 0, 0, 0, 0, 0)
}

func TestVMHooksImpl_DeleteContract(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()

	vmHooks.async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	vmHooks.hooks.DeleteContract(0, 1000000, 0, 0, 0)
}

func TestVMHooksImpl_AsyncCall(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()

	vmHooks.async.On("RegisterLegacyAsyncCall", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	vmHooks.hooks.AsyncCall(0, 0, 0, 0)
}

func TestVMHooksImpl_GetArgumentLength(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	runtime.On("Arguments").Return([][]byte{[]byte("arg1"), []byte("argument2")})

	ret := hooks.GetArgumentLength(1)
	require.Equal(t, int32(len("argument2")), ret)
}

func TestVMHooksImpl_GetArgument(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	runtime.On("Arguments").Return([][]byte{[]byte("arg1"), []byte("argument2")})

	ret := hooks.GetArgument(1, 0)
	require.Equal(t, int32(len("argument2")), ret)
}

func TestVMHooksImpl_GetFunction(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	runtime.On("FunctionName").Return("testFunction")

	ret := hooks.GetFunction(0)
	require.Equal(t, int32(len("testFunction")), ret)
}

func TestVMHooksImpl_GetNumArguments(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	runtime.On("Arguments").Return([][]byte{[]byte("arg1"), []byte("argument2")})

	ret := hooks.GetNumArguments()
	require.Equal(t, int32(2), ret)
}

func TestVMHooksImpl_StorageStore(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	storage.On("SetStorage", mock.Anything, mock.Anything).Return(vmhost.StorageAdded, nil)

	ret := hooks.StorageStore(0, 0, 0, 0)
	require.Equal(t, int32(vmhost.StorageAdded), ret)
}

func TestVMHooksImpl_StorageLoadLength(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	storage.On("GetStorageUnmetered", mock.Anything).Return([]byte("data"), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.StorageLoadLength(0, 0)
	require.Equal(t, int32(len("data")), ret)
}

func TestVMHooksImpl_StorageLoad(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	storage.On("GetStorage", mock.Anything).Return([]byte("data"), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.StorageLoad(0, 0, 0)
	require.Equal(t, int32(len("data")), ret)
}

func TestVMHooksImpl_StorageLoadFromAddress(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	storage.On("GetStorageFromAddress", mock.Anything, mock.Anything).Return([]byte("data"), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.StorageLoadFromAddress(0, 0, 0, 0)
	require.Equal(t, int32(len("data")), ret)
}

func TestVMHooksImpl_GetCaller(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{CallerAddr: []byte("caller")}
	runtime.On("GetVMInput").Return(contractCallInput)

	hooks.GetCaller(0)
}

func TestVMHooksImpl_CheckNoPayment(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, _, _, _, _, _, _, _ := createTestVMHooksWithSetMetering()

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{
		CallValue:     big.NewInt(0),
		ESDTTransfers: []*vmcommon.ESDTTransfer{},
	}
	runtime.On("GetVMInput").Return(contractCallInput)

	hooks.CheckNoPayment()
}

func TestVMHooksImpl_GetCallValue(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, _, _, _, _, _, _, _ := createTestVMHooksWithSetMetering()

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{
		CallValue:     big.NewInt(100),
		ESDTTransfers: []*vmcommon.ESDTTransfer{},
	}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetCallValue(0)
	require.NotEqual(t, int32(-1), ret)
}

func TestVMHooksImpl_GetESDTValue(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{
		CallValue:     big.NewInt(0),
		ESDTTransfers: []*vmcommon.ESDTTransfer{{ESDTValue: big.NewInt(100)}},
	}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetESDTValue(0)
	require.NotEqual(t, int32(-1), ret)
}

func TestVMHooksImpl_GetESDTTokenName(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{ESDTTransfers: []*vmcommon.ESDTTransfer{{ESDTTokenName: []byte("token-name")}}}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetESDTTokenName(0)
	require.NotEqual(t, int32(-1), ret)
}

func TestVMHooksImpl_GetESDTTokenNonce(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{ESDTTransfers: []*vmcommon.ESDTTransfer{{ESDTTokenNonce: 123}}}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetESDTTokenNonce()
	require.Equal(t, int64(123), ret)
}

func TestVMHooksImpl_GetCurrentESDTNFTNonce(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, _, storage := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	storage.On("GetStorageFromAddress", mock.Anything, mock.Anything).Return(big.NewInt(123).Bytes(), uint32(0), false, nil)
	storage.On("UseGasForStorageLoad", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.GetCurrentESDTNFTNonce(0, 0, 0)
	require.Equal(t, int64(123), ret)
}

func TestVMHooksImpl_GetESDTTokenType(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{ESDTTransfers: []*vmcommon.ESDTTransfer{{ESDTTokenType: 1}}}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetESDTTokenType()
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_GetNumESDTTransfers(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{ESDTTransfers: []*vmcommon.ESDTTransfer{{}, {}}}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetNumESDTTransfers()
	require.Equal(t, int32(2), ret)
}

func TestVMHooksImpl_GetCallValueTokenName(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	contractCallInput := &vmcommon.ContractCallInput{}
	contractCallInput.VMInput = vmcommon.VMInput{ESDTTransfers: []*vmcommon.ESDTTransfer{{ESDTTokenName: []byte("token-name"), ESDTValue: big.NewInt(100)}}}
	runtime.On("GetVMInput").Return(contractCallInput)

	ret := hooks.GetCallValueTokenName(0, 0)
	require.Equal(t, int32(len("token-name")), ret)
}

func TestVMHooksImpl_IsReservedFunctionName(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	managedType := &mockery.MockManagedTypesContext{}
	host.On("ManagedTypes").Return(managedType)

	managedType.On("GetBytes", mock.Anything).Return([]byte("init"), nil)
	runtime.On("IsReservedFunctionName", "init").Return(true)

	ret := hooks.IsReservedFunctionName(0)
	require.Equal(t, int32(1), ret)
}

func TestVMHooksImpl_WriteLog(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	runtime.On("GetContextAddress").Return([]byte("address"))
	output.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return()

	hooks.WriteLog(0, 0, 0, 0)
}

func TestVMHooksImpl_WriteEventLog(t *testing.T) {
	t.Parallel()
	hooks, _, runtime, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	runtime.On("GetContextAddress").Return([]byte("address"))
	output.On("WriteLog", mock.Anything, mock.Anything, mock.Anything).Return()

	hooks.WriteEventLog(0, 0, 0, 0, 0)
}

func TestVMHooksImpl_GetBlockTimestamp(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)

	blockchain.On("CurrentTimeStamp").Return(uint64(12345))

	ret := hooks.GetBlockTimestamp()
	require.Equal(t, int64(12345), ret)
}

func TestVMHooksImpl_GetBlockTimestampMs(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)

	blockchain.On("CurrentTimeStampMs").Return(uint64(12345000))

	ret := hooks.GetBlockTimestampMs()
	require.Equal(t, int64(12345000), ret)
}

func TestVMHooksImpl_GetBlockNonce(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)

	blockchain.On("CurrentNonce").Return(uint64(123))

	ret := hooks.GetBlockNonce()
	require.Equal(t, int64(123), ret)
}

func TestVMHooksImpl_GetBlockRound(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()

	vmHooks.blockchain.On("CurrentRound").Return(uint64(456))

	ret := vmHooks.hooks.GetBlockRound()
	require.Equal(t, int64(456), ret)
}

func TestVMHooksImpl_GetBlockEpoch(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()

	vmHooks.blockchain.On("CurrentEpoch").Return(uint32(789))

	ret := vmHooks.hooks.GetBlockEpoch()
	require.Equal(t, int64(789), ret)
}

func TestVMHooksImpl_GetBlockRandomSeed(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()

	vmHooks.blockchain.On("CurrentRandomSeed").Return([]byte("random-seed"))

	vmHooks.hooks.GetBlockRandomSeed(0)
}

func TestVMHooksImpl_GetStateRootHash(t *testing.T) {
	t.Parallel()
	vmHooks := createHooksWithBaseSetup()

	vmHooks.blockchain.On("GetStateRootHash").Return([]byte("state-root-hash"))

	vmHooks.hooks.GetStateRootHash(0)
}

func TestVMHooksImpl_GetPrevBlockTimestamp(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)

	blockchain.On("LastTimeStamp").Return(uint64(12345))

	ret := hooks.GetPrevBlockTimestamp()
	require.Equal(t, int64(12345), ret)
}

func TestVMHooksImpl_GetPrevBlockTimestampMs(t *testing.T) {
	t.Parallel()
	hooks, host, _, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)

	blockchain.On("LastTimeStampMs").Return(uint64(12345000))

	ret := hooks.GetPrevBlockTimestampMs()
	require.Equal(t, int64(12345000), ret)
}

func TestVMHooksImpl_GetPrevBlockNonce(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockchain.On("LastNonce").Return(uint64(123))

	ret := hooks.GetPrevBlockNonce()
	require.Equal(t, int64(123), ret)
}

func TestVMHooksImpl_GetPrevBlockRound(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockchain.On("LastRound").Return(uint64(456))

	ret := hooks.GetPrevBlockRound()
	require.Equal(t, int64(456), ret)
}

func TestVMHooksImpl_GetPrevBlockEpoch(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockchain.On("LastEpoch").Return(uint32(789))

	ret := hooks.GetPrevBlockEpoch()
	require.Equal(t, int64(789), ret)
}

func TestVMHooksImpl_GetPrevBlockRandomSeed(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()

	blockchain.On("LastRandomSeed").Return([]byte("random-seed"))

	hooks.GetPrevBlockRandomSeed(0)
}

func TestVMHooksImpl_GetBlockRoundTimeMs(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, _, _ := createTestVMHooksFull()

	ret := hooks.GetBlockRoundTimeMs()
	require.Equal(t, int64(6000), ret)
}

func TestVMHooksImpl_EpochStartBlockTimestampMs(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _ := createTestVMHooks()

	ret := hooks.EpochStartBlockTimestampMs()
	require.Equal(t, int64(12345000), ret)
}

func TestVMHooksImpl_EpochStartBlockNonce(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()
	blockchain.On("EpochStartBlockNonce").Return(uint64(123))

	ret := hooks.EpochStartBlockNonce()
	require.Equal(t, int64(123), ret)
}

func TestVMHooksImpl_EpochStartBlockRound(t *testing.T) {
	t.Parallel()
	hooks, _, _, _, _, _, blockchain, _ := createTestVMHooksFull()
	blockchain.On("EpochStartBlockRound").Return(uint64(456))

	ret := hooks.EpochStartBlockRound()
	require.Equal(t, int64(456), ret)
}

func TestVMHooksImpl_Finish(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)

	output.On("Finish", mock.Anything).Return()

	hooks.Finish(0, 0)
}

func TestVMHooksImpl_ExecuteOnSameContext(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	metering.On("BoundGasLimit", mock.Anything).Return(uint64(100))
	host.On("AreInSameShard", mock.Anything, mock.Anything).Return(true)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("caller"))
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{})
	host.On("ExecuteOnSameContext", mock.Anything).Return(nil)

	ret := hooks.ExecuteOnSameContext(0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ExecuteOnDestContext(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, _, _, _, _, _, async, _ := createTestVMHooksWithSetMetering()
	host.On("AreInSameShard", mock.Anything, mock.Anything).Return(true)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("IsBuiltinFunctionCall", mock.Anything).Return(false)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("caller"))
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{})

	host.On("ExecuteOnDestContext", mock.Anything).Return(&vmcommon.VMOutput{}, true, nil)
	host.On("CompleteLogEntriesWithCallType", mock.Anything, mock.Anything).Return()

	async.On("SetAsyncArgumentsForCall", mock.Anything).Return()
	async.On("CompleteChildConditional", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.ExecuteOnDestContext(0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_ExecuteReadOnly(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, _, _, _, _, _, async, _ := createTestVMHooksWithSetMetering()
	host.On("AreInSameShard", mock.Anything, mock.Anything).Return(true)
	host.On("IsBuiltinFunctionName", mock.Anything).Return(false)
	host.On("IsBuiltinFunctionCall", mock.Anything).Return(false)

	runtime.On("GetContextAddress").Return([]byte("sender"))
	host.On("ExecuteOnDestContext", mock.Anything).Return(&vmcommon.VMOutput{}, true, nil)
	runtime.On("ReadOnly").Return(false)
	runtime.On("SetReadOnly", mock.Anything).Return()
	runtime.On("GetOriginalCallerAddress").Return([]byte("address"))
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{})

	async.On("SetAsyncArgumentsForCall", mock.Anything).Return()
	async.On("CompleteChildConditional", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	ret := hooks.ExecuteReadOnly(0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_CreateContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	metering.On("BoundGasLimit", mock.Anything).Return(uint64(1000000))

	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{})
	host.On("CreateNewContract", mock.Anything, mock.Anything).Return([]byte("new-address"), nil)

	ret := hooks.CreateContract(0, 0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_DeployFromSourceContract(t *testing.T) {
	t.Parallel()
	hooks, host, runtime, metering, _, _ := createTestVMHooks()
	metering.On("UseGasBounded", mock.Anything).Return(nil)
	metering.On("BoundGasLimit", mock.Anything).Return(uint64(1000000))

	blockchain := &mockery.MockBlockchainContext{}
	host.On("Blockchain").Return(blockchain)
	runtime.On("GetContextAddress").Return([]byte("sender"))
	runtime.On("GetOriginalCallerAddress").Return([]byte("original-caller"))
	runtime.On("GetVMInput").Return(&vmcommon.ContractCallInput{})
	blockchain.On("GetCode", mock.Anything).Return([]byte("code"), nil)
	host.On("CreateNewContract", mock.Anything, mock.Anything).Return([]byte("new-address"), nil)

	ret := hooks.DeployFromSourceContract(0, 0, 0, 0, 0, 0, 0, 0)
	require.Equal(t, int32(0), ret)
}

func TestVMHooksImpl_GetNumReturnData(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	output.On("ReturnData").Return([][]byte{[]byte("data1"), []byte("data2")})

	ret := hooks.GetNumReturnData()
	require.Equal(t, int32(2), ret)
}

func TestVMHooksImpl_GetReturnDataSize(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	output.On("ReturnData").Return([][]byte{[]byte("data1"), []byte("data2")})

	ret := hooks.GetReturnDataSize(1)
	require.Equal(t, int32(len("data2")), ret)
}

func TestVMHooksImpl_GetReturnData(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	output.On("ReturnData").Return([][]byte{[]byte("data1"), []byte("data2")})

	ret := hooks.GetReturnData(1, 0)
	require.Equal(t, int32(len("data2")), ret)
}

func TestVMHooksImpl_CleanReturnData(t *testing.T) {
	t.Parallel()
	hooks, _, _, metering, output, _ := createTestVMHooks()
	metering.On("UseGasBoundedAndAddTracedGas", mock.Anything, mock.Anything).Return(nil)

	output.On("ClearReturnData").Return()

	hooks.CleanReturnData()
	output.AssertCalled(t, "ClearReturnData")
}
