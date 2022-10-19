package main

import (
	"fmt"
	"go/token"
	"os"

	eapigen "github.com/ElrondNetwork/wasm-vm/arwen/elrondapi/generate"
)

const pathToElrondApiPackage = "./"

func initEIMetadata() *eapigen.EIMetadata {
	return &eapigen.EIMetadata{
		Groups: []*eapigen.EIGroup{
			{SourcePath: "elrondei.go", Name: "Main"},
			{SourcePath: "managedei.go", Name: "Managed"},
			{SourcePath: "bigFloatOps.go", Name: "BigFloat"},
			{SourcePath: "bigIntOps.go", Name: "BigInt"},
			{SourcePath: "manBufOps.go", Name: "ManagedBuffer"},
			{SourcePath: "smallIntOps.go", Name: "SmallInt"},
			{SourcePath: "cryptoei.go", Name: "Crypto"},
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

	out1, err := os.Create(pathToElrondApiPackage + "../../executor/vmHooks.go")
	if err != nil {
		panic(err)
	}
	defer out1.Close()
	eapigen.WriteEIInterface(out1, eiMetadata)

	out2, err := os.Create(pathToElrondApiPackage + "../../wasmer/wasmerImportsCgo.go")
	if err != nil {
		panic(err)
	}
	defer out2.Close()
	eapigen.WriteCAPIFunctions(out2, eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))
}
