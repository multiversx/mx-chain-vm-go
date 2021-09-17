#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"
#include "../elrond/test_utils.h"
#include "../elrond/bigFloat.h"

// byte gobEncodedFloat1[] = {1, 10, 0, 0, 0, 100, 0, 0, 0, 108, 136, 217, 65, 19, 144, 71, 160, 0};
// // = 173476272346174583562347456134583.6134671346713451345 

void init() {}

void BigFloatNewFromPartsTest() {
    int bigFloatHandle;
    int reps = int64getArgument(0);
    for (int i = 0; i < reps; i++) {
        bigFloatHandle = bigFloatNewFromParts(i,i,-i-1);
    }
    int64finish(bigFloatHandle);
}

void BigFloatNewFromFracTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    for (int i = 0; i < reps; i++) {
        bigFloatHandle = bigFloatNewFromFrac(reps+i,i+1);
    }
    int64finish(bigFloatHandle);
}

void BigFloatNewFromSciTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    int exponent = int64getArgument(1);
    for (int i = 0; i < reps; i++) {
        bigFloatHandle = bigFloatNewFromSci(-reps-i, exponent);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatAddTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle1 = bigFloatNewFromFrac(0,1);
    int bigFloatHandle2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1, bigFloatHandle1);
    bigFloatGetArgument(2, bigFloatHandle2);
    for (int i = 0; i < reps; i++) {
        bigFloatAdd(bigFloatHandle1,bigFloatHandle1,bigFloatHandle2);
    }
    bigFloatFinish(bigFloatHandle1);
}

void BigFloatSubTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle1 = bigFloatNewFromFrac(0,1);
    int bigFloatHandle2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandle1);
    bigFloatGetArgument(2,bigFloatHandle2);
    for (int i = 0; i < reps; i++) {
        bigFloatSub(bigFloatHandle1,bigFloatHandle1,bigFloatHandle2);
    }
    bigFloatFinish(bigFloatHandle1);
}

void BigFloatMulTest() {
    int reps, bigFloatHandle;
    reps = int64getArgument(0);
    bigFloatHandle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandle);
    for (int i = 0; i < reps; i++) {
        bigFloatMul(bigFloatHandle,bigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatDivTest() {
    int reps = int64getArgument(0);
    int bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    int bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1, bigFloatHandleOp1);
    bigFloatGetArgument(2, bigFloatHandleOp2);
    for (int i = 0; i < reps; i++) {
        bigFloatDiv(bigFloatHandleOp1,bigFloatHandleOp1,bigFloatHandleOp2);
    }
    bigFloatFinish(bigFloatHandleOp1);
}

void BigFloatTruncateTest() {
    int reps = int64getArgument(0);
    int bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    int bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,bigFloatHandleOp1);
    bigFloatGetArgument(2,bigFloatHandleOp2);
    for (int i = 0; i < reps; i++) {
        bigFloatTruncate(bigFloatHandleOp1, 0);
        bigFloatSub(bigFloatHandleOp1, bigFloatHandleOp1, bigFloatHandleOp2);
    }
    bigFloatFinish(bigFloatHandleOp1);
}

void BigFloatAbsTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    int absbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatAbs(absbigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(absbigFloatHandle);
}

void BigFloatNegTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    int negbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1, bigFloatHandle);
        bigFloatNeg(negbigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(negbigFloatHandle);
}

void BigFloatCmpTest() {
    int reps = int64getArgument(0);
    int bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    int bigFloatHandleOp2 = bigFloatNewFromFrac(0,1);
    int result, argsCounter = 1;
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(argsCounter,bigFloatHandleOp1);
        bigFloatGetArgument(argsCounter+1,bigFloatHandleOp2);
        argsCounter += 2;
        result = bigFloatCmp(bigFloatHandleOp1,bigFloatHandleOp2);
    }
    int64finish(result);
}

void BigFloatSignTest() {
    int result;
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        result = bigFloatSign(bigFloatHandle);
    }
    int64finish(result);
}

void BigFloatCloneTest() {
    int result;
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    int copybigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatClone(copybigFloatHandle, bigFloatHandle);
    }
    bigFloatFinish(copybigFloatHandle);
}

void BigFloatSqrtTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    int resultbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatSqrt(resultbigFloatHandle,bigFloatHandle);
    }
    bigFloatFinish(resultbigFloatHandle);
}

void BigFloatPowTest() {
    int reps = int64getArgument(0);
    int bigFloatHandleOp1 = bigFloatNewFromFrac(0,1);
    int resultbigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandleOp1);
        bigFloatPow(resultbigFloatHandle,bigFloatHandleOp1,i);
    }
    bigFloatFinish(resultbigFloatHandle);
}

void BigFloatFloorTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    int resultbigFloatHandle = bigIntNew(0);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatFloor(bigFloatHandle,resultbigFloatHandle);
    }
    bigIntFinishUnsigned(resultbigFloatHandle);
}

void BigFloatCeilTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    int resultbigFloatHandle = bigIntNew(0);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        bigFloatCeil(bigFloatHandle,resultbigFloatHandle);
    }
    bigIntFinishUnsigned(resultbigFloatHandle);
}

void BigFloatIsIntTest() {
    int result;
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetArgument(i+1,bigFloatHandle);
        result = bigFloatIsInt(bigFloatHandle);
    }
    int64finish(result);
}

void BigFloatSetInt64Test() {
    int value;
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        value = int64getArgument(i + 1);
        bigFloatSetInt64(bigFloatHandle,value);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatSetBigIntTest() {
    int bigIntHandle;
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigIntGetUnsignedArgument(i+1,bigIntHandle);
        bigFloatSetBigInt(bigFloatHandle,bigIntHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatGetConstPiTest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetConstPi(bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}

void BigFloatGetConstETest() {
    int reps = int64getArgument(0);
    int bigFloatHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++) {
        bigFloatGetConstE(bigFloatHandle);
    }
    bigFloatFinish(bigFloatHandle);
}
