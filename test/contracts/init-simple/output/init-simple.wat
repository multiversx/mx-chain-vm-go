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
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    local.get 0
    i32.const 0
    i32.load offset=1031 align=1
    i32.store offset=7 align=1
    local.get 0
    i32.const 0
    i64.load offset=1024 align=1
    i64.store
    local.get 0
    i32.const 10
    call 1
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66576))
  (export "memory" (memory 0))
  (export "init" (func 2))
  (export "dummy" (func 3))
  (data (;0;) (i32.const 1024) "dummy text\00"))
