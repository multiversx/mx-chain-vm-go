package executor

// OpcodeVersion specifies which set of opcodes should be used by the VM.
//
// [OpcodeVersionV1] is the legacy opcode set that does not include support
// for bulk memory operations. [OpcodeVersionV2] extends it by adding bulk
// memory instructions, such as MemoryCopy and MemoryFill, while dropping
// these categories from the whitelist, since no contract built with the
// current SDK toolchain can emit any of them:
//
//   - legacy exception handling: Catch, CatchAll, Delegate, Rethrow, Throw, Try
//   - reference types: RefFunc, RefIsNull, RefNull, TypedSelect
//   - table manipulation: TableGet, TableGrow, TableInit, TableSet, TableSize
type OpcodeVersion uint32

const (
	// OpcodeVersionV1 is the legacy opcode set without bulk memory support.
	//
	// Use this for modules or environments that were compiled or designed
	// before bulk memory operations were introduced, or when compatibility
	// with older tooling is required.
	OpcodeVersionV1 OpcodeVersion = iota

	// OpcodeVersionV2 is the opcode set with bulk memory support.
	//
	// Use this for modules that rely on bulk memory operations like
	// MemoryCopy and MemoryFill, or when targeting newer runtimes that
	// support these instructions.
	//
	// Relative to OpcodeVersionV1, OpcodeVersionV2 also drops these
	// categories from the whitelist, since no contract built with today's
	// tooling can produce them:
	//
	//   - legacy exception handling: Catch, CatchAll, Delegate, Rethrow, Throw, Try
	//   - reference types: RefFunc, RefIsNull, RefNull, TypedSelect
	//   - table manipulation: TableGet, TableGrow, TableInit, TableSet, TableSize
	//
	// Note: OpcodeVersionV2 does not add support for memory.init and data.drop.
	OpcodeVersionV2
)
