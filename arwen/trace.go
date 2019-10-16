package arwen

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Trace represents a temporary storage used for data useful in debugging a smart contract.
type Trace struct {
}

var globalTrace = Trace{}

// PutVMOutput saves the VMOutput to a json file.
func (trace *Trace) PutVMOutput(scAddress []byte, vmOutput *vmcommon.VMOutput) {
	scAddressEncoded := hex.EncodeToString(scAddress)
	fileName := fmt.Sprintf("%s.json", scAddressEncoded)
	path := trace.createTracePath(fileName)
	fmt.Printf("Trace.PutVMOutput: save to file %s\n", path)

	serialized, _ := json.MarshalIndent(vmOutput, "", "\t")
	err := ioutil.WriteFile(path, serialized, 0644)

	if err != nil {
		fmt.Printf("Trace.PutVMOutput: could not save file, %s\n", err.Error())
	}
}

func (trace *Trace) createTracePath(fileName string) string {
	folder := filepath.Join(".", "Trace")
	os.MkdirAll(folder, os.ModePerm)

	path := filepath.Join(folder, fileName)
	return path
}
