package arwen

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Trace represents a temporary storage used for data useful in debugging the VM and smart contracts.
type Trace struct {
}

var globalTrace = Trace{}

// PutVMOutput saves the VMOutput to a JSON file, ./trace/smart-contracts/[scAddress]/vmOutput_[timestamp].json
// If any error occurs, it will be sent to logger.
func (trace *Trace) PutVMOutput(scAddress []byte, vmOutput *vmcommon.VMOutput) {
	scAddressEncoded := hex.EncodeToString(scAddress)
	folder := prepareTraceFolder("smart-contracts", scAddressEncoded)
	saveToJSON(folder, "vmOutput", vmOutput)
}

// prepareTraceFolder creates a full path of a trace sub-folder and ensures its existence.
// The result is the full path to the sub-folder.
// If any error occurs in the creation of the sub-folder, it will be sent to logger.
func prepareTraceFolder(folderParts ...string) string {
	parentFolder := filepath.Join(".", "trace")
	subFolder := filepath.Join(folderParts...)
	fullFolderPath := filepath.Join(parentFolder, subFolder)
	err := os.MkdirAll(fullFolderPath, os.ModePerm)

	if err != nil {
		log.Printf("trace.prepareTraceFolder: could not create folder %s. %s\n", fullFolderPath, err.Error())
	}

	return fullFolderPath
}

// saveToJSON creates a file at the specified path (folder and fileNamePrefix), containing a JSON representation of the value parameter.
// It returns no error. If any error occurs, it will be sent to logger and handled silently.
func saveToJSON(parentFolder string, fileNamePrefix string, value interface{}) {
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s.json", fileNamePrefix, timestamp)
	path := filepath.Join(parentFolder, fileName)
	serialized := serializeToJSON(value)
	err := ioutil.WriteFile(path, serialized, 0644)

	if err != nil {
		log.Printf("trace.saveToJSON: could not save file %s. %s\n", path, err.Error())
	} else {
		log.Printf("trace.saveToJSON: saved file %s\n", path)
	}
}

// serializeToJSON creates a JSON representation of the value parameter.
// The JSON representation is pretty formatted.
// It returns no error. If any error occurs, it will be sent to logger, and the JSON representation will be void (empty).
func serializeToJSON(value interface{}) []byte {
	serialized, err := json.MarshalIndent(value, "", "\t")

	if err != nil {
		log.Printf("trace.serializeToJson: Could not serialize value: %s", err.Error())
		serialized = []byte{}
	}

	return serialized
}
