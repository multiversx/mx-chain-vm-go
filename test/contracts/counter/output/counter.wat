(module
  (type $t0 (func (param i32 i32 i64) (result i32)))
  (type $t1 (func (param i32 i32) (result i64)))
  (type $t2 (func (param i64)))
  (type $t3 (func))
  (import "env" "int64storageStore" (func $env.int64storageStore (type $t0)))
  (import "env" "int64storageLoad" (func $env.int64storageLoad (type $t1)))
  (import "env" "int64finish" (func $env.int64finish (type $t2)))
  (func $init (type $t3)
    i32.const 1024
    i32.const 7
    i64.const 1
    call $env.int64storageStore
    drop)
  (func $increment (type $t3)
    (local $l0 i64)
    i32.const 1024
    i32.const 7
    i32.const 1024
    i32.const 7
    call $env.int64storageLoad
    i64.const 1
    i64.add
    local.tee $l0
    call $env.int64storageStore
    drop
    local.get $l0
    call $env.int64finish)
  (func $decrement (type $t3)
    (local $l0 i64)
    i32.const 1024
    i32.const 7
    i32.const 1024
    i32.const 7
    call $env.int64storageLoad
    i64.const -1
    i64.add
    local.tee $l0
    call $env.int64storageStore
    drop
    local.get $l0
    call $env.int64finish)
  (func $get (type $t3)
    i32.const 1024
    i32.const 7
    call $env.int64storageLoad
    call $env.int64finish)
  (table $T0 1 1 funcref)
  (memory $memory 2)
  (global $g0 (mut i32) (i32.const 66576))
  (export "memory" (memory 0))
  (export "init" (func $init))
  (export "increment" (func $increment))
  (export "decrement" (func $decrement))
  (export "get" (func $get))
  (data $d0 (i32.const 1024) "COUNTER\00"))
