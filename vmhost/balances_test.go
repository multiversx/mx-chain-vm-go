package vmhost

import (
	"crypto/elliptic"
	"io"
	"math/big"
	"testing"

	"github.com/multiversx/mx-chain-vm-common-go"
	"github.com/stretchr/testify/require"
)

func TestCheckBalances_BaseCurrency_Balanced(t *testing.T) {
	output := &vmcommon.VMOutput{
		OutputAccounts: map[string]*vmcommon.OutputAccount{
			"sender": {
				Address:      []byte("sender"),
				BalanceDelta: big.NewInt(-100),
			},
			"receiver": {
				Address:      []byte("receiver"),
				BalanceDelta: big.NewInt(100),
			},
		},
	}

	err := checkBaseCurrency(output)
	require.Nil(t, err)
}

func TestCheckBalances_BaseCurrency_Unbalanced(t *testing.T) {
	output := &vmcommon.VMOutput{
		OutputAccounts: map[string]*vmcommon.OutputAccount{
			"sender": {
				Address:      []byte("sender"),
				BalanceDelta: big.NewInt(-100),
			},
			"receiver": {
				Address:      []byte("receiver"),
				BalanceDelta: big.NewInt(99),
			},
		},
	}

	err := checkBaseCurrency(output)
	require.Nil(t, err)
}

func TestCheckBalances_BackTransfers_Empty(t *testing.T) {
	managedContext := &managedTypesContextStub{
		backTransfers: []*vmcommon.ESDTTransfer{},
		backValue:     big.NewInt(0),
	}

	err := checkBackTransfers(nil, managedContext)
	require.Nil(t, err)
}

func TestCheckBalances_BackTransfers_NonEmpty(t *testing.T) {
	managedContext := &managedTypesContextStub{
		backTransfers: []*vmcommon.ESDTTransfer{
			{},
		},
		backValue: big.NewInt(0),
	}

	err := checkBackTransfers(nil, managedContext)
	require.Nil(t, err)
}

type managedTypesContextStub struct {
	backTransfers []*vmcommon.ESDTTransfer
	backValue     *big.Int
}

func (m *managedTypesContextStub) GetBackTransfers() ([]*vmcommon.ESDTTransfer, *big.Int) {
	return m.backTransfers, m.backValue
}

func (m *managedTypesContextStub) AddBackTransfers(value *big.Int, transfers []*vmcommon.ESDTTransfer, index uint32) {
}
func (m *managedTypesContextStub) GetRandReader() io.Reader {
	return nil
}
func (m *managedTypesContextStub) InitState() {
}
func (m *managedTypesContextStub) PushState() {
}
func (m *managedTypesContextStub) PopBackTransferIfAsyncCallBack(vmInput *vmcommon.ContractCallInput) {
}
func (m *managedTypesContextStub) PopSetActiveState() {
}
func (m *managedTypesContextStub) PopDiscard() {
}
func (m *managedTypesContextStub) ClearStateStack() {
}
func (m *managedTypesContextStub) IsInterfaceNil() bool {
	return false
}
func (m *managedTypesContextStub) ConsumeGasForBigIntCopy(values ...*big.Int) error {
	return nil
}
func (m *managedTypesContextStub) ConsumeGasForThisIntNumberOfBytes(byteLen int) error {
	return nil
}
func (m *managedTypesContextStub) ConsumeGasForBytes(bytes []byte) error {
	return nil
}
func (m *managedTypesContextStub) ConsumeGasForThisBigIntNumberOfBytes(byteLen *big.Int) error {
	return nil
}
func (m *managedTypesContextStub) ConsumeGasForBigFloatCopy(values ...*big.Float) error {
	return nil
}
func (m *managedTypesContextStub) GetBigIntOrCreate(handle int32) *big.Int {
	return nil
}
func (m *managedTypesContextStub) GetBigInt(handle int32) (*big.Int, error) {
	return nil, nil
}
func (m *managedTypesContextStub) GetTwoBigInt(handle1 int32, handle2 int32) (*big.Int, *big.Int, error) {
	return nil, nil, nil
}
func (m *managedTypesContextStub) NewBigInt(value *big.Int) int32 {
	return 0
}
func (m *managedTypesContextStub) NewBigIntFromInt64(int64Value int64) int32 {
	return 0
}
func (m *managedTypesContextStub) BigFloatPrecIsNotValid(precision uint) bool {
	return false
}
func (m *managedTypesContextStub) BigFloatExpIsNotValid(exponent int) bool {
	return false
}
func (m *managedTypesContextStub) EncodedBigFloatIsNotValid(encodedBigFloat []byte) bool {
	return false
}
func (m *managedTypesContextStub) GetBigFloatOrCreate(handle int32) (*big.Float, error) {
	return nil, nil
}
func (m *managedTypesContextStub) GetBigFloat(handle int32) (*big.Float, error) {
	return nil, nil
}
func (m *managedTypesContextStub) GetTwoBigFloats(handle1 int32, handle2 int32) (*big.Float, *big.Float, error) {
	return nil, nil, nil
}
func (m *managedTypesContextStub) PutBigFloat(value *big.Float) (int32, error) {
	return 0, nil
}
func (m *managedTypesContextStub) GetEllipticCurve(handle int32) (*elliptic.CurveParams, error) {
	return nil, nil
}
func (m *managedTypesContextStub) PutEllipticCurve(curve *elliptic.CurveParams) int32 {
	return 0
}
func (m *managedTypesContextStub) GetEllipticCurveSizeOfField(ecHandle int32) int32 {
	return 0
}
func (m *managedTypesContextStub) Get100xCurveGasCostMultiplier(ecHandle int32) int32 {
	return 0
}
func (m *managedTypesContextStub) GetScalarMult100xCurveGasCostMultiplier(ecHandle int32) int32 {
	return 0
}
func (m *managedTypesContextStub) GetUCompressed100xCurveGasCostMultiplier(ecHandle int32) int32 {
	return 0
}
func (m *managedTypesContextStub) GetPrivateKeyByteLengthEC(ecHandle int32) int32 {
	return 0
}
func (m *managedTypesContextStub) NewManagedBuffer() int32 {
	return 0
}
func (m *managedTypesContextStub) NewManagedBufferFromBytes(bytes []byte) int32 {
	return 0
}
func (m *managedTypesContextStub) SetBytes(mBufferHandle int32, bytes []byte) {
}
func (m *managedTypesContextStub) GetBytes(mBufferHandle int32) ([]byte, error) {
	return nil, nil
}
func (m *managedTypesContextStub) AppendBytes(mBufferHandle int32, bytes []byte) bool {
	return false
}
func (m *managedTypesContextStub) GetLength(mBufferHandle int32) int32 {
	return 0
}
func (m *managedTypesContextStub) GetSlice(mBufferHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error) {
	return nil, nil
}
func (m *managedTypesContextStub) DeleteSlice(mBufferHandle int32, startPosition int32, lengthOfSlice int32) ([]byte, error) {
	return nil, nil
}
func (m *managedTypesContextStub) InsertSlice(mBufferHandle int32, startPosition int32, slice []byte) ([]byte, error) {
	return nil, nil
}
func (m *managedTypesContextStub) ReadManagedVecOfManagedBuffers(managedVecHandle int32) ([][]byte, uint64, error) {
	return nil, 0, nil
}
func (m *managedTypesContextStub) WriteManagedVecOfManagedBuffers(data [][]byte, destinationHandle int32) error {
	return nil
}
func (m *managedTypesContextStub) NewManagedMap() int32 {
	return 0
}
func (m *managedTypesContextStub) ManagedMapPut(mMapHandle int32, keyHandle int32, valueHandle int32) error {
	return nil
}
func (m *managedTypesContextStub) ManagedMapGet(mMapHandle int32, keyHandle int32, outValueHandle int32) error {
	return nil
}
func (m *managedTypesContextStub) ManagedMapRemove(mMapHandle int32, keyHandle int32, outValueHandle int32) error {
	return nil
}
func (m *managedTypesContextStub) ManagedMapContains(mMapHandle int32, keyHandle int32) (bool, error) {
	return false, nil
}
