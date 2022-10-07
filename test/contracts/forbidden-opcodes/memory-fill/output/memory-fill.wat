(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
  (func $main (type $void)
    (memory.fill
      (i32.const 42)    ;; destination address
      (i32.const 5)     ;; byte value to set
      (i32.const 7)     ;; size of memory region in bytes
    )
    (i64.const 0)
    (call $int64finish)
  )
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
)
