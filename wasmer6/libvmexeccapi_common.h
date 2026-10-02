#ifndef VM_EXEC_COMMON_H
#define VM_EXEC_COMMON_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * The value of `vm_exec_hook_result_t::early_exit` that means "the hook returned normally".
 *
 * A dedicated sentinel is needed because `0` is a real early-exit code: the Rust VM uses
 * `ASYNC_CALL_EARLY_EXIT_CODE = 0` and the Go VM has the same async-call case, so "code 0"
 * and "no code" have to stay distinguishable.
 */
#define VM_EXEC_NO_EARLY_EXIT UINT64_MAX

/**
 * Result of a C API call.
 *
 * `VM_EXEC_OK` is zero and every failure is negative, so the functions that return a
 * count (the memory loads) can share the same code space: a non-negative result is a
 * count, a negative one is one of these.
 */
typedef enum vm_exec_status_t {
  /**
   * The call succeeded.
   */
  VM_EXEC_OK = 0,
  /**
   * The call failed; the error message was written to the provided error buffer.
   */
  VM_EXEC_ERROR = -1,
  /**
   * A pointer argument that must not be null was null. No message is produced, because
   * a null pointer may well be the error buffer itself.
   */
  VM_EXEC_NULL_POINTER = -2,
  /**
   * A hook context was used outside the host call it belongs to. See
   * [`crate::capi_vm_hook_context`].
   */
  VM_EXEC_STALE_CONTEXT = -3,
  /**
   * A panic was caught at the boundary. Rust panics must never unwind into C, so every
   * entry point traps them and reports this instead.
   */
  VM_EXEC_PANIC = -4,
} vm_exec_status_t;

/**
 * How a call ended.
 *
 * This replaces the legacy arrangement, where a call returned a bare ok/error and the VM
 * then polled a wasm global to find out what had actually happened.
 */
typedef enum vm_exec_call_result_kind_t {
  /**
   * The function returned normally.
   */
  VM_EXEC_CALL_OK = 0,
  /**
   * No exported function by that name.
   */
  VM_EXEC_CALL_FUNCTION_NOT_FOUND = 1,
  /**
   * The instance trapped for a reason that is neither a VM hook's doing nor a
   * middleware breakpoint. The message is in the error buffer.
   */
  VM_EXEC_CALL_RUNTIME_ERROR = 2,
  /**
   * A VM hook asked for execution to stop. `early_exit_code` says which; the message, if
   * any, is in the error buffer, and is informational — the VM has already recorded its
   * own.
   */
  VM_EXEC_CALL_EARLY_EXIT = 3,
  /**
   * A middleware stopped execution. `breakpoint` says which.
   */
  VM_EXEC_CALL_BREAKPOINT = 4,
} vm_exec_call_result_kind_t;

/**
 * FFI-safe mirror of `OpcodeCost`.
 *
 * This is a distinct type rather than a cast of `OpcodeCost` itself, so that
 * the general-purpose struct is free to change shape without touching the C
 * ABI, and vice versa.
 */
typedef struct vm_exec_opcode_cost_t {
  uint32_t opcode_block;
  uint32_t opcode_br;
  uint32_t opcode_brif;
  uint32_t opcode_brtable;
  uint32_t opcode_call;
  uint32_t opcode_callindirect;
  uint32_t opcode_catch;
  uint32_t opcode_catchall;
  uint32_t opcode_delegate;
  uint32_t opcode_drop;
  uint32_t opcode_else;
  uint32_t opcode_end;
  uint32_t opcode_globalget;
  uint32_t opcode_globalset;
  uint32_t opcode_i32add;
  uint32_t opcode_i32and;
  uint32_t opcode_i32clz;
  uint32_t opcode_i32const;
  uint32_t opcode_i32ctz;
  uint32_t opcode_i32divs;
  uint32_t opcode_i32divu;
  uint32_t opcode_i32eq;
  uint32_t opcode_i32eqz;
  uint32_t opcode_i32extend16s;
  uint32_t opcode_i32extend8s;
  uint32_t opcode_i32ges;
  uint32_t opcode_i32geu;
  uint32_t opcode_i32gts;
  uint32_t opcode_i32gtu;
  uint32_t opcode_i32les;
  uint32_t opcode_i32leu;
  uint32_t opcode_i32load;
  uint32_t opcode_i32load16s;
  uint32_t opcode_i32load16u;
  uint32_t opcode_i32load8s;
  uint32_t opcode_i32load8u;
  uint32_t opcode_i32lts;
  uint32_t opcode_i32ltu;
  uint32_t opcode_i32mul;
  uint32_t opcode_i32ne;
  uint32_t opcode_i32or;
  uint32_t opcode_i32popcnt;
  uint32_t opcode_i32rems;
  uint32_t opcode_i32remu;
  uint32_t opcode_i32rotl;
  uint32_t opcode_i32rotr;
  uint32_t opcode_i32shl;
  uint32_t opcode_i32shrs;
  uint32_t opcode_i32shru;
  uint32_t opcode_i32store;
  uint32_t opcode_i32store16;
  uint32_t opcode_i32store8;
  uint32_t opcode_i32sub;
  uint32_t opcode_i32wrapi64;
  uint32_t opcode_i32xor;
  uint32_t opcode_i64add;
  uint32_t opcode_i64and;
  uint32_t opcode_i64clz;
  uint32_t opcode_i64const;
  uint32_t opcode_i64ctz;
  uint32_t opcode_i64divs;
  uint32_t opcode_i64divu;
  uint32_t opcode_i64eq;
  uint32_t opcode_i64eqz;
  uint32_t opcode_i64extend16s;
  uint32_t opcode_i64extend32s;
  uint32_t opcode_i64extend8s;
  uint32_t opcode_i64extendi32s;
  uint32_t opcode_i64extendi32u;
  uint32_t opcode_i64ges;
  uint32_t opcode_i64geu;
  uint32_t opcode_i64gts;
  uint32_t opcode_i64gtu;
  uint32_t opcode_i64les;
  uint32_t opcode_i64leu;
  uint32_t opcode_i64load;
  uint32_t opcode_i64load16s;
  uint32_t opcode_i64load16u;
  uint32_t opcode_i64load32s;
  uint32_t opcode_i64load32u;
  uint32_t opcode_i64load8s;
  uint32_t opcode_i64load8u;
  uint32_t opcode_i64lts;
  uint32_t opcode_i64ltu;
  uint32_t opcode_i64mul;
  uint32_t opcode_i64ne;
  uint32_t opcode_i64or;
  uint32_t opcode_i64popcnt;
  uint32_t opcode_i64rems;
  uint32_t opcode_i64remu;
  uint32_t opcode_i64rotl;
  uint32_t opcode_i64rotr;
  uint32_t opcode_i64shl;
  uint32_t opcode_i64shrs;
  uint32_t opcode_i64shru;
  uint32_t opcode_i64store;
  uint32_t opcode_i64store16;
  uint32_t opcode_i64store32;
  uint32_t opcode_i64store8;
  uint32_t opcode_i64sub;
  uint32_t opcode_i64xor;
  uint32_t opcode_if;
  uint32_t opcode_localget;
  uint32_t opcode_localset;
  uint32_t opcode_localtee;
  uint32_t opcode_loop;
  uint32_t opcode_memorygrow;
  uint32_t opcode_memorysize;
  uint32_t opcode_nop;
  uint32_t opcode_reffunc;
  uint32_t opcode_refisnull;
  uint32_t opcode_refnull;
  uint32_t opcode_rethrow;
  uint32_t opcode_return;
  uint32_t opcode_returncall;
  uint32_t opcode_returncallindirect;
  uint32_t opcode_select;
  uint32_t opcode_tableget;
  uint32_t opcode_tablegrow;
  uint32_t opcode_tableinit;
  uint32_t opcode_tableset;
  uint32_t opcode_tablesize;
  uint32_t opcode_throw;
  uint32_t opcode_try;
  uint32_t opcode_typedselect;
  uint32_t opcode_unreachable;
  uint32_t opcode_unwind;
  uint32_t opcode_memorycopy;
  uint32_t opcode_memoryfill;
  uint32_t opcode_memorycopyperbyte;
  uint32_t opcode_memoryfillperbyte;
  uint32_t opcode_localallocate;
} vm_exec_opcode_cost_t;

/**
 * FFI-safe mirror of the new `CompilationOptions`.
 *
 * Four fields, against the legacy seven: `gas_limit` moved to the call, and metering and
 * runtime breakpoints are no longer optional. Fields that are `usize` on the Rust side are
 * fixed to `u64` here, since `usize` has no guaranteed width across host platforms.
 */
typedef struct vm_exec_compilation_options_t {
  uint64_t unmetered_locals;
  uint64_t max_memory_grow;
  uint64_t max_memory_grow_delta;
  bool opcode_trace;
} vm_exec_compilation_options_t;

/**
 * Everything an executor is configured with.
 *
 * The opcode cost table is referenced rather than inlined, because it is large and the
 * caller already holds one. It is read during the call and not retained.
 */
typedef struct vm_exec_executor_config_t {
  /**
   * Which opcode set and whitelist to apply; see `OpcodeVersion` on the Rust side.
   */
  int32_t opcode_version;
  /**
   * The per-opcode gas costs. Must not be null.
   */
  const struct vm_exec_opcode_cost_t *opcode_cost;
} vm_exec_executor_config_t;

/**
 * The out-struct of `vm_exec_w6_instance_call`.
 */
typedef struct vm_exec_call_result_t {
  enum vm_exec_call_result_kind_t kind;
  /**
   * Meaningful only for `VM_EXEC_CALL_EARLY_EXIT`.
   */
  uint64_t early_exit_code;
  /**
   * Meaningful only for `VM_EXEC_CALL_BREAKPOINT`. Values match `BreakpointValue`:
   * 4 = out of gas, 5 = memory limit.
   */
  uint64_t breakpoint;
} vm_exec_call_result_t;

/**
 * What every VM hook returns.
 *
 * One 16-byte struct covers every hook shape — `void`, `i32`, `i64` — and is returned in
 * registers under both SysV AMD64 (RAX:RDX) and AAPCS64 (X0:X1), so the early-exit channel
 * costs no memory traffic. A cgo `//export` function may return a C struct by value, which
 * is what makes this shape available to the Go side.
 *
 * The early-exit *message* does not travel in this struct. A pointer cannot: cgo forbids
 * returning an unpinned Go pointer, and a fixed buffer cannot work because the messages are
 * unbounded user text. The Go VM does not need it to travel, because the hook has already
 * recorded the message in its own output state by the time it returns.
 */
typedef struct vm_exec_hook_result_t {
  /**
   * The hook's own return value, widened to 64 bits. Ignored for a `void` hook, and for
   * any hook that reports an early exit.
   */
  int64_t value;
  /**
   * The early-exit code, or `VM_EXEC_NO_EARLY_EXIT`.
   */
  uint64_t early_exit;
} vm_exec_hook_result_t;

/**
 * What a VM hook is handed as its first argument.
 *
 * Only `go_vm_hooks` is part of the layout the Go side mirrors; the rest is private to
 * Rust, and the whole pointer is what travels back into the state functions.
 */
typedef struct vm_exec_hook_context_t {
  /**
   * The caller's VM hooks handle for this instance. Read by the caller, never by Rust.
   *
   * An opaque word rather than a pointer, deliberately: Rust only stores it and hands it
   * back, and the Go side needs it to be an integer. cgo forbids storing a Go pointer in
   * C memory, so what Go puts here is a `runtime/cgo.Handle`, which is a `uintptr` by
   * construction. Typing the field as a pointer would force Go to convert it back with
   * `unsafe.Pointer(uintptr(...))`, exactly the pattern `go vet` rejects.
   */
  uintptr_t go_vm_hooks;
  /**
   * Erased `&mut S` for the backend's instance state; null while no call is active.
   */
  void *state;
  /**
   * Erased `&RefCell<CapiVMHookContextPool>`, for the token check.
   */
  const void *pool;
  /**
   * This activation's call token; `0` while no call is active.
   */
  uint64_t token;
} vm_exec_hook_context_t;

/**
 * The Go VM's hook implementations, as a table of C function pointers.
 *
 * Copied into Rust at executor creation; see CapiVMHookTable.
 */
typedef struct vm_exec_vm_hook_c_func_pointers {
  struct vm_exec_hook_result_t (*get_gas_left_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_sc_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_owner_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_shard_of_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset);
  struct vm_exec_hook_result_t (*is_smart_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset);
  struct vm_exec_hook_result_t (*signal_error_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t message_offset, int32_t message_length);
  struct vm_exec_hook_result_t (*get_external_balance_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_block_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t nonce, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_esdt_balance_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len, int64_t nonce, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_esdt_nft_name_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len, int64_t nonce);
  struct vm_exec_hook_result_t (*get_esdt_nft_attribute_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len, int64_t nonce);
  struct vm_exec_hook_result_t (*get_esdt_nft_uri_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len, int64_t nonce);
  struct vm_exec_hook_result_t (*get_esdt_token_data_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len, int64_t nonce, int32_t value_handle, int32_t properties_offset, int32_t hash_offset, int32_t name_offset, int32_t attributes_offset, int32_t creator_offset, int32_t royalties_handle, int32_t uris_offset);
  struct vm_exec_hook_result_t (*get_esdt_local_roles_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t token_id_handle);
  struct vm_exec_hook_result_t (*validate_token_identifier_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t token_id_handle);
  struct vm_exec_hook_result_t (*transfer_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t value_offset, int32_t data_offset, int32_t length);
  struct vm_exec_hook_result_t (*transfer_value_execute_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t value_offset, int64_t gas_limit, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*transfer_esdt_execute_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t token_id_offset, int32_t token_id_len, int32_t value_offset, int64_t gas_limit, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*transfer_esdt_nft_execute_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t token_id_offset, int32_t token_id_len, int32_t value_offset, int64_t nonce, int64_t gas_limit, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*multi_transfer_esdt_nft_execute_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t num_token_transfers, int32_t token_transfers_args_length_offset, int32_t token_transfer_data_offset, int64_t gas_limit, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*create_async_call_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t value_offset, int32_t data_offset, int32_t data_length, int32_t success_offset, int32_t success_length, int32_t error_offset, int32_t error_length, int64_t gas, int64_t extra_gas_for_callback);
  struct vm_exec_hook_result_t (*set_async_context_callback_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t callback, int32_t callback_length, int32_t data, int32_t data_length, int64_t gas);
  struct vm_exec_hook_result_t (*upgrade_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int64_t gas_limit, int32_t value_offset, int32_t code_offset, int32_t code_metadata_offset, int32_t length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*upgrade_from_source_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int64_t gas_limit, int32_t value_offset, int32_t source_contract_address_offset, int32_t code_metadata_offset, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*delete_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int64_t gas_limit, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*async_call_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_offset, int32_t value_offset, int32_t data_offset, int32_t length);
  struct vm_exec_hook_result_t (*get_argument_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id);
  struct vm_exec_hook_result_t (*get_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id, int32_t arg_offset);
  struct vm_exec_hook_result_t (*get_function_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t function_offset);
  struct vm_exec_hook_result_t (*get_num_arguments_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*storage_store_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int32_t data_offset, int32_t data_length);
  struct vm_exec_hook_result_t (*storage_load_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*storage_load_from_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t key_offset, int32_t key_length, int32_t data_offset);
  struct vm_exec_hook_result_t (*storage_load_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int32_t data_offset);
  struct vm_exec_hook_result_t (*set_storage_lock_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int64_t lock_timestamp);
  struct vm_exec_hook_result_t (*get_storage_lock_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*is_storage_locked_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*clear_storage_lock_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*get_caller_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset);
  struct vm_exec_hook_result_t (*check_no_payment_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_call_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_esdt_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_esdt_value_by_index_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset, int32_t index);
  struct vm_exec_hook_result_t (*get_esdt_token_name_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset);
  struct vm_exec_hook_result_t (*get_esdt_token_name_by_index_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_offset, int32_t index);
  struct vm_exec_hook_result_t (*get_esdt_token_nonce_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_esdt_token_nonce_by_index_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t index);
  struct vm_exec_hook_result_t (*get_current_esdt_nft_nonce_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len);
  struct vm_exec_hook_result_t (*get_esdt_token_type_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_esdt_token_type_by_index_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t index);
  struct vm_exec_hook_result_t (*get_num_esdt_transfers_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_call_value_token_name_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t call_value_offset, int32_t token_name_offset);
  struct vm_exec_hook_result_t (*get_call_value_token_name_by_index_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t call_value_offset, int32_t token_name_offset, int32_t index);
  struct vm_exec_hook_result_t (*is_reserved_function_name_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t name_handle);
  struct vm_exec_hook_result_t (*write_log_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_pointer, int32_t data_length, int32_t topic_ptr, int32_t num_topics);
  struct vm_exec_hook_result_t (*write_event_log_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t num_topics, int32_t topic_lengths_offset, int32_t topic_offset, int32_t data_offset, int32_t data_length);
  struct vm_exec_hook_result_t (*get_block_timestamp_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_block_timestamp_ms_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_block_nonce_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_block_round_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_block_epoch_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_block_random_seed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t pointer);
  struct vm_exec_hook_result_t (*get_state_root_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t pointer);
  struct vm_exec_hook_result_t (*get_prev_block_timestamp_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_prev_block_timestamp_ms_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_prev_block_nonce_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_prev_block_round_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_prev_block_epoch_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_prev_block_random_seed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t pointer);
  struct vm_exec_hook_result_t (*get_block_round_time_ms_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*epoch_start_block_timestamp_ms_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*epoch_start_block_nonce_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*epoch_start_block_round_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*finish_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t pointer, int32_t length);
  struct vm_exec_hook_result_t (*execute_on_same_context_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas_limit, int32_t address_offset, int32_t value_offset, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*execute_on_dest_context_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas_limit, int32_t address_offset, int32_t value_offset, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*execute_read_only_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas_limit, int32_t address_offset, int32_t function_offset, int32_t function_length, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*create_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas_limit, int32_t value_offset, int32_t code_offset, int32_t code_metadata_offset, int32_t length, int32_t result_offset, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*deploy_from_source_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas_limit, int32_t value_offset, int32_t source_contract_address_offset, int32_t code_metadata_offset, int32_t result_address_offset, int32_t num_arguments, int32_t arguments_length_offset, int32_t data_offset);
  struct vm_exec_hook_result_t (*get_num_return_data_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*get_return_data_size_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_id);
  struct vm_exec_hook_result_t (*get_return_data_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_id, int32_t data_offset);
  struct vm_exec_hook_result_t (*clean_return_data_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*delete_from_return_data_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_id);
  struct vm_exec_hook_result_t (*get_original_tx_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset);
  struct vm_exec_hook_result_t (*get_current_tx_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset);
  struct vm_exec_hook_result_t (*get_prev_tx_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset);
  struct vm_exec_hook_result_t (*managed_sc_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*managed_owner_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*managed_caller_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*managed_get_original_caller_addr_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*managed_get_relayer_addr_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*managed_signal_error_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t err_handle);
  struct vm_exec_hook_result_t (*managed_write_log_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t topics_handle, int32_t data_handle);
  struct vm_exec_hook_result_t (*managed_get_original_tx_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_get_state_root_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_get_block_random_seed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_get_prev_block_random_seed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_get_return_data_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t result_id, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_get_multi_esdt_call_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t multi_call_value_handle);
  struct vm_exec_hook_result_t (*managed_get_all_transfers_call_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t transfer_call_values_list_handle);
  struct vm_exec_hook_result_t (*managed_get_back_transfers_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t esdt_transfers_value_handle, int32_t egld_value_handle);
  struct vm_exec_hook_result_t (*managed_get_esdt_balance_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t token_id_handle, int64_t nonce, int32_t value_handle);
  struct vm_exec_hook_result_t (*managed_get_esdt_token_data_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t token_id_handle, int64_t nonce, int32_t value_handle, int32_t properties_handle, int32_t hash_handle, int32_t name_handle, int32_t attributes_handle, int32_t creator_handle, int32_t royalties_handle, int32_t uris_handle);
  struct vm_exec_hook_result_t (*managed_get_esdt_token_type_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t token_id_handle, int64_t nonce, int32_t type_handle);
  struct vm_exec_hook_result_t (*managed_async_call_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_handle, int32_t value_handle, int32_t function_handle, int32_t arguments_handle);
  struct vm_exec_hook_result_t (*managed_create_async_call_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_handle, int32_t value_handle, int32_t function_handle, int32_t arguments_handle, int32_t success_offset, int32_t success_length, int32_t error_offset, int32_t error_length, int64_t gas, int64_t extra_gas_for_callback, int32_t callback_closure_handle);
  struct vm_exec_hook_result_t (*managed_get_callback_closure_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t callback_closure_handle);
  struct vm_exec_hook_result_t (*managed_upgrade_from_source_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_handle, int64_t gas, int32_t value_handle, int32_t address_handle, int32_t code_metadata_handle, int32_t arguments_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_upgrade_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_handle, int64_t gas, int32_t value_handle, int32_t code_handle, int32_t code_metadata_handle, int32_t arguments_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_delete_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_handle, int64_t gas_limit, int32_t arguments_handle);
  struct vm_exec_hook_result_t (*managed_deploy_from_source_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas, int32_t value_handle, int32_t address_handle, int32_t code_metadata_handle, int32_t arguments_handle, int32_t result_address_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_create_contract_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas, int32_t value_handle, int32_t code_handle, int32_t code_metadata_handle, int32_t arguments_handle, int32_t result_address_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_execute_read_only_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas, int32_t address_handle, int32_t function_handle, int32_t arguments_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_execute_on_same_context_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas, int32_t address_handle, int32_t value_handle, int32_t function_handle, int32_t arguments_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_execute_on_dest_context_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas, int32_t address_handle, int32_t value_handle, int32_t function_handle, int32_t arguments_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_execute_on_dest_context_with_error_return_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t gas, int32_t address_handle, int32_t value_handle, int32_t function_handle, int32_t arguments_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*managed_multi_transfer_esdt_nft_execute_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dst_handle, int32_t token_transfers_handle, int64_t gas_limit, int32_t function_handle, int32_t arguments_handle);
  struct vm_exec_hook_result_t (*managed_multi_transfer_esdt_nft_execute_with_return_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dst_handle, int32_t token_transfers_handle, int64_t gas_limit, int32_t function_handle, int32_t arguments_handle);
  struct vm_exec_hook_result_t (*managed_multi_transfer_esdt_nft_execute_by_user_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t user_handle, int32_t dst_handle, int32_t token_transfers_handle, int64_t gas_limit, int32_t function_handle, int32_t arguments_handle);
  struct vm_exec_hook_result_t (*managed_transfer_value_execute_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dst_handle, int32_t value_handle, int64_t gas_limit, int32_t function_handle, int32_t arguments_handle);
  struct vm_exec_hook_result_t (*managed_is_esdt_frozen_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t token_id_handle, int64_t nonce);
  struct vm_exec_hook_result_t (*managed_is_esdt_limited_transfer_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t token_id_handle);
  struct vm_exec_hook_result_t (*managed_is_esdt_paused_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t token_id_handle);
  struct vm_exec_hook_result_t (*managed_buffer_to_hex_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t source_handle, int32_t dest_handle);
  struct vm_exec_hook_result_t (*managed_get_code_metadata_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t response_handle);
  struct vm_exec_hook_result_t (*managed_get_code_hash_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t code_hash_handle);
  struct vm_exec_hook_result_t (*managed_is_builtin_function_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t function_name_handle);
  struct vm_exec_hook_result_t (*big_float_new_from_parts_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t integral_part, int32_t fractional_part, int32_t exponent);
  struct vm_exec_hook_result_t (*big_float_new_from_frac_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t numerator, int64_t denominator);
  struct vm_exec_hook_result_t (*big_float_new_from_sci_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t significand, int64_t exponent);
  struct vm_exec_hook_result_t (*big_float_add_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_float_sub_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_float_mul_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_float_div_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_float_neg_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_clone_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_cmp_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_float_abs_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_sign_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_sqrt_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_pow_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle, int32_t exponent);
  struct vm_exec_hook_result_t (*big_float_floor_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_big_int_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_ceil_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_big_int_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_truncate_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t dest_big_int_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_set_int64_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int64_t value);
  struct vm_exec_hook_result_t (*big_float_is_int_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_float_set_big_int_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t big_int_handle);
  struct vm_exec_hook_result_t (*big_float_get_const_pi_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_float_get_const_e_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_get_unsigned_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_get_signed_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_storage_store_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int32_t source_handle);
  struct vm_exec_hook_result_t (*big_int_storage_load_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_get_call_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_get_esdt_call_value_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination);
  struct vm_exec_hook_result_t (*big_int_get_esdt_call_value_by_index_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t index);
  struct vm_exec_hook_result_t (*big_int_get_external_balance_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t result);
  struct vm_exec_hook_result_t (*big_int_get_esdt_external_balance_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_offset, int32_t token_id_offset, int32_t token_id_len, int64_t nonce, int32_t result_handle);
  struct vm_exec_hook_result_t (*big_int_new_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t small_value);
  struct vm_exec_hook_result_t (*big_int_unsigned_byte_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t reference_handle);
  struct vm_exec_hook_result_t (*big_int_signed_byte_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t reference_handle);
  struct vm_exec_hook_result_t (*big_int_get_unsigned_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t reference_handle, int32_t byte_offset);
  struct vm_exec_hook_result_t (*big_int_get_signed_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t reference_handle, int32_t byte_offset);
  struct vm_exec_hook_result_t (*big_int_set_unsigned_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t byte_offset, int32_t byte_length);
  struct vm_exec_hook_result_t (*big_int_set_signed_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t byte_offset, int32_t byte_length);
  struct vm_exec_hook_result_t (*big_int_is_int64_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_get_int64_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle);
  struct vm_exec_hook_result_t (*big_int_set_int64_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int64_t value);
  struct vm_exec_hook_result_t (*big_int_add_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_sub_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_mul_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_tdiv_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_tmod_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_ediv_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_emod_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_sqrt_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_int_pow_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_log2_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t op1_handle);
  struct vm_exec_hook_result_t (*big_int_abs_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_int_neg_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_int_sign_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_int_cmp_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_not_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle);
  struct vm_exec_hook_result_t (*big_int_and_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_or_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_xor_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op1_handle, int32_t op2_handle);
  struct vm_exec_hook_result_t (*big_int_shr_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle, int32_t bits);
  struct vm_exec_hook_result_t (*big_int_shl_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t op_handle, int32_t bits);
  struct vm_exec_hook_result_t (*big_int_finish_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t reference_handle);
  struct vm_exec_hook_result_t (*big_int_finish_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t reference_handle);
  struct vm_exec_hook_result_t (*big_int_to_string_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t big_int_handle, int32_t destination_handle);
  struct vm_exec_hook_result_t (*mbuffer_new_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*mbuffer_new_from_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset, int32_t data_length);
  struct vm_exec_hook_result_t (*mbuffer_get_length_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle);
  struct vm_exec_hook_result_t (*mbuffer_get_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t result_offset);
  struct vm_exec_hook_result_t (*mbuffer_get_byte_slice_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t source_handle, int32_t starting_position, int32_t slice_length, int32_t result_offset);
  struct vm_exec_hook_result_t (*mbuffer_copy_byte_slice_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t source_handle, int32_t starting_position, int32_t slice_length, int32_t destination_handle);
  struct vm_exec_hook_result_t (*mbuffer_eq_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle1, int32_t m_buffer_handle2);
  struct vm_exec_hook_result_t (*mbuffer_set_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t data_offset, int32_t data_length);
  struct vm_exec_hook_result_t (*mbuffer_set_byte_slice_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t starting_position, int32_t data_length, int32_t data_offset);
  struct vm_exec_hook_result_t (*mbuffer_append_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t accumulator_handle, int32_t data_handle);
  struct vm_exec_hook_result_t (*mbuffer_append_bytes_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t accumulator_handle, int32_t data_offset, int32_t data_length);
  struct vm_exec_hook_result_t (*mbuffer_to_big_int_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t big_int_handle);
  struct vm_exec_hook_result_t (*mbuffer_to_big_int_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t big_int_handle);
  struct vm_exec_hook_result_t (*mbuffer_from_big_int_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t big_int_handle);
  struct vm_exec_hook_result_t (*mbuffer_from_big_int_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t big_int_handle);
  struct vm_exec_hook_result_t (*mbuffer_to_small_int_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle);
  struct vm_exec_hook_result_t (*mbuffer_to_small_int_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle);
  struct vm_exec_hook_result_t (*mbuffer_from_small_int_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int64_t value);
  struct vm_exec_hook_result_t (*mbuffer_from_small_int_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int64_t value);
  struct vm_exec_hook_result_t (*mbuffer_to_big_float_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t big_float_handle);
  struct vm_exec_hook_result_t (*mbuffer_from_big_float_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_buffer_handle, int32_t big_float_handle);
  struct vm_exec_hook_result_t (*mbuffer_storage_store_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t source_handle);
  struct vm_exec_hook_result_t (*mbuffer_storage_load_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t destination_handle);
  struct vm_exec_hook_result_t (*mbuffer_storage_load_from_address_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t address_handle, int32_t key_handle, int32_t destination_handle);
  struct vm_exec_hook_result_t (*mbuffer_get_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id, int32_t destination_handle);
  struct vm_exec_hook_result_t (*mbuffer_finish_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t source_handle);
  struct vm_exec_hook_result_t (*mbuffer_set_random_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t destination_handle, int32_t length);
  struct vm_exec_hook_result_t (*managed_map_new_func_ptr)(struct vm_exec_hook_context_t *ctx);
  struct vm_exec_hook_result_t (*managed_map_put_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_map_handle, int32_t key_handle, int32_t value_handle);
  struct vm_exec_hook_result_t (*managed_map_get_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_map_handle, int32_t key_handle, int32_t out_value_handle);
  struct vm_exec_hook_result_t (*managed_map_remove_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_map_handle, int32_t key_handle, int32_t out_value_handle);
  struct vm_exec_hook_result_t (*managed_map_contains_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t m_map_handle, int32_t key_handle);
  struct vm_exec_hook_result_t (*small_int_get_unsigned_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id);
  struct vm_exec_hook_result_t (*small_int_get_signed_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id);
  struct vm_exec_hook_result_t (*small_int_finish_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t value);
  struct vm_exec_hook_result_t (*small_int_finish_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t value);
  struct vm_exec_hook_result_t (*small_int_storage_store_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int64_t value);
  struct vm_exec_hook_result_t (*small_int_storage_store_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int64_t value);
  struct vm_exec_hook_result_t (*small_int_storage_load_unsigned_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*small_int_storage_load_signed_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*int64get_argument_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t id);
  struct vm_exec_hook_result_t (*int64finish_func_ptr)(struct vm_exec_hook_context_t *ctx, int64_t value);
  struct vm_exec_hook_result_t (*int64storage_store_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int64_t value);
  struct vm_exec_hook_result_t (*int64storage_load_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length);
  struct vm_exec_hook_result_t (*sha256_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset, int32_t length, int32_t result_offset);
  struct vm_exec_hook_result_t (*managed_sha256_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t input_handle, int32_t output_handle);
  struct vm_exec_hook_result_t (*keccak256_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset, int32_t length, int32_t result_offset);
  struct vm_exec_hook_result_t (*managed_keccak256_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t input_handle, int32_t output_handle);
  struct vm_exec_hook_result_t (*ripemd160_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset, int32_t length, int32_t result_offset);
  struct vm_exec_hook_result_t (*managed_ripemd160_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t input_handle, int32_t output_handle);
  struct vm_exec_hook_result_t (*verify_bls_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t message_offset, int32_t message_length, int32_t sig_offset);
  struct vm_exec_hook_result_t (*managed_verify_bls_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle);
  struct vm_exec_hook_result_t (*verify_ed25519_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t message_offset, int32_t message_length, int32_t sig_offset);
  struct vm_exec_hook_result_t (*managed_verify_ed25519_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle);
  struct vm_exec_hook_result_t (*verify_custom_secp256k1_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int32_t message_offset, int32_t message_length, int32_t sig_offset, int32_t hash_type);
  struct vm_exec_hook_result_t (*managed_verify_custom_secp256k1_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle, int32_t hash_type);
  struct vm_exec_hook_result_t (*verify_secp256k1_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_offset, int32_t key_length, int32_t message_offset, int32_t message_length, int32_t sig_offset);
  struct vm_exec_hook_result_t (*managed_verify_secp256k1_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle);
  struct vm_exec_hook_result_t (*encode_secp256k1_der_signature_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t r_offset, int32_t r_length, int32_t s_offset, int32_t s_length, int32_t sig_offset);
  struct vm_exec_hook_result_t (*managed_encode_secp256k1_der_signature_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t r_handle, int32_t s_handle, int32_t sig_handle);
  struct vm_exec_hook_result_t (*add_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t fst_point_xhandle, int32_t fst_point_yhandle, int32_t snd_point_xhandle, int32_t snd_point_yhandle);
  struct vm_exec_hook_result_t (*double_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t point_xhandle, int32_t point_yhandle);
  struct vm_exec_hook_result_t (*is_on_curve_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t ec_handle, int32_t point_xhandle, int32_t point_yhandle);
  struct vm_exec_hook_result_t (*scalar_base_mult_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t data_offset, int32_t length);
  struct vm_exec_hook_result_t (*managed_scalar_base_mult_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t data_handle);
  struct vm_exec_hook_result_t (*scalar_mult_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t point_xhandle, int32_t point_yhandle, int32_t data_offset, int32_t length);
  struct vm_exec_hook_result_t (*managed_scalar_mult_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t point_xhandle, int32_t point_yhandle, int32_t data_handle);
  struct vm_exec_hook_result_t (*marshal_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_pair_handle, int32_t y_pair_handle, int32_t ec_handle, int32_t result_offset);
  struct vm_exec_hook_result_t (*managed_marshal_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_pair_handle, int32_t y_pair_handle, int32_t ec_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*marshal_compressed_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_pair_handle, int32_t y_pair_handle, int32_t ec_handle, int32_t result_offset);
  struct vm_exec_hook_result_t (*managed_marshal_compressed_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_pair_handle, int32_t y_pair_handle, int32_t ec_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*unmarshal_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t data_offset, int32_t length);
  struct vm_exec_hook_result_t (*managed_unmarshal_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t data_handle);
  struct vm_exec_hook_result_t (*unmarshal_compressed_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t data_offset, int32_t length);
  struct vm_exec_hook_result_t (*managed_unmarshal_compressed_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_result_handle, int32_t y_result_handle, int32_t ec_handle, int32_t data_handle);
  struct vm_exec_hook_result_t (*generate_key_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_pub_key_handle, int32_t y_pub_key_handle, int32_t ec_handle, int32_t result_offset);
  struct vm_exec_hook_result_t (*managed_generate_key_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t x_pub_key_handle, int32_t y_pub_key_handle, int32_t ec_handle, int32_t result_handle);
  struct vm_exec_hook_result_t (*create_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_offset, int32_t data_length);
  struct vm_exec_hook_result_t (*managed_create_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t data_handle);
  struct vm_exec_hook_result_t (*get_curve_length_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t ec_handle);
  struct vm_exec_hook_result_t (*get_priv_key_byte_length_ec_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t ec_handle);
  struct vm_exec_hook_result_t (*elliptic_curve_get_values_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t ec_handle, int32_t field_order_handle, int32_t base_point_order_handle, int32_t eq_constant_handle, int32_t x_base_point_handle, int32_t y_base_point_handle);
  struct vm_exec_hook_result_t (*managed_verify_secp256r1_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle);
  struct vm_exec_hook_result_t (*managed_verify_blssignature_share_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle);
  struct vm_exec_hook_result_t (*managed_verify_blsaggregated_signature_func_ptr)(struct vm_exec_hook_context_t *ctx, int32_t key_handle, int32_t message_handle, int32_t sig_handle);
} vm_exec_vm_hook_c_func_pointers;

/**
 * A buffer of this size holds every message the executor currently produces.
 *
 * Exported so the Go side has one place to take the number from rather than picking its
 * own.
 */
extern const int VM_EXEC_ERROR_BUF_RECOMMENDED_LEN;

#endif  /* VM_EXEC_COMMON_H */
