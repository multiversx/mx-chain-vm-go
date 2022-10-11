(module
  (type (;0;) (func (param i64)))
  (type (;1;) (func))
  (import "env" "int64finish" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    global.get 1
    i64.const 1
    i64.add
    global.set 1

    global.get 2
    i64.const 1
    i64.add
    global.set 2

    global.get 3
    i64.const 1
    i64.add
    global.set 3

    global.get 4
    i64.const 1
    i64.add
    global.set 4

    global.get 0
    call 0
    global.get 1
    call 0
    global.get 2
    call 0
    global.get 3
    call 0
    global.get 4
    call 0
  )
  (memory (;0;) 2)
  (global (;0;) (mut i64) (i64.const 0))
  (global (;1;) (mut i64) (i64.const 1))
  (global (;2;) (mut i64) (i64.const 2))
  (global (;3;) (mut i64) (i64.const 4))
  (global (;4;) (mut i64) (i64.const 6))
  (export "memory" (memory 0))
  (export "increment_globals" (func 1))
)
