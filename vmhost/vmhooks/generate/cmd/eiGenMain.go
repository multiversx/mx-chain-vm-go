package main

import (
	"errors"
	"fmt"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	eapigen "github.com/multiversx/mx-chain-vm-go/vmhost/vmhooks/generate"
)

const pathToApiPackage = "./"
const pathToRustRepoConfigFile = "wasm-vm-executor-rs-path.txt"

// Until we merge the `feat/wasmer2`, there are some files that are not supposed to be generated.
const wasmer2Branch = true

func initEIMetadata() *eapigen.EIMetadata {
	return &eapigen.EIMetadata{
		Groups: []*eapigen.EIGroup{
			{SourcePath: "baseOps.go", Name: "Main"},
			{SourcePath: "managedei.go", Name: "Managed"},
			{SourcePath: "bigFloatOps.go", Name: "BigFloat"},
			{SourcePath: "bigIntOps.go", Name: "BigInt"},
			{SourcePath: "manBufOps.go", Name: "ManagedBuffer"},
			{SourcePath: "manMapOps.go", Name: "ManagedMap"},
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
	err := eapigen.ReadAndParseEIMetadata(fset, pathToApiPackage, eiMetadata)
	if err != nil {
		panic(err)
	}

	writeVMHooks(eiMetadata)
	writeVMHooksWrapper(eiMetadata)
	writeWasmer1ImportsCgo(eiMetadata)
	if wasmer2Branch {
		writeWasmer2ImportsCgo(eiMetadata)
		writeWasmer2Names(eiMetadata)
	}

	writeNamesForMockExecutor(eiMetadata)

	tryCreateRustOutputDirectory()

	writeRustVMHooksTrait(eiMetadata)
	writeRustCapiVMHooks(eiMetadata)
	writeRustCapiVMHooksPointers(eiMetadata)
	writeRustWasmerImports(eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))

	if wasmer2Branch {
		writeExecutorOpcodeCosts()
		writeWasmer2OpcodeCost()
	}
	writeWASMOpcodeCostFuncHelpers()
	writeWASMOpcodeCostConfigHelpers()
	writeOpcodeCostFuncHelpers()
	writeRustOpcodeCost()
	writeRustWasmerMeteringHelpers()

	fmt.Println("Generated code for opcodes and metering helpers.")

	tryCopyFilesToRustExecutorRepo()
}

func writeVMHooks(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../executor/vmHooks.go")
	defer out.Close()
	eapigen.WriteEIInterface(out, eiMetadata)
}

func writeVMHooksWrapper(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../executor/wrapper/wrapperVMHooks.go")
	defer out.Close()
	eapigen.WriteVMHooksWrapper(out, eiMetadata)
}

func writeWasmer1ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer/wasmerImportsCgo.go")
	defer out.Close()
	eapigen.WriteWasmer1Cgo(out, eiMetadata)
}

func writeWasmer2ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer2/wasmer2ImportsCgo.go")
	defer out.Close()
	eapigen.WriteWasmer2Cgo(out, eiMetadata)
}

func writeWasmer2Names(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer2/wasmer2Names.go")
	defer out.Close()
	eapigen.WriteNames(out, "wasmer2", eiMetadata)
}

func writeNamesForMockExecutor(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../mock/context/executorMockFunc.go")
	defer out.Close()
	eapigen.WriteNames(out, "mock", eiMetadata)
}

func tryCreateRustOutputDirectory() {
	outputDirPath := filepath.Join(pathToApiPackage, "generate/cmd/output")
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
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/vm_hooks.rs")
	defer out.Close()
	eapigen.WriteRustVMHooksTrait(out, eiMetadata)
}

func writeRustCapiVMHooks(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/capi_vm_hook.rs")
	defer out.Close()
	eapigen.WriteRustCapiVMHooks(out, eiMetadata)
}

func writeRustCapiVMHooksPointers(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/capi_vm_hook_pointers.rs")
	defer out.Close()
	eapigen.WriteRustCapiVMHooksPointers(out, eiMetadata)
}

func writeRustWasmerImports(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/wasmer_imports.rs")
	defer out.Close()
	eapigen.WriteRustWasmerImports(out, eiMetadata)
}

func writeExecutorOpcodeCosts() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../executor/gasCostWASM.go")
	defer out.Close()
	eapigen.WriteExecutorOpcodeCost(out)
}

func writeWASMOpcodeCostFuncHelpers() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/FillGasMap_WASMOpcodeCosts.txt")
	defer out.Close()
	eapigen.WriteWASMOpcodeCostFuncHelpers(out)
}

func writeWASMOpcodeCostConfigHelpers() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/config.txt")
	defer out.Close()
	eapigen.WriteWASMOpcodeCostConfigHelpers(out)
}

func writeWasmer2OpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer2/opcodeCost.go")
	defer out.Close()
	eapigen.WriteWasmer2OpcodeCost(out)
}

func writeOpcodeCostFuncHelpers() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/extractOpcodeCost.txt")
	defer out.Close()
	eapigen.WriteOpcodeCostFuncHelpers(out)
}

func writeRustOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/opcode_cost.rs")
	defer out.Close()
	eapigen.WriteRustOpcodeCost(out)
}

func writeRustWasmerMeteringHelpers() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/wasmer_metering_helpers.rs")
	defer out.Close()
	eapigen.WriteRustWasmerMeteringHelpers(out)
}

func tryCopyFilesToRustExecutorRepo() {
	fullPathToRustRepoConfigFile := filepath.Join(pathToApiPackage, "generate/cmd/", pathToRustRepoConfigFile)
	contentBytes, err := os.ReadFile(fullPathToRustRepoConfigFile)
	if err != nil {
		// this feature is optional
		fmt.Println("Rust files not copied to wasm-vm-executor-rs. Add a wasm-vm-executor-rs-path.txt with the path to enable feature.")
		return
	}
	rustExecutorPath := strings.Trim(string(contentBytes), " \n\t")

	fmt.Printf("Copying generated Rust files to '%s':\n", rustExecutorPath)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/vm_hooks.rs"),
		filepath.Join(rustExecutorPath, "exec-service/src/vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "exec-service/src/opcode_cost.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_vm_hook.rs"),
		filepath.Join(rustExecutorPath, "exec-c-api/src/capi_vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_vm_hook_pointers.rs"),
		filepath.Join(rustExecutorPath, "exec-c-api/src/capi_vm_hook_pointers.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/wasmer_imports.rs"),
		filepath.Join(rustExecutorPath, "exec-service-wasmer/src/wasmer_imports.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/wasmer_metering_helpers.rs"),
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
