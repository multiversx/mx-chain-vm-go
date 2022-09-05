
#if !defined(WASMER_H_MACROS)

#define WASMER_H_MACROS

// Define the `ARCH_X86_X64` constant.
#if defined(MSVC) && defined(_M_AMD64)
#  define ARCH_X86_64
#elif (defined(GCC) || defined(__GNUC__) || defined(__clang__)) && defined(__x86_64__)
#  define ARCH_X86_64
#endif

// Compatibility with non-Clang compilers.
#if !defined(__has_attribute)
#  define __has_attribute(x) 0
#endif

// Compatibility with non-Clang compilers.
#if !defined(__has_declspec_attribute)
#  define __has_declspec_attribute(x) 0
#endif

// Define the `DEPRECATED` macro.
#if defined(GCC) || defined(__GNUC__) || __has_attribute(deprecated)
#  define DEPRECATED(message) __attribute__((deprecated(message)))
#elif defined(MSVC) || __has_declspec_attribute(deprecated)
#  define DEPRECATED(message) __declspec(deprecated(message))
#endif

#endif // WASMER_H_MACROS


#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * The `wasmer_result_t` enum is a type that represents either a
 * success, or a failure.
 */
typedef enum {
  /**
   * Represents a success.
   */
  VM_EXEC_OK = 1,
  /**
   * Represents a failure.
   */
  VM_EXEC_ERROR = 2,
} vm_exec_result_t;

/**
 * Represents all possibles WebAssembly value types.
 *
 * See `wasmer_value_t` to get a complete example.
 */
enum vm_exec_value_tag {
  /**
   * Represents the `i32` WebAssembly type.
   */
  VM_EXEC_VALUE_I32,
  /**
   * Represents the `i64` WebAssembly type.
   */
  VM_EXEC_VALUE_I64,
};
typedef uint32_t vm_exec_value_tag;

typedef struct {

} vm_exec_import_func_t;

typedef struct {
  const uint8_t *bytes;
  uint32_t bytes_len;
} vm_exec_byte_array;

typedef struct {
  vm_exec_byte_array module_name;
  vm_exec_byte_array import_name;
  const vm_exec_import_func_t *import_func;
} vm_exec_import_t;

int vm_exec_execution_info_length(void);

int vm_exec_execution_info_message(char *dest_buffer, int dest_buffer_len);

/**
 * Frees memory for the given Func
 */
void vm_exec_import_func_destroy(vm_exec_import_func_t *func);

/**
 * Creates new host function, aka imported function. `func` is a
 * function pointer, where the first argument is the famous `vm::Ctx`
 * (in Rust), or `wasmer_instance_context_t` (in C). All arguments
 * must be typed with compatible WebAssembly native types:
 *
 * | WebAssembly type | C/C++ type |
 * | ---------------- | ---------- |
 * | `i32`            | `int32_t`  |
 * | `i64`            | `int64_t`  |
 * | `f32`            | `float`    |
 * | `f64`            | `double`   |
 *
 * The function pointer must have a lifetime greater than the
 * WebAssembly instance lifetime.
 *
 * The caller owns the object and should call
 * `wasmer_import_func_destroy` to free it.
 */
vm_exec_import_func_t *vm_exec_import_func_new(void (*func)(void *data),
                                               const vm_exec_value_tag *params,
                                               unsigned int params_len,
                                               const vm_exec_value_tag *returns,
                                               unsigned int returns_len);

/**
 * Gets the length in bytes of the last error if any.
 *
 * This can be used to dynamically allocate a buffer with the correct number of
 * bytes needed to store a message.
 */
int vm_exec_last_error_length(void);

/**
 * Gets the last error message if any into the provided buffer
 * `buffer` up to the given `length`.
 *
 * The `length` parameter must be large enough to store the last
 * error message. Ideally, the value should come from
 * `wasmer_last_error_length()`.
 *
 * The function returns the length of the string in bytes, `-1` if an
 * error occurs. Potential errors are:
 *
 *  * The buffer is a null pointer,
 *  * The buffer is too smal to hold the error message.
 *
 * Note: The error message always has a trailing null character.
 * ```
 */
int vm_exec_last_error_message(char *dest_buffer, int dest_buffer_len);

vm_exec_result_t vm_exec_set_imports(vm_exec_import_t *imports, unsigned int imports_len);

void vm_exec_set_sigsegv_passthrough(void);
