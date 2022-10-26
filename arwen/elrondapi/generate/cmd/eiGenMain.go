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

	writeVMHooks(eiMetadata)
	writeWasmer1ImportsCgo(eiMetadata)
	writeWasmer2ImportsCgo(eiMetadata)
	writeRustVMHooksTrait(eiMetadata)
	writeRustVMHooksPointers(eiMetadata)
	writeRustWasmerImports(eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))

	writeRustOpcodeCost()
	writeRustWasmerMeteringHelpers()

	fmt.Println("Generated code for opcodes and metering helpers")
}

func writeVMHooks(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(pathToElrondApiPackage + "../../executor/vmHooks.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteEIInterface(out, eiMetadata)
}

func writeWasmer1ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(pathToElrondApiPackage + "../../wasmer/wasmerImportsCgo.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteWasmer1Cgo(out, eiMetadata)
}

func writeWasmer2ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(pathToElrondApiPackage + "../../wasmer2/wasmer2ImportsCgo.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteWasmer2Cgo(out, eiMetadata)
}

func writeRustVMHooksTrait(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/vm_hooks.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustVMHooksTrait(out, eiMetadata)
}

func writeRustVMHooksPointers(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/vm_hook_pointers.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustVMHooksPointers(out, eiMetadata)
}

func writeRustWasmerImports(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/wasmer_imports.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustWasmerImports(out, eiMetadata)
}

func writeRustOpcodeCost() {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/opcodes_cost.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustOpcodeCost(out)
}

func writeRustWasmerMeteringHelpers() {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/wasmer_metering_helpers.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustWasmerMeteringHelpers(out)
}
