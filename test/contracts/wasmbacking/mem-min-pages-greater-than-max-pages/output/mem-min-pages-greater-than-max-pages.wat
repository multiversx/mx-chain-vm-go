;; compile with:
;; wat2wasm --no-check
(module
  (type $void (func))
  (func $main (type $void))
  (memory $mem 1 0)
  (export "memory" (memory $mem))
  (export "main" (func $main))
)
