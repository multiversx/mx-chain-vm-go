(module
  (type (;0;) (func (param i64)))
  (type (;1;) (func (param i32 i32)))
  (type (;2;) (func))
  (import "env" "int64finish" (func (;0;) (type 0)))
  (import "env" "finish" (func (;1;) (type 1)))
  (func (;2;) (type 2)
    i64.const 42
    call 0)
  (func (;3;) (type 2)
    i32.const 1024
    i32.const 10
    call 1)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66576))
  (export "memory" (memory 0))
  (export "init" (func 2))
  (export "dummy" (func 3))
  (data (;0;) (i32.const 1024) "finish0000"))
