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
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/logger"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.VMExecutionHandler = (*ArwenDriver)(nil)

// ArwenDriver manages the execution of the Arwen process
type ArwenDriver struct {
	nodeLogger     logger.Logger
	blockchainHook vmcommon.BlockchainHook
	config         Config
	vmType         []byte
	blockGasLimit  uint64
	gasSchedule    map[string]map[string]uint64

	arwenInputRead   *os.File
	arwenInputWrite  *os.File
	arwenOutputRead  *os.File
	arwenOutputWrite *os.File
	arwenLogRead     *os.File
	arwenLogWrite    *os.File
	command          *exec.Cmd
	part             *NodePart
}

// NewArwenDriver creates a new driver
func NewArwenDriver(
	nodeLogger logger.Logger,
	blockchainHook vmcommon.BlockchainHook,
	config Config,
	vmType []byte,
	blockGasLimit uint64,
	gasSchedule map[string]map[string]uint64,
) (*ArwenDriver, error) {
	driver := &ArwenDriver{
		nodeLogger:     nodeLogger,
		blockchainHook: blockchainHook,
		config:         config,
		vmType:         vmType,
		blockGasLimit:  blockGasLimit,
		gasSchedule:    gasSchedule,
	}

	err := driver.startArwen()
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func (driver *ArwenDriver) startArwen() error {
	driver.nodeLogger.Info("ArwenDriver.startArwen()")

	driver.resetPipeStreams()

	arwenPath, err := driver.getArwenPath()
	if err != nil {
		return err
	}

	arguments, err := common.PrepareArguments(common.Arguments{
		VMType:        driver.vmType,
		BlockGasLimit: driver.blockGasLimit,
		GasSchedule:   driver.gasSchedule,
		LogLevel:      logger.LogDebug,
	})
	if err != nil {
		return err
	}

	driver.command = exec.Command(arwenPath, arguments...)
	driver.command.ExtraFiles = []*os.File{driver.arwenInputRead, driver.arwenOutputWrite, driver.arwenLogWrite}

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

	driver.part, err = NewNodePart(
		driver.nodeLogger,
		driver.arwenOutputRead,
		driver.arwenInputWrite,
		driver.blockchainHook,
		driver.config,
	)
	if err != nil {
		return err
	}

	driver.continuouslyCopyArwenLogs(arwenStdout, arwenStderr)

	return nil
}

func (driver *ArwenDriver) getArwenPath() (string, error) {
	arwenPath := os.Getenv("ARWEN_PATH")
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
	closeFile(driver.arwenLogRead)
	closeFile(driver.arwenLogWrite)

	var err error

	driver.arwenInputRead, driver.arwenInputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.arwenOutputRead, driver.arwenOutputWrite, err = os.Pipe()
	if err != nil {
		return err
	}

	driver.arwenLogRead, driver.arwenLogWrite, err = os.Pipe()
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

// IsClosed checks whether the Arwen process is closed
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

// RunSmartContractCreate sends a deploy request to Arwen and waits for the output
func (driver *ArwenDriver) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	driver.nodeLogger.Trace("RunSmartContractCreate")

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

// RunSmartContractCall sends an execution request to Arwen and waits for the output
func (driver *ArwenDriver) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	driver.nodeLogger.Trace("RunSmartContractCall", "sc", input.RecipientAddr)

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

// DiagnoseWait sends a diagnose message to Arwen
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
	arwenLog := driver.arwenLogRead

	go func() {
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			driver.nodeLogger.Info("ARWEN-OUT", "line", line)
		}
	}()

	go func() {
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			driver.nodeLogger.Error("ARWEN-ERR", "line", line)
		}
	}()

	go func() {
		for {
			err := logger.ReceiveLogThroughPipe(driver.nodeLogger, arwenLog)
			if err != nil {
				driver.nodeLogger.Error("ReceiveLogThroughPipe error", "err", err)
				break
			}
		}
	}()
}
