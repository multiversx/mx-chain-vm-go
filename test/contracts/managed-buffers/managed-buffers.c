#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"
#include "../elrond/test_utils.h"
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
    mBufferAppend(mBufferHandle2,mBuffer3, sizeof(mBuffer3));
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
    int randomBufferHandle = mBufferNewRandom(100);
    if (mBufferGetLength(randomBufferHandle) != 100) ok = 1;
    
    finishResult(ok);
}

void mBufferNewTest() {
    int reps;
    int64getArgument(reps);
    for (int i = 0; i < reps; i++)
    {
        int handle = mBufferNew();
    }
}

void mBufferNewFromBytesTest() {
    int reps, lengthOfBuffer;
    int64getArgument(reps);
    int64getArgument(lengthOfBuffer);
    for (int i = 0; i < reps; i++)
    {
        int handle = mBufferNewFromBytes(mBuffer2,lengthOfBuffer);
    }
}

void mBufferNewRandomTest() {
    int reps;
    int64getArgument(reps);
    for (int i = 0; i < reps; i++)
    {
        int handle = mBufferNewRandom(i+1);
    }
}

void mBufferGetLengthTest() {
    int reps;
    int64getArgument(reps);
    for (int i = 0; i < reps; i++)
    {
        int handle = mBufferNewRandom(i+1);
        int length = mBufferGetLength(handle);
    }
}

void mBufferGetBytesTest() {
    int reps;
    int64getArgument(reps);
    byte returnDataBuffer[255];
    for (int i = 0; i < reps; i++)
    {
        int handle = mBufferNewRandom(i+1);
        mBufferGetBytes(handle, returnDataBuffer);
    }  
}

void mBufferAppendTest() {
    int reps;
    int64getArgument(reps);
    byte returnDataBuffer[255];
    int mBufferHandle = mBufferNew(); 
    for (int i = 0; i < reps; i++)
    {
        int handle = mBufferNewRandom(i+1);
        mBufferGetBytes(handle, returnDataBuffer);
        mBufferAppend(mBufferHandle,returnDataBuffer,i+1);
    }
}

void mBufferToBigIntUnsignedTest() {
    int reps;
    int64getArgument(reps);
    int bigIntHandle;
    for (int i = 0; i < reps; i++)
    {
        int mBufferHandle = mBufferNewRandom(i+1);
        mBufferToBigIntUnsigned(mBufferHandle,bigIntHandle);
    }
}

void mBufferToBigIntSignedTest() {
    int reps;
    int64getArgument(reps);
    int bigIntHandle;
    for (int i = 0; i < reps; i++)
    {
        int mBufferHandle = mBufferNewRandom(i+1);
        mBufferToBigIntSigned(mBufferHandle,bigIntHandle);
    }
}

void mBufferFromBigIntUnsignedTest() {
    int reps;
    int64getArgument(reps);
    int bigIntHandle;
    for (int i = 0; i < reps; i++)
    {
        int mBufferHandle = mBufferNewRandom(i+1);
        mBufferToBigIntUnsigned(mBufferHandle,bigIntHandle);
        mBufferFromBigIntUnsigned(mBufferHandle,bigIntHandle);
    }
}

void mBufferFromBigIntSignedTest() {
    int reps;
    int64getArgument(reps);
    int bigIntHandle;
    for (int i = 0; i < reps; i++)
    {
        int mBufferHandle = mBufferNewRandom(i+1);
        mBufferToBigIntSigned(mBufferHandle,bigIntHandle);
        mBufferFromBigIntSigned(mBufferHandle,bigIntHandle);
    }
}

void mBufferStorageStoreTest() {
    int reps;
    int64getArgument(reps);
    for (int i = 0; i < reps; i++)
    {
        int keyHandle = mBufferNewRandom(5);
        int mBufferHandle = mBufferNewRandom(i+1);
        mBufferStorageStore(keyHandle,mBufferHandle);
    }  
}

void mBufferStorageLoadTest() {
    int reps;
    int64getArgument(reps);
    for (int i = 0; i < reps; i++)
    {
        int keyHandle = mBufferNewRandom(5);
        int mBufferHandle = mBufferNewRandom(i+1);
        mBufferStorageLoad(keyHandle,mBufferHandle);
    }  
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
