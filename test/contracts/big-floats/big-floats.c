#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"
#include "../elrond/test_utils.h"
#include "../elrond/bigFloat.h"

// byte gobEncodedFloat1[] = {1, 10, 0, 0, 0, 100, 0, 0, 0, 108, 136, 217, 65, 19, 144, 71, 160, 0};
// // = 173476272346174583562347456134583.6134671346713451345 

void init() {}

void BigFloatNewTest() {
    int bigFloatHandle;
    int reps = int64getArgument(0);
    for (int i = 0; i < reps; i++) {
        bigFloatHandle = bigFloatNew(i,i,-i-1);
    }
    int64finish(bigFloatHandle);
}

void BigFloatNewFromFracTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    for (int i = 0; i < reps; i++)
    {
        bigFloatHandle = bigFloatNewFromFrac(reps+i,i+1);
    }
    int64finish(bigFloatHandle);
}

void BigFloatAddTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1, bigFloatHandle);
    for (int i = 0; i < reps; i++)
    {
        bigFloatAdd(bigFloatHandle,bigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatSubTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandle);
    for (int i = 0; i < reps; i++)
    {
        bigFloatSub(bigFloatHandle,bigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatMulTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandle);
    for (int i = 0; i < reps; i++)
    {
        bigFloatMul(bigFloatHandle,bigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatDivTest() {
    int reps, bigFloatHandleOp1, bigFloatHandleOp2;
    reps = int64getArgument(0);
    bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1, bigFloatHandleOp1);
    bigFloatGetArgument(2, bigFloatHandleOp2);
    for (int i = 0; i < reps; i++)
    {
        bigFloatDiv(bigFloatHandleOp1,bigFloatHandleOp1,bigFloatHandleOp2);
    }
    bigFloatFinish(bigFloatHandleOp1);
}

void BigFloatTruncateTest() {
    int reps, bigFloatHandleOp1, bigFloatHandleOp2;
    reps = int64getArgument(0);
    bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandleOp1);
    bigFloatGetArgument(2,bigFloatHandleOp2);
    for (int i = 0; i < reps; i++)
    {
        bigFloatTruncate(bigFloatHandleOp1);
        bigFloatSub(bigFloatHandleOp1, bigFloatHandleOp1, bigFloatHandleOp2);
    }
    bigFloatFinish(bigFloatHandleOp1);
}

void BigFloatModTest() {
    int reps, bigFloatHandleOp1, bigFloatHandleOp2;
    reps = int64getArgument(0);
    bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandleOp1);
    bigFloatGetArgument(2,bigFloatHandleOp2);
    for (int i = 0; i < reps; i++)
    {
        bigFloatMod(2,bigFloatHandleOp1,bigFloatHandleOp2);
        bigFloatSub(bigFloatHandleOp1,bigFloatHandleOp1,bigFloatHandleOp2);
    }
    bigFloatFinish(2);
}


void BigFloatAbsTest() {
    int reps, bigFloatHandle, absbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    absbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatAbs(absbigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(absbigFloatHandle);
}

void BigFloatNegTest() {
    int reps, bigFloatHandle, negbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    negbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1, bigFloatHandle);
        bigFloatNeg(negbigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(negbigFloatHandle);
}

void BigFloatCmpTest() {
    int reps, bigFloatHandleOp1, bigFloatHandleOp2;
    reps = int64getArgument(0);
    bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    int result, argsCounter = 1;
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(argsCounter,bigFloatHandleOp1);
        bigFloatGetArgument(argsCounter+1,bigFloatHandleOp2);
        argsCounter += 2;
        result = bigFloatCmp(bigFloatHandleOp1,bigFloatHandleOp2);
    }
    int64finish(result);
}

void BigFloatSignTest() {
    int reps, bigFloatHandle, result;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        result = bigFloatSign(bigFloatHandle);
    }
    int64finish(result);
}

void BigFloatCloneTest() {
    int reps, bigFloatHandle, result;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    int copybigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatClone(copybigFloatHandle, bigFloatHandle);
    }
    bigFloatFinish(copybigFloatHandle);
}

void BigFloatSqrtTest() {
    int reps, bigFloatHandle, resultbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    resultbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatSqrt(resultbigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(resultbigFloatHandle);
}

void BigFloatLog2Test() {
    int reps, bigFloatHandle, result;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        result = bigFloatLog2(bigFloatHandle);
    }
    int64finish(result);
}

void BigFloatPowTest() {
    int reps, bigFloatHandleOp1, resultbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    resultbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandleOp1);
        bigFloatPow(resultbigFloatHandle,bigFloatHandleOp1,i);
    }
    bigFloatFinish(resultbigFloatHandle);
}

void BigFloatFloorTest() {
    int reps, bigFloatHandle, resultbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    resultbigFloatHandle = bigIntNew(0);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatFloor(bigFloatHandle,resultbigFloatHandle);
    }
    bigIntFinishUnsigned(resultbigFloatHandle);
}

void BigFloatCeilTest() {
    int reps, bigFloatHandle, resultbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    resultbigFloatHandle = bigIntNew(0);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatCeil(bigFloatHandle,resultbigFloatHandle);
    }
    bigIntFinishUnsigned(resultbigFloatHandle);
}

void BigFloatIsIntTest() {
    int reps, bigFloatHandle, result;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+1,bigFloatHandle);
        result = bigFloatIsInt(bigFloatHandle);
    }
    int64finish(result);
}

void BigFloatSetInt64Test() {
    int reps, bigFloatHandle, value;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        value = int64getArgument(i + 1);
        bigFloatSetInt64(bigFloatHandle,value);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatSetBigIntTest() {
    int reps, bigFloatHandle, bigIntbigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigIntGetUnsignedArgument(i+1,bigIntbigFloatHandle);
        bigFloatSetBigInt(bigFloatHandle,bigIntbigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatGetConstPiTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetConstPi(bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatGetConstETest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetConstE(bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}
