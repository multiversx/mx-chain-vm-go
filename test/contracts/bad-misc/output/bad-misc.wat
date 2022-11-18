(module
  (type (;0;) (func (param i32 i32 i32 i32)))
  (type (;1;) (func (param i32)))
  (type (;2;) (func (param i64 i32) (result i32)))
  (type (;3;) (func (param i64) (result i32)))
  (type (;4;) (func (param i32 i32 i32) (result i32)))
  (type (;5;) (func))
  (type (;6;) (func (param i64) (result i64)))
  (import "env" "writeLog" (func (;0;) (type 0)))
  (import "env" "getSCAddress" (func (;1;) (type 1)))
  (import "env" "getBlockHash" (func (;2;) (type 2)))
  (import "env" "bigIntNew" (func (;3;) (type 3)))
  (import "env" "bigIntStorageStoreUnsigned" (func (;4;) (type 4)))
  (import "env" "bigIntFinishUnsigned" (func (;5;) (type 1)))
  (func (;6;) (type 5)
    i32.const 0
    i32.const 42
    i32.store8 offset=2147484671)
  (func (;7;) (type 5)
    (local i32)
    local.get 0
    local.get 0
    local.get 0
    local.get 0
    call 0)
  (func (;8;) (type 5)
    i32.const 2147483647
    call 1)
  (func (;9;) (type 5)
    i64.const 0
    i32.const 2147483647
    call 2
    drop)
  (func (;10;) (type 5)
    i64.const 2147483647
    i32.const 0
    call 2
    drop)
  (func (;11;) (type 5)
    i64.const 1
    i32.const 2147483647
    call 2
    drop)
  (func (;12;) (type 5)
    i32.const 0
    i32.const 1
    i32.const 0
    i32.const -1
    call 0)
  (func (;13;) (type 5)
    i32.const 0
    i32.const -1
    i32.const 0
    i32.const 0
    call 0)
  (func (;14;) (type 5)
    i32.const 2147483647
    i32.const 0
    i32.const 0
    i32.const 0
    call 0)
  (func (;15;) (type 5)
    i32.const 0
    i32.const 0
    i32.const 2147483647
    i32.const 500
    call 0)
  (func (;16;) (type 5)
    i32.const 1056
    i32.const 32
    i64.const 100
    call 3
    i32.const 42
    i32.add
    call 4
    drop)
  (func (;17;) (type 6) (param i64) (result i64)
    block  ;; label = @1
      local.get 0
      i64.const 1
      i64.and
      i64.eqz
      br_if 0 (;@1;)
      local.get 0
      i64.const 3
      i64.shl
      i64.const 1
      i64.or
      call 17
      local.get 0
      i64.add
      local.get 0
      i64.const 1
      i64.shl
      i64.const 1
      i64.or
      call 17
      i64.add
      return
    end
    i64.const 42)
  (func (;18;) (type 5)
    i64.const 1
    call 17
    call 3
    call 5)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66608))
  (export "memory" (memory 0))
  (export "memoryFault" (func 6))
  (export "divideByZero" (func 7))
  (export "badGetOwner1" (func 8))
  (export "badGetBlockHash1" (func 9))
  (export "badGetBlockHash2" (func 10))
  (export "badGetBlockHash3" (func 11))
  (export "badWriteLog1" (func 12))
  (export "badWriteLog2" (func 13))
  (export "badWriteLog3" (func 14))
  (export "badWriteLog4" (func 15))
  (export "badBigIntStorageStore1" (func 16))
  (export "badRecursive" (func 18))
  (data (;0;) (i32.const 1024) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
  (data (;1;) (i32.const 1056) "test\00"))
