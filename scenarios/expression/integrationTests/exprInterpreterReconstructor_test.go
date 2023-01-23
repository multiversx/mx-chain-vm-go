package mandosjsontest

import (
	"encoding/hex"
	"testing"

	mei "github.com/multiversx/mx-chain-vm-go/scenarios/expression/interpreter"
	mer "github.com/multiversx/mx-chain-vm-go/scenarios/expression/reconstructor"
	fr "github.com/multiversx/mx-chain-vm-go/scenarios/fileresolver"
	oj "github.com/multiversx/mx-chain-vm-go/scenarios/orderedjson"
	"github.com/stretchr/testify/require"
)

func interpreter() mei.ExprInterpreter {
	return mei.ExprInterpreter{
		VMType: &[2]byte{'V', 'M'},
	}
}

func reconstructor() mer.ExprReconstructor {
	return mer.ExprReconstructor{}
}

func TestEmpty(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestBool(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("true")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = ei.InterpretString("false")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
}

func TestString(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("``abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)
	require.Equal(t, "str:abcdefg", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("``")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
	require.Equal(t, "str:", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("```")
	require.Nil(t, err)
	require.Equal(t, []byte("`"), result)
	require.Equal(t, "str:`", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("`` ")
	require.Nil(t, err)
	require.Equal(t, []byte(" "), result)
	require.Equal(t, "str: ", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("''abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)
	require.Equal(t, "str:abcdefg", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("''")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
	require.Equal(t, "str:", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("'''")
	require.Nil(t, err)
	require.Equal(t, []byte("'"), result)
	require.Equal(t, "str:'", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("'' ")
	require.Nil(t, err)
	require.Equal(t, []byte(" "), result)
	require.Equal(t, "str: ", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("''``")
	require.Nil(t, err)
	require.Equal(t, []byte("``"), result)
	require.Equal(t, "str:``", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("``''")
	require.Nil(t, err)
	require.Equal(t, []byte("''"), result)
	require.Equal(t, "str:''", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("str:abcdefg")
	require.Nil(t, err)
	require.Equal(t, []byte("abcdefg"), result)
	require.Equal(t, "str:abcdefg", er.Reconstruct(result, mer.StrHint))

	result, err = ei.InterpretString("str:")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
	require.Equal(t, "str:", er.Reconstruct(result, mer.StrHint))
}

func TestAddress(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("address:")
	require.Nil(t, err)
	require.Equal(t, []byte("________________________________"), result)
	require.Equal(t, "address:", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:a")
	require.Nil(t, err)
	require.Equal(t, []byte("a_______________________________"), result)
	require.Equal(t, "address:a", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:a\x05")
	require.Nil(t, err)
	require.Equal(t, []byte("a\x05______________________________"), result)
	require.Equal(t, "address:a\x05", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:an_address")
	require.Nil(t, err)
	require.Equal(t, []byte("an_address______________________"), result)
	require.Equal(t, "address:an_address", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:1234567890123456789012345678901\x01")
	require.Nil(t, err)
	require.Equal(t, []byte("1234567890123456789012345678901\x01"), result)
	require.Equal(t, "address:1234567890123456789012345678901#01", er.Reconstruct(result, mer.AddressHint))

	// trims excess
	result, err = ei.InterpretString("address:1234567890123456789012345678901\x013")
	require.Nil(t, err)
	require.Equal(t, []byte("1234567890123456789012345678901\x01"), result)
	require.Equal(t, "address:1234567890123456789012345678901#01", er.Reconstruct(result, mer.AddressHint))
}

func TestAddressWithShardId(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("address:#05")
	require.Nil(t, err)
	require.Equal(t, []byte("_______________________________\x05"), result)
	require.Equal(t, "address:#05", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:a#bb")
	require.Nil(t, err)
	require.Equal(t, []byte("a______________________________\xbb"), result)
	require.Equal(t, "address:a#bb", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:an_address#99")
	require.Nil(t, err)
	require.Equal(t, []byte("an_address_____________________\x99"), result)
	require.Equal(t, "address:an_address#99", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("address:1234567890123456789012345678901#66")
	require.Nil(t, err)
	require.Equal(t, []byte("1234567890123456789012345678901\x66"), result)
	require.Equal(t, "address:1234567890123456789012345678901#66", er.Reconstruct(result, mer.AddressHint))

	// trims excess
	result, err = ei.InterpretString("address:12345678901234567890123456789012#66")
	require.Nil(t, err)
	require.Equal(t, []byte("1234567890123456789012345678901\x66"), result)
	require.Equal(t, "address:1234567890123456789012345678901#66", er.Reconstruct(result, mer.AddressHint))
}

func TestSCAddress(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("sc:a")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00VMa_____________________"), result)
	require.Equal(t, "sc:a", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("sc:123456789012345678912s")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00VM123456789012345678912s"), result)
	require.Equal(t, "sc:123456789012345678912#73", er.Reconstruct(result, mer.AddressHint))

	// trims excess
	result, err = ei.InterpretString("sc:123456789012345678912sx")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00VM123456789012345678912s"), result)
	require.Equal(t, "sc:123456789012345678912#73", er.Reconstruct(result, mer.AddressHint))
}

func TestSCAddressWithShardId(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("sc:a#44")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00VMa____________________\x44"), result)
	require.Equal(t, "sc:a#44", er.Reconstruct(result, mer.AddressHint))

	result, err = ei.InterpretString("sc:123456789012345678912#88")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00VM123456789012345678912\x88"), result)
	require.Equal(t, "sc:123456789012345678912#88", er.Reconstruct(result, mer.AddressHint))

	// trims excess
	result, err = ei.InterpretString("sc:123456789012345678912x#88")
	require.Nil(t, err)
	require.Equal(t, []byte("\x00\x00\x00\x00\x00\x00\x00\x00VM123456789012345678912\x88"), result)
	require.Equal(t, "sc:123456789012345678912#88", er.Reconstruct(result, mer.AddressHint))
}

func TestUnsignedNumber(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x12, 0x34}, result)

	result, err = ei.InterpretString("0x")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
	require.Equal(t, "0", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)
	require.Equal(t, "0", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("12")
	require.Nil(t, err)
	require.Equal(t, []byte{12}, result)
	require.Equal(t, "12", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("256")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x00}, result)
	require.Equal(t, "256", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("0b1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)
	require.Equal(t, "1", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0x05}, result)
}

func TestSignedNumber(t *testing.T) {
	ei := interpreter()
	er := reconstructor()

	result, err := ei.InterpretString("-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = ei.InterpretString("255")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)
	require.Equal(t, "255", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("+255")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)
	require.Equal(t, "255", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)
	require.Equal(t, "255", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("+0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)
	require.Equal(t, "255", er.Reconstruct(result, mer.NumberHint))

	result, err = ei.InterpretString("-256")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0x00}, result)

	result, err = ei.InterpretString("-0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0xfb}, result)
}

func TestUnsignedFixedWidth(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("u8:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00}, result)

	result, err = ei.InterpretString("u16:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00}, result)

	result, err = ei.InterpretString("u32:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)

	result, err = ei.InterpretString("u64:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)

	result, err = ei.InterpretString("u16:0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x12, 0x34}, result)

	result, err = ei.InterpretString("u32:0x1234")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x12, 0x34}, result)

	result, err = ei.InterpretString("u16:256")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x00}, result)

	result, err = ei.InterpretString("u8:0b1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01}, result)

	result, err = ei.InterpretString("u64:0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05}, result)
}

func TestSignedFixedWidth(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("i8:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00}, result)

	result, err = ei.InterpretString("i16:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00}, result)

	result, err = ei.InterpretString("i32:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)

	result, err = ei.InterpretString("i64:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, result)

	result, err = ei.InterpretString("i8:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	result, err = ei.InterpretString("i16:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff}, result)

	result, err = ei.InterpretString("i32:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff}, result)

	result, err = ei.InterpretString("i64:-1")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, result)

	result, err = ei.InterpretString("i8:255") // not completely ok, but we'll let this be for now
	require.Nil(t, err)                        // it could be argued that this should fail
	require.Equal(t, []byte{0xff}, result)     // it is however, consistent with the rest of the format

	_, err = ei.InterpretString("i8:+255")
	require.NotNil(t, err)

	result, err = ei.InterpretString("i16:+255")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = ei.InterpretString("i8:0xff") // same as for i8:255
	require.Nil(t, err)
	require.Equal(t, []byte{0xff}, result)

	_, err = ei.InterpretString("i8:+0xff")
	require.NotNil(t, err)

	result, err = ei.InterpretString("i16:+0xff")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0xff}, result)

	result, err = ei.InterpretString("i64:-256")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}, result)

	result, err = ei.InterpretString("i8:-0b101")
	require.Nil(t, err)
	require.Equal(t, []byte{0xfb}, result)
}

func TestBigFloat(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("bigfloat:6297134613497.34523924564572445")
	require.Nil(t, err)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x12, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("6297134613497.34523924564572445")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("-6297134613497.34523924564572445")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x0b, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("bigfloat:6_297_134_613_497.34523924564572445")
	require.Nil(t, err)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x12, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("bigfloat:6,297,134,613,497.34523924564572445")
	require.Nil(t, err)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x12, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("bigfloat:0x010a000000350000002bb7454f187f2b1000")
	require.Nil(t, err)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x12, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("bigfloat:0x01,0a,00,00,00,35,00,00,00,2b,b7,45,4f,18,7f,2b,10,00")
	require.Nil(t, err)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x12, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

	result, err = ei.InterpretString("bigfloat:0x01_0a_00_00_00_35_00_00_00_2b_b7_45_4f_18_7f_2b_10_00")
	require.Nil(t, err)
	require.Equal(t, []byte{0x0, 0x0, 0x0, 0x12, 0x01, 0x0a, 0x00, 0x00, 0x00, 0x35, 0x00, 0x00, 0x00, 0x2b, 0xb7, 0x45, 0x4f, 0x18, 0x7f, 0x2b, 0x10, 0x00}, result)

}

func TestBigUint(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("biguint:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x00}, result)

	result, err = ei.InterpretString("biguint:1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x01, 0x01}, result)

	result, err = ei.InterpretString("biguint:27,80")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x02, 0x0A, 0xDC}, result)

	result, err = ei.InterpretString("biguint:27_80_86")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x03, 0x04, 0x3E, 0x46}, result)

	result, err = ei.InterpretString("biguint:1")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x01, 0x01}, result)

	result, err = ei.InterpretString("biguint:0x01FF")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x02, 0x01, 0xFF}, result)

	result, err = ei.InterpretString("biguint:0x01,FF,CB")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x03, 0x01, 0xFF, 0xCB}, result)

	result, err = ei.InterpretString("biguint:0x01_FF_CB")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x03, 0x01, 0xFF, 0xCB}, result)

	// should be positive
	_, err = ei.InterpretString("biguint:-0x01")
	require.NotNil(t, err)

	result, err = ei.InterpretString("nested:0x0102030405")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05}, result)

	// accepts any argument
	result, err = ei.InterpretString("nested:-0x01")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x00, 0x00, 0x01, 0xFF}, result)
}

func TestConcat(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("0x01|5")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x05}, result)

	result, err = ei.InterpretString("|||0x01|5||||")
	require.Nil(t, err)
	require.Equal(t, []byte{0x01, 0x05}, result)

	result, err = ei.InterpretString("|")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = ei.InterpretString("|||||||")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = ei.InterpretString("|0")
	require.Nil(t, err)
	require.Equal(t, []byte{}, result)

	result, err = ei.InterpretString("``a|``b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = ei.InterpretString("``a|str:b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = ei.InterpretString("``a|0x62")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = ei.InterpretString("0x61|``b")
	require.Nil(t, err)
	require.Equal(t, []byte("ab"), result)

	result, err = ei.InterpretString("i16:0x61|u32:5")
	require.Nil(t, err)
	require.Equal(t, []byte{0x00, 0x61, 0x00, 0x00, 0x00, 0x05}, result)

	result, err = ei.InterpretString("i64:-1|u8:0")
	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00}, result)
}

func TestKeccak256(t *testing.T) {
	ei := interpreter()
	result, err := ei.InterpretString("keccak256:0x01|5")
	require.Nil(t, err)
	expected, _ := mei.Keccak256([]byte{0x01, 0x05})
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:|||0x01|5||||")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte{0x01, 0x05})
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:|")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:|||||||")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:|0")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte{})
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:``a|``b")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:``a|0x62")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:0x61|``b")
	require.Nil(t, err)
	expected, _ = mei.Keccak256([]byte("ab"))
	require.Equal(t, expected, result)

	// some values from the old ERC20 tests
	result, err = ei.InterpretString("keccak256:1|0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("19efaebcc296cffac396adb4a60d54c05eff43926a6072498a618e943908efe1")
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:1|0x7777777777777777777707777777777777777777777777177777777777771234")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("a3da7395b9df9b4a0ad4ce2fd40d2db4c5b231dbc2a19ce9bafcbc2233dc1b0a")
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:1|0x5555555555555555555505555555555555555555555555155555555555551234")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("648147902a606bf61e05b8b9d828540be393187d2c12a271b45315628f8b05b9")
	require.Equal(t, expected, result)

	result, err = ei.InterpretString("keccak256:2|0x7777777777777777777707777777777777777777777777177777777777771234|0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000")
	require.Nil(t, err)
	expected, _ = hex.DecodeString("e314ce9b5b28a5927ee30ba28b67ee27ad8779e1101baf4224590c8f1e287891")
	require.Equal(t, expected, result)

}

func TestFile(t *testing.T) {
	ei := mei.ExprInterpreter{
		FileResolver: fr.NewDefaultFileResolver(),
	}
	result, err := ei.InterpretString("file:../../json/integrationTests/exampleFile.txt")
	require.Nil(t, err)
	require.Equal(t, []byte("hello!"), result)
}

func TestInterpretSubTree1(t *testing.T) {
	ei := interpreter()
	jobj, err := oj.ParseOrderedJSON([]byte(`
		["''part1", "''part2"]
	`))
	require.Nil(t, err)
	result, err := ei.InterpretSubTree(jobj)
	require.Nil(t, err)
	require.Equal(t, []byte("part1part2"), result)
}

func TestInterpretSubTree2(t *testing.T) {
	ei := interpreter()
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
	result, err := ei.InterpretSubTree(jobj)
	require.Nil(t, err)
	expected := []byte{0x00, 0x00, 0x00, 0x05}
	expected = append(expected, []byte("field2elem1")...)
	expected = append(expected, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}...)
	expected = append(expected, []byte("field2elem3a")...)
	expected = append(expected, []byte("field2elem3b")...)
	require.Equal(t, expected, result)
}
