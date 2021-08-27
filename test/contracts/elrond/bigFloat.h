#ifndef _BIGFLOAT_H_
#define _BIGFLOAT_H_

#include "types.h"


int bigFloatNew(int intBase, int subIntBase, int exponent);
int bigFloatNewFromFrac(long long numerator, long long denominator);

void bigFloatAdd(int destinationHandle, int op1Handle, int op2Handle);
void bigFloatSub(int destinationHandle, int op1Handle, int op2Handle);
void bigFloatMul(int destinationHandle, int op1Handle, int op2Handle);
void bigFloatDiv(int destinationHandle, int op1Handle, int op2Handle);
void bigFloatRoundDiv(int destinationHandle, int op1Handle, int op2Handle);
void bigFloatMod(int destinationHandle, int op1Handle, int op2Handle);

void bigFloatAbs(int destinationHandle, int opHandle);
void bigFloatNeg(int destinationHandle, int opHandle);
int	bigFloatCmp(int op1Handle, int op2Handle);
int	bigFloatSign(int opHandle);
void bigFloatCopy(int destinationHandle, int opHandle);
void bigFloatSqrt(int destinationHandle, int opHandle);
int	bigFloatLog2(int opHandle);
void bigFloatPow(int destinationHandle, int op1Handle, int op2Handle);

void bigFloatFloor(int opHandle, int bigIntHandle);
void bigFloatCeil(int opHandle, int bigIntHandle);

int	bigFloatIsInt(int opHandle);
void bigFloatSetInt64(int destinationHandle, long long value);
void bigFloatSetBigInt(int destinationHandle, int bigIntHandle);

void bigFloatGetConstPi(int destinationHandle);
void bigFloatGetConstE(int destinationHandle);

void bigFloatSetBytes(int destinationHandle, byte* dataOffset, int dataLength);
void bigFloatGetBytes(int destinationHandle, byte* dataOffset);

void bigFloatFinish(int referenceHandle);
void bigFloatGetArgument(int id, int destinationHandle);

#endif
