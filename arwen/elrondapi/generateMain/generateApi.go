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
			&eapigen.EIGroup{SourcePath: "bigFloatOps.go"},
			&eapigen.EIGroup{SourcePath: "bigIntOps.go"},
			&eapigen.EIGroup{SourcePath: "elrondei.go"},
			&eapigen.EIGroup{SourcePath: "generateOps.go"},
			&eapigen.EIGroup{SourcePath: "managedei.go"},
			&eapigen.EIGroup{SourcePath: "manBufOps.go"},
			&eapigen.EIGroup{SourcePath: "smallIntOps.go"},
			&eapigen.EIGroup{SourcePath: "cryptoei.go"},
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
	eapigen.WriteEIInterface(out1, eiMetadata)

	out2, err := os.Create(pathToElrondApiPackage + "../../wasmer/wasmerImportsCgo.go")
	if err != nil {
		panic(err)
	}
	defer out2.Close()
	eapigen.WriteCAPIFunctions(out2, eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))
}
