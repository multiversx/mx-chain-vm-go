(module
  (type (;0;) (func (param i64)))
  (type (;1;) (func))
  (import "env" "int64finish" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    global.get 2
    call 0
  )
  (memory (;0;) 2)  
  (export "memory" (memory 0))
  (global (;0;) (mut i64) (i64.const 42))
  (global (;1;) (mut i64) (i64.const 43))
  (export "getglobal" (func 1))
)
