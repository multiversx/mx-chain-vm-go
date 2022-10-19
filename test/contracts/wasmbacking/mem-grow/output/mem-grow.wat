(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
	(func $main (type $void)
		i32.const 5
		memory.grow
		drop
    memory.size
    i64.extend_i32_u
		call $int64finish
	)
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
)
