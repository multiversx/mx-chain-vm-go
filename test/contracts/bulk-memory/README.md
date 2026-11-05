# bulk-memory

Contract used to measure the gas cost of the bulk memory opcodes, `memory.copy` and `memory.fill`.

Every endpoint takes the number of bytes to copy/fill as argument 0. Since the rest of the
code path is identical no matter the argument, the gas difference between two calls is exactly
`(bytes1 - bytes2) * <opcode>PerByte`, which is what the tests in
`vmhost/hosttest/bulkMemoryGas_test.go` assert.

Built directly from WAT, without `wasm-opt`, so that the opcodes stay exactly as written:

```
./build.sh
```
