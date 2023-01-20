package wasmer

import "github.com/ElrondNetwork/wasm-vm/executor"

// OpcodeCount is the number of opcodes that we account for when setting gas costs.
const opcodeCount = 448

// opcodes list
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

func toOpcodeCostsArray(opcodeCostsStruct *executor.WASMOpcodeCost) [opcodeCount]uint32 {
	opcodeCosts := [opcodeCount]uint32{}

	opcodeCosts[OpcodeUnreachable] = opcodeCostsStruct.Unreachable
	opcodeCosts[OpcodeNop] = opcodeCostsStruct.Nop
	opcodeCosts[OpcodeBlock] = opcodeCostsStruct.Block
	opcodeCosts[OpcodeLoop] = opcodeCostsStruct.Loop
	opcodeCosts[OpcodeIf] = opcodeCostsStruct.If
	opcodeCosts[OpcodeElse] = opcodeCostsStruct.Else
	opcodeCosts[OpcodeEnd] = opcodeCostsStruct.End
	opcodeCosts[OpcodeBr] = opcodeCostsStruct.Br
	opcodeCosts[OpcodeBrIf] = opcodeCostsStruct.BrIf
	opcodeCosts[OpcodeBrTable] = opcodeCostsStruct.BrTable
	opcodeCosts[OpcodeReturn] = opcodeCostsStruct.Return
	opcodeCosts[OpcodeCall] = opcodeCostsStruct.Call
	opcodeCosts[OpcodeCallIndirect] = opcodeCostsStruct.CallIndirect
	opcodeCosts[OpcodeDrop] = opcodeCostsStruct.Drop
	opcodeCosts[OpcodeSelect] = opcodeCostsStruct.Select
	opcodeCosts[OpcodeTypedSelect] = opcodeCostsStruct.TypedSelect
	opcodeCosts[OpcodeLocalGet] = opcodeCostsStruct.LocalGet
	opcodeCosts[OpcodeLocalSet] = opcodeCostsStruct.LocalSet
	opcodeCosts[OpcodeLocalTee] = opcodeCostsStruct.LocalTee
	opcodeCosts[OpcodeGlobalGet] = opcodeCostsStruct.GlobalGet
	opcodeCosts[OpcodeGlobalSet] = opcodeCostsStruct.GlobalSet
	opcodeCosts[OpcodeI32Load] = opcodeCostsStruct.I32Load
	opcodeCosts[OpcodeI64Load] = opcodeCostsStruct.I64Load
	opcodeCosts[OpcodeF32Load] = opcodeCostsStruct.F32Load
	opcodeCosts[OpcodeF64Load] = opcodeCostsStruct.F64Load
	opcodeCosts[OpcodeI32Load8S] = opcodeCostsStruct.I32Load8S
	opcodeCosts[OpcodeI32Load8U] = opcodeCostsStruct.I32Load8U
	opcodeCosts[OpcodeI32Load16S] = opcodeCostsStruct.I32Load16S
	opcodeCosts[OpcodeI32Load16U] = opcodeCostsStruct.I32Load16U
	opcodeCosts[OpcodeI64Load8S] = opcodeCostsStruct.I64Load8S
	opcodeCosts[OpcodeI64Load8U] = opcodeCostsStruct.I64Load8U
	opcodeCosts[OpcodeI64Load16S] = opcodeCostsStruct.I64Load16S
	opcodeCosts[OpcodeI64Load16U] = opcodeCostsStruct.I64Load16U
	opcodeCosts[OpcodeI64Load32S] = opcodeCostsStruct.I64Load32S
	opcodeCosts[OpcodeI64Load32U] = opcodeCostsStruct.I64Load32U
	opcodeCosts[OpcodeI32Store] = opcodeCostsStruct.I32Store
	opcodeCosts[OpcodeI64Store] = opcodeCostsStruct.I64Store
	opcodeCosts[OpcodeF32Store] = opcodeCostsStruct.F32Store
	opcodeCosts[OpcodeF64Store] = opcodeCostsStruct.F64Store
	opcodeCosts[OpcodeI32Store8] = opcodeCostsStruct.I32Store8
	opcodeCosts[OpcodeI32Store16] = opcodeCostsStruct.I32Store16
	opcodeCosts[OpcodeI64Store8] = opcodeCostsStruct.I64Store8
	opcodeCosts[OpcodeI64Store16] = opcodeCostsStruct.I64Store16
	opcodeCosts[OpcodeI64Store32] = opcodeCostsStruct.I64Store32
	opcodeCosts[OpcodeMemorySize] = opcodeCostsStruct.MemorySize
	opcodeCosts[OpcodeMemoryGrow] = opcodeCostsStruct.MemoryGrow
	opcodeCosts[OpcodeI32Const] = opcodeCostsStruct.I32Const
	opcodeCosts[OpcodeI64Const] = opcodeCostsStruct.I64Const
	opcodeCosts[OpcodeF32Const] = opcodeCostsStruct.F32Const
	opcodeCosts[OpcodeF64Const] = opcodeCostsStruct.F64Const
	opcodeCosts[OpcodeRefNull] = opcodeCostsStruct.RefNull
	opcodeCosts[OpcodeRefIsNull] = opcodeCostsStruct.RefIsNull
	opcodeCosts[OpcodeRefFunc] = opcodeCostsStruct.RefFunc
	opcodeCosts[OpcodeI32Eqz] = opcodeCostsStruct.I32Eqz
	opcodeCosts[OpcodeI32Eq] = opcodeCostsStruct.I32Eq
	opcodeCosts[OpcodeI32Ne] = opcodeCostsStruct.I32Ne
	opcodeCosts[OpcodeI32LtS] = opcodeCostsStruct.I32LtS
	opcodeCosts[OpcodeI32LtU] = opcodeCostsStruct.I32LtU
	opcodeCosts[OpcodeI32GtS] = opcodeCostsStruct.I32GtS
	opcodeCosts[OpcodeI32GtU] = opcodeCostsStruct.I32GtU
	opcodeCosts[OpcodeI32LeS] = opcodeCostsStruct.I32LeS
	opcodeCosts[OpcodeI32LeU] = opcodeCostsStruct.I32LeU
	opcodeCosts[OpcodeI32GeS] = opcodeCostsStruct.I32GeS
	opcodeCosts[OpcodeI32GeU] = opcodeCostsStruct.I32GeU
	opcodeCosts[OpcodeI64Eqz] = opcodeCostsStruct.I64Eqz
	opcodeCosts[OpcodeI64Eq] = opcodeCostsStruct.I64Eq
	opcodeCosts[OpcodeI64Ne] = opcodeCostsStruct.I64Ne
	opcodeCosts[OpcodeI64LtS] = opcodeCostsStruct.I64LtS
	opcodeCosts[OpcodeI64LtU] = opcodeCostsStruct.I64LtU
	opcodeCosts[OpcodeI64GtS] = opcodeCostsStruct.I64GtS
	opcodeCosts[OpcodeI64GtU] = opcodeCostsStruct.I64GtU
	opcodeCosts[OpcodeI64LeS] = opcodeCostsStruct.I64LeS
	opcodeCosts[OpcodeI64LeU] = opcodeCostsStruct.I64LeU
	opcodeCosts[OpcodeI64GeS] = opcodeCostsStruct.I64GeS
	opcodeCosts[OpcodeI64GeU] = opcodeCostsStruct.I64GeU
	opcodeCosts[OpcodeF32Eq] = opcodeCostsStruct.F32Eq
	opcodeCosts[OpcodeF32Ne] = opcodeCostsStruct.F32Ne
	opcodeCosts[OpcodeF32Lt] = opcodeCostsStruct.F32Lt
	opcodeCosts[OpcodeF32Gt] = opcodeCostsStruct.F32Gt
	opcodeCosts[OpcodeF32Le] = opcodeCostsStruct.F32Le
	opcodeCosts[OpcodeF32Ge] = opcodeCostsStruct.F32Ge
	opcodeCosts[OpcodeF64Eq] = opcodeCostsStruct.F64Eq
	opcodeCosts[OpcodeF64Ne] = opcodeCostsStruct.F64Ne
	opcodeCosts[OpcodeF64Lt] = opcodeCostsStruct.F64Lt
	opcodeCosts[OpcodeF64Gt] = opcodeCostsStruct.F64Gt
	opcodeCosts[OpcodeF64Le] = opcodeCostsStruct.F64Le
	opcodeCosts[OpcodeF64Ge] = opcodeCostsStruct.F64Ge
	opcodeCosts[OpcodeI32Clz] = opcodeCostsStruct.I32Clz
	opcodeCosts[OpcodeI32Ctz] = opcodeCostsStruct.I32Ctz
	opcodeCosts[OpcodeI32Popcnt] = opcodeCostsStruct.I32Popcnt
	opcodeCosts[OpcodeI32Add] = opcodeCostsStruct.I32Add
	opcodeCosts[OpcodeI32Sub] = opcodeCostsStruct.I32Sub
	opcodeCosts[OpcodeI32Mul] = opcodeCostsStruct.I32Mul
	opcodeCosts[OpcodeI32DivS] = opcodeCostsStruct.I32DivS
	opcodeCosts[OpcodeI32DivU] = opcodeCostsStruct.I32DivU
	opcodeCosts[OpcodeI32RemS] = opcodeCostsStruct.I32RemS
	opcodeCosts[OpcodeI32RemU] = opcodeCostsStruct.I32RemU
	opcodeCosts[OpcodeI32And] = opcodeCostsStruct.I32And
	opcodeCosts[OpcodeI32Or] = opcodeCostsStruct.I32Or
	opcodeCosts[OpcodeI32Xor] = opcodeCostsStruct.I32Xor
	opcodeCosts[OpcodeI32Shl] = opcodeCostsStruct.I32Shl
	opcodeCosts[OpcodeI32ShrS] = opcodeCostsStruct.I32ShrS
	opcodeCosts[OpcodeI32ShrU] = opcodeCostsStruct.I32ShrU
	opcodeCosts[OpcodeI32Rotl] = opcodeCostsStruct.I32Rotl
	opcodeCosts[OpcodeI32Rotr] = opcodeCostsStruct.I32Rotr
	opcodeCosts[OpcodeI64Clz] = opcodeCostsStruct.I64Clz
	opcodeCosts[OpcodeI64Ctz] = opcodeCostsStruct.I64Ctz
	opcodeCosts[OpcodeI64Popcnt] = opcodeCostsStruct.I64Popcnt
	opcodeCosts[OpcodeI64Add] = opcodeCostsStruct.I64Add
	opcodeCosts[OpcodeI64Sub] = opcodeCostsStruct.I64Sub
	opcodeCosts[OpcodeI64Mul] = opcodeCostsStruct.I64Mul
	opcodeCosts[OpcodeI64DivS] = opcodeCostsStruct.I64DivS
	opcodeCosts[OpcodeI64DivU] = opcodeCostsStruct.I64DivU
	opcodeCosts[OpcodeI64RemS] = opcodeCostsStruct.I64RemS
	opcodeCosts[OpcodeI64RemU] = opcodeCostsStruct.I64RemU
	opcodeCosts[OpcodeI64And] = opcodeCostsStruct.I64And
	opcodeCosts[OpcodeI64Or] = opcodeCostsStruct.I64Or
	opcodeCosts[OpcodeI64Xor] = opcodeCostsStruct.I64Xor
	opcodeCosts[OpcodeI64Shl] = opcodeCostsStruct.I64Shl
	opcodeCosts[OpcodeI64ShrS] = opcodeCostsStruct.I64ShrS
	opcodeCosts[OpcodeI64ShrU] = opcodeCostsStruct.I64ShrU
	opcodeCosts[OpcodeI64Rotl] = opcodeCostsStruct.I64Rotl
	opcodeCosts[OpcodeI64Rotr] = opcodeCostsStruct.I64Rotr
	opcodeCosts[OpcodeF32Abs] = opcodeCostsStruct.F32Abs
	opcodeCosts[OpcodeF32Neg] = opcodeCostsStruct.F32Neg
	opcodeCosts[OpcodeF32Ceil] = opcodeCostsStruct.F32Ceil
	opcodeCosts[OpcodeF32Floor] = opcodeCostsStruct.F32Floor
	opcodeCosts[OpcodeF32Trunc] = opcodeCostsStruct.F32Trunc
	opcodeCosts[OpcodeF32Nearest] = opcodeCostsStruct.F32Nearest
	opcodeCosts[OpcodeF32Sqrt] = opcodeCostsStruct.F32Sqrt
	opcodeCosts[OpcodeF32Add] = opcodeCostsStruct.F32Add
	opcodeCosts[OpcodeF32Sub] = opcodeCostsStruct.F32Sub
	opcodeCosts[OpcodeF32Mul] = opcodeCostsStruct.F32Mul
	opcodeCosts[OpcodeF32Div] = opcodeCostsStruct.F32Div
	opcodeCosts[OpcodeF32Min] = opcodeCostsStruct.F32Min
	opcodeCosts[OpcodeF32Max] = opcodeCostsStruct.F32Max
	opcodeCosts[OpcodeF32Copysign] = opcodeCostsStruct.F32Copysign
	opcodeCosts[OpcodeF64Abs] = opcodeCostsStruct.F64Abs
	opcodeCosts[OpcodeF64Neg] = opcodeCostsStruct.F64Neg
	opcodeCosts[OpcodeF64Ceil] = opcodeCostsStruct.F64Ceil
	opcodeCosts[OpcodeF64Floor] = opcodeCostsStruct.F64Floor
	opcodeCosts[OpcodeF64Trunc] = opcodeCostsStruct.F64Trunc
	opcodeCosts[OpcodeF64Nearest] = opcodeCostsStruct.F64Nearest
	opcodeCosts[OpcodeF64Sqrt] = opcodeCostsStruct.F64Sqrt
	opcodeCosts[OpcodeF64Add] = opcodeCostsStruct.F64Add
	opcodeCosts[OpcodeF64Sub] = opcodeCostsStruct.F64Sub
	opcodeCosts[OpcodeF64Mul] = opcodeCostsStruct.F64Mul
	opcodeCosts[OpcodeF64Div] = opcodeCostsStruct.F64Div
	opcodeCosts[OpcodeF64Min] = opcodeCostsStruct.F64Min
	opcodeCosts[OpcodeF64Max] = opcodeCostsStruct.F64Max
	opcodeCosts[OpcodeF64Copysign] = opcodeCostsStruct.F64Copysign
	opcodeCosts[OpcodeI32WrapI64] = opcodeCostsStruct.I32WrapI64
	opcodeCosts[OpcodeI32TruncF32S] = opcodeCostsStruct.I32TruncF32S
	opcodeCosts[OpcodeI32TruncF32U] = opcodeCostsStruct.I32TruncF32U
	opcodeCosts[OpcodeI32TruncF64S] = opcodeCostsStruct.I32TruncF64S
	opcodeCosts[OpcodeI32TruncF64U] = opcodeCostsStruct.I32TruncF64U
	opcodeCosts[OpcodeI64ExtendI32S] = opcodeCostsStruct.I64ExtendI32S
	opcodeCosts[OpcodeI64ExtendI32U] = opcodeCostsStruct.I64ExtendI32U
	opcodeCosts[OpcodeI64TruncF32S] = opcodeCostsStruct.I64TruncF32S
	opcodeCosts[OpcodeI64TruncF32U] = opcodeCostsStruct.I64TruncF32U
	opcodeCosts[OpcodeI64TruncF64S] = opcodeCostsStruct.I64TruncF64S
	opcodeCosts[OpcodeI64TruncF64U] = opcodeCostsStruct.I64TruncF64U
	opcodeCosts[OpcodeF32ConvertI32S] = opcodeCostsStruct.F32ConvertI32S
	opcodeCosts[OpcodeF32ConvertI32U] = opcodeCostsStruct.F32ConvertI32U
	opcodeCosts[OpcodeF32ConvertI64S] = opcodeCostsStruct.F32ConvertI64S
	opcodeCosts[OpcodeF32ConvertI64U] = opcodeCostsStruct.F32ConvertI64U
	opcodeCosts[OpcodeF32DemoteF64] = opcodeCostsStruct.F32DemoteF64
	opcodeCosts[OpcodeF64ConvertI32S] = opcodeCostsStruct.F64ConvertI32S
	opcodeCosts[OpcodeF64ConvertI32U] = opcodeCostsStruct.F64ConvertI32U
	opcodeCosts[OpcodeF64ConvertI64S] = opcodeCostsStruct.F64ConvertI64S
	opcodeCosts[OpcodeF64ConvertI64U] = opcodeCostsStruct.F64ConvertI64U
	opcodeCosts[OpcodeF64PromoteF32] = opcodeCostsStruct.F64PromoteF32
	opcodeCosts[OpcodeI32ReinterpretF32] = opcodeCostsStruct.I32ReinterpretF32
	opcodeCosts[OpcodeI64ReinterpretF64] = opcodeCostsStruct.I64ReinterpretF64
	opcodeCosts[OpcodeF32ReinterpretI32] = opcodeCostsStruct.F32ReinterpretI32
	opcodeCosts[OpcodeF64ReinterpretI64] = opcodeCostsStruct.F64ReinterpretI64
	opcodeCosts[OpcodeI32Extend8S] = opcodeCostsStruct.I32Extend8S
	opcodeCosts[OpcodeI32Extend16S] = opcodeCostsStruct.I32Extend16S
	opcodeCosts[OpcodeI64Extend8S] = opcodeCostsStruct.I64Extend8S
	opcodeCosts[OpcodeI64Extend16S] = opcodeCostsStruct.I64Extend16S
	opcodeCosts[OpcodeI64Extend32S] = opcodeCostsStruct.I64Extend32S
	opcodeCosts[OpcodeI32TruncSatF32S] = opcodeCostsStruct.I32TruncSatF32S
	opcodeCosts[OpcodeI32TruncSatF32U] = opcodeCostsStruct.I32TruncSatF32U
	opcodeCosts[OpcodeI32TruncSatF64S] = opcodeCostsStruct.I32TruncSatF64S
	opcodeCosts[OpcodeI32TruncSatF64U] = opcodeCostsStruct.I32TruncSatF64U
	opcodeCosts[OpcodeI64TruncSatF32S] = opcodeCostsStruct.I64TruncSatF32S
	opcodeCosts[OpcodeI64TruncSatF32U] = opcodeCostsStruct.I64TruncSatF32U
	opcodeCosts[OpcodeI64TruncSatF64S] = opcodeCostsStruct.I64TruncSatF64S
	opcodeCosts[OpcodeI64TruncSatF64U] = opcodeCostsStruct.I64TruncSatF64U
	opcodeCosts[OpcodeMemoryInit] = opcodeCostsStruct.MemoryInit
	opcodeCosts[OpcodeDataDrop] = opcodeCostsStruct.DataDrop
	opcodeCosts[OpcodeMemoryCopy] = opcodeCostsStruct.MemoryCopy
	opcodeCosts[OpcodeMemoryFill] = opcodeCostsStruct.MemoryFill
	opcodeCosts[OpcodeTableInit] = opcodeCostsStruct.TableInit
	opcodeCosts[OpcodeElemDrop] = opcodeCostsStruct.ElemDrop
	opcodeCosts[OpcodeTableCopy] = opcodeCostsStruct.TableCopy
	opcodeCosts[OpcodeTableFill] = opcodeCostsStruct.TableFill
	opcodeCosts[OpcodeTableGet] = opcodeCostsStruct.TableGet
	opcodeCosts[OpcodeTableSet] = opcodeCostsStruct.TableSet
	opcodeCosts[OpcodeTableGrow] = opcodeCostsStruct.TableGrow
	opcodeCosts[OpcodeTableSize] = opcodeCostsStruct.TableSize
	opcodeCosts[OpcodeAtomicNotify] = opcodeCostsStruct.AtomicNotify
	opcodeCosts[OpcodeI32AtomicWait] = opcodeCostsStruct.I32AtomicWait
	opcodeCosts[OpcodeI64AtomicWait] = opcodeCostsStruct.I64AtomicWait
	opcodeCosts[OpcodeAtomicFence] = opcodeCostsStruct.AtomicFence
	opcodeCosts[OpcodeI32AtomicLoad] = opcodeCostsStruct.I32AtomicLoad
	opcodeCosts[OpcodeI64AtomicLoad] = opcodeCostsStruct.I64AtomicLoad
	opcodeCosts[OpcodeI32AtomicLoad8U] = opcodeCostsStruct.I32AtomicLoad8U
	opcodeCosts[OpcodeI32AtomicLoad16U] = opcodeCostsStruct.I32AtomicLoad16U
	opcodeCosts[OpcodeI64AtomicLoad8U] = opcodeCostsStruct.I64AtomicLoad8U
	opcodeCosts[OpcodeI64AtomicLoad16U] = opcodeCostsStruct.I64AtomicLoad16U
	opcodeCosts[OpcodeI64AtomicLoad32U] = opcodeCostsStruct.I64AtomicLoad32U
	opcodeCosts[OpcodeI32AtomicStore] = opcodeCostsStruct.I32AtomicStore
	opcodeCosts[OpcodeI64AtomicStore] = opcodeCostsStruct.I64AtomicStore
	opcodeCosts[OpcodeI32AtomicStore8] = opcodeCostsStruct.I32AtomicStore8
	opcodeCosts[OpcodeI32AtomicStore16] = opcodeCostsStruct.I32AtomicStore16
	opcodeCosts[OpcodeI64AtomicStore8] = opcodeCostsStruct.I64AtomicStore8
	opcodeCosts[OpcodeI64AtomicStore16] = opcodeCostsStruct.I64AtomicStore16
	opcodeCosts[OpcodeI64AtomicStore32] = opcodeCostsStruct.I64AtomicStore32
	opcodeCosts[OpcodeI32AtomicRmwAdd] = opcodeCostsStruct.I32AtomicRmwAdd
	opcodeCosts[OpcodeI64AtomicRmwAdd] = opcodeCostsStruct.I64AtomicRmwAdd
	opcodeCosts[OpcodeI32AtomicRmw8AddU] = opcodeCostsStruct.I32AtomicRmw8AddU
	opcodeCosts[OpcodeI32AtomicRmw16AddU] = opcodeCostsStruct.I32AtomicRmw16AddU
	opcodeCosts[OpcodeI64AtomicRmw8AddU] = opcodeCostsStruct.I64AtomicRmw8AddU
	opcodeCosts[OpcodeI64AtomicRmw16AddU] = opcodeCostsStruct.I64AtomicRmw16AddU
	opcodeCosts[OpcodeI64AtomicRmw32AddU] = opcodeCostsStruct.I64AtomicRmw32AddU
	opcodeCosts[OpcodeI32AtomicRmwSub] = opcodeCostsStruct.I32AtomicRmwSub
	opcodeCosts[OpcodeI64AtomicRmwSub] = opcodeCostsStruct.I64AtomicRmwSub
	opcodeCosts[OpcodeI32AtomicRmw8SubU] = opcodeCostsStruct.I32AtomicRmw8SubU
	opcodeCosts[OpcodeI32AtomicRmw16SubU] = opcodeCostsStruct.I32AtomicRmw16SubU
	opcodeCosts[OpcodeI64AtomicRmw8SubU] = opcodeCostsStruct.I64AtomicRmw8SubU
	opcodeCosts[OpcodeI64AtomicRmw16SubU] = opcodeCostsStruct.I64AtomicRmw16SubU
	opcodeCosts[OpcodeI64AtomicRmw32SubU] = opcodeCostsStruct.I64AtomicRmw32SubU
	opcodeCosts[OpcodeI32AtomicRmwAnd] = opcodeCostsStruct.I32AtomicRmwAnd
	opcodeCosts[OpcodeI64AtomicRmwAnd] = opcodeCostsStruct.I64AtomicRmwAnd
	opcodeCosts[OpcodeI32AtomicRmw8AndU] = opcodeCostsStruct.I32AtomicRmw8AndU
	opcodeCosts[OpcodeI32AtomicRmw16AndU] = opcodeCostsStruct.I32AtomicRmw16AndU
	opcodeCosts[OpcodeI64AtomicRmw8AndU] = opcodeCostsStruct.I64AtomicRmw8AndU
	opcodeCosts[OpcodeI64AtomicRmw16AndU] = opcodeCostsStruct.I64AtomicRmw16AndU
	opcodeCosts[OpcodeI64AtomicRmw32AndU] = opcodeCostsStruct.I64AtomicRmw32AndU
	opcodeCosts[OpcodeI32AtomicRmwOr] = opcodeCostsStruct.I32AtomicRmwOr
	opcodeCosts[OpcodeI64AtomicRmwOr] = opcodeCostsStruct.I64AtomicRmwOr
	opcodeCosts[OpcodeI32AtomicRmw8OrU] = opcodeCostsStruct.I32AtomicRmw8OrU
	opcodeCosts[OpcodeI32AtomicRmw16OrU] = opcodeCostsStruct.I32AtomicRmw16OrU
	opcodeCosts[OpcodeI64AtomicRmw8OrU] = opcodeCostsStruct.I64AtomicRmw8OrU
	opcodeCosts[OpcodeI64AtomicRmw16OrU] = opcodeCostsStruct.I64AtomicRmw16OrU
	opcodeCosts[OpcodeI64AtomicRmw32OrU] = opcodeCostsStruct.I64AtomicRmw32OrU
	opcodeCosts[OpcodeI32AtomicRmwXor] = opcodeCostsStruct.I32AtomicRmwXor
	opcodeCosts[OpcodeI64AtomicRmwXor] = opcodeCostsStruct.I64AtomicRmwXor
	opcodeCosts[OpcodeI32AtomicRmw8XorU] = opcodeCostsStruct.I32AtomicRmw8XorU
	opcodeCosts[OpcodeI32AtomicRmw16XorU] = opcodeCostsStruct.I32AtomicRmw16XorU
	opcodeCosts[OpcodeI64AtomicRmw8XorU] = opcodeCostsStruct.I64AtomicRmw8XorU
	opcodeCosts[OpcodeI64AtomicRmw16XorU] = opcodeCostsStruct.I64AtomicRmw16XorU
	opcodeCosts[OpcodeI64AtomicRmw32XorU] = opcodeCostsStruct.I64AtomicRmw32XorU
	opcodeCosts[OpcodeI32AtomicRmwXchg] = opcodeCostsStruct.I32AtomicRmwXchg
	opcodeCosts[OpcodeI64AtomicRmwXchg] = opcodeCostsStruct.I64AtomicRmwXchg
	opcodeCosts[OpcodeI32AtomicRmw8XchgU] = opcodeCostsStruct.I32AtomicRmw8XchgU
	opcodeCosts[OpcodeI32AtomicRmw16XchgU] = opcodeCostsStruct.I32AtomicRmw16XchgU
	opcodeCosts[OpcodeI64AtomicRmw8XchgU] = opcodeCostsStruct.I64AtomicRmw8XchgU
	opcodeCosts[OpcodeI64AtomicRmw16XchgU] = opcodeCostsStruct.I64AtomicRmw16XchgU
	opcodeCosts[OpcodeI64AtomicRmw32XchgU] = opcodeCostsStruct.I64AtomicRmw32XchgU
	opcodeCosts[OpcodeI32AtomicRmwCmpxchg] = opcodeCostsStruct.I32AtomicRmwCmpxchg
	opcodeCosts[OpcodeI64AtomicRmwCmpxchg] = opcodeCostsStruct.I64AtomicRmwCmpxchg
	opcodeCosts[OpcodeI32AtomicRmw8CmpxchgU] = opcodeCostsStruct.I32AtomicRmw8CmpxchgU
	opcodeCosts[OpcodeI32AtomicRmw16CmpxchgU] = opcodeCostsStruct.I32AtomicRmw16CmpxchgU
	opcodeCosts[OpcodeI64AtomicRmw8CmpxchgU] = opcodeCostsStruct.I64AtomicRmw8CmpxchgU
	opcodeCosts[OpcodeI64AtomicRmw16CmpxchgU] = opcodeCostsStruct.I64AtomicRmw16CmpxchgU
	opcodeCosts[OpcodeI64AtomicRmw32CmpxchgU] = opcodeCostsStruct.I64AtomicRmw32CmpxchgU
	opcodeCosts[OpcodeV128Load] = opcodeCostsStruct.V128Load
	opcodeCosts[OpcodeV128Store] = opcodeCostsStruct.V128Store
	opcodeCosts[OpcodeV128Const] = opcodeCostsStruct.V128Const
	opcodeCosts[OpcodeI8x16Splat] = opcodeCostsStruct.I8x16Splat
	opcodeCosts[OpcodeI8x16ExtractLaneS] = opcodeCostsStruct.I8x16ExtractLaneS
	opcodeCosts[OpcodeI8x16ExtractLaneU] = opcodeCostsStruct.I8x16ExtractLaneU
	opcodeCosts[OpcodeI8x16ReplaceLane] = opcodeCostsStruct.I8x16ReplaceLane
	opcodeCosts[OpcodeI16x8Splat] = opcodeCostsStruct.I16x8Splat
	opcodeCosts[OpcodeI16x8ExtractLaneS] = opcodeCostsStruct.I16x8ExtractLaneS
	opcodeCosts[OpcodeI16x8ExtractLaneU] = opcodeCostsStruct.I16x8ExtractLaneU
	opcodeCosts[OpcodeI16x8ReplaceLane] = opcodeCostsStruct.I16x8ReplaceLane
	opcodeCosts[OpcodeI32x4Splat] = opcodeCostsStruct.I32x4Splat
	opcodeCosts[OpcodeI32x4ExtractLane] = opcodeCostsStruct.I32x4ExtractLane
	opcodeCosts[OpcodeI32x4ReplaceLane] = opcodeCostsStruct.I32x4ReplaceLane
	opcodeCosts[OpcodeI64x2Splat] = opcodeCostsStruct.I64x2Splat
	opcodeCosts[OpcodeI64x2ExtractLane] = opcodeCostsStruct.I64x2ExtractLane
	opcodeCosts[OpcodeI64x2ReplaceLane] = opcodeCostsStruct.I64x2ReplaceLane
	opcodeCosts[OpcodeF32x4Splat] = opcodeCostsStruct.F32x4Splat
	opcodeCosts[OpcodeF32x4ExtractLane] = opcodeCostsStruct.F32x4ExtractLane
	opcodeCosts[OpcodeF32x4ReplaceLane] = opcodeCostsStruct.F32x4ReplaceLane
	opcodeCosts[OpcodeF64x2Splat] = opcodeCostsStruct.F64x2Splat
	opcodeCosts[OpcodeF64x2ExtractLane] = opcodeCostsStruct.F64x2ExtractLane
	opcodeCosts[OpcodeF64x2ReplaceLane] = opcodeCostsStruct.F64x2ReplaceLane
	opcodeCosts[OpcodeI8x16Eq] = opcodeCostsStruct.I8x16Eq
	opcodeCosts[OpcodeI8x16Ne] = opcodeCostsStruct.I8x16Ne
	opcodeCosts[OpcodeI8x16LtS] = opcodeCostsStruct.I8x16LtS
	opcodeCosts[OpcodeI8x16LtU] = opcodeCostsStruct.I8x16LtU
	opcodeCosts[OpcodeI8x16GtS] = opcodeCostsStruct.I8x16GtS
	opcodeCosts[OpcodeI8x16GtU] = opcodeCostsStruct.I8x16GtU
	opcodeCosts[OpcodeI8x16LeS] = opcodeCostsStruct.I8x16LeS
	opcodeCosts[OpcodeI8x16LeU] = opcodeCostsStruct.I8x16LeU
	opcodeCosts[OpcodeI8x16GeS] = opcodeCostsStruct.I8x16GeS
	opcodeCosts[OpcodeI8x16GeU] = opcodeCostsStruct.I8x16GeU
	opcodeCosts[OpcodeI16x8Eq] = opcodeCostsStruct.I16x8Eq
	opcodeCosts[OpcodeI16x8Ne] = opcodeCostsStruct.I16x8Ne
	opcodeCosts[OpcodeI16x8LtS] = opcodeCostsStruct.I16x8LtS
	opcodeCosts[OpcodeI16x8LtU] = opcodeCostsStruct.I16x8LtU
	opcodeCosts[OpcodeI16x8GtS] = opcodeCostsStruct.I16x8GtS
	opcodeCosts[OpcodeI16x8GtU] = opcodeCostsStruct.I16x8GtU
	opcodeCosts[OpcodeI16x8LeS] = opcodeCostsStruct.I16x8LeS
	opcodeCosts[OpcodeI16x8LeU] = opcodeCostsStruct.I16x8LeU
	opcodeCosts[OpcodeI16x8GeS] = opcodeCostsStruct.I16x8GeS
	opcodeCosts[OpcodeI16x8GeU] = opcodeCostsStruct.I16x8GeU
	opcodeCosts[OpcodeI32x4Eq] = opcodeCostsStruct.I32x4Eq
	opcodeCosts[OpcodeI32x4Ne] = opcodeCostsStruct.I32x4Ne
	opcodeCosts[OpcodeI32x4LtS] = opcodeCostsStruct.I32x4LtS
	opcodeCosts[OpcodeI32x4LtU] = opcodeCostsStruct.I32x4LtU
	opcodeCosts[OpcodeI32x4GtS] = opcodeCostsStruct.I32x4GtS
	opcodeCosts[OpcodeI32x4GtU] = opcodeCostsStruct.I32x4GtU
	opcodeCosts[OpcodeI32x4LeS] = opcodeCostsStruct.I32x4LeS
	opcodeCosts[OpcodeI32x4LeU] = opcodeCostsStruct.I32x4LeU
	opcodeCosts[OpcodeI32x4GeS] = opcodeCostsStruct.I32x4GeS
	opcodeCosts[OpcodeI32x4GeU] = opcodeCostsStruct.I32x4GeU
	opcodeCosts[OpcodeF32x4Eq] = opcodeCostsStruct.F32x4Eq
	opcodeCosts[OpcodeF32x4Ne] = opcodeCostsStruct.F32x4Ne
	opcodeCosts[OpcodeF32x4Lt] = opcodeCostsStruct.F32x4Lt
	opcodeCosts[OpcodeF32x4Gt] = opcodeCostsStruct.F32x4Gt
	opcodeCosts[OpcodeF32x4Le] = opcodeCostsStruct.F32x4Le
	opcodeCosts[OpcodeF32x4Ge] = opcodeCostsStruct.F32x4Ge
	opcodeCosts[OpcodeF64x2Eq] = opcodeCostsStruct.F64x2Eq
	opcodeCosts[OpcodeF64x2Ne] = opcodeCostsStruct.F64x2Ne
	opcodeCosts[OpcodeF64x2Lt] = opcodeCostsStruct.F64x2Lt
	opcodeCosts[OpcodeF64x2Gt] = opcodeCostsStruct.F64x2Gt
	opcodeCosts[OpcodeF64x2Le] = opcodeCostsStruct.F64x2Le
	opcodeCosts[OpcodeF64x2Ge] = opcodeCostsStruct.F64x2Ge
	opcodeCosts[OpcodeV128Not] = opcodeCostsStruct.V128Not
	opcodeCosts[OpcodeV128And] = opcodeCostsStruct.V128And
	opcodeCosts[OpcodeV128AndNot] = opcodeCostsStruct.V128AndNot
	opcodeCosts[OpcodeV128Or] = opcodeCostsStruct.V128Or
	opcodeCosts[OpcodeV128Xor] = opcodeCostsStruct.V128Xor
	opcodeCosts[OpcodeV128Bitselect] = opcodeCostsStruct.V128Bitselect
	opcodeCosts[OpcodeI8x16Neg] = opcodeCostsStruct.I8x16Neg
	opcodeCosts[OpcodeI8x16AnyTrue] = opcodeCostsStruct.I8x16AnyTrue
	opcodeCosts[OpcodeI8x16AllTrue] = opcodeCostsStruct.I8x16AllTrue
	opcodeCosts[OpcodeI8x16Shl] = opcodeCostsStruct.I8x16Shl
	opcodeCosts[OpcodeI8x16ShrS] = opcodeCostsStruct.I8x16ShrS
	opcodeCosts[OpcodeI8x16ShrU] = opcodeCostsStruct.I8x16ShrU
	opcodeCosts[OpcodeI8x16Add] = opcodeCostsStruct.I8x16Add
	opcodeCosts[OpcodeI8x16AddSaturateS] = opcodeCostsStruct.I8x16AddSaturateS
	opcodeCosts[OpcodeI8x16AddSaturateU] = opcodeCostsStruct.I8x16AddSaturateU
	opcodeCosts[OpcodeI8x16Sub] = opcodeCostsStruct.I8x16Sub
	opcodeCosts[OpcodeI8x16SubSaturateS] = opcodeCostsStruct.I8x16SubSaturateS
	opcodeCosts[OpcodeI8x16SubSaturateU] = opcodeCostsStruct.I8x16SubSaturateU
	opcodeCosts[OpcodeI8x16MinS] = opcodeCostsStruct.I8x16MinS
	opcodeCosts[OpcodeI8x16MinU] = opcodeCostsStruct.I8x16MinU
	opcodeCosts[OpcodeI8x16MaxS] = opcodeCostsStruct.I8x16MaxS
	opcodeCosts[OpcodeI8x16MaxU] = opcodeCostsStruct.I8x16MaxU
	opcodeCosts[OpcodeI8x16Mul] = opcodeCostsStruct.I8x16Mul
	opcodeCosts[OpcodeI16x8Neg] = opcodeCostsStruct.I16x8Neg
	opcodeCosts[OpcodeI16x8AnyTrue] = opcodeCostsStruct.I16x8AnyTrue
	opcodeCosts[OpcodeI16x8AllTrue] = opcodeCostsStruct.I16x8AllTrue
	opcodeCosts[OpcodeI16x8Shl] = opcodeCostsStruct.I16x8Shl
	opcodeCosts[OpcodeI16x8ShrS] = opcodeCostsStruct.I16x8ShrS
	opcodeCosts[OpcodeI16x8ShrU] = opcodeCostsStruct.I16x8ShrU
	opcodeCosts[OpcodeI16x8Add] = opcodeCostsStruct.I16x8Add
	opcodeCosts[OpcodeI16x8AddSaturateS] = opcodeCostsStruct.I16x8AddSaturateS
	opcodeCosts[OpcodeI16x8AddSaturateU] = opcodeCostsStruct.I16x8AddSaturateU
	opcodeCosts[OpcodeI16x8Sub] = opcodeCostsStruct.I16x8Sub
	opcodeCosts[OpcodeI16x8SubSaturateS] = opcodeCostsStruct.I16x8SubSaturateS
	opcodeCosts[OpcodeI16x8SubSaturateU] = opcodeCostsStruct.I16x8SubSaturateU
	opcodeCosts[OpcodeI16x8Mul] = opcodeCostsStruct.I16x8Mul
	opcodeCosts[OpcodeI16x8MinS] = opcodeCostsStruct.I16x8MinS
	opcodeCosts[OpcodeI16x8MinU] = opcodeCostsStruct.I16x8MinU
	opcodeCosts[OpcodeI16x8MaxS] = opcodeCostsStruct.I16x8MaxS
	opcodeCosts[OpcodeI16x8MaxU] = opcodeCostsStruct.I16x8MaxU
	opcodeCosts[OpcodeI32x4Neg] = opcodeCostsStruct.I32x4Neg
	opcodeCosts[OpcodeI32x4AnyTrue] = opcodeCostsStruct.I32x4AnyTrue
	opcodeCosts[OpcodeI32x4AllTrue] = opcodeCostsStruct.I32x4AllTrue
	opcodeCosts[OpcodeI32x4Shl] = opcodeCostsStruct.I32x4Shl
	opcodeCosts[OpcodeI32x4ShrS] = opcodeCostsStruct.I32x4ShrS
	opcodeCosts[OpcodeI32x4ShrU] = opcodeCostsStruct.I32x4ShrU
	opcodeCosts[OpcodeI32x4Add] = opcodeCostsStruct.I32x4Add
	opcodeCosts[OpcodeI32x4Sub] = opcodeCostsStruct.I32x4Sub
	opcodeCosts[OpcodeI32x4Mul] = opcodeCostsStruct.I32x4Mul
	opcodeCosts[OpcodeI32x4MinS] = opcodeCostsStruct.I32x4MinS
	opcodeCosts[OpcodeI32x4MinU] = opcodeCostsStruct.I32x4MinU
	opcodeCosts[OpcodeI32x4MaxS] = opcodeCostsStruct.I32x4MaxS
	opcodeCosts[OpcodeI32x4MaxU] = opcodeCostsStruct.I32x4MaxU
	opcodeCosts[OpcodeI64x2Neg] = opcodeCostsStruct.I64x2Neg
	opcodeCosts[OpcodeI64x2AnyTrue] = opcodeCostsStruct.I64x2AnyTrue
	opcodeCosts[OpcodeI64x2AllTrue] = opcodeCostsStruct.I64x2AllTrue
	opcodeCosts[OpcodeI64x2Shl] = opcodeCostsStruct.I64x2Shl
	opcodeCosts[OpcodeI64x2ShrS] = opcodeCostsStruct.I64x2ShrS
	opcodeCosts[OpcodeI64x2ShrU] = opcodeCostsStruct.I64x2ShrU
	opcodeCosts[OpcodeI64x2Add] = opcodeCostsStruct.I64x2Add
	opcodeCosts[OpcodeI64x2Sub] = opcodeCostsStruct.I64x2Sub
	opcodeCosts[OpcodeI64x2Mul] = opcodeCostsStruct.I64x2Mul
	opcodeCosts[OpcodeF32x4Abs] = opcodeCostsStruct.F32x4Abs
	opcodeCosts[OpcodeF32x4Neg] = opcodeCostsStruct.F32x4Neg
	opcodeCosts[OpcodeF32x4Sqrt] = opcodeCostsStruct.F32x4Sqrt
	opcodeCosts[OpcodeF32x4Add] = opcodeCostsStruct.F32x4Add
	opcodeCosts[OpcodeF32x4Sub] = opcodeCostsStruct.F32x4Sub
	opcodeCosts[OpcodeF32x4Mul] = opcodeCostsStruct.F32x4Mul
	opcodeCosts[OpcodeF32x4Div] = opcodeCostsStruct.F32x4Div
	opcodeCosts[OpcodeF32x4Min] = opcodeCostsStruct.F32x4Min
	opcodeCosts[OpcodeF32x4Max] = opcodeCostsStruct.F32x4Max
	opcodeCosts[OpcodeF64x2Abs] = opcodeCostsStruct.F64x2Abs
	opcodeCosts[OpcodeF64x2Neg] = opcodeCostsStruct.F64x2Neg
	opcodeCosts[OpcodeF64x2Sqrt] = opcodeCostsStruct.F64x2Sqrt
	opcodeCosts[OpcodeF64x2Add] = opcodeCostsStruct.F64x2Add
	opcodeCosts[OpcodeF64x2Sub] = opcodeCostsStruct.F64x2Sub
	opcodeCosts[OpcodeF64x2Mul] = opcodeCostsStruct.F64x2Mul
	opcodeCosts[OpcodeF64x2Div] = opcodeCostsStruct.F64x2Div
	opcodeCosts[OpcodeF64x2Min] = opcodeCostsStruct.F64x2Min
	opcodeCosts[OpcodeF64x2Max] = opcodeCostsStruct.F64x2Max
	opcodeCosts[OpcodeI32x4TruncSatF32x4S] = opcodeCostsStruct.I32x4TruncSatF32x4S
	opcodeCosts[OpcodeI32x4TruncSatF32x4U] = opcodeCostsStruct.I32x4TruncSatF32x4U
	opcodeCosts[OpcodeI64x2TruncSatF64x2S] = opcodeCostsStruct.I64x2TruncSatF64x2S
	opcodeCosts[OpcodeI64x2TruncSatF64x2U] = opcodeCostsStruct.I64x2TruncSatF64x2U
	opcodeCosts[OpcodeF32x4ConvertI32x4S] = opcodeCostsStruct.F32x4ConvertI32x4S
	opcodeCosts[OpcodeF32x4ConvertI32x4U] = opcodeCostsStruct.F32x4ConvertI32x4U
	opcodeCosts[OpcodeF64x2ConvertI64x2S] = opcodeCostsStruct.F64x2ConvertI64x2S
	opcodeCosts[OpcodeF64x2ConvertI64x2U] = opcodeCostsStruct.F64x2ConvertI64x2U
	opcodeCosts[OpcodeV8x16Swizzle] = opcodeCostsStruct.V8x16Swizzle
	opcodeCosts[OpcodeV8x16Shuffle] = opcodeCostsStruct.V8x16Shuffle
	opcodeCosts[OpcodeV8x16LoadSplat] = opcodeCostsStruct.V8x16LoadSplat
	opcodeCosts[OpcodeV16x8LoadSplat] = opcodeCostsStruct.V16x8LoadSplat
	opcodeCosts[OpcodeV32x4LoadSplat] = opcodeCostsStruct.V32x4LoadSplat
	opcodeCosts[OpcodeV64x2LoadSplat] = opcodeCostsStruct.V64x2LoadSplat
	opcodeCosts[OpcodeI8x16NarrowI16x8S] = opcodeCostsStruct.I8x16NarrowI16x8S
	opcodeCosts[OpcodeI8x16NarrowI16x8U] = opcodeCostsStruct.I8x16NarrowI16x8U
	opcodeCosts[OpcodeI16x8NarrowI32x4S] = opcodeCostsStruct.I16x8NarrowI32x4S
	opcodeCosts[OpcodeI16x8NarrowI32x4U] = opcodeCostsStruct.I16x8NarrowI32x4U
	opcodeCosts[OpcodeI16x8WidenLowI8x16S] = opcodeCostsStruct.I16x8WidenLowI8x16S
	opcodeCosts[OpcodeI16x8WidenHighI8x16S] = opcodeCostsStruct.I16x8WidenHighI8x16S
	opcodeCosts[OpcodeI16x8WidenLowI8x16U] = opcodeCostsStruct.I16x8WidenLowI8x16U
	opcodeCosts[OpcodeI16x8WidenHighI8x16U] = opcodeCostsStruct.I16x8WidenHighI8x16U
	opcodeCosts[OpcodeI32x4WidenLowI16x8S] = opcodeCostsStruct.I32x4WidenLowI16x8S
	opcodeCosts[OpcodeI32x4WidenHighI16x8S] = opcodeCostsStruct.I32x4WidenHighI16x8S
	opcodeCosts[OpcodeI32x4WidenLowI16x8U] = opcodeCostsStruct.I32x4WidenLowI16x8U
	opcodeCosts[OpcodeI32x4WidenHighI16x8U] = opcodeCostsStruct.I32x4WidenHighI16x8U
	opcodeCosts[OpcodeI16x8Load8x8S] = opcodeCostsStruct.I16x8Load8x8S
	opcodeCosts[OpcodeI16x8Load8x8U] = opcodeCostsStruct.I16x8Load8x8U
	opcodeCosts[OpcodeI32x4Load16x4S] = opcodeCostsStruct.I32x4Load16x4S
	opcodeCosts[OpcodeI32x4Load16x4U] = opcodeCostsStruct.I32x4Load16x4U
	opcodeCosts[OpcodeI64x2Load32x2S] = opcodeCostsStruct.I64x2Load32x2S
	opcodeCosts[OpcodeI64x2Load32x2U] = opcodeCostsStruct.I64x2Load32x2U
	opcodeCosts[OpcodeI8x16RoundingAverageU] = opcodeCostsStruct.I8x16RoundingAverageU
	opcodeCosts[OpcodeI16x8RoundingAverageU] = opcodeCostsStruct.I16x8RoundingAverageU
	opcodeCosts[OpcodeLocalAllocate] = opcodeCostsStruct.LocalAllocate
	// LocalsUnmetered, MaxMemoryGrow and MaxMemoryGrowDelta are not added to the
	// opcodeCosts array; the values will be sent to Wasmer as compilation
	// options instead

	return opcodeCosts
}
