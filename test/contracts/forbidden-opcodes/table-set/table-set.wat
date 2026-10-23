(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
  (table $tab 1 funcref)
  (elem (i32.const 0) $main)
  (func $main (type $void)
    (table.set $tab
      (i32.const 0)             ;; index into the table
      (ref.func $main)          ;; function reference to store
    )
    (i64.const 0)
    (call $int64finish)
  )
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
)
