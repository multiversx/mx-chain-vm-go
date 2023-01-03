package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"

	twos "github.com/ElrondNetwork/big-int-util/twos-complement"
)

const numberOfDataSets = 10000

var positiveEncodedBigFloatPrefix = [...]byte{1, 10, 0, 0, 0, 53}
var negativeEncodedBigFloatPrefix = [...]byte{1, 11, 0, 0, 0, 53}
var positiveEncodedBigFloatForPowPrefix = [...]byte{1, 10, 0, 0, 53, 0, 0, 0}
var negativeEncodedBigFloatForPowPrefix = [...]byte{1, 11, 0, 0, 53, 0, 0, 0}

// contract endpoint names
const (
	bfToManagedBufferName = "BigFloatToManagedBufferTest"
	bfFromPartsName       = "BigFloatNewFromPartsTest"
	bfFromFracName        = "BigFloatNewFromFracTest"
	bfFromSciName         = "BigFloatNewFromSciTest"
	bfAddName             = "BigFloatAddTest"
	bfSubName             = "BigFloatSubTest"
	bfMulName             = "BigFloatMulTest"
	bfDivName             = "BigFloatDivTest"
	bfTruncName           = "BigFloatTruncateTest"
	bfAbsName             = "BigFloatAbsTest"
	bfNegName             = "BigFloatNegTest"
	bfCmpName             = "BigFloatCmpTest"
	bfSignName            = "BigFloatSignTest"
	bfCloneName           = "BigFloatCloneTest"
	bfSqrtName            = "BigFloatSqrtTest"
	bfPowName             = "BigFloatPowTest"
	bfFloorName           = "BigFloatFloorTest"
	bfCeilName            = "BigFloatCeilTest"
	bfIsIntName           = "BigFloatIsIntTest"
	bfSetInt64Name        = "BigFloatSetInt64Test"
	bfSetIntName          = "BigFloatSetIntTest"
)

//POW AND SETBIGINT

func main() {
	generateBigFloatFromPartsData()
	generateBigFloatsFromFracData()
	generateBigFloatsFromSciData()
	generateDataForEndpoint(1, bfToManagedBufferName, 4000000)
	generateDataForEndpoint(2, bfAddName, 4000000)
	generateDataForEndpoint(2, bfSubName, 4000000)
	generateDataForEndpoint(2, bfMulName, 4000000)
	generateDataForEndpoint(2, bfDivName, 4000000)
	generateDataForEndpoint(1, bfTruncName, 4000000)
	generateDataForEndpoint(1, bfAbsName, 1390000)
	generateDataForEndpoint(1, bfNegName, 1390000)
	generateDataForEndpoint(2, bfCmpName, 4000000)
	generateDataForEndpoint(1, bfSignName, 1390000)
	generateDataForEndpoint(1, bfCloneName, 1390000)
	generateDataForEndpoint(1, bfSqrtName, 4000000)
	generateDataForBigFloatPow()
	generateDataForEndpoint(1, bfFloorName, 1390000)
	generateDataForEndpoint(1, bfCeilName, 1390000)
	generateDataForEndpoint(1, bfIsIntName, 4000000)
	generateDataForEndpoint(1, bfSetInt64Name, 4000000)
	generateBigFloatsSetBigInt()
	generateBigFloatsSetInt64()

}

func generateBigFloatFromPartsData() {
	file, _ := os.Create(bfFromPartsName + ".data")
	for i := 0; i < numberOfDataSets; i++ {
		_, _ = file.WriteString("BigFloatNewFromPartsTest@")

		// integralPart
		integralPart := rand.Intn(math.MaxInt32 - 1)
		bigIntegralPart := big.NewInt(0).SetInt64(int64(integralPart))
		if rand.Intn(2) == 1 {
			bigIntegralPart.Neg(bigIntegralPart)
		}
		hexEncodedIntegralPart := hex.EncodeToString(twos.ToBytes(bigIntegralPart))
		_, _ = file.WriteString(hexEncodedIntegralPart + "@")
		// fractionalPart
		fractionalPart := rand.Intn(math.MaxInt32 - 1)
		bigFractionalPart := big.NewInt(0).SetInt64(int64(fractionalPart))
		hexEncodedFractionalPart := hex.EncodeToString(twos.ToBytes(bigFractionalPart))
		_, _ = file.WriteString(hexEncodedFractionalPart + "@")
		// exponent
		exponent := rand.Intn(400)
		bigExponent := big.NewInt(0).SetInt64(int64(exponent))
		validExponent := big.NewInt(0).Neg(bigExponent)
		hexEncodedExponent := hex.EncodeToString(validExponent.Bytes())
		_, _ = file.WriteString(hexEncodedExponent + ":30000" + "\n")

	}
	defer func() {
		_ = file.Close()
	}()
}

func generateBigFloatsFromFracData() {
	file, _ := os.Create(bfFromFracName + ".data")

	for i := 0; i < numberOfDataSets; i++ {
		_, _ = file.WriteString("BigFloatNewFromFracTest@")

		// numerator
		numeratorPart := rand.Intn(math.MaxInt64 - 1)
		bigNumeratorPart := big.NewInt(0).SetInt64(int64(numeratorPart))
		if rand.Intn(2) == 1 {
			bigNumeratorPart.Neg(bigNumeratorPart)
		}
		hexEncodedNumerator := hex.EncodeToString(bigNumeratorPart.Bytes())
		_, _ = file.WriteString(hexEncodedNumerator + "@")
		// denominator
		denominatorPart := rand.Intn(math.MaxInt64 - 1)
		bigDenominatorPart := big.NewInt(0).SetInt64(int64(denominatorPart))
		hexEncodedDenominator := hex.EncodeToString(bigDenominatorPart.Bytes())
		_, _ = file.WriteString(hexEncodedDenominator + ":30000" + "\n")

	}
	defer func() {
		_ = file.Close()
	}()
}

func generateBigFloatsFromSciData() {
	file, _ := os.Create(bfFromSciName + ".data")

	for i := 0; i < numberOfDataSets; i++ {
		_, _ = file.WriteString("BigFloatNewFromSciTest@")

		// significand
		significandPart := rand.Intn(math.MaxInt64 - 1)
		bigSignificandPart := big.NewInt(0).SetInt64(int64(significandPart))
		if rand.Intn(2) == 1 {
			bigSignificandPart.Neg(bigSignificandPart)
		}
		hexEncodedSignificand := hex.EncodeToString(bigSignificandPart.Bytes())
		_, _ = file.WriteString(hexEncodedSignificand + "@")
		// exponent
		exponentPart := rand.Intn(400)
		bigExponentPart := big.NewInt(0).SetInt64(int64(exponentPart))
		hexEncodedExponent := hex.EncodeToString(bigExponentPart.Bytes())
		_, _ = file.WriteString(hexEncodedExponent + ":30000" + "\n")
	}

	defer func() {
		_ = file.Close()
	}()
}

func generateBigFloatsSetBigInt() {
	file, _ := os.Create(bfSetIntName + ".data")

	for i := 0; i < numberOfDataSets; i++ {
		_, _ = file.WriteString("BigFloatSetBigIntTest@")
		numberOfBytes := rand.Intn(200)
		bigIntBytes := make([]byte, numberOfBytes)
		rand.Read(bigIntBytes)

		hexEncodedBytes := hex.EncodeToString(bigIntBytes)
		_, _ = file.WriteString(hexEncodedBytes + ":4000000" + "\n")
	}
	defer func() {
		_ = file.Close()
	}()
}

func generateBigFloatsSetInt64() {
	file, _ := os.Create(bfSetInt64Name + ".data")

	for i := 0; i < numberOfDataSets; i++ {
		_, _ = file.WriteString("BigFloatSetInt64Test@")
		smallValue := rand.Intn(math.MaxInt64 - 1)
		bigIntVal := big.NewInt(int64(smallValue))
		if rand.Intn(2) == 1 {
			bigIntVal.Neg(bigIntVal)
		}
		argumentBytes := twos.ToBytes(bigIntVal)
		hexEncodedArgument := hex.EncodeToString(argumentBytes)
		_, _ = file.WriteString(hexEncodedArgument + ":4000000" + "\n")
	}
	defer func() {
		_ = file.Close()
	}()
}

func generateDataForBigFloatPow() {
	file, _ := os.Create(bfPowName + ".data")
	for i := 0; i < 1000; i++ {
		_, _ = file.WriteString("BigFloatPowTest@")
		hexEncodedFloat := generateHexEncodedBigFloatForPow()
		_, _ = file.WriteString(hexEncodedFloat + "@")

		bytes, _ := hex.DecodeString(hexEncodedFloat)
		floatVal := big.NewFloat(0)
		_ = floatVal.GobDecode(bytes)
		intOp := big.NewInt(0)
		floatVal.Int(intOp)

		//exponent
		exponentBytes := make([]byte, 1)
		rand.Read(exponentBytes)
		bigExponent := big.NewInt(0).SetBytes(exponentBytes)
		if rand.Intn(2) == 1 {
			bigExponent.Neg(bigExponent)
		}
		exponentBytes, _ = twos.ToBytesOfLength(bigExponent, 4)

		//gas cost
		lengthOfResult := big.NewInt(0).Div(big.NewInt(0).Mul(bigExponent, big.NewInt(int64(intOp.BitLen()))), big.NewInt(8))
		gasForPow := big.NewInt(0).Mul(lengthOfResult, big.NewInt(1000))
		//fmt.Println(gasForPow)
		gasToUse := uint64(math.MaxUint64)
		if gasForPow.Cmp(big.NewInt(0).SetUint64(math.MaxUint64)) < 0 {
			gasToUse = gasForPow.Uint64() + 10000
		}
		hexEncodedExponent := hex.EncodeToString(exponentBytes)
		_, _ = file.WriteString(hexEncodedExponent)
		_, _ = file.WriteString(":" + strconv.Itoa(int(gasToUse)) + "\n")
	}
	defer func() {
		_ = file.Close()
	}()
}

func generateDataForEndpoint(numberOfBigFloats int, endpointName string, gasLimit int) {
	fileName := fmt.Sprintf("%s.data", endpointName)
	file, _ := os.Create(fileName)

	for i := 0; i < numberOfDataSets; i++ {
		_, _ = file.WriteString(endpointName)
		for j := 0; j < numberOfBigFloats; j++ {
			bigFloatValue := generateHexEncodedBigFloat()
			_, _ = file.WriteString("@" + bigFloatValue)
		}
		_, _ = file.WriteString(":" + strconv.Itoa(gasLimit) + "\n")
	}
	defer func() {
		_ = file.Close()
	}()
}

func generateHexEncodedBigFloat() string {
	encodedBigFloat := make([]byte, 0)
	if rand.Intn(2) == 1 {
		encodedBigFloat = append(encodedBigFloat, positiveEncodedBigFloatPrefix[:]...)
	} else {
		encodedBigFloat = append(encodedBigFloat, negativeEncodedBigFloatPrefix[:]...)
	}
	randomExponentAndMantissa := make([]byte, 12)
	rand.Read(randomExponentAndMantissa)
	encodedBigFloat = append(encodedBigFloat, randomExponentAndMantissa...)
	hexEncodedBigFloat := hex.EncodeToString(encodedBigFloat)
	return hexEncodedBigFloat
}

func generateHexEncodedBigFloatForPow() string {
	encodedBigFloat := make([]byte, 0)
	if rand.Intn(2) == 1 {
		encodedBigFloat = append(encodedBigFloat, positiveEncodedBigFloatForPowPrefix[:]...)
	} else {
		encodedBigFloat = append(encodedBigFloat, negativeEncodedBigFloatForPowPrefix[:]...)
	}
	randomExponentAndMantissa := make([]byte, 9)
	rand.Read(randomExponentAndMantissa)
	encodedBigFloat = append(encodedBigFloat, randomExponentAndMantissa...)
	hexEncodedBigFloat := hex.EncodeToString(encodedBigFloat)
	return hexEncodedBigFloat
}
