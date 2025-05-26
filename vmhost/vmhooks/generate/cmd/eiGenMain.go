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
	writeWasmer2ImportsCgo(eiMetadata)
	writeWasmer2Names(eiMetadata)

	writeNamesForMockExecutor(eiMetadata)

	tryCreateRustOutputDirectory()

	writeRustVMHooksNames(eiMetadata)
	writeRustVMHooksTrait(eiMetadata)
	writeRustVMHooksLegacyTrait(eiMetadata)
	writeRustVMHooksLegacyAdapter(eiMetadata)
	writeRustCapiVMHooks(eiMetadata)
	writeRustCapiVMHooksPointers(eiMetadata)
	writeRustWasmerProdImports(eiMetadata)
	writeRustWasmerExperimentalImports(eiMetadata)
	writeRustVHDispatcherLegacy(eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))

	writeExecutorOpcodeCosts()
	writeWasmer2OpcodeCost()
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

func writeRustVMHooksNames(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/ei_1_5.rs")
	defer out.Close()
	eapigen.WriteRustHookNames(out, eiMetadata)
}

func writeRustVMHooksLegacyTrait(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/vm_hooks.rs")
	defer out.Close()
	eapigen.WriteRustVMHooksLegacyTrait(out, eiMetadata)
}

func writeRustVMHooksLegacyAdapter(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/vm_hooks_legacy_adapter.rs")
	defer out.Close()
	eapigen.WriteRustVMHooksLegacyAdapter(out, eiMetadata)
}

func writeRustVMHooksTrait(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/vm_hooks_new.rs")
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

func writeRustWasmerProdImports(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/wasmer_imports.rs")
	defer out.Close()
	eapigen.WriteRustWasmerImports(out, eiMetadata)
}

func writeRustWasmerExperimentalImports(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/we_imports.rs")
	defer out.Close()
	eapigen.WriteRustWasmer5Imports(out, eiMetadata)
}

func writeRustVHDispatcherLegacy(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/vh_dispatcher_legacy.rs")
	defer out.Close()
	eapigen.WriteRustVHDispatcherLegacy(out, eiMetadata)
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
		filepath.Join(rustExecutorPath, "vm-executor/src/vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/vm_hooks_new.rs"),
		filepath.Join(rustExecutorPath, "vm-executor/src/new_traits/vm_hooks_new.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/vm_hooks_legacy_adapter.rs"),
		filepath.Join(rustExecutorPath, "vm-executor/src/new_traits/vm_hooks_legacy_adapter.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "vm-executor/src/opcode_cost.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_vm_hook.rs"),
		filepath.Join(rustExecutorPath, "c-api/src/capi_vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_vm_hook_pointers.rs"),
		filepath.Join(rustExecutorPath, "c-api/src/capi_vm_hook_pointers.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/wasmer_imports.rs"),
		filepath.Join(rustExecutorPath, "vm-executor-wasmer/src/wasmer_imports.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/we_imports.rs"),
		filepath.Join(rustExecutorPath, "vm-executor-experimental/src/we_imports.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/wasmer_metering_helpers.rs"),
		filepath.Join(rustExecutorPath, "vm-executor-wasmer/src/wasmer_metering_helpers.rs"),
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
