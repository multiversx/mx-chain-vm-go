;; Regression fixture for warm-instance reuse: after Instance::reset the
;; static-memory guard region must still be in place and the memory base
;; must not move. Build with: wat2wasm memory-reset.wat -o output/memory-reset.wasm
(module
  (memory (export "memory") 1 20)
  (data (i32.const 0) "\01\02\03\04")

  (func (export "init"))
  (func (export "callBack"))

  ;; Grow by three pages and write a marker into the last new page.
  (func (export "growAndDirty")
    (drop (memory.grow (i32.const 3)))
    (i32.store (i32.const 196608) (i32.const 0x5eed5eed)))

  ;; Always in bounds; a successful call that leaves the warm instance cached.
  (func (export "touchMinimum")
    (drop (i32.load (i32.const 0))))

  ;; Touch the page that growAndDirty added (offset 3 pages).
  (func (export "loadGrownPage")
    (drop (i32.load (i32.const 196608))))

  ;; Touch an address deep inside the 4 GiB static bound (3 GiB).
  (func (export "loadHigh")
    (drop (i32.load (i32.const 0xC0000000))))

  ;; Touch the last aligned word of the 32-bit address space.
  (func (export "loadLast")
    (drop (i32.load (i32.const 0xFFFFFFFC)))))
