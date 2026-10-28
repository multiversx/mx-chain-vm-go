package executor

import "fmt"

// BreakpointValueNew is the reason the executor itself interrupted execution.
//
// Counterpart of the Rust `BreakpointValue` (new_traits/breakpoint_value_new.rs). The New
// suffix separates it from the legacy, six-valued vmhost.BreakpointValue, which it does not
// clash with in Go only because the two live in different packages; the suffix goes away with
// the old type.
//
// Unlike the legacy breakpoint, this only covers the interruptions raised by the executor's
// own middlewares. Everything the high-level VM signals travels as an [EarlyExit] instead,
// and reaches the VM through [InstanceCallResult].
//
// The values are the ones the metering and opcode control middlewares write into the
// instance's breakpoint global, and are the same numbers as the matching
// vmhost.BreakpointValue constants.
type BreakpointValueNew uint64

const (
	// BreakpointNoneNew means the lack of a breakpoint.
	BreakpointNoneNew BreakpointValueNew = 0

	// BreakpointOutOfGasNew means execution stopped because gas was exhausted.
	BreakpointOutOfGasNew BreakpointValueNew = 4

	// BreakpointMemoryLimitNew means execution stopped because of over-allocation of WASM
	// memory.
	BreakpointMemoryLimitNew BreakpointValueNew = 5
)

// String returns the human-readable name of a BreakpointValueNew.
func (breakpoint BreakpointValueNew) String() string {
	switch breakpoint {
	case BreakpointNoneNew:
		return "BreakpointNone"
	case BreakpointOutOfGasNew:
		return "BreakpointOutOfGas"
	case BreakpointMemoryLimitNew:
		return "BreakpointMemoryLimit"
	default:
		return "unknown breakpoint"
	}
}

// ParseBreakpointValueNew converts the raw contents of an instance's breakpoint global.
//
// The inverse of the constants above: a value the executor's middlewares never write is an
// error rather than a silently accepted breakpoint.
func ParseBreakpointValueNew(value uint64) (BreakpointValueNew, error) {
	switch BreakpointValueNew(value) {
	case BreakpointNoneNew, BreakpointOutOfGasNew, BreakpointMemoryLimitNew:
		return BreakpointValueNew(value), nil
	default:
		return BreakpointNoneNew, fmt.Errorf("unknown breakpoint value: %d", value)
	}
}
