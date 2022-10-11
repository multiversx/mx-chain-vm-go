(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
  (func $main (type $void)
    (memory.init 0    ;; init from data segment 0
      (i32.const 42)  ;; destination address
      (i32.const 0)   ;; offset into the source segment
      (i32.const 2)   ;; size of memory region in bytes
    )
    (i64.const 0)
    (call $int64finish)
  )
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
  (data "ok")         ;; data segment 0
)
