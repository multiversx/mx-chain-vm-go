package mock

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	"github.com/ElrondNetwork/arwen-wasm-vm/math"
	"github.com/ElrondNetwork/elrond-go/core/vmcommon"
	"github.com/stretchr/testify/require"
)

type HostGasTester struct {
	DefaultCodeLength uint64
	GasUsed           uint64
	GasUsedForAPI     uint64
	Host              arwen.VMHost
	TB                testing.TB
}

func NewHostGasTester(tb testing.TB, host arwen.VMHost, defaultCodeLength uint64) *HostGasTester {
	return &HostGasTester{
		DefaultCodeLength: defaultCodeLength,
		GasUsed:           0,
		Host:              host,
		TB:                tb,
	}
}

func (hgt *HostGasTester) UseGasForContractCode() {
	retrieveCodeCost := uint64(1)
	hgt.GasUsed = math.AddUint64(hgt.GasUsed, hgt.DefaultCodeLength)
	hgt.GasUsed = math.AddUint64(hgt.GasUsed, retrieveCodeCost)
}

func (hgt *HostGasTester) UseGas(gas uint64) {
	hgt.GasUsed = math.AddUint64(hgt.GasUsed, gas)
	err := hgt.Host.Metering().UseGasBounded(gas)
	require.Nil(hgt.TB, err)
}

func (hgt *HostGasTester) UseGasForAPI() {
	hgt.UseGas(1)
	hgt.GasUsedForAPI = math.AddUint64(hgt.GasUsedForAPI, 1)
}

func (hgt *HostGasTester) UseGasForLastFinish() {
	returnData := hgt.Host.Output().ReturnData()
	if len(returnData) == 0 {
		return
	}

	lastFinish := returnData[len(returnData)-1]
	hgt.UseGas(uint64(len(lastFinish)))
}

func (hgt *HostGasTester) Validate(vmInput *vmcommon.ContractCallInput, vmOutput *vmcommon.VMOutput) {
	gasIn := vmInput.GasProvided
	gasOut := math.AddUint64(vmOutput.GasRemaining, hgt.GasUsed)

	fmt.Println(hgt.GasUsed)

	// gasIn and gasOut are cast to int for require.Equal() so that they will be
	// printed to stdout as decimals, and not as hexadecimals.
	require.Equal(hgt.TB, int(gasIn), int(gasOut))
}
