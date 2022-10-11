(module
  (type (;0;) (func (param i64)))
  (type (;1;) (func))
  (import "env" "int64finish" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    global.get 0
    i64.const 1
    i64.add
    global.set 0

    global.get 0
    call 0
  )
  (memory (;0;) 2)
  (global (;0;) (mut i64) (i64.const 66560))
  (export "memory" (memory 0))
  (export "getglobal" (func 1))
)
