(module
  (type $void (func))
  (type $finish_type(func (param i32 i32)))
  (import "env" "finish" (func $finish (type $finish_type)))
	(func $main (type $void)
    (local $offset i32)
    i32.const 1024
    local.set $offset

    local.get $offset
    local.get $offset
    i32.load8_u
    i32.const 1 ;; makes 'n' into 'o'
    i32.add
    i32.store8 

    local.get $offset
    i32.const 1
    i32.add 
    local.get $offset
    i32.const 1
    i32.add 
    i32.load8_u
    i32.const 1 ;; makes 'j' into 'k'
    i32.add
    i32.store8

    local.get $offset
    i32.const 2
		call $finish
	)
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
  (data (i32.const 1024) "nj")
)
