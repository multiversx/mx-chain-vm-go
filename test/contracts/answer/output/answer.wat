(module
  (type (;0;) (func (param i64)))
  (type (;1;) (func))
  (import "env" "int64finish" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    (local i64 i64 i64 i64)
    i64.const 42
    call 0)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66560))
  (export "memory" (memory 0))
  (export "answer" (func 1)))
