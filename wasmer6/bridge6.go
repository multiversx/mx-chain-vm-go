package wasmer6

// #cgo LDFLAGS: -Wl,-rpath,${SRCDIR} -L${SRCDIR}
// #cgo linux,amd64 LDFLAGS:-lvmexeccapi_w6
// #cgo linux,arm64 LDFLAGS:-lvmexeccapi_w6_arm
// #cgo darwin,amd64 LDFLAGS:-lvmexeccapi_w6
// #cgo darwin,arm64 LDFLAGS:-lvmexeccapi_w6_arm
// #include <stdlib.h>
// #include "./libvmexeccapi_w6.h"
//
import "C"

import (
	"unsafe"
)

type cExecutorT = C.struct_vm_exec_w6_executor_t
type cInstanceT = C.struct_vm_exec_w6_instance_t
type cStatusT = C.enum_vm_exec_w6_status_t
type cHookContextT = C.vm_exec_w6_hook_context_t
type cOpcodeCostT = C.struct_vm_exec_opcode_cost_t
type cVMHookPointers = C.vm_exec_w6_vm_hook_c_func_pointers

// --- Executor lifecycle ------------------------------------------------------

func cNewExecutor(
	executorOut **cExecutorT,
	config *C.struct_vm_exec_w6_executor_config_t,
	hookTable *cVMHookPointers,
	buf errBuf,
) cStatusT {
	return C.vm_exec_w6_new_executor(executorOut, config, hookTable, buf.ptr(), errBufLen)
}

func cExecutorDestroy(executor *cExecutorT) {
	C.vm_exec_w6_executor_destroy(executor)
}

// --- Instance lifecycle ------------------------------------------------------

func cNewInstance(
	executor *cExecutorT,
	instanceOut **cInstanceT,
	wasmBytes []byte,
	options *C.struct_vm_exec_w6_compilation_options_t,
	goVMHooks uintptr,
	buf errBuf,
) cStatusT {
	return C.vm_exec_w6_new_instance(
		executor,
		instanceOut,
		(*C.uint8_t)(unsafe.Pointer(&wasmBytes[0])),
		C.uint32_t(len(wasmBytes)),
		options,
		C.uintptr_t(goVMHooks),
		buf.ptr(),
		errBufLen,
	)
}

func cInstanceDestroy(instance *cInstanceT) {
	C.vm_exec_w6_instance_destroy(instance)
}

// --- Outer instance operations -----------------------------------------------

func cInstanceCall(
	instance *cInstanceT,
	funcName string,
	pointsLimit uint64,
	resultOut *C.struct_vm_exec_w6_call_result_t,
	buf errBuf,
) cStatusT {
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	return C.vm_exec_w6_instance_call(
		instance,
		cFuncName,
		C.uint64_t(pointsLimit),
		resultOut,
		buf.ptr(),
		errBufLen,
	)
}

func cInstanceGetPointsUsed(instance *cInstanceT, pointsOut *C.uint64_t, buf errBuf) cStatusT {
	return C.vm_exec_w6_instance_get_points_used(instance, pointsOut, buf.ptr(), errBufLen)
}

func cInstanceCheckSignatures(instance *cInstanceT) C.int {
	return C.vm_exec_w6_instance_check_signatures(instance)
}

func cInstanceHasFunction(instance *cInstanceT, funcName string) C.int {
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))

	return C.vm_exec_w6_instance_has_function(instance, cFuncName)
}

func cInstanceExportedFunctionNamesLength(instance *cInstanceT) C.int {
	return C.vm_exec_w6_instance_exported_function_names_length(instance)
}

func cInstanceExportedFunctionNames(instance *cInstanceT, dest []byte) C.int {
	return C.vm_exec_w6_instance_exported_function_names(
		instance,
		(*C.char)(unsafe.Pointer(&dest[0])),
		C.int(len(dest)),
	)
}

func cInstanceImportedFunctionNamesLength(instance *cInstanceT) C.int {
	return C.vm_exec_w6_instance_imported_function_names_length(instance)
}

func cInstanceImportedFunctionNames(instance *cInstanceT, dest []byte) C.int {
	return C.vm_exec_w6_instance_imported_function_names(
		instance,
		(*C.char)(unsafe.Pointer(&dest[0])),
		C.int(len(dest)),
	)
}

// --- Instance state operations, reachable only from inside a hook ------------

func cStateMemLoad(ctx *cHookContextT, memPtr int32, dest []byte, buf errBuf) C.int64_t {
	return C.vm_exec_w6_state_mem_load(
		ctx,
		C.int32_t(memPtr),
		(*C.uint8_t)(unsafe.Pointer(&dest[0])),
		C.int32_t(len(dest)),
		buf.ptr(),
		errBufLen,
	)
}

func cStateMemStore(ctx *cHookContextT, memPtr int32, src []byte, buf errBuf) cStatusT {
	return C.vm_exec_w6_state_mem_store(
		ctx,
		C.int32_t(memPtr),
		(*C.uint8_t)(unsafe.Pointer(&src[0])),
		C.int32_t(len(src)),
		buf.ptr(),
		errBufLen,
	)
}

func cStateMemLength(ctx *cHookContextT, lengthOut *C.uint64_t, buf errBuf) cStatusT {
	return C.vm_exec_w6_state_mem_length(ctx, lengthOut, buf.ptr(), errBufLen)
}

func cStateMemGrow(
	ctx *cHookContextT,
	byNumPages uint32,
	prevPagesOut *C.uint32_t,
	buf errBuf,
) cStatusT {
	return C.vm_exec_w6_state_mem_grow(ctx, C.uint32_t(byNumPages), prevPagesOut, buf.ptr(), errBufLen)
}

func cStateGetPointsUsed(ctx *cHookContextT, pointsOut *C.uint64_t, buf errBuf) cStatusT {
	return C.vm_exec_w6_state_get_points_used(ctx, pointsOut, buf.ptr(), errBufLen)
}

func cStateSetPointsUsed(ctx *cHookContextT, points uint64, buf errBuf) cStatusT {
	return C.vm_exec_w6_state_set_points_used(ctx, C.uint64_t(points), buf.ptr(), errBufLen)
}
