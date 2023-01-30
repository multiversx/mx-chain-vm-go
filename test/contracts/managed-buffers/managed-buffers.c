#include "../mxvm/context.h"
#include "../mxvm/bigInt.h"
#include "../mxvm/types.h"
#include "../mxvm/test_utils.h"
byte mBuffer1[] = {
        0xff, 0x2a, 0x26, 0x5f, 0x8b, 0xcb, 0xdc, 0xaf, 
        0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24, 
        0xcb, 0x40, 0xd6, 0x4a, 0x50, 0x1f, 0xba, 0x9c, 
        0x11, 0x84, 0x7b, 0x28, 0x96, 0x5b, 0xc7, 0x37
    };
byte mBuffer2[] = {
        0xff, 0x2a, 0x26, 0x5f, 0x8b, 0xcb, 0xdc, 0xaf, 
        0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24, 
        0xcb, 0x40, 0xd6, 0x4a, 0x50, 0x1f, 0xba, 0x9c, 
        0x11, 0x84, 0x7b, 0x28, 0x96, 0x5b, 0xc7, 0x37,
        0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24, 
        0xcb, 0x40, 0xd6, 0x4a, 0x50, 0x1f, 0xba, 0x9c, 
        0x11, 0x84, 0x7b, 0x28, 0x96, 0x5b, 0xc7, 0x37,
        0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24,
    };
byte mBuffer3[] = {  // this is mBuffer2 - mBuffer1
        0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24, 
        0xcb, 0x40, 0xd6, 0x4a, 0x50, 0x1f, 0xba, 0x9c, 
        0x11, 0x84, 0x7b, 0x28, 0x96, 0x5b, 0xc7, 0x37,
        0xd5, 0x85, 0x19, 0x14, 0x1e, 0x57, 0x81, 0x24,
    };

byte mBufferKey[] = "mBuffer";

void mBufferMethod();
int verifyIfBuffersAreEqual(int handle1, int handle2);
int verifyBytesMBufferAndBigInt(int bigIntHandle, int mBufferHandle, int isSigned);
int byteArraysAreEqual(byte firstArray[], byte secondArray[], int length);

void init() {}

void mBufferMethod() {
    //Basic functionalities
    int mBufferHandle1 = mBufferNew();
    int mBufferHandle2 = mBufferNewFromBytes(mBuffer1, sizeof(mBuffer1));
    mBufferSetBytes(mBufferHandle1, mBuffer1, sizeof(mBuffer1));
    int ok = 0;

    if (mBufferHandle1 != 0 || mBufferHandle2 != 1) ok = 1;

    if (mBufferGetLength(mBufferHandle1) != mBufferGetLength(mBufferHandle2)) ok = 1;

    mBufferSetBytes(mBufferHandle1, mBuffer2, sizeof(mBuffer2));
    mBufferSetBytes(mBufferHandle2, mBuffer1, sizeof(mBuffer1));
    mBufferAppendBytes(mBufferHandle2,mBuffer3, sizeof(mBuffer3));
    if (verifyIfBuffersAreEqual(mBufferHandle1, mBufferHandle2)==1) ok = 1;

    int bigIntHandle1 = bigIntNew(0);
    int bigIntHandle2 = bigIntNew(0);
    if (bigIntHandle1 != 0 || bigIntHandle2 != 1) ok = 1;

    // To/From BigInts functionalities
    mBufferToBigIntUnsigned(mBufferHandle1, bigIntHandle1);
    if(verifyBytesMBufferAndBigInt(bigIntHandle1,mBufferHandle1, 0)==1) ok = 1;
    mBufferToBigIntSigned(mBufferHandle1, bigIntHandle2);
    if (verifyBytesMBufferAndBigInt(bigIntHandle2,mBufferHandle1, 1) != 0) ok = 1;

    int mBufferHandle3 = mBufferNew();
    mBufferFromBigIntUnsigned(mBufferHandle3, bigIntHandle1);
    int mBufferHandle4 = mBufferNew();
    mBufferFromBigIntSigned(mBufferHandle4, bigIntHandle2);
    if( verifyBytesMBufferAndBigInt(bigIntHandle1,mBufferHandle3, 0) != 0) ok = 1;
    if( verifyBytesMBufferAndBigInt(bigIntHandle2,mBufferHandle4, 1) != 0) ok = 1;
    if( verifyIfBuffersAreEqual(mBufferHandle1,mBufferHandle4) != 0) ok = 1;
    if( verifyIfBuffersAreEqual(mBufferHandle1,mBufferHandle3) != 0) ok = 1;

    // Storage
    int storageKeyLength = sizeof(mBufferKey) - 1;
    int keyHandle = mBufferNewFromBytes(mBufferKey,storageKeyLength);
    if( mBufferStorageStore(keyHandle, mBufferHandle4) != 0) ok = 1;
    int mBufferHandle5 = mBufferNew();
    if( mBufferStorageLoad(keyHandle, mBufferHandle5) != 0) ok = 1;
    if( verifyIfBuffersAreEqual(mBufferHandle4,mBufferHandle5) != 0) ok = 1;

    // Finish
    if( mBufferFinish(mBufferHandle4) != 0) ok = 1;
    int lengthReturnData = getReturnDataSize(0);
    int lengthOfBuffer = mBufferGetLength(mBufferHandle4);
    if (lengthReturnData!=lengthOfBuffer) ok = 1;
    byte returnDataBuffer[255];
    getReturnData(0,returnDataBuffer);
    mBufferSetBytes(mBufferHandle5,returnDataBuffer,lengthReturnData);
    if ( verifyIfBuffersAreEqual(mBufferHandle4,mBufferHandle5) != 0) ok = 1;

    //Random
    int randomBufferHandle = mBufferNew();
    int result = mBufferSetRandom(randomBufferHandle, 100);
    if (mBufferGetLength(randomBufferHandle) != 100 || result == 1) ok = 1;
    
    finishResult(ok);
}

void mBufferNewTest() {
    int reps, handle;
    reps = int64getArgument(0);
    for (int i = 0; i < reps; i++)
    {
        handle = mBufferNew();
    }
    int64finish(handle);
}

void mBufferNewFromBytesTest() {
    int reps, lengthOfBuffer, handle;
    reps = int64getArgument(0);
    lengthOfBuffer = int64getArgument(1);
    for (int i = 0; i < reps; i++)
    {
        handle = mBufferNewFromBytes(mBuffer2,lengthOfBuffer);      
    }
    mBufferFinish(handle);
}

void mBufferSetRandomTest() {
    int reps, result;
    reps = int64getArgument(0);
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle, reps);
    }
    mBufferFinish(randomBufferHandle);
}

void mBufferGetLengthTest() {
    int reps, result, length;
    reps = int64getArgument(0);
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,i+1);
        length = mBufferGetLength(randomBufferHandle);
    }
    int64finish(length);
}

void mBufferGetBytesTest() {
    int reps, result;
    reps = int64getArgument(0);
    byte returnDataBuffer[255];
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferGetBytes(randomBufferHandle, returnDataBuffer);
    }
    mBufferFinish(randomBufferHandle);
    finish(returnDataBuffer,reps);
}

void mBufferSetByteSliceTest() {
		int startPos = int64getArgument(0);
		int copyLen = int64getArgument(1);

		if (copyLen > 36) {
			byte msg[] = "max 36 bytes to copy";
			signalError(msg, 20);
		}

		byte sourceBytes[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890";
		byte destBytes[] = "abcdefghijklmnopqrstuvwxyz";

		int buffer = mBufferNew();
		mBufferSetBytes(buffer, destBytes, 26);

		int result;
		result = mBufferSetByteSlice(buffer, startPos, copyLen, sourceBytes);
		mBufferFinish(buffer);
}

void mBufferAppendTest() {
    int reps, result;
    reps = int64getArgument(0);
    byte returnDataBuffer[255];
    int randomBufferHandle = mBufferNew();
    int handle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferGetBytes(randomBufferHandle, returnDataBuffer);
        mBufferAppendBytes(handle,returnDataBuffer,reps);
    }
    mBufferFinish(handle);
}

void mBufferToBigIntUnsignedTest() {
    int reps, result, bigIntHandle;
    reps = int64getArgument(0);
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferToBigIntUnsigned(randomBufferHandle,bigIntHandle);
    }
    mBufferFinish(randomBufferHandle);
    bigIntFinishUnsigned(bigIntHandle);
}

void mBufferToBigIntSignedTest() {
    int reps, result, bigIntHandle;
    reps = int64getArgument(0);
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferToBigIntSigned(randomBufferHandle,bigIntHandle);
    }
    mBufferFinish(randomBufferHandle);
    bigIntFinishSigned(bigIntHandle);
}

void mBufferFromBigIntUnsignedTest() {
    int reps, result, bigIntHandle;
    reps = int64getArgument(0);
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferToBigIntUnsigned(randomBufferHandle,bigIntHandle);
        mBufferFromBigIntUnsigned(randomBufferHandle,bigIntHandle);
    }
    mBufferFinish(randomBufferHandle);
    bigIntFinishUnsigned(bigIntHandle);
}

void mBufferFromBigIntSignedTest() {
    int reps, result;
    reps = int64getArgument(0);
    int bigIntHandle;
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferToBigIntSigned(randomBufferHandle,bigIntHandle);
        mBufferFromBigIntSigned(randomBufferHandle,bigIntHandle);
    }
    mBufferFinish(randomBufferHandle);
    bigIntFinishSigned(bigIntHandle);
}

void mBufferStorageStoreTest() {
    int reps, result;
    reps = int64getArgument(0);
    int randomBufferHandle = mBufferNew();
    int randomKeyHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomKeyHandle,5);
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferStorageStore(randomKeyHandle,randomBufferHandle);        
    }
    mBufferFinish(randomBufferHandle);
    mBufferFinish(randomKeyHandle);
}

void mBufferStorageLoadTest() {
    int reps, result;
    reps = int64getArgument(0);
    int randomKeyHandle = mBufferNew();
    int randomBufferHandle = mBufferNew();
    for (int i = 0; i < reps; i++)
    {
        result = mBufferSetRandom(randomKeyHandle, 5);
        result = mBufferSetRandom(randomBufferHandle,reps);
        mBufferStorageStore(randomKeyHandle,randomBufferHandle);
        mBufferStorageLoad(randomKeyHandle,randomBufferHandle);
    }
    mBufferFinish(randomBufferHandle);
    mBufferFinish(randomKeyHandle);    
}

int verifyIfBuffersAreEqual(int handle1, int handle2) {
    byte firstBuffer[255];
    int length1 = mBufferGetLength(handle1);
    byte secondBuffer[255];
    int length2 = mBufferGetLength(handle2);
    if (length1!=length2)
        return 1;

    mBufferGetBytes(handle1, firstBuffer);
    mBufferGetBytes(handle2, secondBuffer);
    return byteArraysAreEqual(firstBuffer,secondBuffer,length1);
}

int verifyBytesMBufferAndBigInt(int bigIntHandle, int mBufferHandle, int isSigned) {
    byte bufferBytes[255];
    byte bigIntBytes[255];
    int mBufferLength;
    int bigIntByteLength;

    mBufferLength = mBufferGetLength(mBufferHandle);
    if (isSigned!=0) bigIntByteLength = bigIntSignedByteLength(bigIntHandle);
    else bigIntByteLength = bigIntUnsignedByteLength(bigIntHandle);

    if (mBufferLength!=bigIntByteLength)
        return 1;

    mBufferGetBytes(mBufferHandle, bufferBytes);
    if (isSigned == 0) { bigIntGetUnsignedBytes(bigIntHandle, bigIntBytes); }
    else { bigIntGetSignedBytes(bigIntHandle, bigIntBytes); }
    
    return byteArraysAreEqual(bufferBytes, bigIntBytes, mBufferLength);
    return 0;
}

int byteArraysAreEqual(byte firstArray[], byte secondArray[], int length) {
    for (int i = 0; i < length; i++)
        if (firstArray[i] != secondArray[i])
            { return 1; }
    return 0;
}
