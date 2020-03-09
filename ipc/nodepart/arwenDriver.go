package nodepart

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/common"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.VMExecutionHandler = (*ArwenDriver)(nil)

// ArwenDriver is
type ArwenDriver struct {
	nodeLogger     common.NodeLogger
	blockchainHook vmcommon.BlockchainHook
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
	nodeLogger common.NodeLogger,
	blockchainHook vmcommon.BlockchainHook,
	vmType []byte,
	blockGasLimit uint64,
	gasSchedule map[string]map[string]uint64,
) (*ArwenDriver, error) {
	driver := &ArwenDriver{
		nodeLogger:     nodeLogger,
		blockchainHook: blockchainHook,
		vmType:         vmType,
		blockGasLimit:  blockGasLimit,
		gasSchedule:    gasSchedule,
	}

	err := driver.startArwen()
	return driver, err
}

func (driver *ArwenDriver) startArwen() error {
	driver.nodeLogger.Info("ArwenDriver.startArwen()")

	driver.resetPipeStreams()

	arwenPath, err := driver.getArwenPath()
	if err != nil {
		return err
	}

	arguments, err := common.PrepareArguments(driver.vmType, driver.blockGasLimit, driver.gasSchedule)
	if err != nil {
		return err
	}

	driver.command = exec.Command(arwenPath, arguments...)
	driver.command.ExtraFiles = []*os.File{driver.arwenInputRead, driver.arwenOutputWrite}

	arwenStdout, err := driver.command.StdoutPipe()
	if err != nil {
		return err
	}

	arwenStderr, err := driver.command.StderrPipe()
	if err != nil {
		return err
	}

	err = driver.command.Start()
	if err != nil {
		return err
	}

	driver.part, err = NewNodePart(driver.arwenOutputRead, driver.arwenInputWrite, driver.blockchainHook)
	if err != nil {
		return err
	}

	driver.continuouslyCopyArwenLogs(arwenStdout, arwenStderr)

	return nil
}

func (driver *ArwenDriver) getArwenPath() (string, error) {
	arwenPath := os.Getenv("ARWEN_PATH")
	driver.nodeLogger.Info("ARWEN_PATH environment variable", "ARWEN_PATH", arwenPath)

	if fileExists(arwenPath) {
		return arwenPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	arwenPath = path.Join(cwd, "arwen")
	if fileExists(arwenPath) {
		return arwenPath, nil
	}

	return "", common.ErrArwenNotFound
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func (driver *ArwenDriver) resetPipeStreams() error {
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
	if !driver.IsClosed() {
		return nil
	}

	err := driver.startArwen()
	return err
}

// IsClosed returns
func (driver *ArwenDriver) IsClosed() bool {
	pid := driver.command.Process.Pid
	process, err := os.FindProcess(pid)
	if err != nil {
		return true
	}
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return true
	}
	return false
}

// RunSmartContractCreate creates
func (driver *ArwenDriver) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	driver.nodeLogger.Info("RunSmartContractCreate")

	err := driver.restartArwenIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractDeployRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.stopArwen()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	return typedResponse.VMOutput, response.GetError()
}

// RunSmartContractCall calls
func (driver *ArwenDriver) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	driver.nodeLogger.Info("RunSmartContractCall", "sc", input.RecipientAddr)

	err := driver.restartArwenIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractCallRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.nodeLogger.Error("RunSmartContractCall", "err", err)
		driver.stopArwen()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	return typedResponse.VMOutput, response.GetError()
}

// DiagnoseWait calls
func (driver *ArwenDriver) DiagnoseWait(milliseconds uint32) error {
	err := driver.restartArwenIfNecessary()
	if err != nil {
		return common.WrapCriticalError(err)
	}

	request := common.NewMessageDiagnoseWaitRequest(milliseconds)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.nodeLogger.Error("RunSmartContractCall", "err", err)
		driver.stopArwen()
		return common.WrapCriticalError(err)
	}

	return response.GetError()
}

// TODO: func OnRoundEnded -> triggers Arwen restart.

// Close stops Arwen
func (driver *ArwenDriver) Close() error {
	err := driver.stopArwen()
	if err != nil {
		return err
	}

	return nil
}

func (driver *ArwenDriver) stopArwen() error {
	err := driver.command.Process.Kill()
	driver.command.Process.Wait()
	if err != nil {
		driver.nodeLogger.Error("stopArwen error=%s", err)
	}

	return err
}

func (driver *ArwenDriver) continuouslyCopyArwenLogs(arwenStdout io.Reader, arwenStderr io.Reader) {
	stdoutReader := bufio.NewReader(arwenStdout)
	stderrReader := bufio.NewReader(arwenStderr)

	go func() {
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			driver.nodeLogger.Info(line)
		}
	}()

	go func() {
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			driver.nodeLogger.Error(line)
		}
	}()
}
