(module
  (type (;0;) (func (result i32)))
  (type (;1;) (func (param i32 i32)))
  (type (;2;) (func (param i32 i32) (result i32)))
  (type (;3;) (func))
  (import "env" "getNumArguments" (func (;0;) (type 0)))
  (import "env" "finish" (func (;1;) (type 1)))
  (import "env" "getArgument" (func (;2;) (type 2)))
  (import "env" "signalError" (func (;3;) (type 1)))
  (func (;4;) (type 3)
    (local i32 i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 0
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1032
        i64.store offset=24
        local.get 0
        i32.const 0
        i64.load offset=1024
        i64.store offset=16
        local.get 0
        i32.const 16
        i32.add
        i32.const 15
        call 1
        br 1 (;@1;)
      end
      local.get 0
      i32.const 0
      i32.store8 offset=15
      i32.const 0
      local.get 0
      i32.const 15
      i32.add
      call 2
      drop
      block  ;; label = @2
        local.get 0
        i32.load8_u offset=15
        local.tee 1
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1032
        i64.store offset=24
        local.get 0
        i32.const 0
        i64.load offset=1024
        i64.store offset=16
        local.get 0
        i32.const 16
        i32.add
        i32.const 15
        call 1
        local.get 0
        i32.load8_u offset=15
        local.set 1
      end
      block  ;; label = @2
        local.get 1
        i32.const 255
        i32.and
        i32.const 1
        i32.ne
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1046 align=1
        i64.store offset=22 align=2
        local.get 0
        i32.const 0
        i64.load offset=1040 align=1
        i64.store offset=16
        local.get 0
        i32.const 16
        i32.add
        i32.const 13
        call 3
        local.get 0
        i32.load8_u offset=15
        local.set 1
      end
      local.get 1
      i32.const 255
      i32.and
      i32.const 2
      i32.ne
      br_if 0 (;@1;)
      local.get 0
      i32.const 16
      i32.add
      i32.const 4
      i32.add
      i32.const 0
      i32.load8_u offset=1058
      i32.store8
      local.get 0
      i32.const 0
      i32.load offset=1054 align=1
      i32.store offset=16
      loop  ;; label = @2
        local.get 0
        i32.const 16
        i32.add
        i32.const 4
        call 1
        br 0 (;@2;)
      end
    end
    local.get 0
    i32.const 32
    i32.add
    global.set 0)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66608))
  (export "memory" (memory 0))
  (export "init" (func 4))
  (data (;0;) (i32.const 1024) "init successful\00don't do this\00loop\00"))
