(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
  (table $tab 2 funcref)
  (elem $elems func $main)
  (func $main (type $void)
    (table.init $tab $elems
      (i32.const 0)   ;; destination index in the table
      (i32.const 0)   ;; offset into the source element segment
      (i32.const 1)   ;; number of entries to copy
    )
    (i64.const 0)
    (call $int64finish)
  )
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
)
