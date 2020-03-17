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
	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/marshaling"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.VMExecutionHandler = (*ArwenDriver)(nil)

// ArwenDriver manages the execution of the Arwen process
type ArwenDriver struct {
	nodeLogger          logger.Logger
	blockchainHook      vmcommon.BlockchainHook
	arwenArguments      common.ArwenArguments
	config              Config
	logsMarshalizer     marshaling.Marshalizer
	messagesMarshalizer marshaling.Marshalizer

	arwenInitRead    *os.File
	arwenInitWrite   *os.File
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
	arwenArguments common.ArwenArguments,
	config Config,
) (*ArwenDriver, error) {
	driver := &ArwenDriver{
		nodeLogger:          nodeLogger,
		blockchainHook:      blockchainHook,
		arwenArguments:      arwenArguments,
		config:              config,
		logsMarshalizer:     marshaling.CreateMarshalizer(arwenArguments.LogsMarshalizer),
		messagesMarshalizer: marshaling.CreateMarshalizer(arwenArguments.MessagesMarshalizer),
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

	driver.command = exec.Command(arwenPath)
	driver.command.ExtraFiles = []*os.File{driver.arwenInitRead, driver.arwenInputRead, driver.arwenOutputWrite, driver.arwenLogWrite}

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

	err = common.SendArwenArguments(driver.arwenInitWrite, driver.arwenArguments)
	if err != nil {
		return err
	}

	driver.part, err = NewNodePart(
		driver.nodeLogger,
		driver.arwenOutputRead,
		driver.arwenInputWrite,
		driver.blockchainHook,
		driver.config,
		driver.messagesMarshalizer,
	)
	if err != nil {
		return err
	}

	driver.continuouslyCopyArwenLogs(arwenStdout, arwenStderr)

	return nil
}

func (driver *ArwenDriver) getArwenPath() (string, error) {
	arwenPath := os.Getenv(common.EnvVarArwenPath)
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
	closeFile(driver.arwenInitRead)
	closeFile(driver.arwenInitWrite)
	closeFile(driver.arwenInputRead)
	closeFile(driver.arwenInputWrite)
	closeFile(driver.arwenOutputRead)
	closeFile(driver.arwenOutputWrite)
	closeFile(driver.arwenLogRead)
	closeFile(driver.arwenLogWrite)

	var err error

	driver.arwenInitRead, driver.arwenInitWrite, err = os.Pipe()
	if err != nil {
		return err
	}

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

// RestartArwenIfNecessary restarts Arwen if the process is closed
// TODO: This has to be called on Node's behalf when a critical error is encountered while processing a smart contract transaction.
// The basic idea is that the node should not wait for Arwen to restart,
// but process other types of transactions, and only when Arwen is ready should go with the smart contract transactions.
func (driver *ArwenDriver) RestartArwenIfNecessary() error {
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

	request := common.NewMessageContractDeployRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.stopArwen()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	vmOutput, err := typedResponse.VMOutput, response.GetError()
	if err != nil {
		return nil, err
	}

	return vmOutput, nil
}

// RunSmartContractCall sends an execution request to Arwen and waits for the output
func (driver *ArwenDriver) RunSmartContractCall(input *vmcommon.ContractCallInput) (*vmcommon.VMOutput, error) {
	driver.nodeLogger.Trace("RunSmartContractCall", "sc", input.RecipientAddr)

	request := common.NewMessageContractCallRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.nodeLogger.Error("RunSmartContractCall", "err", err)
		driver.stopArwen()
		return nil, common.WrapCriticalError(err)
	}

	typedResponse := response.(*common.MessageContractResponse)
	vmOutput, err := typedResponse.VMOutput, response.GetError()
	if err != nil {
		return nil, err
	}

	return vmOutput, nil
}

// DiagnoseWait sends a diagnose message to Arwen
func (driver *ArwenDriver) DiagnoseWait(milliseconds uint32) error {
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
			// TODO: refactor to struct / component
			err := logger.ReceiveLogThroughPipe(driver.nodeLogger, arwenLog, driver.logsMarshalizer)
			if err != nil {
				driver.nodeLogger.Error("ReceiveLogThroughPipe error", "err", err)
				break
			}
		}
	}()
}
