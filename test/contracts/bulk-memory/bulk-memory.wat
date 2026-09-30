;; Contract used to measure the gas cost of the bulk memory opcodes.
;;
;; Each endpoint takes the number of bytes to copy/fill as argument 0, so that the
;; same code path can be executed with different sizes. Everything apart from the
;; bulk memory opcode itself costs the same on every call, which means that the
;; difference in gas between two calls is exactly
;;     (bytes1 - bytes2) * <opcode>PerByte
;;
;; Memory is 2 pages, the second one is the source of the copy, so that source and
;; destination never overlap, for sizes up to 65536.
(module
  (type $void (func))
  (type $getArgument (func (param i32) (result i64)))

  (import "env" "int64getArgument" (func $int64getArgument (type $getArgument)))

  (memory $mem 2)

  (func $init (type $void))

  ;; number of bytes to copy comes from argument 0
  (func $memoryCopy (type $void)
    (memory.copy
      (i32.const 0)      ;; destination address, first page
      (i32.const 65536)  ;; source address, second page
      (i32.wrap_i64 (call $int64getArgument (i32.const 0)))
    )
  )

  ;; number of bytes to fill comes from argument 0
  (func $memoryFill (type $void)
    (memory.fill
      (i32.const 0)      ;; destination address, first page
      (i32.const 7)      ;; value to fill with
      (i32.wrap_i64 (call $int64getArgument (i32.const 0)))
    )
  )

  ;; same as $memoryCopy, but the size goes through a local first, to check that the
  ;; metering injection puts the size operand back on the stack unchanged
  (func $memoryCopyViaLocal (type $void) (local $size i32)
    (local.set $size (i32.wrap_i64 (call $int64getArgument (i32.const 0))))
    (memory.copy
      (i32.const 0)      ;; destination address, first page
      (i32.const 65536)  ;; source address, second page
      (local.get $size)
    )
  )

  ;; same as $memoryFill, but the size goes through a local first
  (func $memoryFillViaLocal (type $void) (local $size i32)
    (local.set $size (i32.wrap_i64 (call $int64getArgument (i32.const 0))))
    (memory.fill
      (i32.const 0)      ;; destination address, first page
      (i32.const 7)      ;; value to fill with
      (local.get $size)
    )
  )

  ;; everything the *ViaLocal endpoints do, except the bulk memory opcode and its three
  ;; operands, so that the base cost of the opcode can be isolated by subtraction. Also
  ;; proves that decoding the argument costs the same regardless of its value.
  (func $noBulkMemory (type $void) (local $size i32)
    (local.set $size (i32.wrap_i64 (call $int64getArgument (i32.const 0))))
  )

  (export "memory" (memory $mem))
  (export "init" (func $init))
  (export "memoryCopy" (func $memoryCopy))
  (export "memoryFill" (func $memoryFill))
  (export "memoryCopyViaLocal" (func $memoryCopyViaLocal))
  (export "memoryFillViaLocal" (func $memoryFillViaLocal))
  (export "noBulkMemory" (func $noBulkMemory))
)
