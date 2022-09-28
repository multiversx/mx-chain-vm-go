(module
  (type (;0;) (func (param i64)))
  (type (;1;) (func))
  (import "env" "int64finish" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    i64.const 42
    call 0
  )
  (memory (;0;) 2)
  (export "memory" (memory 0))
  (export "getnumber" (func 1))
)
