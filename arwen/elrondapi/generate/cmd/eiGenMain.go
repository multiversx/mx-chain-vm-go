package main

import (
	"errors"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	eapigen "github.com/ElrondNetwork/wasm-vm/arwen/elrondapi/generate"
)

const pathToElrondApiPackage = "./"
const pathToRustRepoConfigFile = "wasm-vm-executor-rs-path.txt"

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
	writeWasmer2Names(eiMetadata)

	tryCreateRustOutputDirectory()

	writeRustVMHooksTrait(eiMetadata)
	writeRustCapiVMHooks(eiMetadata)
	writeRustCapiVMHooksPointers(eiMetadata)
	writeRustWasmerImports(eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))

	writeWASMOpcodeCost()
	writeRustWasmerOpcodeCost()
	writeRustWasmerMeteringHelpers()

	fmt.Println("Generated code for opcodes and metering helpers.")

	tryCopyFilesToRustExecutorRepo()
}

func writeVMHooks(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "../../executor/vmHooks.go"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteEIInterface(out, eiMetadata)
}

func writeWasmer1ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "../../wasmer/wasmerImportsCgo.go"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteWasmer1Cgo(out, eiMetadata)
}

func writeWasmer2ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "../../wasmer2/wasmer2ImportsCgo.go"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteWasmer2Cgo(out, eiMetadata)
}

func writeWasmer2Names(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "../../wasmer2/wasmer2Names.go"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteNames(out, eiMetadata)
}

func tryCreateRustOutputDirectory() {
	outputDirPath := filepath.Join(pathToElrondApiPackage, "generate/cmd/output")
	if _, err := os.Stat(outputDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(outputDirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created output directory.")
		return
	}
	fmt.Println("Output directory already exists.")
}

func writeRustVMHooksTrait(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "generate/cmd/output/vm_hooks.rs"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustVMHooksTrait(out, eiMetadata)
}

func writeRustCapiVMHooks(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "generate/cmd/output/capi_vm_hook.rs"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustCapiVMHooks(out, eiMetadata)
}

func writeRustCapiVMHooksPointers(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "generate/cmd/output/capi_vm_hook_pointers.rs"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustCapiVMHooksPointers(out, eiMetadata)
}

func writeRustWasmerImports(eiMetadata *eapigen.EIMetadata) {
	out, err := os.Create(filepath.Join(pathToElrondApiPackage, "generate/cmd/output/wasmer_imports.rs"))
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustWasmerImports(out, eiMetadata)
}

func writeWASMOpcodeCost() {
	out, err := os.Create(pathToElrondApiPackage + "../../config/gasCostWASM.go")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteWASMOpcodeCost(out)
}

func writeRustWasmerOpcodeCost() {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/output/opcode_cost.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustWasmerOpcodeCost(out)
}

func writeRustWasmerMeteringHelpers() {
	out, err := os.Create(pathToElrondApiPackage + "generate/cmd/output/wasmer_metering_helpers.rs")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	eapigen.WriteRustWasmerMeteringHelpers(out)
}

func tryCopyFilesToRustExecutorRepo() {
	fullPathToRustRepoConfigFile := filepath.Join(pathToElrondApiPackage, "generate/cmd/", pathToRustRepoConfigFile)
	contentBytes, err := ioutil.ReadFile(fullPathToRustRepoConfigFile)
	if err != nil {
		// this feature is optional
		fmt.Println("Rust files not copied to wasm-vm-executor-rs. Add a wasm-vm-executor-rs-path.txt with the path to enable feature.")
		return
	}
	rustExecutorPath := strings.Trim(string(contentBytes), " \n\t")

	fmt.Printf("Copying generated Rust files to '%s':\n", rustExecutorPath)
	copyFile(
		filepath.Join(pathToElrondApiPackage, "generate/cmd/output/vm_hooks.rs"),
		filepath.Join(rustExecutorPath, "exec-service/src/vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToElrondApiPackage, "generate/cmd/output/opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "exec-service/src/opcode_cost.rs"),
	)
	copyFile(
		filepath.Join(pathToElrondApiPackage, "generate/cmd/output/capi_vm_hook.rs"),
		filepath.Join(rustExecutorPath, "exec-c-api/src/capi_vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToElrondApiPackage, "generate/cmd/output/capi_vm_hook_pointers.rs"),
		filepath.Join(rustExecutorPath, "exec-c-api/src/capi_vm_hook_pointers.rs"),
	)
	copyFile(
		filepath.Join(pathToElrondApiPackage, "generate/cmd/output/wasmer_imports.rs"),
		filepath.Join(rustExecutorPath, "exec-service-wasmer/src/wasmer_imports.rs"),
	)
	copyFile(
		filepath.Join(pathToElrondApiPackage, "generate/cmd/output/wasmer_metering_helpers.rs"),
		filepath.Join(rustExecutorPath, "exec-service-wasmer/src/wasmer_metering_helpers.rs"),
	)
}

func copyFile(from, to string) {
	fmt.Printf("    %s -> %s\n", from, to)

	// Open original file
	original, err := os.Open(from)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	// Create new file
	new, err := os.Create(to)
	if err != nil {
		log.Fatal(err)
	}
	defer new.Close()

	//This will copy
	_, err = io.Copy(new, original)
	if err != nil {
		log.Fatal(err)
	}
}
