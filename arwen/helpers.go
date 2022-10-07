package arwen

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"unsafe"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/wasm-vm/crypto"
	"github.com/ElrondNetwork/wasm-vm/math"
	"github.com/ElrondNetwork/wasm-vm/wasmer"
	"github.com/pelletier/go-toml"
)

// Zero is the big integer 0
var Zero = big.NewInt(0)

// One is the big integer 1
var One = big.NewInt(1)

// CustomStorageKey generates a storage key of a specific type.
func CustomStorageKey(keyType string, associatedKey []byte) []byte {
	return append([]byte(keyType), associatedKey...)
}

// BooleanToInt returns 1 if the given bool is true, 0 otherwise
func BooleanToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// GuardedGetBytesSlice extracts a subslice from a given slice, guarding against overstepping the bounds.
func GuardedGetBytesSlice(data []byte, offset int32, length int32) ([]byte, error) {
	dataLength := uint32(len(data))
	isOffsetTooSmall := offset < 0
	isOffsetTooLarge := uint32(offset) > dataLength
	requestedEnd := math.AddInt32(offset, length)
	isRequestedEndTooLarge := uint32(requestedEnd) > dataLength
	isLengthNegative := length < 0

	if isOffsetTooSmall || isOffsetTooLarge || isRequestedEndTooLarge {
		return nil, fmt.Errorf("GuardedGetBytesSlice: bad bounds")
	}

	if isLengthNegative {
		return nil, fmt.Errorf("GuardedGetBytesSlice: negative length")
	}

	result := data[offset:requestedEnd]
	return result, nil
}

// PadBytesLeft adds a specified number of zeros to the left of a byte slice.
func PadBytesLeft(data []byte, size int) []byte {
	if data == nil {
		return nil
	}
	if len(data) == 0 {
		return []byte{}
	}
	padSize := math.SubInt(size, len(data))
	if padSize <= 0 {
		return data
	}

	paddedBytes := make([]byte, padSize)
	paddedBytes = append(paddedBytes, data...)
	return paddedBytes
}

// InverseBytes reverses the order of a byte slice.
func InverseBytes(data []byte) []byte {
	length := len(data)
	invBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		invBytes[length-i-1] = data[i]
	}
	return invBytes
}

// GetSCCode retrieves the bytecode of a WASM module from a file.
func GetSCCode(fileName string) []byte {
	code, err := ioutil.ReadFile(filepath.Clean(fileName))
	if err != nil {
		panic(fmt.Sprintf("GetSCCode(): %s", fileName))
	}
	return code
}

// GetTestSCCode retrieves the bytecode of a WASM testing module.
func GetTestSCCode(scName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + scName + ".wasm"
	return GetSCCode(pathToSC)
}

// GetTestSCCodeModule retrieves the bytecode of a WASM testing contract, given
// a specific name of the WASM module
func GetTestSCCodeModule(scName string, moduleName string, prefixToTestSCs string) []byte {
	pathToSC := prefixToTestSCs + "test/contracts/" + scName + "/output/" + moduleName + ".wasm"
	return GetSCCode(pathToSC)
}

// OpenFile method opens the file from given path - does not close the file
func OpenFile(relativePath string) (*os.File, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		fmt.Printf("cannot create absolute path for the provided file: %s", err.Error())
		return nil, err
	}
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	return f, nil
}

// LoadTomlFileToMap opens and decodes a toml file as a map[string]interface{}
func LoadTomlFileToMap(relativePath string) (map[string]interface{}, error) {
	f, err := OpenFile(relativePath)
	if err != nil {
		return nil, err
	}

	fileinfo, err := f.Stat()
	if err != nil {
		fmt.Printf("cannot stat file: %s", err.Error())
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = f.Read(buffer)
	if err != nil {
		fmt.Printf("cannot read from file: %s", err.Error())
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Printf("cannot close file: %s", err.Error())
		}
	}()

	loadedTree, err := toml.Load(string(buffer))
	if err != nil {
		fmt.Printf("cannot interpret file contents as toml: %s", err.Error())
		return nil, err
	}

	loadedMap := loadedTree.ToMap()

	return loadedMap, nil
}

// SetPlainLoggerFormatter configures the logger to output only ASCII characters
func SetPlainLoggerFormatter() {
	logger.ClearLogObservers()
	_ = logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
}

// SetLoggingForTests configures the logger package with *:TRACE and enabled logger names
func SetLoggingForTests() {
	SetLoggingForTestsWithLogger("*")
}

// SetLoggingForTestsWithLogger configures the logger package with a certain logger
func SetLoggingForTestsWithLogger(loggerName string) {
	_ = logger.SetLogLevel(fmt.Sprintf("*:NONE,%s:TRACE", loggerName))
	logger.ToggleCorrelation(false)
	logger.ToggleLoggerName(true)
}

// DisableLoggingForTests sets log level to *:NONE
func DisableLoggingForTests() {
	_ = logger.SetLogLevel("*:NONE")
}

// U64ToLEB128 encodes an uint64 using LEB128 (Little Endian Base 128), used in WASM bytecode
// See https://en.wikipedia.org/wiki/LEB128
// Copied from https://github.com/filecoin-project/go-leb128/blob/master/leb128.go
func U64ToLEB128(n uint64) (out []byte) {
	more := true
	for more {
		b := byte(n & 0x7F)
		n >>= 7
		if n == 0 {
			more = false
		} else {
			b |= 0x80
		}
		out = append(out, b)
	}
	return
}

// IfNil tests if the provided interface pointer or underlying object is nil
func IfNil(checker nilInterfaceChecker) bool {
	if checker == nil {
		return true
	}
	return checker.IsInterfaceNil()
}

type nilInterfaceChecker interface {
	IsInterfaceNil() bool
}

// GetVMHost returns the vm Context from the vm context map
func GetVMHost(vmHostPtr unsafe.Pointer) VMHost {
	instCtx := wasmer.IntoInstanceContext(vmHostPtr)
	var ptr = *(*uintptr)(instCtx.Data())
	return *(*VMHost)(unsafe.Pointer(ptr))
}

// GetBlockchainContext returns the blockchain context
func GetBlockchainContext(vmHostPtr unsafe.Pointer) BlockchainContext {
	return GetVMHost(vmHostPtr).Blockchain()
}

// GetRuntimeContext returns the runtime context
func GetRuntimeContext(vmHostPtr unsafe.Pointer) RuntimeContext {
	host := GetVMHost(vmHostPtr)
	return host.Runtime()
}

// GetCryptoContext returns the crypto context
func GetCryptoContext(vmHostPtr unsafe.Pointer) crypto.VMCrypto {
	return GetVMHost(vmHostPtr).Crypto()
}

// GetManagedTypesContext returns the big int context
func GetManagedTypesContext(vmHostPtr unsafe.Pointer) ManagedTypesContext {
	return GetVMHost(vmHostPtr).ManagedTypes()
}

// GetOutputContext returns the output context
func GetOutputContext(vmHostPtr unsafe.Pointer) OutputContext {
	return GetVMHost(vmHostPtr).Output()
}

// GetMeteringContext returns the metering context
func GetMeteringContext(vmHostPtr unsafe.Pointer) MeteringContext {
	return GetVMHost(vmHostPtr).Metering()
}

// GetStorageContext returns the storage context
func GetStorageContext(vmHostPtr unsafe.Pointer) StorageContext {
	return GetVMHost(vmHostPtr).Storage()
}

// WithFault handles an error, taking into account whether it should completely
// fail the execution of a contract or not.
func WithFault(err error, vmHostPtr unsafe.Pointer, failExecution bool) bool {
	runtime := GetVMHost(vmHostPtr)
	return WithFaultAndHost(runtime, err, failExecution)
}

// WithFaultAndHost fails the execution with the provided error
func WithFaultAndHost(host VMHost, err error, failExecution bool) bool {
	if err == nil {
		return false
	}

	if failExecution {
		runtime := host.Runtime()
		metering := host.Metering()
		metering.UseGas(metering.GasLeft())
		runtime.FailExecution(err)
	}

	return true
}
