(module
  (type (;0;) (func (param i32) (result i64)))
  (type (;1;) (func (param i64)))
  (type (;2;) (func))
  (import "env" "int64getArgument" (func (;0;) (type 0)))
  (import "env" "int64finish" (func (;1;) (type 1)))
  (func (;2;) (type 2)
    (local i64 i64)
    block  ;; label = @1
      i32.const 0
      call 0
      local.tee 0
      i64.const 1
      i64.lt_s
      br_if 0 (;@1;)
      i64.const 0
      local.set 1
      loop  ;; label = @2
				i32.const 1
				memory.grow
				i64.extend_i32_s
        call 1
        local.get 0
        local.get 1
        i64.const 1
        i64.add
        local.tee 1
        i64.ne
        br_if 0 (;@2;)
      end
    end)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66560))
  (export "memory" (memory 0))
  (export "memGrow" (func 2)))
