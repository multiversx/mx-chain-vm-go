(module
  (type (;0;) (func (param i32 i32 i64) (result i32)))
  (type (;1;) (func (param i32 i32) (result i64)))
  (type (;2;) (func (param i64)))
  (type (;3;) (func))
  (import "env" "int64storageStore" (func (;0;) (type 0)))
  (import "env" "int64storageLoad" (func (;1;) (type 1)))
  (import "env" "int64finish" (func (;2;) (type 2)))
  (func (;3;) (type 3)
    (local i32 i32 i64)
    i32.const 1024
    local.set 0
    i32.const 7
    local.set 1
    i64.const 1
    local.set 2
    local.get 0
    local.get 1
    local.get 2
    call 0
    drop
    return)
  (func (;4;) (type 3)
    (local i32 i32 i64)
    i32.const 1024
    local.set 0
    i32.const 7
    local.set 1
    i64.const 1
    local.set 2
    local.get 0
    local.get 1
    local.get 2
    call 0
    drop
    return)
  (func (;5;) (type 3)
    (local i32 i32 i32 i32 i32 i64 i64 i64 i64 i64 i64 i32 i32)
    global.get 0
    local.set 0
    i32.const 16
    local.set 1
    local.get 0
    local.get 1
    i32.sub
    local.set 2
    local.get 2
    global.set 0
    i32.const 1024
    local.set 3
    i32.const 7
    local.set 4
    local.get 3
    local.get 4
    call 1
    local.set 5
    local.get 2
    local.get 5
    i64.store offset=8
    local.get 2
    i64.load offset=8
    local.set 6
    i64.const 1
    local.set 7
    local.get 6
    local.get 7
    i64.add
    local.set 8
    local.get 2
    local.get 8
    i64.store offset=8
    local.get 2
    i64.load offset=8
    local.set 9
    local.get 3
    local.get 4
    local.get 9
    call 0
    drop
    local.get 2
    i64.load offset=8
    local.set 10
    local.get 10
    call 2
    i32.const 16
    local.set 11
    local.get 2
    local.get 11
    i32.add
    local.set 12
    local.get 12
    global.set 0
    return)
  (func (;6;) (type 3)
    (local i32 i32 i32 i32 i32 i64 i64 i64 i64 i64 i64 i32 i32)
    global.get 0
    local.set 0
    i32.const 16
    local.set 1
    local.get 0
    local.get 1
    i32.sub
    local.set 2
    local.get 2
    global.set 0
    i32.const 1024
    local.set 3
    i32.const 7
    local.set 4
    local.get 3
    local.get 4
    call 1
    local.set 5
    local.get 2
    local.get 5
    i64.store offset=8
    local.get 2
    i64.load offset=8
    local.set 6
    i64.const -1
    local.set 7
    local.get 6
    local.get 7
    i64.add
    local.set 8
    local.get 2
    local.get 8
    i64.store offset=8
    local.get 2
    i64.load offset=8
    local.set 9
    local.get 3
    local.get 4
    local.get 9
    call 0
    drop
    local.get 2
    i64.load offset=8
    local.set 10
    local.get 10
    call 2
    i32.const 16
    local.set 11
    local.get 2
    local.get 11
    i32.add
    local.set 12
    local.get 12
    global.set 0
    return)
  (func (;7;) (type 3)
    (local i32 i32 i32 i32 i32 i64 i64 i32 i32)
    global.get 0
    local.set 0
    i32.const 16
    local.set 1
    local.get 0
    local.get 1
    i32.sub
    local.set 2
    local.get 2
    global.set 0
    i32.const 1024
    local.set 3
    i32.const 7
    local.set 4
    local.get 3
    local.get 4
    call 1
    local.set 5
    local.get 2
    local.get 5
    i64.store offset=8
    local.get 2
    i64.load offset=8
    local.set 6
    local.get 6
    call 2
    i32.const 16
    local.set 7
    local.get 2
    local.get 7
    i32.add
    local.set 8
    local.get 8
    global.set 0
    return)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66576))
  (export "memory" (memory 0))
  (export "init" (func 3))
  (export "upgrade" (func 4))
  (export "increment" (func 5))
  (export "decrement" (func 6))
  (export "get" (func 7))
  (data (;0;) (i32.const 1024) "COUNTER\00"))
