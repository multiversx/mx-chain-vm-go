package mandosvalueinterpreter

import (
	"encoding/hex"
	"testing"

	fr "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/json/fileresolver"
	oj "github.com/ElrondNetwork/arwen-wasm-vm/mandos-go/orderedjson"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestBool(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("true")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = vi.InterpretString("false")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestString(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("``abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)

	result, err = vi.InterpretString("``")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("```")
	require.Nil(t, err)
	require.Equal(t, []byte("`"), result)

	result, err = vi.InterpretString("`` ")
	require.Nil(t, err)
	require.Equal(t, []byte(" "), result)

	result, err = vi.InterpretString("''abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)

	result, err = vi.InterpretString("''")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("'''")
	require.Nil(t, err)
	require.Equal(t, []byte("'"), result)

	result, err = vi.InterpretString("'' ")
	require.Nil(t, err)
	require.Equal(t, []byte(" "), result)

	result, err = vi.InterpretString("''``")
	require.Nil(t, err)
	require.Equal(t, []byte("``"), result)

	result, err = vi.InterpretString("``''")
	require.Nil(t, err)
	require.Equal(t, []byte("''"), result)

	result, err = vi.InterpretString("str:abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)

	result, err = vi.InterpretString("str:")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestAddress(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("address:")
	require.Nil(t, err)
	require.Equal(t, []byte("________________________________"), result)

	result, err = vi.InterpretString("address:a")
	require.Nil(t, err)
	require.Equal(t, []byte("a_______________________________"), result)

	result, err = vi.InterpretString("address:an_address")
	require.Nil(t, err)
	require.Equal(t, []byte("an_address______________________"), result)

	result, err = vi.InterpretString("address:12345678901234567890123456789012")
	require.Nil(t, err)
	require.Equal(t, []byte("12345678901234567890123456789012"), result)

	// trims excess
	result, err = vi.InterpretString("address:123456789012345678901234567890123")
	require.Nil(t, err)
	require.Equal(t, []byte("12345678901234567890123456789012"), result)
}

func TestSCAddress(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("sc:a")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00a_______________________"), result)

	result, err = vi.InterpretString("sc:123456789012345678901234")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00123456789012345678901234"), result)

	// trims excess
	result, err = vi.InterpretString("sc:123456789012345678901234x")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00123456789012345678901234"), result)
}

func TestUnsignedNumber(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x12, 0x34}, result)

	result, err = vi.InterpretString("0x")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("12")
	require.Nil(t, err)
	require.Equal(t, []byte{12}, result)

	result, err = vi.InterpretString("256")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x00}, result)

	result, err = vi.InterpretString("0b1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = vi.InterpretString("0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0x05}, result)
}

func TestSignedNumber(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = vi.InterpretString("255")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = vi.InterpretString("+255")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = vi.InterpretString("0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = vi.InterpretString("+0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = vi.InterpretString("-256")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0x00}, result)

	result, err = vi.InterpretString("-0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0xfb}, result)
}

func TestUnsignedFixedWidth(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("u8:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00}, result)

	result, err = vi.InterpretString("u16:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00}, result)

	result, err = vi.InterpretString("u32:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)

	result, err = vi.InterpretString("u64:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)

	result, err = vi.InterpretString("u16:0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x12, 0x34}, result)

	result, err = vi.InterpretString("u32:0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x12, 0x34}, result)

	result, err = vi.InterpretString("u16:256")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x00}, result)

	result, err = vi.InterpretString("u8:0b1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = vi.InterpretString("u64:0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05}, result)
}

func TestSignedFixedWidth(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("i8:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00}, result)

	result, err = vi.InterpretString("i16:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00}, result)

	result, err = vi.InterpretString("i32:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)

	result, err = vi.InterpretString("i64:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)

	result, err = vi.InterpretString("i8:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = vi.InterpretString("i16:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff}, result)

	result, err = vi.InterpretString("i32:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff}, result)

	result, err = vi.InterpretString("i64:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, result)

	result, err = vi.InterpretString("i8:255") // not completely ok, but we'll let this be for now
	require.Nil(t, err)                        // it could be argued that this should fail
	require.Equal(t, []byte{0xff}, result)     // it is however, consistent with the rest of the format

	result, err = vi.InterpretString("i8:+255")
	require.NotNil(t, err)

	result, err = vi.InterpretString("i16:+255")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = vi.InterpretString("i8:0xff") // same as for i8:255
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = vi.InterpretString("i8:+0xff")
	require.NotNil(t, err)

	result, err = vi.InterpretString("i16:+0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = vi.InterpretString("i64:-256")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}, result)

	result, err = vi.InterpretString("i8:-0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0xfb}, result)
}

func TestBigUint(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("biguint:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)

	result, err = vi.InterpretString("biguint:1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x01, 0x01}, result)

	result, err = vi.InterpretString("biguint:0x01FF")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x02, 0x01, 0xFF}, result)

	// should be positive
	result, err = vi.InterpretString("biguint:-0x01")
	require.NotNil(t, err)

	result, err = vi.InterpretString("nested:0x0102030405")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05}, result)

	// accepts any argument
	result, err = vi.InterpretString("nested:-0x01")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x01, 0xFF}, result)
}

func TestConcat(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("0x01|5")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x05}, result)

	result, err = vi.InterpretString("|||0x01|5||||")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x05}, result)

	result, err = vi.InterpretString("|")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("|||||||")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("|0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = vi.InterpretString("``a|``b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = vi.InterpretString("``a|str:b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = vi.InterpretString("``a|0x62")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = vi.InterpretString("0x61|``b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = vi.InterpretString("i16:0x61|u32:5")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x61, 0x00, 0x00, 0x00, 0x05}, result)

	result, err = vi.InterpretString("i64:-1|u8:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}, result)
}

func TestKeccak256(t *testing.T) {
	vi := ValueInterpreter{}
	result, err := vi.InterpretString("keccak256:0x01|5")
	require.Nil(t, err)
	expected, _ := keccak256([]byte{0x01, 0x05})
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:|||0x01|5||||")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{0x01, 0x05})
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:|")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:|||||||")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:|0")
	require.Nil(t, err)
	expected, _ = keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:``a|``b")
	require.Nil(t, err)
	expected, _ = keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:``a|0x62")
	require.Nil(t, err)
	expected, _ = keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:0x61|``b")
	require.Nil(t, err)
	expected, _ = keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	// some values from the old ERC20 tests
	result, err = vi.InterpretString("keccak256:1|0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("19efaebcc296cffac396adb4a60d54c05eff43926a6072498a618e943908efe1")
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:1|0x7777777777777777777707777777777777777777777777177777777777771234")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("a3da7395b9df9b4a0ad4ce2fd40d2db4c5b231dbc2a19ce9bafcbc2233dc1b0a")
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:1|0x5555555555555555555505555555555555555555555555155555555555551234")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("648147902a606bf61e05b8b9d828540be393187d2c12a271b45315628f8b05b9")
	require.Equal(t, expected, result)

	result, err = vi.InterpretString("keccak256:2|0x7777777777777777777707777777777777777777777777177777777777771234|0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("e314ce9b5b28a5927ee30ba28b67ee27ad8779e1101baf4224590c8f1e287891")
	require.Equal(t, expected, result)

}

func TestFile(t *testing.T) {
	vi := ValueInterpreter{
		FileResolver: fr.NewDefaultFileResolver(),
	}
	result, err := vi.InterpretString("file:../integrationTests/exampleFile.txt")
	require.Nil(t, err)
	require.Equal(t, []byte("hello!"), result)
}

func TestInterpretSubTree1(t *testing.T) {
	vi := ValueInterpreter{}
	jobj, err := oj.ParseOrderedJSON([]byte(`
		["''part1", "''part2"]
	`))
	require.Nil(t, err)
	result, err := vi.InterpretSubTree(jobj)
	require.Nil(t, err)
	require.Equal(t, []byte("part1part2"), result)
}

func TestInterpretSubTree2(t *testing.T) {
	vi := ValueInterpreter{}
	jobj, err := oj.ParseOrderedJSON([]byte(`
		{
			"''field1": "u32:5",
			"''field2": [
				"''field2elem1",
				"u64:0",
				["''field2elem3a", "''field2elem3b"]
			]
		}
	`))
	require.Nil(t, err)
	result, err := vi.InterpretSubTree(jobj)
	require.Nil(t, err)
	expected := []byte{0x00, 0x00, 0x00, 0x05}
	expected = append(expected, []byte("field2elem1")...)
	expected = append(expected, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	expected = append(expected, []byte("field2elem3a")...)
	expected = append(expected, []byte("field2elem3b")...)
	require.Equal(t, expected, result)
}
