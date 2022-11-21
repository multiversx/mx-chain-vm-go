package wasmer

import "github.com/ElrondNetwork/wasm-vm/executor"

// OpcodeCount is the number of opcodes that we account for when setting gas costs.
const opcodeCount = 448

const (
	OpcodeUnreachable = iota
	OpcodeNop
	OpcodeBlock
	OpcodeLoop
	OpcodeIf
	OpcodeElse
	OpcodeEnd
	OpcodeBr
	OpcodeBrIf
	OpcodeBrTable
	OpcodeReturn
	OpcodeCall
	OpcodeCallIndirect
	OpcodeDrop
	OpcodeSelect
	OpcodeTypedSelect
	OpcodeLocalGet
	OpcodeLocalSet
	OpcodeLocalTee
	OpcodeGlobalGet
	OpcodeGlobalSet
	OpcodeI32Load
	OpcodeI64Load
	OpcodeF32Load
	OpcodeF64Load
	OpcodeI32Load8S
	OpcodeI32Load8U
	OpcodeI32Load16S
	OpcodeI32Load16U
	OpcodeI64Load8S
	OpcodeI64Load8U
	OpcodeI64Load16S
	OpcodeI64Load16U
	OpcodeI64Load32S
	OpcodeI64Load32U
	OpcodeI32Store
	OpcodeI64Store
	OpcodeF32Store
	OpcodeF64Store
	OpcodeI32Store8
	OpcodeI32Store16
	OpcodeI64Store8
	OpcodeI64Store16
	OpcodeI64Store32
	OpcodeMemorySize
	OpcodeMemoryGrow
	OpcodeI32Const
	OpcodeI64Const
	OpcodeF32Const
	OpcodeF64Const
	OpcodeRefNull
	OpcodeRefIsNull
	OpcodeRefFunc
	OpcodeI32Eqz
	OpcodeI32Eq
	OpcodeI32Ne
	OpcodeI32LtS
	OpcodeI32LtU
	OpcodeI32GtS
	OpcodeI32GtU
	OpcodeI32LeS
	OpcodeI32LeU
	OpcodeI32GeS
	OpcodeI32GeU
	OpcodeI64Eqz
	OpcodeI64Eq
	OpcodeI64Ne
	OpcodeI64LtS
	OpcodeI64LtU
	OpcodeI64GtS
	OpcodeI64GtU
	OpcodeI64LeS
	OpcodeI64LeU
	OpcodeI64GeS
	OpcodeI64GeU
	OpcodeF32Eq
	OpcodeF32Ne
	OpcodeF32Lt
	OpcodeF32Gt
	OpcodeF32Le
	OpcodeF32Ge
	OpcodeF64Eq
	OpcodeF64Ne
	OpcodeF64Lt
	OpcodeF64Gt
	OpcodeF64Le
	OpcodeF64Ge
	OpcodeI32Clz
	OpcodeI32Ctz
	OpcodeI32Popcnt
	OpcodeI32Add
	OpcodeI32Sub
	OpcodeI32Mul
	OpcodeI32DivS
	OpcodeI32DivU
	OpcodeI32RemS
	OpcodeI32RemU
	OpcodeI32And
	OpcodeI32Or
	OpcodeI32Xor
	OpcodeI32Shl
	OpcodeI32ShrS
	OpcodeI32ShrU
	OpcodeI32Rotl
	OpcodeI32Rotr
	OpcodeI64Clz
	OpcodeI64Ctz
	OpcodeI64Popcnt
	OpcodeI64Add
	OpcodeI64Sub
	OpcodeI64Mul
	OpcodeI64DivS
	OpcodeI64DivU
	OpcodeI64RemS
	OpcodeI64RemU
	OpcodeI64And
	OpcodeI64Or
	OpcodeI64Xor
	OpcodeI64Shl
	OpcodeI64ShrS
	OpcodeI64ShrU
	OpcodeI64Rotl
	OpcodeI64Rotr
	OpcodeF32Abs
	OpcodeF32Neg
	OpcodeF32Ceil
	OpcodeF32Floor
	OpcodeF32Trunc
	OpcodeF32Nearest
	OpcodeF32Sqrt
	OpcodeF32Add
	OpcodeF32Sub
	OpcodeF32Mul
	OpcodeF32Div
	OpcodeF32Min
	OpcodeF32Max
	OpcodeF32Copysign
	OpcodeF64Abs
	OpcodeF64Neg
	OpcodeF64Ceil
	OpcodeF64Floor
	OpcodeF64Trunc
	OpcodeF64Nearest
	OpcodeF64Sqrt
	OpcodeF64Add
	OpcodeF64Sub
	OpcodeF64Mul
	OpcodeF64Div
	OpcodeF64Min
	OpcodeF64Max
	OpcodeF64Copysign
	OpcodeI32WrapI64
	OpcodeI32TruncF32S
	OpcodeI32TruncF32U
	OpcodeI32TruncF64S
	OpcodeI32TruncF64U
	OpcodeI64ExtendI32S
	OpcodeI64ExtendI32U
	OpcodeI64TruncF32S
	OpcodeI64TruncF32U
	OpcodeI64TruncF64S
	OpcodeI64TruncF64U
	OpcodeF32ConvertI32S
	OpcodeF32ConvertI32U
	OpcodeF32ConvertI64S
	OpcodeF32ConvertI64U
	OpcodeF32DemoteF64
	OpcodeF64ConvertI32S
	OpcodeF64ConvertI32U
	OpcodeF64ConvertI64S
	OpcodeF64ConvertI64U
	OpcodeF64PromoteF32
	OpcodeI32ReinterpretF32
	OpcodeI64ReinterpretF64
	OpcodeF32ReinterpretI32
	OpcodeF64ReinterpretI64
	OpcodeI32Extend8S
	OpcodeI32Extend16S
	OpcodeI64Extend8S
	OpcodeI64Extend16S
	OpcodeI64Extend32S
	OpcodeI32TruncSatF32S
	OpcodeI32TruncSatF32U
	OpcodeI32TruncSatF64S
	OpcodeI32TruncSatF64U
	OpcodeI64TruncSatF32S
	OpcodeI64TruncSatF32U
	OpcodeI64TruncSatF64S
	OpcodeI64TruncSatF64U
	OpcodeMemoryInit
	OpcodeDataDrop
	OpcodeMemoryCopy
	OpcodeMemoryFill
	OpcodeTableInit
	OpcodeElemDrop
	OpcodeTableCopy
	OpcodeTableFill
	OpcodeTableGet
	OpcodeTableSet
	OpcodeTableGrow
	OpcodeTableSize
	OpcodeAtomicNotify
	OpcodeI32AtomicWait
	OpcodeI64AtomicWait
	OpcodeAtomicFence
	OpcodeI32AtomicLoad
	OpcodeI64AtomicLoad
	OpcodeI32AtomicLoad8U
	OpcodeI32AtomicLoad16U
	OpcodeI64AtomicLoad8U
	OpcodeI64AtomicLoad16U
	OpcodeI64AtomicLoad32U
	OpcodeI32AtomicStore
	OpcodeI64AtomicStore
	OpcodeI32AtomicStore8
	OpcodeI32AtomicStore16
	OpcodeI64AtomicStore8
	OpcodeI64AtomicStore16
	OpcodeI64AtomicStore32
	OpcodeI32AtomicRmwAdd
	OpcodeI64AtomicRmwAdd
	OpcodeI32AtomicRmw8AddU
	OpcodeI32AtomicRmw16AddU
	OpcodeI64AtomicRmw8AddU
	OpcodeI64AtomicRmw16AddU
	OpcodeI64AtomicRmw32AddU
	OpcodeI32AtomicRmwSub
	OpcodeI64AtomicRmwSub
	OpcodeI32AtomicRmw8SubU
	OpcodeI32AtomicRmw16SubU
	OpcodeI64AtomicRmw8SubU
	OpcodeI64AtomicRmw16SubU
	OpcodeI64AtomicRmw32SubU
	OpcodeI32AtomicRmwAnd
	OpcodeI64AtomicRmwAnd
	OpcodeI32AtomicRmw8AndU
	OpcodeI32AtomicRmw16AndU
	OpcodeI64AtomicRmw8AndU
	OpcodeI64AtomicRmw16AndU
	OpcodeI64AtomicRmw32AndU
	OpcodeI32AtomicRmwOr
	OpcodeI64AtomicRmwOr
	OpcodeI32AtomicRmw8OrU
	OpcodeI32AtomicRmw16OrU
	OpcodeI64AtomicRmw8OrU
	OpcodeI64AtomicRmw16OrU
	OpcodeI64AtomicRmw32OrU
	OpcodeI32AtomicRmwXor
	OpcodeI64AtomicRmwXor
	OpcodeI32AtomicRmw8XorU
	OpcodeI32AtomicRmw16XorU
	OpcodeI64AtomicRmw8XorU
	OpcodeI64AtomicRmw16XorU
	OpcodeI64AtomicRmw32XorU
	OpcodeI32AtomicRmwXchg
	OpcodeI64AtomicRmwXchg
	OpcodeI32AtomicRmw8XchgU
	OpcodeI32AtomicRmw16XchgU
	OpcodeI64AtomicRmw8XchgU
	OpcodeI64AtomicRmw16XchgU
	OpcodeI64AtomicRmw32XchgU
	OpcodeI32AtomicRmwCmpxchg
	OpcodeI64AtomicRmwCmpxchg
	OpcodeI32AtomicRmw8CmpxchgU
	OpcodeI32AtomicRmw16CmpxchgU
	OpcodeI64AtomicRmw8CmpxchgU
	OpcodeI64AtomicRmw16CmpxchgU
	OpcodeI64AtomicRmw32CmpxchgU
	OpcodeV128Load
	OpcodeV128Store
	OpcodeV128Const
	OpcodeI8x16Splat
	OpcodeI8x16ExtractLaneS
	OpcodeI8x16ExtractLaneU
	OpcodeI8x16ReplaceLane
	OpcodeI16x8Splat
	OpcodeI16x8ExtractLaneS
	OpcodeI16x8ExtractLaneU
	OpcodeI16x8ReplaceLane
	OpcodeI32x4Splat
	OpcodeI32x4ExtractLane
	OpcodeI32x4ReplaceLane
	OpcodeI64x2Splat
	OpcodeI64x2ExtractLane
	OpcodeI64x2ReplaceLane
	OpcodeF32x4Splat
	OpcodeF32x4ExtractLane
	OpcodeF32x4ReplaceLane
	OpcodeF64x2Splat
	OpcodeF64x2ExtractLane
	OpcodeF64x2ReplaceLane
	OpcodeI8x16Eq
	OpcodeI8x16Ne
	OpcodeI8x16LtS
	OpcodeI8x16LtU
	OpcodeI8x16GtS
	OpcodeI8x16GtU
	OpcodeI8x16LeS
	OpcodeI8x16LeU
	OpcodeI8x16GeS
	OpcodeI8x16GeU
	OpcodeI16x8Eq
	OpcodeI16x8Ne
	OpcodeI16x8LtS
	OpcodeI16x8LtU
	OpcodeI16x8GtS
	OpcodeI16x8GtU
	OpcodeI16x8LeS
	OpcodeI16x8LeU
	OpcodeI16x8GeS
	OpcodeI16x8GeU
	OpcodeI32x4Eq
	OpcodeI32x4Ne
	OpcodeI32x4LtS
	OpcodeI32x4LtU
	OpcodeI32x4GtS
	OpcodeI32x4GtU
	OpcodeI32x4LeS
	OpcodeI32x4LeU
	OpcodeI32x4GeS
	OpcodeI32x4GeU
	OpcodeF32x4Eq
	OpcodeF32x4Ne
	OpcodeF32x4Lt
	OpcodeF32x4Gt
	OpcodeF32x4Le
	OpcodeF32x4Ge
	OpcodeF64x2Eq
	OpcodeF64x2Ne
	OpcodeF64x2Lt
	OpcodeF64x2Gt
	OpcodeF64x2Le
	OpcodeF64x2Ge
	OpcodeV128Not
	OpcodeV128And
	OpcodeV128AndNot
	OpcodeV128Or
	OpcodeV128Xor
	OpcodeV128Bitselect
	OpcodeI8x16Neg
	OpcodeI8x16AnyTrue
	OpcodeI8x16AllTrue
	OpcodeI8x16Shl
	OpcodeI8x16ShrS
	OpcodeI8x16ShrU
	OpcodeI8x16Add
	OpcodeI8x16AddSaturateS
	OpcodeI8x16AddSaturateU
	OpcodeI8x16Sub
	OpcodeI8x16SubSaturateS
	OpcodeI8x16SubSaturateU
	OpcodeI8x16MinS
	OpcodeI8x16MinU
	OpcodeI8x16MaxS
	OpcodeI8x16MaxU
	OpcodeI8x16Mul
	OpcodeI16x8Neg
	OpcodeI16x8AnyTrue
	OpcodeI16x8AllTrue
	OpcodeI16x8Shl
	OpcodeI16x8ShrS
	OpcodeI16x8ShrU
	OpcodeI16x8Add
	OpcodeI16x8AddSaturateS
	OpcodeI16x8AddSaturateU
	OpcodeI16x8Sub
	OpcodeI16x8SubSaturateS
	OpcodeI16x8SubSaturateU
	OpcodeI16x8Mul
	OpcodeI16x8MinS
	OpcodeI16x8MinU
	OpcodeI16x8MaxS
	OpcodeI16x8MaxU
	OpcodeI32x4Neg
	OpcodeI32x4AnyTrue
	OpcodeI32x4AllTrue
	OpcodeI32x4Shl
	OpcodeI32x4ShrS
	OpcodeI32x4ShrU
	OpcodeI32x4Add
	OpcodeI32x4Sub
	OpcodeI32x4Mul
	OpcodeI32x4MinS
	OpcodeI32x4MinU
	OpcodeI32x4MaxS
	OpcodeI32x4MaxU
	OpcodeI64x2Neg
	OpcodeI64x2AnyTrue
	OpcodeI64x2AllTrue
	OpcodeI64x2Shl
	OpcodeI64x2ShrS
	OpcodeI64x2ShrU
	OpcodeI64x2Add
	OpcodeI64x2Sub
	OpcodeI64x2Mul
	OpcodeF32x4Abs
	OpcodeF32x4Neg
	OpcodeF32x4Sqrt
	OpcodeF32x4Add
	OpcodeF32x4Sub
	OpcodeF32x4Mul
	OpcodeF32x4Div
	OpcodeF32x4Min
	OpcodeF32x4Max
	OpcodeF64x2Abs
	OpcodeF64x2Neg
	OpcodeF64x2Sqrt
	OpcodeF64x2Add
	OpcodeF64x2Sub
	OpcodeF64x2Mul
	OpcodeF64x2Div
	OpcodeF64x2Min
	OpcodeF64x2Max
	OpcodeI32x4TruncSatF32x4S
	OpcodeI32x4TruncSatF32x4U
	OpcodeI64x2TruncSatF64x2S
	OpcodeI64x2TruncSatF64x2U
	OpcodeF32x4ConvertI32x4S
	OpcodeF32x4ConvertI32x4U
	OpcodeF64x2ConvertI64x2S
	OpcodeF64x2ConvertI64x2U
	OpcodeV8x16Swizzle
	OpcodeV8x16Shuffle
	OpcodeV8x16LoadSplat
	OpcodeV16x8LoadSplat
	OpcodeV32x4LoadSplat
	OpcodeV64x2LoadSplat
	OpcodeI8x16NarrowI16x8S
	OpcodeI8x16NarrowI16x8U
	OpcodeI16x8NarrowI32x4S
	OpcodeI16x8NarrowI32x4U
	OpcodeI16x8WidenLowI8x16S
	OpcodeI16x8WidenHighI8x16S
	OpcodeI16x8WidenLowI8x16U
	OpcodeI16x8WidenHighI8x16U
	OpcodeI32x4WidenLowI16x8S
	OpcodeI32x4WidenHighI16x8S
	OpcodeI32x4WidenLowI16x8U
	OpcodeI32x4WidenHighI16x8U
	OpcodeI16x8Load8x8S
	OpcodeI16x8Load8x8U
	OpcodeI32x4Load16x4S
	OpcodeI32x4Load16x4U
	OpcodeI64x2Load32x2S
	OpcodeI64x2Load32x2U
	OpcodeI8x16RoundingAverageU
	OpcodeI16x8RoundingAverageU
	OpcodeLocalAllocate
)

func toOpcodeCostsArray(opcode_costs_struct *executor.WASMOpcodeCost) [opcodeCount]uint32 {
	opcode_costs := [opcodeCount]uint32{}

	opcode_costs[OpcodeUnreachable] = opcode_costs_struct.Unreachable
	opcode_costs[OpcodeNop] = opcode_costs_struct.Nop
	opcode_costs[OpcodeBlock] = opcode_costs_struct.Block
	opcode_costs[OpcodeLoop] = opcode_costs_struct.Loop
	opcode_costs[OpcodeIf] = opcode_costs_struct.If
	opcode_costs[OpcodeElse] = opcode_costs_struct.Else
	opcode_costs[OpcodeEnd] = opcode_costs_struct.End
	opcode_costs[OpcodeBr] = opcode_costs_struct.Br
	opcode_costs[OpcodeBrIf] = opcode_costs_struct.BrIf
	opcode_costs[OpcodeBrTable] = opcode_costs_struct.BrTable
	opcode_costs[OpcodeReturn] = opcode_costs_struct.Return
	opcode_costs[OpcodeCall] = opcode_costs_struct.Call
	opcode_costs[OpcodeCallIndirect] = opcode_costs_struct.CallIndirect
	opcode_costs[OpcodeDrop] = opcode_costs_struct.Drop
	opcode_costs[OpcodeSelect] = opcode_costs_struct.Select
	opcode_costs[OpcodeTypedSelect] = opcode_costs_struct.TypedSelect
	opcode_costs[OpcodeLocalGet] = opcode_costs_struct.LocalGet
	opcode_costs[OpcodeLocalSet] = opcode_costs_struct.LocalSet
	opcode_costs[OpcodeLocalTee] = opcode_costs_struct.LocalTee
	opcode_costs[OpcodeGlobalGet] = opcode_costs_struct.GlobalGet
	opcode_costs[OpcodeGlobalSet] = opcode_costs_struct.GlobalSet
	opcode_costs[OpcodeI32Load] = opcode_costs_struct.I32Load
	opcode_costs[OpcodeI64Load] = opcode_costs_struct.I64Load
	opcode_costs[OpcodeF32Load] = opcode_costs_struct.F32Load
	opcode_costs[OpcodeF64Load] = opcode_costs_struct.F64Load
	opcode_costs[OpcodeI32Load8S] = opcode_costs_struct.I32Load8S
	opcode_costs[OpcodeI32Load8U] = opcode_costs_struct.I32Load8U
	opcode_costs[OpcodeI32Load16S] = opcode_costs_struct.I32Load16S
	opcode_costs[OpcodeI32Load16U] = opcode_costs_struct.I32Load16U
	opcode_costs[OpcodeI64Load8S] = opcode_costs_struct.I64Load8S
	opcode_costs[OpcodeI64Load8U] = opcode_costs_struct.I64Load8U
	opcode_costs[OpcodeI64Load16S] = opcode_costs_struct.I64Load16S
	opcode_costs[OpcodeI64Load16U] = opcode_costs_struct.I64Load16U
	opcode_costs[OpcodeI64Load32S] = opcode_costs_struct.I64Load32S
	opcode_costs[OpcodeI64Load32U] = opcode_costs_struct.I64Load32U
	opcode_costs[OpcodeI32Store] = opcode_costs_struct.I32Store
	opcode_costs[OpcodeI64Store] = opcode_costs_struct.I64Store
	opcode_costs[OpcodeF32Store] = opcode_costs_struct.F32Store
	opcode_costs[OpcodeF64Store] = opcode_costs_struct.F64Store
	opcode_costs[OpcodeI32Store8] = opcode_costs_struct.I32Store8
	opcode_costs[OpcodeI32Store16] = opcode_costs_struct.I32Store16
	opcode_costs[OpcodeI64Store8] = opcode_costs_struct.I64Store8
	opcode_costs[OpcodeI64Store16] = opcode_costs_struct.I64Store16
	opcode_costs[OpcodeI64Store32] = opcode_costs_struct.I64Store32
	opcode_costs[OpcodeMemorySize] = opcode_costs_struct.MemorySize
	opcode_costs[OpcodeMemoryGrow] = opcode_costs_struct.MemoryGrow
	opcode_costs[OpcodeI32Const] = opcode_costs_struct.I32Const
	opcode_costs[OpcodeI64Const] = opcode_costs_struct.I64Const
	opcode_costs[OpcodeF32Const] = opcode_costs_struct.F32Const
	opcode_costs[OpcodeF64Const] = opcode_costs_struct.F64Const
	opcode_costs[OpcodeRefNull] = opcode_costs_struct.RefNull
	opcode_costs[OpcodeRefIsNull] = opcode_costs_struct.RefIsNull
	opcode_costs[OpcodeRefFunc] = opcode_costs_struct.RefFunc
	opcode_costs[OpcodeI32Eqz] = opcode_costs_struct.I32Eqz
	opcode_costs[OpcodeI32Eq] = opcode_costs_struct.I32Eq
	opcode_costs[OpcodeI32Ne] = opcode_costs_struct.I32Ne
	opcode_costs[OpcodeI32LtS] = opcode_costs_struct.I32LtS
	opcode_costs[OpcodeI32LtU] = opcode_costs_struct.I32LtU
	opcode_costs[OpcodeI32GtS] = opcode_costs_struct.I32GtS
	opcode_costs[OpcodeI32GtU] = opcode_costs_struct.I32GtU
	opcode_costs[OpcodeI32LeS] = opcode_costs_struct.I32LeS
	opcode_costs[OpcodeI32LeU] = opcode_costs_struct.I32LeU
	opcode_costs[OpcodeI32GeS] = opcode_costs_struct.I32GeS
	opcode_costs[OpcodeI32GeU] = opcode_costs_struct.I32GeU
	opcode_costs[OpcodeI64Eqz] = opcode_costs_struct.I64Eqz
	opcode_costs[OpcodeI64Eq] = opcode_costs_struct.I64Eq
	opcode_costs[OpcodeI64Ne] = opcode_costs_struct.I64Ne
	opcode_costs[OpcodeI64LtS] = opcode_costs_struct.I64LtS
	opcode_costs[OpcodeI64LtU] = opcode_costs_struct.I64LtU
	opcode_costs[OpcodeI64GtS] = opcode_costs_struct.I64GtS
	opcode_costs[OpcodeI64GtU] = opcode_costs_struct.I64GtU
	opcode_costs[OpcodeI64LeS] = opcode_costs_struct.I64LeS
	opcode_costs[OpcodeI64LeU] = opcode_costs_struct.I64LeU
	opcode_costs[OpcodeI64GeS] = opcode_costs_struct.I64GeS
	opcode_costs[OpcodeI64GeU] = opcode_costs_struct.I64GeU
	opcode_costs[OpcodeF32Eq] = opcode_costs_struct.F32Eq
	opcode_costs[OpcodeF32Ne] = opcode_costs_struct.F32Ne
	opcode_costs[OpcodeF32Lt] = opcode_costs_struct.F32Lt
	opcode_costs[OpcodeF32Gt] = opcode_costs_struct.F32Gt
	opcode_costs[OpcodeF32Le] = opcode_costs_struct.F32Le
	opcode_costs[OpcodeF32Ge] = opcode_costs_struct.F32Ge
	opcode_costs[OpcodeF64Eq] = opcode_costs_struct.F64Eq
	opcode_costs[OpcodeF64Ne] = opcode_costs_struct.F64Ne
	opcode_costs[OpcodeF64Lt] = opcode_costs_struct.F64Lt
	opcode_costs[OpcodeF64Gt] = opcode_costs_struct.F64Gt
	opcode_costs[OpcodeF64Le] = opcode_costs_struct.F64Le
	opcode_costs[OpcodeF64Ge] = opcode_costs_struct.F64Ge
	opcode_costs[OpcodeI32Clz] = opcode_costs_struct.I32Clz
	opcode_costs[OpcodeI32Ctz] = opcode_costs_struct.I32Ctz
	opcode_costs[OpcodeI32Popcnt] = opcode_costs_struct.I32Popcnt
	opcode_costs[OpcodeI32Add] = opcode_costs_struct.I32Add
	opcode_costs[OpcodeI32Sub] = opcode_costs_struct.I32Sub
	opcode_costs[OpcodeI32Mul] = opcode_costs_struct.I32Mul
	opcode_costs[OpcodeI32DivS] = opcode_costs_struct.I32DivS
	opcode_costs[OpcodeI32DivU] = opcode_costs_struct.I32DivU
	opcode_costs[OpcodeI32RemS] = opcode_costs_struct.I32RemS
	opcode_costs[OpcodeI32RemU] = opcode_costs_struct.I32RemU
	opcode_costs[OpcodeI32And] = opcode_costs_struct.I32And
	opcode_costs[OpcodeI32Or] = opcode_costs_struct.I32Or
	opcode_costs[OpcodeI32Xor] = opcode_costs_struct.I32Xor
	opcode_costs[OpcodeI32Shl] = opcode_costs_struct.I32Shl
	opcode_costs[OpcodeI32ShrS] = opcode_costs_struct.I32ShrS
	opcode_costs[OpcodeI32ShrU] = opcode_costs_struct.I32ShrU
	opcode_costs[OpcodeI32Rotl] = opcode_costs_struct.I32Rotl
	opcode_costs[OpcodeI32Rotr] = opcode_costs_struct.I32Rotr
	opcode_costs[OpcodeI64Clz] = opcode_costs_struct.I64Clz
	opcode_costs[OpcodeI64Ctz] = opcode_costs_struct.I64Ctz
	opcode_costs[OpcodeI64Popcnt] = opcode_costs_struct.I64Popcnt
	opcode_costs[OpcodeI64Add] = opcode_costs_struct.I64Add
	opcode_costs[OpcodeI64Sub] = opcode_costs_struct.I64Sub
	opcode_costs[OpcodeI64Mul] = opcode_costs_struct.I64Mul
	opcode_costs[OpcodeI64DivS] = opcode_costs_struct.I64DivS
	opcode_costs[OpcodeI64DivU] = opcode_costs_struct.I64DivU
	opcode_costs[OpcodeI64RemS] = opcode_costs_struct.I64RemS
	opcode_costs[OpcodeI64RemU] = opcode_costs_struct.I64RemU
	opcode_costs[OpcodeI64And] = opcode_costs_struct.I64And
	opcode_costs[OpcodeI64Or] = opcode_costs_struct.I64Or
	opcode_costs[OpcodeI64Xor] = opcode_costs_struct.I64Xor
	opcode_costs[OpcodeI64Shl] = opcode_costs_struct.I64Shl
	opcode_costs[OpcodeI64ShrS] = opcode_costs_struct.I64ShrS
	opcode_costs[OpcodeI64ShrU] = opcode_costs_struct.I64ShrU
	opcode_costs[OpcodeI64Rotl] = opcode_costs_struct.I64Rotl
	opcode_costs[OpcodeI64Rotr] = opcode_costs_struct.I64Rotr
	opcode_costs[OpcodeF32Abs] = opcode_costs_struct.F32Abs
	opcode_costs[OpcodeF32Neg] = opcode_costs_struct.F32Neg
	opcode_costs[OpcodeF32Ceil] = opcode_costs_struct.F32Ceil
	opcode_costs[OpcodeF32Floor] = opcode_costs_struct.F32Floor
	opcode_costs[OpcodeF32Trunc] = opcode_costs_struct.F32Trunc
	opcode_costs[OpcodeF32Nearest] = opcode_costs_struct.F32Nearest
	opcode_costs[OpcodeF32Sqrt] = opcode_costs_struct.F32Sqrt
	opcode_costs[OpcodeF32Add] = opcode_costs_struct.F32Add
	opcode_costs[OpcodeF32Sub] = opcode_costs_struct.F32Sub
	opcode_costs[OpcodeF32Mul] = opcode_costs_struct.F32Mul
	opcode_costs[OpcodeF32Div] = opcode_costs_struct.F32Div
	opcode_costs[OpcodeF32Min] = opcode_costs_struct.F32Min
	opcode_costs[OpcodeF32Max] = opcode_costs_struct.F32Max
	opcode_costs[OpcodeF32Copysign] = opcode_costs_struct.F32Copysign
	opcode_costs[OpcodeF64Abs] = opcode_costs_struct.F64Abs
	opcode_costs[OpcodeF64Neg] = opcode_costs_struct.F64Neg
	opcode_costs[OpcodeF64Ceil] = opcode_costs_struct.F64Ceil
	opcode_costs[OpcodeF64Floor] = opcode_costs_struct.F64Floor
	opcode_costs[OpcodeF64Trunc] = opcode_costs_struct.F64Trunc
	opcode_costs[OpcodeF64Nearest] = opcode_costs_struct.F64Nearest
	opcode_costs[OpcodeF64Sqrt] = opcode_costs_struct.F64Sqrt
	opcode_costs[OpcodeF64Add] = opcode_costs_struct.F64Add
	opcode_costs[OpcodeF64Sub] = opcode_costs_struct.F64Sub
	opcode_costs[OpcodeF64Mul] = opcode_costs_struct.F64Mul
	opcode_costs[OpcodeF64Div] = opcode_costs_struct.F64Div
	opcode_costs[OpcodeF64Min] = opcode_costs_struct.F64Min
	opcode_costs[OpcodeF64Max] = opcode_costs_struct.F64Max
	opcode_costs[OpcodeF64Copysign] = opcode_costs_struct.F64Copysign
	opcode_costs[OpcodeI32WrapI64] = opcode_costs_struct.I32WrapI64
	opcode_costs[OpcodeI32TruncF32S] = opcode_costs_struct.I32TruncF32S
	opcode_costs[OpcodeI32TruncF32U] = opcode_costs_struct.I32TruncF32U
	opcode_costs[OpcodeI32TruncF64S] = opcode_costs_struct.I32TruncF64S
	opcode_costs[OpcodeI32TruncF64U] = opcode_costs_struct.I32TruncF64U
	opcode_costs[OpcodeI64ExtendI32S] = opcode_costs_struct.I64ExtendI32S
	opcode_costs[OpcodeI64ExtendI32U] = opcode_costs_struct.I64ExtendI32U
	opcode_costs[OpcodeI64TruncF32S] = opcode_costs_struct.I64TruncF32S
	opcode_costs[OpcodeI64TruncF32U] = opcode_costs_struct.I64TruncF32U
	opcode_costs[OpcodeI64TruncF64S] = opcode_costs_struct.I64TruncF64S
	opcode_costs[OpcodeI64TruncF64U] = opcode_costs_struct.I64TruncF64U
	opcode_costs[OpcodeF32ConvertI32S] = opcode_costs_struct.F32ConvertI32S
	opcode_costs[OpcodeF32ConvertI32U] = opcode_costs_struct.F32ConvertI32U
	opcode_costs[OpcodeF32ConvertI64S] = opcode_costs_struct.F32ConvertI64S
	opcode_costs[OpcodeF32ConvertI64U] = opcode_costs_struct.F32ConvertI64U
	opcode_costs[OpcodeF32DemoteF64] = opcode_costs_struct.F32DemoteF64
	opcode_costs[OpcodeF64ConvertI32S] = opcode_costs_struct.F64ConvertI32S
	opcode_costs[OpcodeF64ConvertI32U] = opcode_costs_struct.F64ConvertI32U
	opcode_costs[OpcodeF64ConvertI64S] = opcode_costs_struct.F64ConvertI64S
	opcode_costs[OpcodeF64ConvertI64U] = opcode_costs_struct.F64ConvertI64U
	opcode_costs[OpcodeF64PromoteF32] = opcode_costs_struct.F64PromoteF32
	opcode_costs[OpcodeI32ReinterpretF32] = opcode_costs_struct.I32ReinterpretF32
	opcode_costs[OpcodeI64ReinterpretF64] = opcode_costs_struct.I64ReinterpretF64
	opcode_costs[OpcodeF32ReinterpretI32] = opcode_costs_struct.F32ReinterpretI32
	opcode_costs[OpcodeF64ReinterpretI64] = opcode_costs_struct.F64ReinterpretI64
	opcode_costs[OpcodeI32Extend8S] = opcode_costs_struct.I32Extend8S
	opcode_costs[OpcodeI32Extend16S] = opcode_costs_struct.I32Extend16S
	opcode_costs[OpcodeI64Extend8S] = opcode_costs_struct.I64Extend8S
	opcode_costs[OpcodeI64Extend16S] = opcode_costs_struct.I64Extend16S
	opcode_costs[OpcodeI64Extend32S] = opcode_costs_struct.I64Extend32S
	opcode_costs[OpcodeI32TruncSatF32S] = opcode_costs_struct.I32TruncSatF32S
	opcode_costs[OpcodeI32TruncSatF32U] = opcode_costs_struct.I32TruncSatF32U
	opcode_costs[OpcodeI32TruncSatF64S] = opcode_costs_struct.I32TruncSatF64S
	opcode_costs[OpcodeI32TruncSatF64U] = opcode_costs_struct.I32TruncSatF64U
	opcode_costs[OpcodeI64TruncSatF32S] = opcode_costs_struct.I64TruncSatF32S
	opcode_costs[OpcodeI64TruncSatF32U] = opcode_costs_struct.I64TruncSatF32U
	opcode_costs[OpcodeI64TruncSatF64S] = opcode_costs_struct.I64TruncSatF64S
	opcode_costs[OpcodeI64TruncSatF64U] = opcode_costs_struct.I64TruncSatF64U
	opcode_costs[OpcodeMemoryInit] = opcode_costs_struct.MemoryInit
	opcode_costs[OpcodeDataDrop] = opcode_costs_struct.DataDrop
	opcode_costs[OpcodeMemoryCopy] = opcode_costs_struct.MemoryCopy
	opcode_costs[OpcodeMemoryFill] = opcode_costs_struct.MemoryFill
	opcode_costs[OpcodeTableInit] = opcode_costs_struct.TableInit
	opcode_costs[OpcodeElemDrop] = opcode_costs_struct.ElemDrop
	opcode_costs[OpcodeTableCopy] = opcode_costs_struct.TableCopy
	opcode_costs[OpcodeTableFill] = opcode_costs_struct.TableFill
	opcode_costs[OpcodeTableGet] = opcode_costs_struct.TableGet
	opcode_costs[OpcodeTableSet] = opcode_costs_struct.TableSet
	opcode_costs[OpcodeTableGrow] = opcode_costs_struct.TableGrow
	opcode_costs[OpcodeTableSize] = opcode_costs_struct.TableSize
	opcode_costs[OpcodeAtomicNotify] = opcode_costs_struct.AtomicNotify
	opcode_costs[OpcodeI32AtomicWait] = opcode_costs_struct.I32AtomicWait
	opcode_costs[OpcodeI64AtomicWait] = opcode_costs_struct.I64AtomicWait
	opcode_costs[OpcodeAtomicFence] = opcode_costs_struct.AtomicFence
	opcode_costs[OpcodeI32AtomicLoad] = opcode_costs_struct.I32AtomicLoad
	opcode_costs[OpcodeI64AtomicLoad] = opcode_costs_struct.I64AtomicLoad
	opcode_costs[OpcodeI32AtomicLoad8U] = opcode_costs_struct.I32AtomicLoad8U
	opcode_costs[OpcodeI32AtomicLoad16U] = opcode_costs_struct.I32AtomicLoad16U
	opcode_costs[OpcodeI64AtomicLoad8U] = opcode_costs_struct.I64AtomicLoad8U
	opcode_costs[OpcodeI64AtomicLoad16U] = opcode_costs_struct.I64AtomicLoad16U
	opcode_costs[OpcodeI64AtomicLoad32U] = opcode_costs_struct.I64AtomicLoad32U
	opcode_costs[OpcodeI32AtomicStore] = opcode_costs_struct.I32AtomicStore
	opcode_costs[OpcodeI64AtomicStore] = opcode_costs_struct.I64AtomicStore
	opcode_costs[OpcodeI32AtomicStore8] = opcode_costs_struct.I32AtomicStore8
	opcode_costs[OpcodeI32AtomicStore16] = opcode_costs_struct.I32AtomicStore16
	opcode_costs[OpcodeI64AtomicStore8] = opcode_costs_struct.I64AtomicStore8
	opcode_costs[OpcodeI64AtomicStore16] = opcode_costs_struct.I64AtomicStore16
	opcode_costs[OpcodeI64AtomicStore32] = opcode_costs_struct.I64AtomicStore32
	opcode_costs[OpcodeI32AtomicRmwAdd] = opcode_costs_struct.I32AtomicRmwAdd
	opcode_costs[OpcodeI64AtomicRmwAdd] = opcode_costs_struct.I64AtomicRmwAdd
	opcode_costs[OpcodeI32AtomicRmw8AddU] = opcode_costs_struct.I32AtomicRmw8AddU
	opcode_costs[OpcodeI32AtomicRmw16AddU] = opcode_costs_struct.I32AtomicRmw16AddU
	opcode_costs[OpcodeI64AtomicRmw8AddU] = opcode_costs_struct.I64AtomicRmw8AddU
	opcode_costs[OpcodeI64AtomicRmw16AddU] = opcode_costs_struct.I64AtomicRmw16AddU
	opcode_costs[OpcodeI64AtomicRmw32AddU] = opcode_costs_struct.I64AtomicRmw32AddU
	opcode_costs[OpcodeI32AtomicRmwSub] = opcode_costs_struct.I32AtomicRmwSub
	opcode_costs[OpcodeI64AtomicRmwSub] = opcode_costs_struct.I64AtomicRmwSub
	opcode_costs[OpcodeI32AtomicRmw8SubU] = opcode_costs_struct.I32AtomicRmw8SubU
	opcode_costs[OpcodeI32AtomicRmw16SubU] = opcode_costs_struct.I32AtomicRmw16SubU
	opcode_costs[OpcodeI64AtomicRmw8SubU] = opcode_costs_struct.I64AtomicRmw8SubU
	opcode_costs[OpcodeI64AtomicRmw16SubU] = opcode_costs_struct.I64AtomicRmw16SubU
	opcode_costs[OpcodeI64AtomicRmw32SubU] = opcode_costs_struct.I64AtomicRmw32SubU
	opcode_costs[OpcodeI32AtomicRmwAnd] = opcode_costs_struct.I32AtomicRmwAnd
	opcode_costs[OpcodeI64AtomicRmwAnd] = opcode_costs_struct.I64AtomicRmwAnd
	opcode_costs[OpcodeI32AtomicRmw8AndU] = opcode_costs_struct.I32AtomicRmw8AndU
	opcode_costs[OpcodeI32AtomicRmw16AndU] = opcode_costs_struct.I32AtomicRmw16AndU
	opcode_costs[OpcodeI64AtomicRmw8AndU] = opcode_costs_struct.I64AtomicRmw8AndU
	opcode_costs[OpcodeI64AtomicRmw16AndU] = opcode_costs_struct.I64AtomicRmw16AndU
	opcode_costs[OpcodeI64AtomicRmw32AndU] = opcode_costs_struct.I64AtomicRmw32AndU
	opcode_costs[OpcodeI32AtomicRmwOr] = opcode_costs_struct.I32AtomicRmwOr
	opcode_costs[OpcodeI64AtomicRmwOr] = opcode_costs_struct.I64AtomicRmwOr
	opcode_costs[OpcodeI32AtomicRmw8OrU] = opcode_costs_struct.I32AtomicRmw8OrU
	opcode_costs[OpcodeI32AtomicRmw16OrU] = opcode_costs_struct.I32AtomicRmw16OrU
	opcode_costs[OpcodeI64AtomicRmw8OrU] = opcode_costs_struct.I64AtomicRmw8OrU
	opcode_costs[OpcodeI64AtomicRmw16OrU] = opcode_costs_struct.I64AtomicRmw16OrU
	opcode_costs[OpcodeI64AtomicRmw32OrU] = opcode_costs_struct.I64AtomicRmw32OrU
	opcode_costs[OpcodeI32AtomicRmwXor] = opcode_costs_struct.I32AtomicRmwXor
	opcode_costs[OpcodeI64AtomicRmwXor] = opcode_costs_struct.I64AtomicRmwXor
	opcode_costs[OpcodeI32AtomicRmw8XorU] = opcode_costs_struct.I32AtomicRmw8XorU
	opcode_costs[OpcodeI32AtomicRmw16XorU] = opcode_costs_struct.I32AtomicRmw16XorU
	opcode_costs[OpcodeI64AtomicRmw8XorU] = opcode_costs_struct.I64AtomicRmw8XorU
	opcode_costs[OpcodeI64AtomicRmw16XorU] = opcode_costs_struct.I64AtomicRmw16XorU
	opcode_costs[OpcodeI64AtomicRmw32XorU] = opcode_costs_struct.I64AtomicRmw32XorU
	opcode_costs[OpcodeI32AtomicRmwXchg] = opcode_costs_struct.I32AtomicRmwXchg
	opcode_costs[OpcodeI64AtomicRmwXchg] = opcode_costs_struct.I64AtomicRmwXchg
	opcode_costs[OpcodeI32AtomicRmw8XchgU] = opcode_costs_struct.I32AtomicRmw8XchgU
	opcode_costs[OpcodeI32AtomicRmw16XchgU] = opcode_costs_struct.I32AtomicRmw16XchgU
	opcode_costs[OpcodeI64AtomicRmw8XchgU] = opcode_costs_struct.I64AtomicRmw8XchgU
	opcode_costs[OpcodeI64AtomicRmw16XchgU] = opcode_costs_struct.I64AtomicRmw16XchgU
	opcode_costs[OpcodeI64AtomicRmw32XchgU] = opcode_costs_struct.I64AtomicRmw32XchgU
	opcode_costs[OpcodeI32AtomicRmwCmpxchg] = opcode_costs_struct.I32AtomicRmwCmpxchg
	opcode_costs[OpcodeI64AtomicRmwCmpxchg] = opcode_costs_struct.I64AtomicRmwCmpxchg
	opcode_costs[OpcodeI32AtomicRmw8CmpxchgU] = opcode_costs_struct.I32AtomicRmw8CmpxchgU
	opcode_costs[OpcodeI32AtomicRmw16CmpxchgU] = opcode_costs_struct.I32AtomicRmw16CmpxchgU
	opcode_costs[OpcodeI64AtomicRmw8CmpxchgU] = opcode_costs_struct.I64AtomicRmw8CmpxchgU
	opcode_costs[OpcodeI64AtomicRmw16CmpxchgU] = opcode_costs_struct.I64AtomicRmw16CmpxchgU
	opcode_costs[OpcodeI64AtomicRmw32CmpxchgU] = opcode_costs_struct.I64AtomicRmw32CmpxchgU
	opcode_costs[OpcodeV128Load] = opcode_costs_struct.V128Load
	opcode_costs[OpcodeV128Store] = opcode_costs_struct.V128Store
	opcode_costs[OpcodeV128Const] = opcode_costs_struct.V128Const
	opcode_costs[OpcodeI8x16Splat] = opcode_costs_struct.I8x16Splat
	opcode_costs[OpcodeI8x16ExtractLaneS] = opcode_costs_struct.I8x16ExtractLaneS
	opcode_costs[OpcodeI8x16ExtractLaneU] = opcode_costs_struct.I8x16ExtractLaneU
	opcode_costs[OpcodeI8x16ReplaceLane] = opcode_costs_struct.I8x16ReplaceLane
	opcode_costs[OpcodeI16x8Splat] = opcode_costs_struct.I16x8Splat
	opcode_costs[OpcodeI16x8ExtractLaneS] = opcode_costs_struct.I16x8ExtractLaneS
	opcode_costs[OpcodeI16x8ExtractLaneU] = opcode_costs_struct.I16x8ExtractLaneU
	opcode_costs[OpcodeI16x8ReplaceLane] = opcode_costs_struct.I16x8ReplaceLane
	opcode_costs[OpcodeI32x4Splat] = opcode_costs_struct.I32x4Splat
	opcode_costs[OpcodeI32x4ExtractLane] = opcode_costs_struct.I32x4ExtractLane
	opcode_costs[OpcodeI32x4ReplaceLane] = opcode_costs_struct.I32x4ReplaceLane
	opcode_costs[OpcodeI64x2Splat] = opcode_costs_struct.I64x2Splat
	opcode_costs[OpcodeI64x2ExtractLane] = opcode_costs_struct.I64x2ExtractLane
	opcode_costs[OpcodeI64x2ReplaceLane] = opcode_costs_struct.I64x2ReplaceLane
	opcode_costs[OpcodeF32x4Splat] = opcode_costs_struct.F32x4Splat
	opcode_costs[OpcodeF32x4ExtractLane] = opcode_costs_struct.F32x4ExtractLane
	opcode_costs[OpcodeF32x4ReplaceLane] = opcode_costs_struct.F32x4ReplaceLane
	opcode_costs[OpcodeF64x2Splat] = opcode_costs_struct.F64x2Splat
	opcode_costs[OpcodeF64x2ExtractLane] = opcode_costs_struct.F64x2ExtractLane
	opcode_costs[OpcodeF64x2ReplaceLane] = opcode_costs_struct.F64x2ReplaceLane
	opcode_costs[OpcodeI8x16Eq] = opcode_costs_struct.I8x16Eq
	opcode_costs[OpcodeI8x16Ne] = opcode_costs_struct.I8x16Ne
	opcode_costs[OpcodeI8x16LtS] = opcode_costs_struct.I8x16LtS
	opcode_costs[OpcodeI8x16LtU] = opcode_costs_struct.I8x16LtU
	opcode_costs[OpcodeI8x16GtS] = opcode_costs_struct.I8x16GtS
	opcode_costs[OpcodeI8x16GtU] = opcode_costs_struct.I8x16GtU
	opcode_costs[OpcodeI8x16LeS] = opcode_costs_struct.I8x16LeS
	opcode_costs[OpcodeI8x16LeU] = opcode_costs_struct.I8x16LeU
	opcode_costs[OpcodeI8x16GeS] = opcode_costs_struct.I8x16GeS
	opcode_costs[OpcodeI8x16GeU] = opcode_costs_struct.I8x16GeU
	opcode_costs[OpcodeI16x8Eq] = opcode_costs_struct.I16x8Eq
	opcode_costs[OpcodeI16x8Ne] = opcode_costs_struct.I16x8Ne
	opcode_costs[OpcodeI16x8LtS] = opcode_costs_struct.I16x8LtS
	opcode_costs[OpcodeI16x8LtU] = opcode_costs_struct.I16x8LtU
	opcode_costs[OpcodeI16x8GtS] = opcode_costs_struct.I16x8GtS
	opcode_costs[OpcodeI16x8GtU] = opcode_costs_struct.I16x8GtU
	opcode_costs[OpcodeI16x8LeS] = opcode_costs_struct.I16x8LeS
	opcode_costs[OpcodeI16x8LeU] = opcode_costs_struct.I16x8LeU
	opcode_costs[OpcodeI16x8GeS] = opcode_costs_struct.I16x8GeS
	opcode_costs[OpcodeI16x8GeU] = opcode_costs_struct.I16x8GeU
	opcode_costs[OpcodeI32x4Eq] = opcode_costs_struct.I32x4Eq
	opcode_costs[OpcodeI32x4Ne] = opcode_costs_struct.I32x4Ne
	opcode_costs[OpcodeI32x4LtS] = opcode_costs_struct.I32x4LtS
	opcode_costs[OpcodeI32x4LtU] = opcode_costs_struct.I32x4LtU
	opcode_costs[OpcodeI32x4GtS] = opcode_costs_struct.I32x4GtS
	opcode_costs[OpcodeI32x4GtU] = opcode_costs_struct.I32x4GtU
	opcode_costs[OpcodeI32x4LeS] = opcode_costs_struct.I32x4LeS
	opcode_costs[OpcodeI32x4LeU] = opcode_costs_struct.I32x4LeU
	opcode_costs[OpcodeI32x4GeS] = opcode_costs_struct.I32x4GeS
	opcode_costs[OpcodeI32x4GeU] = opcode_costs_struct.I32x4GeU
	opcode_costs[OpcodeF32x4Eq] = opcode_costs_struct.F32x4Eq
	opcode_costs[OpcodeF32x4Ne] = opcode_costs_struct.F32x4Ne
	opcode_costs[OpcodeF32x4Lt] = opcode_costs_struct.F32x4Lt
	opcode_costs[OpcodeF32x4Gt] = opcode_costs_struct.F32x4Gt
	opcode_costs[OpcodeF32x4Le] = opcode_costs_struct.F32x4Le
	opcode_costs[OpcodeF32x4Ge] = opcode_costs_struct.F32x4Ge
	opcode_costs[OpcodeF64x2Eq] = opcode_costs_struct.F64x2Eq
	opcode_costs[OpcodeF64x2Ne] = opcode_costs_struct.F64x2Ne
	opcode_costs[OpcodeF64x2Lt] = opcode_costs_struct.F64x2Lt
	opcode_costs[OpcodeF64x2Gt] = opcode_costs_struct.F64x2Gt
	opcode_costs[OpcodeF64x2Le] = opcode_costs_struct.F64x2Le
	opcode_costs[OpcodeF64x2Ge] = opcode_costs_struct.F64x2Ge
	opcode_costs[OpcodeV128Not] = opcode_costs_struct.V128Not
	opcode_costs[OpcodeV128And] = opcode_costs_struct.V128And
	opcode_costs[OpcodeV128AndNot] = opcode_costs_struct.V128AndNot
	opcode_costs[OpcodeV128Or] = opcode_costs_struct.V128Or
	opcode_costs[OpcodeV128Xor] = opcode_costs_struct.V128Xor
	opcode_costs[OpcodeV128Bitselect] = opcode_costs_struct.V128Bitselect
	opcode_costs[OpcodeI8x16Neg] = opcode_costs_struct.I8x16Neg
	opcode_costs[OpcodeI8x16AnyTrue] = opcode_costs_struct.I8x16AnyTrue
	opcode_costs[OpcodeI8x16AllTrue] = opcode_costs_struct.I8x16AllTrue
	opcode_costs[OpcodeI8x16Shl] = opcode_costs_struct.I8x16Shl
	opcode_costs[OpcodeI8x16ShrS] = opcode_costs_struct.I8x16ShrS
	opcode_costs[OpcodeI8x16ShrU] = opcode_costs_struct.I8x16ShrU
	opcode_costs[OpcodeI8x16Add] = opcode_costs_struct.I8x16Add
	opcode_costs[OpcodeI8x16AddSaturateS] = opcode_costs_struct.I8x16AddSaturateS
	opcode_costs[OpcodeI8x16AddSaturateU] = opcode_costs_struct.I8x16AddSaturateU
	opcode_costs[OpcodeI8x16Sub] = opcode_costs_struct.I8x16Sub
	opcode_costs[OpcodeI8x16SubSaturateS] = opcode_costs_struct.I8x16SubSaturateS
	opcode_costs[OpcodeI8x16SubSaturateU] = opcode_costs_struct.I8x16SubSaturateU
	opcode_costs[OpcodeI8x16MinS] = opcode_costs_struct.I8x16MinS
	opcode_costs[OpcodeI8x16MinU] = opcode_costs_struct.I8x16MinU
	opcode_costs[OpcodeI8x16MaxS] = opcode_costs_struct.I8x16MaxS
	opcode_costs[OpcodeI8x16MaxU] = opcode_costs_struct.I8x16MaxU
	opcode_costs[OpcodeI8x16Mul] = opcode_costs_struct.I8x16Mul
	opcode_costs[OpcodeI16x8Neg] = opcode_costs_struct.I16x8Neg
	opcode_costs[OpcodeI16x8AnyTrue] = opcode_costs_struct.I16x8AnyTrue
	opcode_costs[OpcodeI16x8AllTrue] = opcode_costs_struct.I16x8AllTrue
	opcode_costs[OpcodeI16x8Shl] = opcode_costs_struct.I16x8Shl
	opcode_costs[OpcodeI16x8ShrS] = opcode_costs_struct.I16x8ShrS
	opcode_costs[OpcodeI16x8ShrU] = opcode_costs_struct.I16x8ShrU
	opcode_costs[OpcodeI16x8Add] = opcode_costs_struct.I16x8Add
	opcode_costs[OpcodeI16x8AddSaturateS] = opcode_costs_struct.I16x8AddSaturateS
	opcode_costs[OpcodeI16x8AddSaturateU] = opcode_costs_struct.I16x8AddSaturateU
	opcode_costs[OpcodeI16x8Sub] = opcode_costs_struct.I16x8Sub
	opcode_costs[OpcodeI16x8SubSaturateS] = opcode_costs_struct.I16x8SubSaturateS
	opcode_costs[OpcodeI16x8SubSaturateU] = opcode_costs_struct.I16x8SubSaturateU
	opcode_costs[OpcodeI16x8Mul] = opcode_costs_struct.I16x8Mul
	opcode_costs[OpcodeI16x8MinS] = opcode_costs_struct.I16x8MinS
	opcode_costs[OpcodeI16x8MinU] = opcode_costs_struct.I16x8MinU
	opcode_costs[OpcodeI16x8MaxS] = opcode_costs_struct.I16x8MaxS
	opcode_costs[OpcodeI16x8MaxU] = opcode_costs_struct.I16x8MaxU
	opcode_costs[OpcodeI32x4Neg] = opcode_costs_struct.I32x4Neg
	opcode_costs[OpcodeI32x4AnyTrue] = opcode_costs_struct.I32x4AnyTrue
	opcode_costs[OpcodeI32x4AllTrue] = opcode_costs_struct.I32x4AllTrue
	opcode_costs[OpcodeI32x4Shl] = opcode_costs_struct.I32x4Shl
	opcode_costs[OpcodeI32x4ShrS] = opcode_costs_struct.I32x4ShrS
	opcode_costs[OpcodeI32x4ShrU] = opcode_costs_struct.I32x4ShrU
	opcode_costs[OpcodeI32x4Add] = opcode_costs_struct.I32x4Add
	opcode_costs[OpcodeI32x4Sub] = opcode_costs_struct.I32x4Sub
	opcode_costs[OpcodeI32x4Mul] = opcode_costs_struct.I32x4Mul
	opcode_costs[OpcodeI32x4MinS] = opcode_costs_struct.I32x4MinS
	opcode_costs[OpcodeI32x4MinU] = opcode_costs_struct.I32x4MinU
	opcode_costs[OpcodeI32x4MaxS] = opcode_costs_struct.I32x4MaxS
	opcode_costs[OpcodeI32x4MaxU] = opcode_costs_struct.I32x4MaxU
	opcode_costs[OpcodeI64x2Neg] = opcode_costs_struct.I64x2Neg
	opcode_costs[OpcodeI64x2AnyTrue] = opcode_costs_struct.I64x2AnyTrue
	opcode_costs[OpcodeI64x2AllTrue] = opcode_costs_struct.I64x2AllTrue
	opcode_costs[OpcodeI64x2Shl] = opcode_costs_struct.I64x2Shl
	opcode_costs[OpcodeI64x2ShrS] = opcode_costs_struct.I64x2ShrS
	opcode_costs[OpcodeI64x2ShrU] = opcode_costs_struct.I64x2ShrU
	opcode_costs[OpcodeI64x2Add] = opcode_costs_struct.I64x2Add
	opcode_costs[OpcodeI64x2Sub] = opcode_costs_struct.I64x2Sub
	opcode_costs[OpcodeI64x2Mul] = opcode_costs_struct.I64x2Mul
	opcode_costs[OpcodeF32x4Abs] = opcode_costs_struct.F32x4Abs
	opcode_costs[OpcodeF32x4Neg] = opcode_costs_struct.F32x4Neg
	opcode_costs[OpcodeF32x4Sqrt] = opcode_costs_struct.F32x4Sqrt
	opcode_costs[OpcodeF32x4Add] = opcode_costs_struct.F32x4Add
	opcode_costs[OpcodeF32x4Sub] = opcode_costs_struct.F32x4Sub
	opcode_costs[OpcodeF32x4Mul] = opcode_costs_struct.F32x4Mul
	opcode_costs[OpcodeF32x4Div] = opcode_costs_struct.F32x4Div
	opcode_costs[OpcodeF32x4Min] = opcode_costs_struct.F32x4Min
	opcode_costs[OpcodeF32x4Max] = opcode_costs_struct.F32x4Max
	opcode_costs[OpcodeF64x2Abs] = opcode_costs_struct.F64x2Abs
	opcode_costs[OpcodeF64x2Neg] = opcode_costs_struct.F64x2Neg
	opcode_costs[OpcodeF64x2Sqrt] = opcode_costs_struct.F64x2Sqrt
	opcode_costs[OpcodeF64x2Add] = opcode_costs_struct.F64x2Add
	opcode_costs[OpcodeF64x2Sub] = opcode_costs_struct.F64x2Sub
	opcode_costs[OpcodeF64x2Mul] = opcode_costs_struct.F64x2Mul
	opcode_costs[OpcodeF64x2Div] = opcode_costs_struct.F64x2Div
	opcode_costs[OpcodeF64x2Min] = opcode_costs_struct.F64x2Min
	opcode_costs[OpcodeF64x2Max] = opcode_costs_struct.F64x2Max
	opcode_costs[OpcodeI32x4TruncSatF32x4S] = opcode_costs_struct.I32x4TruncSatF32x4S
	opcode_costs[OpcodeI32x4TruncSatF32x4U] = opcode_costs_struct.I32x4TruncSatF32x4U
	opcode_costs[OpcodeI64x2TruncSatF64x2S] = opcode_costs_struct.I64x2TruncSatF64x2S
	opcode_costs[OpcodeI64x2TruncSatF64x2U] = opcode_costs_struct.I64x2TruncSatF64x2U
	opcode_costs[OpcodeF32x4ConvertI32x4S] = opcode_costs_struct.F32x4ConvertI32x4S
	opcode_costs[OpcodeF32x4ConvertI32x4U] = opcode_costs_struct.F32x4ConvertI32x4U
	opcode_costs[OpcodeF64x2ConvertI64x2S] = opcode_costs_struct.F64x2ConvertI64x2S
	opcode_costs[OpcodeF64x2ConvertI64x2U] = opcode_costs_struct.F64x2ConvertI64x2U
	opcode_costs[OpcodeV8x16Swizzle] = opcode_costs_struct.V8x16Swizzle
	opcode_costs[OpcodeV8x16Shuffle] = opcode_costs_struct.V8x16Shuffle
	opcode_costs[OpcodeV8x16LoadSplat] = opcode_costs_struct.V8x16LoadSplat
	opcode_costs[OpcodeV16x8LoadSplat] = opcode_costs_struct.V16x8LoadSplat
	opcode_costs[OpcodeV32x4LoadSplat] = opcode_costs_struct.V32x4LoadSplat
	opcode_costs[OpcodeV64x2LoadSplat] = opcode_costs_struct.V64x2LoadSplat
	opcode_costs[OpcodeI8x16NarrowI16x8S] = opcode_costs_struct.I8x16NarrowI16x8S
	opcode_costs[OpcodeI8x16NarrowI16x8U] = opcode_costs_struct.I8x16NarrowI16x8U
	opcode_costs[OpcodeI16x8NarrowI32x4S] = opcode_costs_struct.I16x8NarrowI32x4S
	opcode_costs[OpcodeI16x8NarrowI32x4U] = opcode_costs_struct.I16x8NarrowI32x4U
	opcode_costs[OpcodeI16x8WidenLowI8x16S] = opcode_costs_struct.I16x8WidenLowI8x16S
	opcode_costs[OpcodeI16x8WidenHighI8x16S] = opcode_costs_struct.I16x8WidenHighI8x16S
	opcode_costs[OpcodeI16x8WidenLowI8x16U] = opcode_costs_struct.I16x8WidenLowI8x16U
	opcode_costs[OpcodeI16x8WidenHighI8x16U] = opcode_costs_struct.I16x8WidenHighI8x16U
	opcode_costs[OpcodeI32x4WidenLowI16x8S] = opcode_costs_struct.I32x4WidenLowI16x8S
	opcode_costs[OpcodeI32x4WidenHighI16x8S] = opcode_costs_struct.I32x4WidenHighI16x8S
	opcode_costs[OpcodeI32x4WidenLowI16x8U] = opcode_costs_struct.I32x4WidenLowI16x8U
	opcode_costs[OpcodeI32x4WidenHighI16x8U] = opcode_costs_struct.I32x4WidenHighI16x8U
	opcode_costs[OpcodeI16x8Load8x8S] = opcode_costs_struct.I16x8Load8x8S
	opcode_costs[OpcodeI16x8Load8x8U] = opcode_costs_struct.I16x8Load8x8U
	opcode_costs[OpcodeI32x4Load16x4S] = opcode_costs_struct.I32x4Load16x4S
	opcode_costs[OpcodeI32x4Load16x4U] = opcode_costs_struct.I32x4Load16x4U
	opcode_costs[OpcodeI64x2Load32x2S] = opcode_costs_struct.I64x2Load32x2S
	opcode_costs[OpcodeI64x2Load32x2U] = opcode_costs_struct.I64x2Load32x2U
	opcode_costs[OpcodeI8x16RoundingAverageU] = opcode_costs_struct.I8x16RoundingAverageU
	opcode_costs[OpcodeI16x8RoundingAverageU] = opcode_costs_struct.I16x8RoundingAverageU
	opcode_costs[OpcodeLocalAllocate] = opcode_costs_struct.LocalAllocate
	// LocalsUnmetered, MaxMemoryGrow and MaxMemoryGrowDelta are not added to the
	// opcode_costs array; the values will be sent to Wasmer as compilation
	// options instead

	return opcode_costs
}
