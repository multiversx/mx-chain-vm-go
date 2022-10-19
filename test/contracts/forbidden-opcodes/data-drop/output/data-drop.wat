(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
  (func $main (type $void)
    (data.drop 0)   ;; drop the data segment 0
    (i64.const 0)
    (call $int64finish)
  )
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
  (data "ok")       ;; data segment 0
)
