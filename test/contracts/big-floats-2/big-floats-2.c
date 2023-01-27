#include "../mxvm/context.h"
#include "../mxvm/bigInt.h"
#include "../mxvm/types.h"
#include "../mxvm/test_utils.h"
#include "../mxvm/bigFloat.h"

void init() {}

void BigFloatToManagedBufferTest()
{
    mBufferGetArgument(0, 0);
    mBufferToBigFloat(0, 0);
}

void BigFloatNewFromPartsTest()
{
    int integralPart = int64getArgument(0);
    int fractionalPart = int64getArgument(1);
    int exponent = int64getArgument(2);
    int bigFloatHandle = bigFloatNewFromParts(integralPart, fractionalPart, exponent);
}

void BigFloatNewFromFracTest()
{
    int numerator = int64getArgument(0);
    int denominator = int64getArgument(1);
    int bigFloatHandle = bigFloatNewFromFrac(numerator, denominator);
}

void BigFloatNewFromSciTest()
{
    int significand = int64getArgument(0);
    int exponent = int64getArgument(1);
    int bigFloatHandle = bigFloatNewFromSci(significand, exponent);
}

void BigFloatAddTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    bigFloatAdd(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
}

void BigFloatSubTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    bigFloatSub(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
}

void BigFloatMulTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    bigFloatMul(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
}

void BigFloatDivTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    bigFloatDiv(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
}

void BigFloatTruncateTest()
{
    int bigFloatHandle1 = 0;
    bigFloatGetArgument(1, bigFloatHandle1);
    bigFloatTruncate(bigFloatHandle1, 0);
}

void BigFloatAbsTest()
{
    int bigFloatHandle = 0;
    int absbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatAbs(absbigFloatHandle, bigFloatHandle);
}

void BigFloatNegTest()
{
    int bigFloatHandle = 0;
    int negbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatNeg(negbigFloatHandle, bigFloatHandle);
}

void BigFloatCmpTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    bigFloatCmp(bigFloatHandle1, bigFloatHandle2);
}

void BigFloatSignTest()
{
    int bigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatSign(bigFloatHandle);
}

void BigFloatCloneTest()
{
    int bigFloatHandle = 0;
    int copybigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatClone(copybigFloatHandle, bigFloatHandle);
}

void BigFloatSqrtTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatSqrt(resultbigFloatHandle, bigFloatHandle);
}

void BigFloatPowTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    int exponent = int64getArgument(1);
    bigFloatPow(resultbigFloatHandle, bigFloatHandle, exponent);
}

void BigFloatFloorTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatFloor(resultbigFloatHandle, bigFloatHandle);
}

void BigFloatCeilTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatCeil(resultbigFloatHandle, bigFloatHandle);
}

void BigFloatIsIntTest()
{
    int bigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    bigFloatIsInt(bigFloatHandle);
}

void BigFloatSetInt64Test()
{
    int bigFloatHandle = 0;
    int value = int64getArgument(0);
    bigFloatSetInt64(bigFloatHandle, value);
}

void BigFloatSetBigIntTest()
{
    int bigIntHandle;
    int bigFloatHandle = 0;
    bigIntGetUnsignedArgument(0, bigIntHandle);
    bigFloatSetBigInt(bigFloatHandle, bigIntHandle);
}

void BigFloatGetConstPiTest()
{
    int bigFloatHandle = 0;
    bigFloatGetConstPi(bigFloatHandle);
}

void BigFloatGetConstETest()
{
    int bigFloatHandle = 0;
    bigFloatGetConstE(bigFloatHandle);
}
