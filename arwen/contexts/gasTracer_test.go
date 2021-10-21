package contexts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGasTracer(t *testing.T) {
	gasTracer := NewEnabledGasTracer()
	require.Equal(t, 0, len(gasTracer.gasTrace))
	require.Equal(t, "", gasTracer.functionNameTraced)
	require.Equal(t, "", gasTracer.scAddress)
	require.False(t, gasTracer.IsInterfaceNil())

	gasTracer.setCurrentFunctionTraced("functionName")
	require.Equal(t, "functionName", gasTracer.functionNameTraced)
	gasTracer.setCurrentScAddressTraced("scAddress")
	require.Equal(t, "scAddress", gasTracer.scAddress)

	scAddress1 := "scAddress1"
	function1 := "function1"

	gasTracer.BeginTrace(scAddress1, function1)
	require.Equal(t, function1, gasTracer.functionNameTraced)
	require.Equal(t, scAddress1, gasTracer.scAddress)
	require.Equal(t, 1, len(gasTracer.gasTrace))
	require.Equal(t, 1, len(gasTracer.gasTrace[scAddress1]))
	require.Equal(t, 1, len(gasTracer.gasTrace[scAddress1][function1]))
	require.Equal(t, uint64(0), gasTracer.gasTrace[scAddress1][function1][0])

	gasTracer.AddToCurrentTrace(uint64(2000))
	require.Equal(t, uint64(2000), gasTracer.gasTrace[scAddress1][function1][0])
	gasTracer.AddToCurrentTrace(uint64(4000))
	require.Equal(t, uint64(6000), gasTracer.gasTrace[scAddress1][function1][0])

	gasTracer.AddTracedGas(scAddress1, function1, uint64(3000))
	require.Equal(t, 2, len(gasTracer.gasTrace[scAddress1][function1]))
	require.Equal(t, uint64(3000), gasTracer.gasTrace[scAddress1][function1][1])

	gasTracer.BeginTrace(scAddress1, function1)
	require.Equal(t, 3, len(gasTracer.gasTrace[scAddress1][function1]))

	function2 := "function2"

	gasTracer.BeginTrace(scAddress1, function2)
	require.Equal(t, 2, len(gasTracer.gasTrace[scAddress1]))
	require.Equal(t, "function2", gasTracer.functionNameTraced)

	scAddress2 := "scAddress2"

	gasTracer.AddTracedGas(scAddress2, function2, uint64(4800))
	require.Equal(t, 2, len(gasTracer.gasTrace))
}
