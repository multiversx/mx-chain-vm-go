(module
  (type (;0;) (func (param i32) (result i64)))
  (type (;1;) (func (param i32 i32)))
  (type (;2;) (func (param i64)))
  (type (;3;) (func))
  (import "env" "int64getArgument" (func (;0;) (type 0)))
  (import "env" "signalError" (func (;1;) (type 1)))
  (import "env" "int64finish" (func (;2;) (type 2)))
  (func (;3;) (type 3)
    (local i32 i64)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      i32.const 0
      call 0
      i64.const -1
      i64.add
      local.tee 1
      i64.const 1
      i64.gt_u
      br_if 0 (;@1;)
      block  ;; label = @2
        block  ;; label = @3
          local.get 1
          i32.wrap_i64
          br_table 0 (;@3;) 1 (;@2;) 0 (;@3;)
        end
        local.get 0
        i32.const 24
        i32.add
        i32.const 0
        i32.load16_u offset=1064 align=1
        i32.store16
        local.get 0
        i32.const 0
        i64.load offset=1056 align=1
        i64.store offset=16
        local.get 0
        i32.const 16
        i32.add
        i32.const 9
        call 1
        local.get 0
        i32.const 0
        i32.load offset=1073 align=1
        i32.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1066 align=1
        i64.store
        local.get 0
        i32.const 10
        call 1
        br 1 (;@1;)
      end
      i32.const 0
      i32.const 42
      i32.store8 offset=2147484671
      i64.const 42
      call 2
    end
    i64.const 100
    call 2
    local.get 0
    i32.const 32
    i32.add
    global.set 0)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66624))
  (export "memory" (memory 0))
  (export "testFunc" (func 3))
  (data (;0;) (i32.const 1024) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
  (data (;1;) (i32.const 1056) "exit here\00exit later\00"))
