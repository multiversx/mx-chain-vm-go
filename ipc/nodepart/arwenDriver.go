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
	driverLogger        logger.Logger
	arwenMainLogger     logger.Logger
	dialogueLogger      logger.Logger
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

	// TODO: Encapsulate to PipeLoggerSinkPart / PipeLoggerSourcePart
	arwenLogRead     *os.File
	arwenLogWrite    *os.File
	dialogueLogRead  *os.File
	dialogueLogWrite *os.File

	counterDeploy uint64
	counterCall   uint64

	command *exec.Cmd
	part    *NodePart
}

// NewArwenDriver creates a new driver
func NewArwenDriver(
	driverLogger logger.Logger,
	arwenMainLogger logger.Logger,
	dialogueLogger logger.Logger,
	blockchainHook vmcommon.BlockchainHook,
	arwenArguments common.ArwenArguments,
	config Config,
) (*ArwenDriver, error) {
	driver := &ArwenDriver{
		driverLogger:        driverLogger,
		arwenMainLogger:     arwenMainLogger,
		dialogueLogger:      dialogueLogger,
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
	driver.driverLogger.Info("ArwenDriver.startArwen()")

	err := driver.resetPipeStreams()
	if err != nil {
		return err
	}

	arwenPath, err := driver.getArwenPath()
	if err != nil {
		return err
	}

	driver.command = exec.Command(arwenPath)
	driver.command.ExtraFiles = []*os.File{
		driver.arwenInitRead,
		driver.arwenInputRead,
		driver.arwenOutputWrite,
		driver.arwenLogWrite,
		driver.dialogueLogWrite,
	}

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
		driver.driverLogger,
		driver.dialogueLogger,
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
	arwenPath, err := driver.getArwenPathInCurrentDirectory()
	if err == nil {
		return arwenPath, nil
	}

	arwenPath, err = driver.getArwenPathFromEnvironment()
	if err == nil {
		return arwenPath, nil
	}

	return "", common.ErrArwenNotFound
}

func (driver *ArwenDriver) getArwenPathInCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	arwenPath := path.Join(cwd, "arwen")
	if fileExists(arwenPath) {
		return arwenPath, nil
	}

	return "", common.ErrArwenNotFound
}

func (driver *ArwenDriver) getArwenPathFromEnvironment() (string, error) {
	arwenPath := os.Getenv(common.EnvVarArwenPath)
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

	// TODO: Encapsulate logger-pipes (see above TODO)
	closeFile(driver.arwenLogRead)
	closeFile(driver.arwenLogWrite)
	closeFile(driver.dialogueLogRead)
	closeFile(driver.dialogueLogWrite)

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

	driver.dialogueLogRead, driver.dialogueLogWrite, err = os.Pipe()
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
	return err != nil
}

// RunSmartContractCreate sends a deploy request to Arwen and waits for the output
func (driver *ArwenDriver) RunSmartContractCreate(input *vmcommon.ContractCreateInput) (*vmcommon.VMOutput, error) {
	driver.counterDeploy++
	driver.driverLogger.Trace("RunSmartContractCreate", "counter", driver.counterDeploy)

	err := driver.RestartArwenIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractDeployRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.driverLogger.Error("RunSmartContractCreate", "err", err)
		_ = driver.Close()
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
	driver.counterCall++
	driver.driverLogger.Trace("RunSmartContractCall", "counter", driver.counterCall, "func", input.Function, "sc", input.RecipientAddr)

	err := driver.RestartArwenIfNecessary()
	if err != nil {
		return nil, common.WrapCriticalError(err)
	}

	request := common.NewMessageContractCallRequest(input)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.driverLogger.Error("RunSmartContractCall", "err", err)
		_ = driver.Close()
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
	err := driver.RestartArwenIfNecessary()
	if err != nil {
		return common.WrapCriticalError(err)
	}

	request := common.NewMessageDiagnoseWaitRequest(milliseconds)
	response, err := driver.part.StartLoop(request)
	if err != nil {
		driver.driverLogger.Error("DiagnoseWait", "err", err)
		_ = driver.Close()
		return common.WrapCriticalError(err)
	}

	return response.GetError()
}

// Close stops Arwen
func (driver *ArwenDriver) Close() error {
	err := driver.stopArwen()
	if err != nil {
		driver.driverLogger.Error("ArwenDriver.Close()", "err", err)
		return err
	}

	return nil
}

func (driver *ArwenDriver) stopArwen() error {
	err := driver.command.Process.Kill()
	if err != nil {
		return err
	}

	_, err = driver.command.Process.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (driver *ArwenDriver) continuouslyCopyArwenLogs(arwenStdout io.Reader, arwenStderr io.Reader) {
	stdoutReader := bufio.NewReader(arwenStdout)
	stderrReader := bufio.NewReader(arwenStderr)
	arwenLog := driver.arwenLogRead
	dialogueLog := driver.dialogueLogRead

	go func() {
		for {
			line, err := stdoutReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			driver.arwenMainLogger.Info("ARWEN-OUT", "line", line)
		}
	}()

	go func() {
		for {
			line, err := stderrReader.ReadString('\n')
			if err != nil {
				break
			}

			line = strings.TrimSpace(line)
			driver.arwenMainLogger.Error("ARWEN-ERR", "line", line)
		}
	}()

	go func() {
		for {
			err := logger.ReceiveLogThroughPipe(driver.arwenMainLogger, arwenLog, driver.logsMarshalizer)
			if err != nil {
				driver.driverLogger.Error("ReceiveLogThroughPipe error (arwenMainLogger)", "err", err)
				break
			}
		}
	}()

	go func() {
		for {
			err := logger.ReceiveLogThroughPipe(driver.dialogueLogger, dialogueLog, driver.logsMarshalizer)
			if err != nil {
				driver.driverLogger.Error("ReceiveLogThroughPipe error (dialogueLogger)", "err", err)
				break
			}
		}
	}()
}
