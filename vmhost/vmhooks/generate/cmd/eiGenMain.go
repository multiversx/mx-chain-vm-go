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
const pathVmExecutorRsFile = "path-to-vm-executor-rs.txt"
const pathToMxSdkRsFile = "path-to-mx-sdk-rs.txt"

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
	writeWasmer2LegacyImportsCgo(eiMetadata)
	writeWasmer2ImportsCgo(eiMetadata)
	writeWasmer7ImportsCgo(eiMetadata)

	writeExecutorVMHookNames(eiMetadata)

	tryCreateRustOutputDirectory()

	writeRustVMHooksNames(eiMetadata)
	writeRustVMHooksTrait(eiMetadata)
	writeRustVMHooksLegacyTrait(eiMetadata)
	writeRustVMHooksLegacyAdapter(eiMetadata)
	writeRustVMHooksSignatures(eiMetadata)
	writeRustCapiVMHooks(eiMetadata)
	writeRustCapiVMHooksPointers(eiMetadata)
	writeRustCapiCommonVMHooks(eiMetadata)
	writeRustCapiCommonVMHooksPointers(eiMetadata)
	writeRustWasmerProdImports(eiMetadata)
	writeRustWasmerExperimentalImports(eiMetadata)
	writeRustVHDispatcherLegacy(eiMetadata)

	fmt.Printf("Generated code for %d executor callback methods.\n", len(eiMetadata.AllFunctions))

	writeExecutorWASMOpcodeCosts()
	writeExecutorOpcodeCost()
	writeWASMOpcodeCostFuncHelpers()
	writeWASMOpcodeCostConfigHelpers()
	writeRustOpcodeCost()
	writeRustCapiLegacyOpcodeCost()
	writeRustCapiOpcodeCost()
	writeRustWasmerOpcodeCost()
	writeOpcodeWhitelisted()
	writeRustWasmerExperimentalOpcodeCost()

	fmt.Println("Generated code for opcodes and metering helpers.")

	tryCopyFilesToMxExecutorRsRepo()
	tryCopyFilesToMxSdkRsRepo()
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

func writeWasmer2LegacyImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer2/wasmer2ImportsCgo.go")
	defer out.Close()
	eapigen.WriteWasmer2Cgo(out, eiMetadata)
}

func writeWasmer2ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer2new/wasmer2ImportsCgo.go")
	defer out.Close()
	eapigen.WriteWasmer2NewCgo(out, eiMetadata)
}

func writeWasmer7ImportsCgo(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../wasmer7/wasmer7ImportsCgo.go")
	defer out.Close()
	eapigen.WriteWasmer7Cgo(out, eiMetadata)
}

// The hook names, shared by every executor. Previously generated once per executor, into
// wasmer2new/wasmer2Names.go and mock/context/executorMockFunc.go.
func writeExecutorVMHookNames(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../executor/vmHookNames.go")
	defer out.Close()
	eapigen.WriteExecutorVMHookNames(out, eiMetadata)
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

func writeRustCapiCommonVMHooks(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/capiCommon_vm_hook.rs")
	defer out.Close()
	eapigen.WriteRustCapiCommonVMHooks(out, eiMetadata)
}

func writeRustCapiCommonVMHooksPointers(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/capiCommon_vm_hook_pointers.rs")
	defer out.Close()
	eapigen.WriteRustCapiCommonVMHooksPointers(out, eiMetadata)
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

func writeRustVMHooksSignatures(eiMetadata *eapigen.EIMetadata) {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/vm_hook_signature_list.rs")
	defer out.Close()
	eapigen.WriteRustVMHooksSignatures(out, eiMetadata)
}

func writeExecutorWASMOpcodeCosts() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../executor/gasCostWASM.go")
	defer out.Close()
	eapigen.WriteExecutorWASMOpcodeCost(out)
}

// The metered subset, shared by every backend. Previously generated once per backend, into
// wasmer2new/opcodeCost.go and wasmer7/opcodeCost.go.
func writeExecutorOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "../../executor/opcodeCost.go")
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

func writeRustOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/opcode_cost.rs")
	defer out.Close()
	eapigen.WriteRustOpcodeCost(out)
}

// Two outputs from one opcode list: the legacy C API's own type, and the shared one the new
// crates publish. Same layout, different names, so both can appear in one binary.
func writeRustCapiLegacyOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/capi_legacy_opcode_cost.rs")
	defer out.Close()
	eapigen.WriteRustCapiOpcodeCost(out, "vm_legacy_opcode_cost_t")
}

func writeRustCapiOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/capi_opcode_cost.rs")
	defer out.Close()
	eapigen.WriteRustCapiOpcodeCost(out, "vm_exec_opcode_cost_t")
}

func writeRustWasmerOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/wasmer_opcode_cost.rs")
	defer out.Close()
	eapigen.WriteRustWasmerOpcodeCost(out)
}

func writeOpcodeWhitelisted() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/opcode_whitelist.rs")
	defer out.Close()
	eapigen.WriteOpcodeWhitelisted(out)
}

func writeRustWasmerExperimentalOpcodeCost() {
	out := eapigen.NewEIGenWriter(pathToApiPackage, "generate/cmd/output/we_opcode_cost.rs")
	defer out.Close()
	eapigen.WriteRustWasmerExperimentalOpcodeCost(out)
}

func tryCopyFilesToMxExecutorRsRepo() {
	fullPathToConfigFile := filepath.Join(pathToApiPackage, "generate/cmd/", pathVmExecutorRsFile)
	contentBytes, err := os.ReadFile(fullPathToConfigFile)
	if err != nil {
		// this feature is optional
		fmt.Printf("Rust files not copied to wasm-vm-executor-rs. Add a %s with the path to enable feature.\n", pathVmExecutorRsFile)
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
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_legacy_opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "c-api-legacy/src/capi_opcode_cost.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_vm_hook.rs"),
		filepath.Join(rustExecutorPath, "c-api-legacy/src/capi_vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_vm_hook_pointers.rs"),
		filepath.Join(rustExecutorPath, "c-api-legacy/src/capi_vm_hook_pointers.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capi_opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "c-api-common/src/capi_opcode_cost.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capiCommon_vm_hook.rs"),
		filepath.Join(rustExecutorPath, "c-api-common/src/capi_vm_hooks.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/capiCommon_vm_hook_pointers.rs"),
		filepath.Join(rustExecutorPath, "c-api-common/src/capi_vm_hook_pointers.rs"),
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
		filepath.Join(pathToApiPackage, "generate/cmd/output/wasmer_opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "vm-executor-wasmer/src/wasmer_opcode_cost.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/we_opcode_cost.rs"),
		filepath.Join(rustExecutorPath, "vm-executor-experimental/src/middlewares/we_opcode_cost.rs"),
	)
}

func tryCopyFilesToMxSdkRsRepo() {
	fullPathToConfigFile := filepath.Join(pathToApiPackage, "generate/cmd/", pathToMxSdkRsFile)
	contentBytes, err := os.ReadFile(fullPathToConfigFile)
	if err != nil {
		// this feature is optional
		fmt.Printf("Rust files not copied to mx-sdk-rs. Add a %s with the path to enable feature.\n", pathToMxSdkRsFile)
		return
	}
	rustExecutorPath := strings.Trim(string(contentBytes), " \n\t")

	fmt.Printf("Copying generated Rust files to '%s':\n", rustExecutorPath)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/opcode_whitelist.rs"),
		filepath.Join(rustExecutorPath, "framework/meta-lib/src/tools/wasm_extractor/opcode_whitelist.rs"),
	)
	copyFile(
		filepath.Join(pathToApiPackage, "generate/cmd/output/vm_hook_signature_list.rs"),
		filepath.Join(rustExecutorPath, "framework/meta-lib/src/ei/vm_hook_signature_list.rs"),
	)
}

func copyFile(from, to string) {
	fmt.Printf("    %s -> %s\n", from, to)

	// Open original file
	original, err := os.Open(from)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = original.Close() }()

	// Create new file
	new, err := os.Create(to)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = new.Close() }()

	//This will copy
	_, err = io.Copy(new, original)
	if err != nil {
		log.Fatal(err)
	}
}
