# VM Hooks Code generator

The code generator generates boilerplate code for both this Go VM, and for [the executor repository](https://github.com/multiversx/mx-vm-executor-rs)

For it to automatically copy files there, create a file called `wasm-vm-executor-rs-path.txt` here, in the `cmd` folder, contianing your local path to that repository, on your disk.

Finally, simply run `go generate` in `vmhost/vmhooks`.

## Inputs

Everything under `input/` is a plain list of names, one per line. These lists are the only
place where opcodes and costs are still declared by hand.

### `wasm_all_opcodes.txt`

The full list of Wasm opcode names, used to generate the opcode cost struct and its helpers.
It is a superset of every other list in this directory.

## The derived cost list

`OpcodeNames.AllCostNames()` concatenates the allowed opcodes of v1, then those of v2, then
the `*PerByte` costs, then the additional costs, skipping the names it has already seen.

That list is the field layout of the `OpcodeCost` struct on both sides of the FFI boundary,
`wasmer2/opcodeCost.go` here and `vm-executor/src/opcode_cost.rs` in the executor, and the
Go struct is cast straight to the C one. Reordering an input file, or inserting a name
anywhere but at the end of one, therefore changes that layout silently. Whenever it does
change, the executor has to be rebuilt and `wasmer2/libvmexeccapi.*` refreshed, or the VM
will read every cost past the first moved field from the wrong place.
