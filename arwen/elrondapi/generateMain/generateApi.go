package main

import (
	"go/token"
	"os"

	eapigen "github.com/ElrondNetwork/wasm-vm/arwen/elrondapi/generate"
)

const pathToElrondApiPackage = "./"

func initEIMetadata() *eapigen.EIMetadata {
	m := make(map[string]*eapigen.EIFileMetadata)
	m["bigFloatOps.go"] = &eapigen.EIFileMetadata{}
	m["bigIntOps.go"] = &eapigen.EIFileMetadata{}
	m["elrondei.go"] = &eapigen.EIFileMetadata{}
	m["generateOps.go"] = &eapigen.EIFileMetadata{}
	m["managedei.go"] = &eapigen.EIFileMetadata{}
	m["manBufOps.go"] = &eapigen.EIFileMetadata{}
	m["smallIntOps.go"] = &eapigen.EIFileMetadata{}
	return &eapigen.EIMetadata{
		FileMap:      m,
		AllFunctions: nil,
	}
}

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	fset := token.NewFileSet() // positions are relative to fset
	eiMetadata := initEIMetadata()
	err := eapigen.ReadAndParseEIMetadata(fset, pathToElrondApiPackage, eiMetadata)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(pathToElrondApiPackage + "../../executor/executorImportsInterface.go")
	if err != nil {
		panic(err)
	}
	eapigen.WriteEIInterface(eiMetadata, out)
	out.Close()
}
