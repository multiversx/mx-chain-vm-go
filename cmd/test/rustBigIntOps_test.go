package main

import (
	"fmt"
	"math/big"
	"testing"

	arwen "github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

func getSmallNumbers() []*big.Int {
	return []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(10),
		big.NewInt(12345),
	}
}

func getLargeNumbers() []*big.Int {
	big1, _ := big.NewInt(0).SetString("18446744073709551615", 10)
	big2, _ := big.NewInt(0).SetString("18446744073709551616", 10)
	big3, _ := big.NewInt(0).SetString("18446744073709551617", 10)

	return []*big.Int{
		big.NewInt(2147483647),
		big.NewInt(2147483648),
		big.NewInt(2147483649),
		big.NewInt(4294967295),
		big.NewInt(4294967296),
		big.NewInt(4294967297),
		big1,
		big2,
		big3,
	}
}

func getPositiveNumbers() []*big.Int {
	var numbers []*big.Int
	numbers = append(numbers, getSmallNumbers()...)
	numbers = append(numbers, getLargeNumbers()...)
	return numbers
}

func getNumbers() []*big.Int {
	positiveNumbers := getPositiveNumbers()
	var numbers []*big.Int
	numbers = append(numbers, positiveNumbers...)
	for _, num := range numbers {
		neg := big.NewInt(0).Neg(num)
		numbers = append(numbers, neg)
	}
	return numbers
}

func unsignedInterpreter(bytes []byte) *big.Int {
	return big.NewInt(0).SetBytes(bytes)
}

func signedInterpreter(bytes []byte) *big.Int {
	return twos.FromBytes(bytes)
}

func appendBinaryOpTestCase(
	testCases []*pureFunctionIO,
	opName string, signed bool,
	arg1, arg2, result []byte,
	expectedStatus vmi.ReturnCode, expectedMessage string,
) []*pureFunctionIO {

	var typeName string
	if signed {
		typeName = "big_int"
	} else {
		typeName = "big_uint"
	}

	expectedResults := [][]byte{}
	if expectedStatus == vmi.Ok {
		expectedResults = [][]byte{result}
	}

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_%s", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_%s_ref", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_assign_%s", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	testCases = append(testCases, &pureFunctionIO{
		functionName:    fmt.Sprintf("%s_assign_%s_ref", opName, typeName),
		arguments:       [][]byte{arg1, arg2},
		expectedStatus:  expectedStatus,
		expectedMessage: expectedMessage,
		expectedResults: expectedResults,
	})

	return testCases
}

func TestBigUintArith(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	var testCases []*pureFunctionIO

	numbers := getPositiveNumbers()

	for _, num1 := range numbers {
		for _, num2 := range numbers {
			bytes1 := num1.Bytes()
			bytes2 := num2.Bytes()

			// add
			sumBytes := big.NewInt(0).Add(num1, num2).Bytes()
			testCases = appendBinaryOpTestCase(testCases,
				"add", false,
				bytes1, bytes2, sumBytes,
				vmi.Ok, "")

			// sub
			diff := big.NewInt(0).Sub(num1, num2)
			if diff.Sign() < 0 {
				testCases = appendBinaryOpTestCase(testCases,
					"sub", false,
					bytes1, bytes2, nil,
					vmi.UserError, "cannot subtract because result would be negative")
			} else {
				testCases = appendBinaryOpTestCase(testCases,
					"sub", false,
					bytes1, bytes2, diff.Bytes(),
					vmi.Ok, "")
			}

			// mul
			mulBytes := big.NewInt(0).Mul(num1, num2).Bytes()
			testCases = appendBinaryOpTestCase(testCases,
				"mul", false,
				bytes1, bytes2, mulBytes,
				vmi.Ok, "")

			// div
			if num2.Sign() == 0 {
				testCases = appendBinaryOpTestCase(testCases,
					"div", false,
					bytes1, bytes2, nil,
					vmi.UserError, arwen.UserErrorDivZero)
			} else {
				divBytes := big.NewInt(0).Quo(num1, num2).Bytes()
				testCases = appendBinaryOpTestCase(testCases,
					"div", false,
					bytes1, bytes2, divBytes,
					vmi.Ok, "")
			}

			// mod
			if num2.Sign() == 0 {
				testCases = appendBinaryOpTestCase(testCases,
					"rem", false,
					bytes1, bytes2, nil,
					vmi.UserError, arwen.UserErrorDivZero)
			} else {
				remBytes := big.NewInt(0).Rem(num1, num2).Bytes()
				testCases = appendBinaryOpTestCase(testCases,
					"rem", false,
					bytes1, bytes2, remBytes,
					vmi.Ok, "")
			}
		}
	}

	logFunc := func(testCaseIndex, testCaseCount int) {
		if testCaseIndex%100 == 0 {
			fmt.Printf("Big uint operation test case %d/%d\n", testCaseIndex, len(testCases))
		}
	}

	pureFunctionTest(t, testCases, unsignedInterpreter, logFunc)
}

func TestBigIntArith(t *testing.T) {
	if testing.Short() {
		t.Skip("this is not a short test")
	}

	var testCases []*pureFunctionIO

	numbers := getNumbers()

	for _, num1 := range numbers {
		for _, num2 := range numbers {
			bytes1 := twos.ToBytes(num1)
			bytes2 := twos.ToBytes(num2)

			// add
			sumBytes := twos.ToBytes(big.NewInt(0).Add(num1, num2))
			testCases = appendBinaryOpTestCase(testCases,
				"add", true,
				bytes1, bytes2, sumBytes,
				vmi.Ok, "")

			// sub
			diffBytes := twos.ToBytes(big.NewInt(0).Sub(num1, num2))
			testCases = appendBinaryOpTestCase(testCases,
				"sub", true,
				bytes1, bytes2, diffBytes,
				vmi.Ok, "")

			// mul
			mulBytes := twos.ToBytes(big.NewInt(0).Mul(num1, num2))
			testCases = appendBinaryOpTestCase(testCases,
				"mul", true,
				bytes1, bytes2, mulBytes,
				vmi.Ok, "")

			// div
			if num2.Sign() == 0 {
				testCases = appendBinaryOpTestCase(testCases,
					"div", true,
					bytes1, bytes2, nil,
					vmi.UserError, arwen.UserErrorDivZero)
			} else {
				divBytes := twos.ToBytes(big.NewInt(0).Quo(num1, num2))
				testCases = appendBinaryOpTestCase(testCases,
					"div", true,
					bytes1, bytes2, divBytes,
					vmi.Ok, "")
			}

			// mod
			if num2.Sign() == 0 {
				testCases = appendBinaryOpTestCase(testCases,
					"rem", true,
					bytes1, bytes2, nil,
					vmi.UserError, arwen.UserErrorDivZero)
			} else {
				remBytes := twos.ToBytes(big.NewInt(0).Rem(num1, num2))
				testCases = appendBinaryOpTestCase(testCases,
					"rem", true,
					bytes1, bytes2, remBytes,
					vmi.Ok, "")
			}
		}
	}

	logFunc := func(testCaseIndex, testCaseCount int) {
		if testCaseIndex%100 == 0 {
			fmt.Printf("Big int operation test case %d/%d\n", testCaseIndex, len(testCases))
		}
	}

	pureFunctionTest(t, testCases, unsignedInterpreter, logFunc)
}
