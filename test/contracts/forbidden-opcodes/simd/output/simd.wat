;; SIMD instructions are not supported at the moment
;; Using any two SIMD instructions just to see the contract is invalid
(module
  (type $void (func))
  (type $finish(func (param i64)))
  (import "env" "int64finish" (func $int64finish (type $finish)))
  (func $main (type $void)
    (i32.const 42)
    (i8x16.splat)           ;; Create vector with identical lanes
    (i64x2.extract_lane 1)  ;; Extract lane 1 as a scalar
    (call $int64finish)
  )
  (memory $mem 1)
  (export "memory" (memory $mem))
  (export "main" (func $main))
)
