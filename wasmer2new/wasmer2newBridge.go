package wasmer2new

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR}
// #cgo linux,amd64 LDFLAGS:-lvmexeccapi_w2
// #cgo linux,arm64 LDFLAGS:-lvmexeccapi_w2_arm
// #cgo darwin,amd64 LDFLAGS:-lvmexeccapi_w2
// #cgo darwin,arm64 LDFLAGS:-lvmexeccapi_w2_arm
// #include <stdlib.h>
// #include "./libvmexeccapi_w2.h"
//
import "C"

import (
	"unsafe"

	"github.com/multiversx/mx-chain-vm-go/wasmercommon"
)

// This file is the only one in the package that calls the library, and therefore the only one
// that deals in its C status type: every function here converts to wasmercommon.Status on the
// way out, so everything above the bridge is written against the shared vocabulary and could
// not accidentally be written against one backend's.

type cExecutorT = C.struct_vm_exec_w2_executor_t
type cInstanceT = C.struct_vm_exec_w2_instance_t
type cStopHandleT = C.struct_vm_exec_w2_stop_handle_t
type cHookContextT = C.vm_exec_hook_context_t
type cOpcodeCostT = C.struct_vm_exec_opcode_cost_t
type cVMHookPointers = C.vm_exec_vm_hook_c_func_pointers

// --- Executor lifecycle ------------------------------------------------------

func cNewExecutor(
	executorOut **cExecutorT,
	config *C.struct_vm_exec_executor_config_t,
	hookTable *cVMHookPointers,
	buf errBuf,
) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_new_executor(executorOut, config, hookTable, bufPtr(buf), errBufLen))
}

func cExecutorDestroy(executor *cExecutorT) {
	C.vm_exec_w2_executor_destroy(executor)
}

// --- Instance lifecycle ------------------------------------------------------

func cNewInstance(
	executor *cExecutorT,
	instanceOut **cInstanceT,
	wasmBytes []byte,
	options *C.struct_vm_exec_compilation_options_t,
	goVMHooks uintptr,
	buf errBuf,
) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_new_instance(
		executor,
		instanceOut,
		(*C.uint8_t)(unsafe.Pointer(&wasmBytes[0])),
		C.uint32_t(len(wasmBytes)),
		options,
		C.uintptr_t(goVMHooks),
		bufPtr(buf),
		errBufLen,
	))
}

func cNewInstanceFromCache(
	executor *cExecutorT,
	instanceOut **cInstanceT,
	cacheBytes []byte,
	options *C.struct_vm_exec_compilation_options_t,
	goVMHooks uintptr,
	buf errBuf,
) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_new_instance_from_cache(
		executor,
		instanceOut,
		(*C.uint8_t)(unsafe.Pointer(&cacheBytes[0])),
		C.uint32_t(len(cacheBytes)),
		options,
		C.uintptr_t(goVMHooks),
		bufPtr(buf),
		errBufLen,
	))
}

func cInstanceDestroy(instance *cInstanceT) {
	C.vm_exec_w2_instance_destroy(instance)
}

func cInstanceStopHandle(
	instance *cInstanceT,
	handleOut **cStopHandleT,
	buf errBuf,
) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_instance_stop_handle(instance, handleOut, bufPtr(buf), errBufLen))
}

func cStopHandleRequestStop(handle *cStopHandleT) {
	C.vm_exec_w2_stop_handle_request_stop(handle)
}

func cStopHandleDestroy(handle *cStopHandleT) {
	C.vm_exec_w2_stop_handle_destroy(handle)
}

func cInstanceReset(instance *cInstanceT, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_instance_reset(instance, bufPtr(buf), errBufLen))
}

// --- Outer instance operations -----------------------------------------------

func cInstanceCacheLength(instance *cInstanceT, lengthOut *C.uint32_t, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_instance_cache_length(instance, lengthOut, bufPtr(buf), errBufLen))
}

func cInstanceCache(instance *cInstanceT, dest []byte, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_instance_cache(
		instance,
		(*C.uint8_t)(unsafe.Pointer(&dest[0])),
		C.uint32_t(len(dest)),
		bufPtr(buf),
		errBufLen,
	))
}

func cInstanceCall(
	instance *cInstanceT,
	funcName string,
	pointsLimit uint64,
	resultOut *C.struct_vm_exec_call_result_t,
	buf errBuf,
) wasmercommon.Status {
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	return toStatus(C.vm_exec_w2_instance_call(
		instance,
		cFuncName,
		C.uint64_t(pointsLimit),
		resultOut,
		bufPtr(buf),
		errBufLen,
	))
}

func cInstanceGetPointsUsed(
	instance *cInstanceT,
	pointsOut *C.uint64_t,
	buf errBuf,
) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_instance_get_points_used(instance, pointsOut, bufPtr(buf), errBufLen))
}

func cInstanceCheckSignatures(instance *cInstanceT) int32 {
	return int32(C.vm_exec_w2_instance_check_signatures(instance))
}

func cInstanceHasFunction(instance *cInstanceT, funcName string) int32 {
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	return int32(C.vm_exec_w2_instance_has_function(instance, cFuncName))
}

func cInstanceExportedFunctionNamesLength(instance *cInstanceT) int32 {
	return int32(C.vm_exec_w2_instance_exported_function_names_length(instance))
}

func cInstanceExportedFunctionNames(instance *cInstanceT, dest []byte) int32 {
	return int32(C.vm_exec_w2_instance_exported_function_names(
		instance,
		(*C.char)(unsafe.Pointer(&dest[0])),
		C.int(len(dest)),
	))
}

func cInstanceImportedFunctionNamesLength(instance *cInstanceT) int32 {
	return int32(C.vm_exec_w2_instance_imported_function_names_length(instance))
}

func cInstanceImportedFunctionNames(instance *cInstanceT, dest []byte) int32 {
	return int32(C.vm_exec_w2_instance_imported_function_names(
		instance,
		(*C.char)(unsafe.Pointer(&dest[0])),
		C.int(len(dest)),
	))
}

// --- Instance state operations, reachable only from inside a hook ------------

func cStateMemLoad(ctx *cHookContextT, memPtr int32, dest []byte, buf errBuf) int64 {
	return int64(C.vm_exec_w2_state_mem_load(
		ctx,
		C.int32_t(memPtr),
		(*C.uint8_t)(unsafe.Pointer(&dest[0])),
		C.int32_t(len(dest)),
		bufPtr(buf),
		errBufLen,
	))
}

func cStateMemStore(ctx *cHookContextT, memPtr int32, src []byte, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_state_mem_store(
		ctx,
		C.int32_t(memPtr),
		(*C.uint8_t)(unsafe.Pointer(&src[0])),
		C.int32_t(len(src)),
		bufPtr(buf),
		errBufLen,
	))
}

func cStateMemLength(ctx *cHookContextT, lengthOut *C.uint64_t, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_state_mem_length(ctx, lengthOut, bufPtr(buf), errBufLen))
}

func cStateMemGrow(
	ctx *cHookContextT,
	byNumPages uint32,
	prevPagesOut *C.uint32_t,
	buf errBuf,
) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_state_mem_grow(
		ctx,
		C.uint32_t(byNumPages),
		prevPagesOut,
		bufPtr(buf),
		errBufLen,
	))
}

func cStateGetPointsUsed(ctx *cHookContextT, pointsOut *C.uint64_t, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_state_get_points_used(ctx, pointsOut, bufPtr(buf), errBufLen))
}

func cStateSetPointsUsed(ctx *cHookContextT, points uint64, buf errBuf) wasmercommon.Status {
	return toStatus(C.vm_exec_w2_state_set_points_used(ctx, C.uint64_t(points), bufPtr(buf), errBufLen))
}
