(module
  (type $void (func))
  (func $main (type $void))
  (memory $mem 10 42)
  (export "memory" (memory $mem))
  (export "empty" (func $main))
)