package vmhost

import (
	"math/big"

	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// AsyncCallResults holds the results of an async call
type AsyncCallResults struct {
	ReturnData     [][]byte
	TokenTransfers []*vmcommon.ESDTTransfer
	ReturnCode     vmcommon.ReturnCode
}

// FinishedAsyncCall holds the final status and results of an async call
type FinishedAsyncCall struct {
	InitialCall *AsyncCallResults
	Callback    *AsyncCallResults
}

// NewFinishedAsyncCallFromVMOutput creates a new FinishedAsyncCall from a VMOutput
func NewFinishedAsyncCallFromVMOutput(vmOutput *vmcommon.VMOutput, esdtParser vmcommon.ESDTTransferParser) *FinishedAsyncCall {
	transfers := make([]*vmcommon.ESDTTransfer, 0)
	for _, acc := range vmOutput.OutputAccounts {
		for _, outTransfer := range acc.OutputTransfers {
			if outTransfer.Value.Cmp(big.NewInt(0)) > 0 {
				continue
			}
			parsed, err := esdtParser.ParseESDTTransfers(outTransfer.SenderAddress, acc.Address, string(outTransfer.Data), nil)
			if err != nil {
				continue
			}
			transfers = append(transfers, parsed.ESDTTransfers...)
		}
	}

	return &FinishedAsyncCall{
		InitialCall: &AsyncCallResults{
			ReturnData:     vmOutput.ReturnData,
			TokenTransfers: transfers,
			ReturnCode:     vmOutput.ReturnCode,
		},
		Callback: nil,
	}
}

// Merge merges two FinishedAsyncCall objects
func (fac *FinishedAsyncCall) Merge(other *FinishedAsyncCall) {
	if other == nil {
		return
	}
	if fac.InitialCall == nil {
		fac.InitialCall = other.InitialCall
	}
	if fac.Callback == nil {
		fac.Callback = other.Callback
	}
}

func (fac *FinishedAsyncCall) ToSerializable() *SerializableFinishedAsyncCall {
	return &SerializableFinishedAsyncCall{
		InitialCall: fac.InitialCall.ToSerializable(),
		Callback:    fac.Callback.ToSerializable(),
	}
}

func (acr *AsyncCallResults) ToSerializable() *SerializableAsyncCallResults {
	if acr == nil {
		return nil
	}
	return &SerializableAsyncCallResults{
		ReturnData: acr.ReturnData,
		ReturnCode: uint32(acr.ReturnCode),
	}
}

func fromSerializableFinishedAsyncCall(sfac *SerializableFinishedAsyncCall) *FinishedAsyncCall {
	return &FinishedAsyncCall{
		InitialCall: fromSerializableAsyncCallResults(sfac.InitialCall),
		Callback:    fromSerializableAsyncCallResults(sfac.Callback),
	}
}

func fromSerializableAsyncCallResults(sacr *SerializableAsyncCallResults) *AsyncCallResults {
	if sacr == nil {
		return nil
	}
	return &AsyncCallResults{
		ReturnData: sacr.ReturnData,
		ReturnCode: vmcommon.ReturnCode(sacr.ReturnCode),
	}
}
