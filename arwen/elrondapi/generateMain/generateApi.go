package main

import (
	"go/token"
	"os"

	eapigen "github.com/ElrondNetwork/wasm-vm/arwen/elrondapi/generate"
)

const pathToElrondApiPackage = "./"

func initEIMetadata() *eapigen.EIMetadata {
	return &eapigen.EIMetadata{
		Groups: []*eapigen.EIGroup{
			{SourcePath: "bigFloatOps.go"},
			{SourcePath: "bigIntOps.go"},
			{SourcePath: "elrondei.go"},
			{SourcePath: "generateOps.go"},
			{SourcePath: "managedei.go"},
			{SourcePath: "manBufOps.go"},
			{SourcePath: "smallIntOps.go"},
			{SourcePath: "../cryptoapi/cryptoei.go"},
		},
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

	out1, err := os.Create(pathToElrondApiPackage + "../../executor/executorImportsInterface.go")
	if err != nil {
		panic(err)
	}
	defer out1.Close()
	eapigen.WriteEIInterface(eiMetadata, out1)

	out2, err := os.Create(pathToElrondApiPackage + "../../wasmer/wasmerImportsCgo.go")
	if err != nil {
		panic(err)
	}
	defer out2.Close()
	eapigen.WriteCAPIFunctions(eiMetadata, out2)
}
