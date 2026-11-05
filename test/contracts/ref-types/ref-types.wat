;; Uses reference types, which the V1 opcode whitelist accepts and the V2 one drops.
;; Used to check that code cached under V1 is not run under V2.
(module
  (type $void (func))
  (memory $mem 2)
  (func $init (type $void))
  (func $isNull (type $void)
    (drop (ref.is_null (ref.null func)))
  )
  (export "memory" (memory $mem))
  (export "init" (func $init))
  (export "isNull" (func $isNull))
)
