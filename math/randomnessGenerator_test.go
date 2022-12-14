package math

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomnessGenerator(t *testing.T) {
	t.Parallel()
	firstStr := "the usage of the same seed must generate 2 equal randomizers"
	firstSeed := []byte(firstStr)
	var firstRandomizer *seedRandReader
	require.True(t, firstRandomizer.IsInterfaceNil())
	firstRandomizer = NewSeedRandReader(firstSeed)
	require.False(t, firstRandomizer.IsInterfaceNil())

	secondStr := "the usage of the same seed must generate 2 equal randomizers"
	secondSeed := []byte(secondStr)
	secondRandomizer := NewSeedRandReader(secondSeed)
	require.Equal(t, firstRandomizer, secondRandomizer)

	thirdStr := "the usage of two different seeds must generate 2 different randomizers"
	thirdSeed := []byte(thirdStr)
	var thirdRandomizer *seedRandReader
	require.True(t, thirdRandomizer.IsInterfaceNil())
	thirdRandomizer = NewSeedRandReader(thirdSeed)
	require.NotEqual(t, firstRandomizer, thirdRandomizer)
	require.NotEqual(t, secondRandomizer, thirdRandomizer)
	require.False(t, thirdRandomizer.IsInterfaceNil())

	a := make([]byte, 100)
	_, _ = firstRandomizer.Read(a)
	require.NotEqual(t, firstRandomizer, secondRandomizer)
	b := make([]byte, 100)
	_, _ = secondRandomizer.Read(b)
	require.Equal(t, a, b)
	require.Equal(t, firstRandomizer, secondRandomizer)
	c := make([]byte, 100)
	_, _ = thirdRandomizer.Read(c)
	require.NotEqual(t, a, c)

	length, err := thirdRandomizer.Read(nil)
	require.Nil(t, err)
	require.Equal(t, 0, length)

	c = make([]byte, 0)
	length, err = thirdRandomizer.Read(c)
	require.Nil(t, err)
	require.Equal(t, 0, length)
}

func TestBackwardsCompatible(t *testing.T) {
	t.Parallel()
	str := "Backwards compatible test string"
	seed := []byte(str)
	randomizer := NewSeedRandReader(seed)
	a := make([]byte, 32)
	_, _ = randomizer.Read(a)
	require.Equal(t, "7459d163b20b5b0269ce2211a2cc061cc9e512fdcbe025b0fa359014f6619ed0", hex.EncodeToString(a))
}
