#include "libvmexeccapi_common.h"

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * Opaque handle to an executor.
 */
typedef struct vm_exec_w6_executor_t {

} vm_exec_w6_executor_t;

/**
 * Opaque handle to an instance.
 */
typedef struct vm_exec_w6_instance_t {

} vm_exec_w6_instance_t;

/**
 * Creates an executor from one configuration struct and one hook table.
 *
 * `hook_table` is copied, so the caller is free to release it as soon as this returns. Every
 * entry must be populated; see `CapiVMHookTable::copy_from`.
 *
 * On success writes the executor to `executor_out` and returns `VM_EXEC_OK`. The caller
 * owns the result and must release it with `vm_exec_w6_executor_destroy`.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_new_executor(struct vm_exec_w6_executor_t **executor_out, const vm_exec_executor_config_t *config, const vm_exec_vm_hook_c_func_pointers *hook_table, char *err_buf, int err_buf_len);

/**
 * Destroys an executor. Does nothing if `executor` is null.
 *
 * Instances created from an executor do not borrow it, so the order of destruction between
 * the two does not matter.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
void vm_exec_w6_executor_destroy(struct vm_exec_w6_executor_t *executor);

/**
 * Creates an instance.
 *
 * `go_vm_hooks` is the caller's VM hooks handle for *this* instance; it is handed back to
 * every VM hook as field 0 of the hook context, and is otherwise opaque here. Whatever it
 * refers to must stay alive for the instance's lifetime. From Go it is a `cgo.Handle`,
 * which is why it is an integer and not a pointer.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_new_instance(struct vm_exec_w6_executor_t *executor, struct vm_exec_w6_instance_t **instance_out, const uint8_t *wasm_bytes, uint32_t wasm_bytes_len, const vm_exec_compilation_options_t *options, uintptr_t go_vm_hooks, char *err_buf, int err_buf_len);

/**
 * Destroys an instance. Does nothing if `instance` is null.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
void vm_exec_w6_instance_destroy(struct vm_exec_w6_instance_t *instance);

/**
 * Calls an exported function, with a gas limit for this call.
 *
 * Always writes `result_out` when it returns `VM_EXEC_OK`. A non-`OK` status means the
 * call could not be attempted at all, which is distinct from a call that was attempted and
 * failed — the latter is reported through `result_out.kind`.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_instance_call(struct vm_exec_w6_instance_t *instance, const char *func_name, uint64_t points_limit, vm_exec_call_result_t *result_out, char *err_buf, int err_buf_len);

/**
 * Returns the number of points (gas) used so far.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_instance_get_points_used(struct vm_exec_w6_instance_t *instance, uint64_t *points_out, char *err_buf, int err_buf_len);

/**
 * Checks that every exported function takes no arguments and returns nothing.
 *
 * Returns 1 if so, 0 if not, or a negative status.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int vm_exec_w6_instance_check_signatures(const struct vm_exec_w6_instance_t *instance);

/**
 * Checks whether the contract exports a function with the given name.
 *
 * Returns 1 if so, 0 if not, or a negative status.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int vm_exec_w6_instance_has_function(const struct vm_exec_w6_instance_t *instance, const char *func_name);

/**
 * Buffer length needed by `vm_exec_w6_instance_exported_function_names`.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int vm_exec_w6_instance_exported_function_names_length(const struct vm_exec_w6_instance_t *instance);

/**
 * Writes the exported function names, pipe-separated, e.g. `"init|endpoint1|endpoint2"`.
 *
 * No ordering is guaranteed. Returns the number of bytes written including the terminator,
 * or `-1` if the buffer is absent or too small.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int vm_exec_w6_instance_exported_function_names(const struct vm_exec_w6_instance_t *instance, char *dest, int dest_len);

/**
 * Buffer length needed by `vm_exec_w6_instance_imported_function_names`.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int vm_exec_w6_instance_imported_function_names_length(const struct vm_exec_w6_instance_t *instance);

/**
 * Writes the names of the VM hooks the contract imports, pipe-separated.
 *
 * All of them at once, rather than a per-name query, because the Go contract validator
 * checks the module against whole sets of allowed names.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int vm_exec_w6_instance_imported_function_names(const struct vm_exec_w6_instance_t *instance, char *dest, int dest_len);

/**
 * Loads from instance memory into a caller-provided buffer, **truncating** at the end of
 * the memory.
 *
 * Returns the number of bytes actually written, which may be fewer than `dest_len`, or a
 * negative status. The truncation is where the Go VM's historical load semantics live: a
 * load whose offset is valid but whose end lies past the end of the memory yields the bytes
 * that are available and no error. Contracts have observed that, so the executor keeps
 * offering it. An offset outside the memory, or a negative length, is still an error.
 *
 * `dest_len` doubles as the requested length, which is what makes the clamping visible to
 * the caller: Go allocates `length` bytes and gets back `min(length, available)`.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
int64_t vm_exec_w6_state_mem_load(vm_exec_hook_context_t *ctx, int32_t mem_ptr, uint8_t *dest, int32_t dest_len, char *err_buf, int err_buf_len);

/**
 * Writes to instance memory. The range must fit entirely inside the memory.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_state_mem_store(vm_exec_hook_context_t *ctx, int32_t mem_ptr, const uint8_t *src, int32_t src_len, char *err_buf, int err_buf_len);

/**
 * Returns the size in bytes of the instance memory.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_state_mem_length(vm_exec_hook_context_t *ctx, uint64_t *length_out, char *err_buf, int err_buf_len);

/**
 * Grows the instance memory by `by_num_pages` pages of 64Kb each.
 *
 * Writes the number of pages the memory had *before* growing to `prev_pages_out`.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_state_mem_grow(vm_exec_hook_context_t *ctx, uint32_t by_num_pages, uint32_t *prev_pages_out, char *err_buf, int err_buf_len);

/**
 * Returns the number of points (gas) used by the instance so far.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_state_get_points_used(vm_exec_hook_context_t *ctx, uint64_t *points_out, char *err_buf, int err_buf_len);

/**
 * Sets the number of points (gas) used by the instance.
 *
 * # Safety
 *
 * C API function, works with raw object pointers.
 */
vm_exec_status_t vm_exec_w6_state_set_points_used(vm_exec_hook_context_t *ctx, uint64_t points, char *err_buf, int err_buf_len);
