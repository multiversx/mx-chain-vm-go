package vmhooks

import (
	"bytes"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	gnarkgroth16 "github.com/consensys/gnark/backend/groth16"
	gnarkplonk "github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/examples/exponentiate"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"
	"github.com/multiversx/mx-chain-crypto-go/curves/bn254"
	"github.com/multiversx/mx-chain-crypto-go/zk/lowLevelFeatures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestManagedVerifyGroth16_Success(t *testing.T) {
	css, err := frontend.Compile(ecc.BLS12_381.ScalarField(), r1cs.NewBuilder, &exponentiate.Circuit{})
	require.Nil(t, err)

	pk, vk, err := gnarkgroth16.Setup(css)
	require.Nil(t, err)

	homework := &exponentiate.Circuit{
		X: 2,
		Y: 16,
		E: 4,
	}
	witness, err := frontend.NewWitness(homework, ecc.BLS12_381.ScalarField())
	require.Nil(t, err)

	proof, err := gnarkgroth16.Prove(css, pk, witness)
	require.Nil(t, err)

	var serializedProof bytes.Buffer
	_, err = proof.WriteTo(&serializedProof)
	require.Nil(t, err)

	var serializedVK bytes.Buffer
	_, err = vk.WriteTo(&serializedVK)
	require.Nil(t, err)

	pubW, err := witness.Public()
	require.Nil(t, err)
	pubWBytes, err := pubW.MarshalBinary()
	require.Nil(t, err)

	vmHooks := createHooksWithBaseSetup()
	hooks := vmHooks.hooks
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return(serializedProof.Bytes(), nil)
	managedType.On("GetBytes", int32(2)).Return(serializedVK.Bytes(), nil)
	managedType.On("GetBytes", int32(3)).Return(pubWBytes, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("IsUnsafeMode").Return(false)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := hooks.ManagedVerifyGroth16(int32(ecc.BLS12_381), 1, 2, 3)

	assert.Equal(t, int32(0), ret)
	runtime.AssertNotCalled(t, "FailExecution", mock.Anything)
}

func TestManagedVerifyPlonk_Success(t *testing.T) {
	css, err := frontend.Compile(ecc.BLS12_381.ScalarField(), scs.NewBuilder, &exponentiate.Circuit{})
	require.Nil(t, err)

	srs, srsLagrange, err := unsafekzg.NewSRS(css)
	require.Nil(t, err)

	pk, vk, err := gnarkplonk.Setup(css, srs, srsLagrange)
	require.Nil(t, err)

	homework := &exponentiate.Circuit{
		X: 2,
		Y: 16,
		E: 4,
	}
	witness, err := frontend.NewWitness(homework, ecc.BLS12_381.ScalarField())
	require.Nil(t, err)

	proof, err := gnarkplonk.Prove(css, pk, witness)
	require.Nil(t, err)

	var serializedProof bytes.Buffer
	_, err = proof.WriteTo(&serializedProof)
	require.Nil(t, err)

	var serializedVK bytes.Buffer
	_, err = vk.WriteTo(&serializedVK)
	require.Nil(t, err)

	pubW, err := witness.Public()
	require.Nil(t, err)
	pubWBytes, err := pubW.MarshalBinary()
	require.Nil(t, err)

	vmHooks := createHooksWithBaseSetup()
	hooks := vmHooks.hooks
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	managedType.On("GetBytes", int32(1)).Return(serializedProof.Bytes(), nil)
	managedType.On("GetBytes", int32(2)).Return(serializedVK.Bytes(), nil)
	managedType.On("GetBytes", int32(3)).Return(pubWBytes, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	runtime.On("IsUnsafeMode").Return(false)
	runtime.On("FailExecution", mock.Anything).Return()

	ret := hooks.ManagedVerifyPlonk(int32(ecc.BLS12_381), 1, 2, 3)

	assert.Equal(t, int32(0), ret)
	runtime.AssertNotCalled(t, "FailExecution", mock.Anything)
}

func TestManagedAddEC_Success(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	hooks := vmHooks.hooks
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	p1, _ := bn254.NewPointG1().Pick()
	p2, _ := bn254.NewPointG1().Pick()
	p3, _ := p1.Add(p2)
	p1Bytes, _ := p1.MarshalBinary()
	p2Bytes, _ := p2.MarshalBinary()
	p3Bytes, _ := p3.MarshalBinary()

	managedType.On("GetBytes", int32(1)).Return(p1Bytes, nil)
	managedType.On("GetBytes", int32(2)).Return(p2Bytes, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("SetBytes", int32(3), p3Bytes).Return()

	ret := hooks.ManagedAddEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(0), ret)
	runtime.AssertNotCalled(t, "FailExecution", mock.Anything)
}

func TestManagedMulEC_Success(t *testing.T) {
	vmHooks := createHooksWithBaseSetup()
	hooks := vmHooks.hooks
	managedType := vmHooks.managedType
	runtime := vmHooks.runtime

	p, _ := bn254.NewPointG1().Pick()
	s, _ := bn254.NewScalar().Pick()
	res, _ := p.Mul(s)
	pBytes, _ := p.MarshalBinary()
	sBytes, _ := s.MarshalBinary()
	resBytes, _ := res.MarshalBinary()

	managedType.On("GetBytes", int32(1)).Return(pBytes, nil)
	managedType.On("GetBytes", int32(2)).Return(sBytes, nil)
	managedType.On("ConsumeGasForBytes", mock.Anything).Return(nil)
	managedType.On("SetBytes", int32(3), resBytes).Return()

	ret := hooks.ManagedMulEC(int32(lowLevelFeatures.BN254), int32(lowLevelFeatures.G1), 1, 2, 3)
	assert.Equal(t, int32(0), ret)
	runtime.AssertNotCalled(t, "FailExecution", mock.Anything)
}
