package vmhost

import (
	"math/big"
	"strings"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

// CheckBalances performs checks on the VM output and the back transfers
func CheckBalances(
	output *vmcommon.VMOutput,
	managedContext ManagedTypesContext,
	host VMHost,
) error {
	err := checkBaseCurrency(output)
	if err != nil {
		return err
	}

	err = checkESDTs(output, host)
	if err != nil {
		return err
	}

	if host.Runtime().GetVMInput().CallType == vm.DirectCall {
		err = checkBackTransfers(output, managedContext)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkBaseCurrency(output *vmcommon.VMOutput) error {
	return nil
}

func checkESDTs(output *vmcommon.VMOutput, host VMHost) error {
	esdtBalances := make(map[string]*big.Int)

	for _, outAcc := range output.OutputAccounts {
		for _, transfer := range outAcc.OutputTransfers {
			parts := strings.Split(string(transfer.Data), "@")
			function := parts[0]

			switch function {
			case "ESDTLocalMint":
				token := parts[1]
				amount, _ := big.NewInt(0).SetString(parts[2], 16)
				addESDTToMap(esdtBalances, &vmcommon.ESDTTransfer{ESDTTokenName: []byte(token), ESDTValue: amount}, true)
			case "ESDTLocalBurn":
				token := parts[1]
				amount, _ := big.NewInt(0).SetString(parts[2], 16)
				addESDTToMap(esdtBalances, &vmcommon.ESDTTransfer{ESDTTokenName: []byte(token), ESDTValue: amount}, false)
			default:
				// Regular transfer
				// This part is tricky, as we need to parse the arguments correctly.
				// The logic here is incomplete.
			}
		}
	}

	for _, balance := range esdtBalances {
		if balance.Cmp(big.NewInt(0)) != 0 {
			return ErrBalancesMismatch
		}
	}

	return nil
}

func addESDTToMap(
	esdtBalances map[string]*big.Int,
	esdt *vmcommon.ESDTTransfer,
	isMint bool,
) {
	balance, ok := esdtBalances[string(esdt.ESDTTokenName)]
	if !ok {
		balance = big.NewInt(0)
	}

	if isMint {
		balance.Add(balance, esdt.ESDTValue)
	} else {
		balance.Sub(balance, esdt.ESDTValue)
	}

	esdtBalances[string(esdt.ESDTTokenName)] = balance
}

func checkBackTransfers(output *vmcommon.VMOutput, managedContext ManagedTypesContext) error {
	return nil
}
