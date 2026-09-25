# VM Hooks Code generator

The code generator generates boilerplate code for both this Go VM, and for [the executor repository](https://github.com/multiversx/mx-vm-executor-rs)

For it to automatically copy files there, create a file called `wasm-vm-executor-rs-path.txt` here, in the `cmd` folder, contianing your local path to that repository, on your disk.

Finally, simply run `go generate` in `vmhost/vmhooks`.

## Inputs

Everything under `input/` is a plain list of names, one per line. These lists are the only
place where opcodes and costs are still declared by hand.

### `wasm_opcodes.txt`

Every entry the gas schedule knows about, whether it is metered today or not. It drives:

- `executor/gasCostWASM.go`, the `WASMOpcodeCost` struct;
- `output/config.txt`, to be pasted into the `[WASMOpcodeCost]` section of `config/config.toml`;
- `output/FillGasMap_WASMOpcodeCosts.txt`, to be pasted into `config/gasSchedule.go`.

Being the gas schedule, it also holds names that are not WASM operators at all
(`LocalAllocate`, `LocalsUnmetered`, the `*PerByte` costs), as well as operators that no
opcode version allows. Dropping a name from here is a breaking change for every gas
schedule TOML.

### `wasmer2_allowed_opcodes_v1.txt`, `wasmer2_allowed_opcodes_v2.txt`

The operators allowed by each `OpcodeVersion`. The names must match the variants of
`wasmparser::Operator`, since they are emitted as match arms into:

- `opcode_whitelist.rs`, the contract validator in mx-sdk-rs;
- `wasmer_opcode_cost.rs` and `we_opcode_cost.rs`, the metering middlewares in the executor.

Each version is a standalone list rather than a delta over the previous one, because an
older version has to keep metering already deployed contracts exactly as it did before. A
name the target no longer understands is filtered out while writing: `Unwind` is dropped
from the whitelist entirely, and commented out in the Wasmer 6 middleware.

`v2` additionally drops 15 operators that `v1` still allows: legacy exception handling
(`Catch`, `CatchAll`, `Delegate`, `Rethrow`, `Throw`, `Try`), reference types (`RefFunc`,
`RefIsNull`, `RefNull`, `TypedSelect`), and table manipulation (`TableGet`, `TableGrow`,
`TableInit`, `TableSet`, `TableSize`). Each group sits behind a WASM proposal
(exception-handling, reference-types, bulk-memory/reference-types) that this SDK's
`wasm32-unknown-unknown` build never turns on, so no contract built with today's tooling
can contain them: shrinking the whitelist removes attack surface from the execution engine
at no functional cost. They stay in `v1` and in `wasm_opcodes.txt`, since `v1` still has to
meter contracts deployed before this change exactly as before, and the cost middlewares
keep a match arm for them regardless of version, as those are unaffected by the whitelist.
Tail calls (`ReturnCall`, `ReturnCallIndirect`) are equally unreachable by this toolchain
today but were deliberately kept allowed in both versions.

### `wasmer2_per_byte_opcodes.txt`

The subset of allowed opcodes whose cost also depends on how many bytes they end up
processing, that is, the bulk memory operators. Each name here gains a second cost,
`<Name>PerByte`, and its match arm becomes `Cost::BulkMemory { base, per_byte }` instead of
`Cost::Base`. The name itself has to appear in an allowed opcode list, and the `*PerByte`
name it implies has to appear in `wasm_opcodes.txt`.

### `wasmer2_additional_costs.txt`

Costs handed over to the executor that are not charged per opcode, so they can never show
up in an allowed opcode list. Currently only `LocalAllocate`, charged once for every local
variable a function declares, beyond the unmetered allowance.

### `wasmer2_opcodes.txt`

Unused, a leftover of the Wasmer 2 migration kept for reference. `wasm_opcodes.txt` is a
superset of it.

## The derived cost list

`OpcodeNames.AllCostNames()` concatenates the allowed opcodes of v1, then those of v2, then
the `*PerByte` costs, then the additional costs, skipping the names it has already seen.

That list is the field layout of the `OpcodeCost` struct on both sides of the FFI boundary,
`wasmer2/opcodeCost.go` here and `vm-executor/src/opcode_cost.rs` in the executor, and the
Go struct is cast straight to the C one. Reordering an input file, or inserting a name
anywhere but at the end of one, therefore changes that layout silently. Whenever it does
change, the executor has to be rebuilt and `wasmer2/libvmexeccapi.*` refreshed, or the VM
will read every cost past the first moved field from the wrong place.
