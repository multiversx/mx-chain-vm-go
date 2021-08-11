package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomnessGenerator(t *testing.T) {
	firstStr := "the usage of the same seed must generate 2 equal randomizers"
	firstSeed := []byte(firstStr)
	var firstRandomizer *seedRandReader
	require.True(t, firstRandomizer.IsInterfaceNil())
	firstRandomizer = NewSeedRandReader(firstSeed)
	require.False(t, firstRandomizer.IsInterfaceNil())

	secondStr := "the usage of the same seed must generate 2 equal randomizers"
	secondSeed := []byte(secondStr)
	var secondRandomizer *seedRandReader
	require.True(t, secondRandomizer.IsInterfaceNil())
	secondRandomizer = NewSeedRandReader(secondSeed)
	require.Equal(t, firstRandomizer, secondRandomizer)
	require.False(t, secondRandomizer.IsInterfaceNil())

	thirdStr := "the usage of two different seeds must generate 2 different randomizers"
	thirdSeed := []byte(thirdStr)
	var thirdRandomizer *seedRandReader
	require.True(t, thirdRandomizer.IsInterfaceNil())
	thirdRandomizer = NewSeedRandReader(thirdSeed)
	require.NotEqual(t, firstRandomizer, thirdRandomizer)
	require.NotEqual(t, secondRandomizer, thirdRandomizer)
	require.False(t, thirdRandomizer.IsInterfaceNil())

	a := make([]byte, 100)
	firstRandomizer.Read(a)
	require.NotEqual(t, firstRandomizer, secondRandomizer)
	b := make([]byte, 100)
	secondRandomizer.Read(b)
	require.Equal(t, a, b)
	require.Equal(t, firstRandomizer, secondRandomizer)
	c := make([]byte, 100)
	thirdRandomizer.Read(c)
	require.NotEqual(t, a, c)
}
