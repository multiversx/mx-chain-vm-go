#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"
#include "../elrond/test_utils.h"
#include "../elrond/bigFloat.h"

// byte gobEncodedFloat1[] = {1, 10, 0, 0, 0, 100, 0, 0, 0, 108, 136, 217, 65, 19, 144, 71, 160, 0};
// // = 173476272346174583562347456134583.6134671346713451345 


void init() {}

void BigFloatNewTest() {
    int handle;
    int reps = int64getArgument(0);
    for (int i = 0; i < reps; i++) {
        handle = bigFloatNew(i,i,-i-1);
    }
    int64finish(handle);
}

void BigFloatNewFromFracTest() {
    int reps, handle;
    reps = int64getArgument(0);
    for (int i = 0; i < reps; i++)
    {
        handle = bigFloatNewFromFrac(reps+i,i+1);
    }
    int64finish(handle);
}

void BigFloatAddTest() {
    int reps, handle;
    reps = int64getArgument(0);
    handle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,handle);
    for (int i = 0; i < reps; i++)
    {
        bigFloatAdd(handle,handle,handle);
    }
    bigFloatFinish(handle);
}

void BigFloatSubTest() {
    int reps, handle;
    reps = int64getArgument(0);
    handle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,handle);
    for (int i = 0; i < reps; i++)
    {
        bigFloatSub(handle,handle,handle);
    }
    bigFloatFinish(handle);
}

void BigFloatMulTest() {
    int reps, handle;
    reps = int64getArgument(0);
    handle = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,handle);
    for (int i = 0; i < reps; i++)
    {
        bigFloatMul(handle,handle,handle);
    }
    bigFloatFinish(handle);
}

void BigFloatDivTest() {
    int reps, handleOp1, handleOp2;
    reps = int64getArgument(0);
    handleOp1 = bigFloatNewFromFrac(0,1);
    handleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,handleOp1);
    bigFloatGetArgument(2,handleOp2);
    for (int i = 0; i < reps; i++)
    {
        bigFloatDiv(handleOp1,handleOp1,handleOp2);
    }
    bigFloatFinish(handleOp1);
}

void BigFloatRoundDivTest() {
    int reps, handleOp1, handleOp2;
    reps = int64getArgument(0);
    handleOp1 = bigFloatNewFromFrac(0,1);
    handleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,handleOp1);
    bigFloatGetArgument(2,handleOp2);
    for (int i = 0; i < reps; i++)
    {
        bigFloatRoundDiv(handleOp1,handleOp1,handleOp2);
    }
    bigFloatFinish(handleOp1);
}


void BigFloatModTest() {
    int reps, handleOp1, handleOp2;
    reps = int64getArgument(0);
    handleOp1 = bigFloatNewFromFrac(0,1);
    handleOp2 = bigFloatNewFromFrac(0,1);
    bigFloatGetArgument(1,handleOp1);
    bigFloatGetArgument(2,handleOp2);
    for (int i = 0; i < reps; i++)
    {
        bigFloatMod(2,handleOp1,handleOp2);
        bigFloatSub(handleOp1,handleOp1,handleOp2);
    }
    bigFloatFinish(2);
}

/*
void BigFloatAbsTest() {
    int reps, handle, absHandle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    absHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatAbs(absHandle,handle);
    }
    bigFloatFinish(absHandle);
}

void BigFloatNegTest() {
    int reps, handle, negHandle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    negHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2, handle);
        bigFloatNeg(negHandle,handle);
    }
    bigFloatFinish(negHandle);
}

void BigFloatCmpTest() {
    int reps, handleOp1, handleOp2;
    reps = int64getArgument(1);
    handleOp1 = bigFloatNewFromFrac(0,1);
    handleOp2 = bigFloatNewFromFrac(0,1);
    int result, argsCounter = 2;
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(argsCounter,handleOp1);
        bigFloatGetArgument(argsCounter+1,handleOp2);
        argsCounter += 2;
        result = bigFloatCmp(handleOp1,handleOp2);
    }
    int64finish(result);
}

void BigFloatSignTest() {
    int reps, handle, result;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        result = bigFloatSign(handle);
    }
    int64finish(result);
}

void BigFloatCopyTest() {
    int reps, handle, result;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    int copyHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatCopy(handle, copyHandle);
    }
    bigFloatFinish(copyHandle);
}

void BigFloatSqrtTest() {
    int reps, handle, resultHandle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    resultHandle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatSqrt(resultHandle,handle);
    }
    bigFloatFinish(resultHandle);
}

void BigFloatLog2Test() {
    int reps, handle, result;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        result = bigFloatLog2(handle);
    }
    int64finish(result);
}

void BigFloatPowTest() {
    int reps, handleOp1, handleOp2, resultHandle;
    reps = int64getArgument(1);
    handleOp1 = bigFloatNewFromFrac(0,1);
    handleOp2 = bigFloatNewFromFrac(0,1);
    resultHandle = bigFloatNewFromFrac(0,1);
    int argsCounter = 2;
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(argsCounter,handleOp1);
        bigFloatGetArgument(argsCounter+1,handleOp2);
        bigFloatPow(resultHandle,handleOp1,handleOp2);
    }
    bigFloatFinish(resultHandle);
}

void BigFloatFloorTest() {
    int reps, handle, resultHandle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    resultHandle = bigIntNew(0);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatFloor(handle,resultHandle);
    }
    bigIntFinishUnsigned(resultHandle);
}

void BigFloatCeilTest() {
    int reps, handle, resultHandle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    resultHandle = bigIntNew(0);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatCeil(handle,resultHandle);
    }
    bigIntFinishUnsigned(resultHandle);
}

void BigFloatIsIntTest() {
    int reps, handle, result;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        result = bigFloatIsInt(handle);
    }
    int64finish(result);
}

void BigFloatSetInt64Test() {
    int reps, handle, value;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        value = int64getArgument(i + 2);
        bigFloatSetInt64(handle,value);
    }
    bigFloatFinish(handle);
}

void BigFloatSetBigIntTest() {
    int reps, handle, bigIntHandle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigIntGetUnsignedArgument(i+2,bigIntHandle);
        bigFloatSetBigInt(handle,bigIntHandle);
    }
    bigFloatFinish(handle);
}

void BigFloatGetConstPiTest() {
    int reps, handle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetConstPi(handle);
    }
}

void BigFloatGetConstETest() {
    int reps, handle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetConstE(handle);
    }
}

void BigFloatSetBytesTest() {
    int reps, handle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    byte buffer[255];
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatGetBytes(handle,buffer);
        bigFloatSetBytes(handle,buffer,18);
    }
}

void BigFloatGetBytesTest() {
    int reps, handle;
    reps = int64getArgument(1);
    handle = bigFloatNewFromFrac(0,1);
    byte buffer[255];
    for (int i = 0; i < reps; i++)
    {
        bigFloatGetArgument(i+2,handle);
        bigFloatGetBytes(handle,buffer);
    }
}*/
