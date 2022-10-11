;; compile with:
;; wat2wasm --no-check
(module
  (type $void (func))
  (func $main (type $void))
  (memory $mem1 1 2)
  (memory $mem2 2 4)
  (memory $mem3 3 6)
  (memory $mem4 4 8)
  (memory $mem5 5 10)
  (export "memory" (memory $mem1))
  (export "main" (func $main))
)
