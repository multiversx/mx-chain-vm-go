package nodepart

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.VMExecutionHandler = (*ArwenDriver)(nil)

// ArwenDriver is
type ArwenDriver struct {
	blockchainHook vmcommon.BlockchainHook
	cryptoHook     vmcommon.CryptoHook
	vmType         []byte
	blockGasLimit  uint64
	gasSchedule    map[string]map[string]uint64

	arwenInputRead   *os.File
	arwenInputWrite  *os.File
	arwenOutputRead  *os.File
	arwenOutputWrite *os.File
	command          *exec.Cmd
	part             *NodePart
}

// NewArwenDriver creates
func NewArwenDriver(
	blockchainHook vmcommon.BlockchainHook,
	cryptoHook vmcommon.CryptoHook,
	vmType []byte,
	blockGasLimit uint64,
	gasSchedule map[string]map[string]uint64,
) (*ArwenDriver, error) {
	driver := &ArwenDriver{
		blockchainHook: blockchainHook,
		cryptoHook:     cryptoHook,
		vmType:         vmType,
		blockGasLimit:  blockGasLimit,
		gasSchedule:    gasSchedule,
	}

	err := driver.startArwen()
	return driver, err
}

func (driver *ArwenDriver) startArwen() error {
	driver.resetStreams()

	driver.command = exec.Command("arwen")
	// TODO: pass vmType, blockGasLimit and gasSchedule when starting Arwen

	driver.command.ExtraFiles = []*os.File{driver.arwenInputRead, driver.arwenOutputWrite}
	err := driver.command.Start()
	if err != nil {
		return err
	}

	driver.part, err = NewNodePart(driver.arwenOutputRead, driver.arwenInputWrite, driver.blockchainHook, driver.cryptoHook)
	if err != nil {
		return err
	}

	return nil
}

func (driver *ArwenDriver) resetStreams() error {
	closeFile(driver.arwenInputRead)
	closeFile(driver.arwenInputWrite)
	closeFile(driver.arwenOutputRead)
	closeFile(driver.arwenOutputWrite)

	var err error

	driver.arwenInputRead, driver.arwenInputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.arwenOutputRead, driver.arwenOutputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	return nil
}

func closeFile(file *os.File) {
	if file != nil {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot close file.\n")
		}
	}
}

func (driver *ArwenDriver) restartArwenIfNecessary() error {
	if !driver.command.ProcessState.Exited() {
		return nil
	}

	err := driver.startArwen()
	return err
}

// RunSmartContractCreate creates
func (driver *ArwenDriver) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	request := &common.ContractRequest{
		Action:      "Deploy",
		CreateInput: input,
	}

	response, err := driver.part.StartLoop(request)
	if err != nil {
		// TODO: if critical error, restart.
	}

	return response.VMOutput, err
}

// RunSmartContractCall calls
func (driver *ArwenDriver) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	request := &common.ContractRequest{
		Action:    "Call",
		CallInput: input,
	}

	response, err := driver.part.StartLoop(request)
	if err != nil {
		// TODO: if critical error, restart.
	}

	return response.VMOutput, err
}

// func OnRoundEnded -> triggers Arwen restart.
