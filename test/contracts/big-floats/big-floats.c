#include "../mxvm/context.h"
#include "../mxvm/bigInt.h"
#include "../mxvm/types.h"
#include "../mxvm/test_utils.h"
#include "../mxvm/bigFloat.h"

void init() {}

void BigFloatToManagedBufferTest()
{
    mBufferGetArgument(0, 0);
    for (int i = 0; i < 150000; i++)
    {
        mBufferToBigFloat(0, 0);
    }
}

void BigFloatNewFromPartsTest()
{
    int integralPart = int64getArgument(0);
    int fractionalPart = int64getArgument(1);
    int exponent = int64getArgument(2);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatNewFromParts(integralPart, fractionalPart, exponent);
    }    
}

void BigFloatNewFromFracTest()
{
    int numerator = int64getArgument(0);
    int denominator = int64getArgument(1);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatNewFromFrac(numerator, denominator);
    }
}

void BigFloatNewFromSciTest()
{
    int significand = int64getArgument(0);
    int exponent = int64getArgument(1);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatNewFromSci(significand, exponent);
    }
}

void BigFloatAddTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    for (int i = 0; i < 98000; i++)
    {
        bigFloatAdd(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
    }
}

void BigFloatSubTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    for (int i = 0; i < 98000; i++)
    {
        bigFloatSub(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
    }   
}

void BigFloatMulTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    for (int i = 0; i < 70000; i++)
    {
        bigFloatMul(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
    }
}

void BigFloatDivTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    for (int i = 0; i < 70000; i++)
    {
        bigFloatDiv(bigFloatHandle1, bigFloatHandle1, bigFloatHandle2);
    }
}

void BigFloatTruncateTest()
{
    int bigFloatHandle1 = 0;
    bigFloatGetArgument(1, bigFloatHandle1);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatTruncate(bigFloatHandle1, 0);
    }
}

void BigFloatAbsTest()
{
    int bigFloatHandle = 0;
    int absbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 140000; i++)
    {
          bigFloatAbs(absbigFloatHandle, bigFloatHandle);  
    }    
}

void BigFloatNegTest()
{
    int bigFloatHandle = 0;
    int negbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 140000; i++)
    {
        bigFloatNeg(negbigFloatHandle, bigFloatHandle);
    }
}

void BigFloatCmpTest()
{
    int bigFloatHandle1 = 0;
    int bigFloatHandle2 = 1;
    bigFloatGetArgument(0, bigFloatHandle1);
    bigFloatGetArgument(1, bigFloatHandle2);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatCmp(bigFloatHandle1, bigFloatHandle2);
    }
}

void BigFloatSignTest()
{
    int bigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 140000; i++)
    {
        bigFloatSign(bigFloatHandle);
    }
}

void BigFloatCloneTest()
{
    int bigFloatHandle = 0;
    int copybigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 140000; i++)
    {
        bigFloatClone(copybigFloatHandle, bigFloatHandle);
    }    
}

void BigFloatSqrtTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatSqrt(resultbigFloatHandle, bigFloatHandle);    
    }  
}

void BigFloatPowTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 1;
    bigFloatGetArgument(0, bigFloatHandle);
    int exponent = int64getArgument(1);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatPow(resultbigFloatHandle, bigFloatHandle, exponent);
    }   
}

void BigFloatFloorTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 210000; i++)
    {
        bigFloatFloor(resultbigFloatHandle, bigFloatHandle);    
    }    
}

void BigFloatCeilTest()
{
    int bigFloatHandle = 0;
    int resultbigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatCeil(resultbigFloatHandle, bigFloatHandle);
    }    
}

void BigFloatIsIntTest()
{
    int bigFloatHandle = 0;
    bigFloatGetArgument(0, bigFloatHandle);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatIsInt(bigFloatHandle);
    }
}

void BigFloatSetInt64Test()
{
    int bigFloatHandle = 0;
    int value = int64getArgument(0);
    for (int i = 0; i < 910000; i++)
    {
        bigFloatSetInt64(bigFloatHandle, value);    
    }    
}

void BigFloatSetBigIntTest()
{
    int bigIntHandle;
    int bigFloatHandle = 0;
    bigIntGetUnsignedArgument(0, bigIntHandle);
    for (int i = 0; i < 10000; i++)
    {
        bigFloatSetBigInt(bigFloatHandle, bigIntHandle);    
    }    
}

void BigFloatGetConstPiTest()
{
    int bigFloatHandle = 0;
    for (int i = 0; i < 910000; i++)
    {
        bigFloatGetConstPi(bigFloatHandle);    
    }    
}

void BigFloatGetConstETest()
{
    int bigFloatHandle = 0;
    for (int i = 0; i < 910000; i++)
    {
        bigFloatGetConstE(bigFloatHandle);
    }    
}
