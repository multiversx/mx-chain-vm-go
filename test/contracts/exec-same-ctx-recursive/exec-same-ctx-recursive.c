#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/types.h"

byte selfAddress[] = "parentSC........................";
byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte functionName[] = "callRecursive";
byte arguments[1] = {0};
int argumentsLengths[1] = {1};

byte recursiveIterationCounterKey[] = "recursiveIterationCounter.......";
byte recursiveIterationBigCounterKey[] = "recursiveIterationBigCounter....";
bigInt bigIntIterationCounterID = 16;

void intTo3String(int value, byte *string, int startPos);
void finishResult(int result);
void incrementIterCounter();
void incrementBigIntCounter();
void storeIterationNumber(byte iteration, byte prefix);
void finishIterationNumber(byte iteration, byte prefix);

void recursiveMethodA();

void callRecursive() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);

  finishIterationNumber(iteration, 'R');
  storeIterationNumber(iteration, 'R');
  incrementIterCounter();
  incrementBigIntCounter();

  // Run next iteration.
  if (iteration > 0) {
    arguments[0] = iteration - 1;
    int result = executeOnSameContext(
        10000,
        selfAddress,
        executeValue,
        functionName,
        13,
        1,
        (byte*)argumentsLengths,
        arguments
    );

    finishResult(result);
  } else {
    bigIntStorageStoreUnsigned(recursiveIterationBigCounterKey, bigIntIterationCounterID);
  }
}

void callRecursiveMutualMethods() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);
  if (iteration < 2) {
		byte message[] = "need number of recursive calls >= 2";
		signalError(message, 25);
  }

  byte startMsg[] = "start recursive mutual calls";
  finish(startMsg, 28);

  byte endMsg[] = "end recursive mutual calls";
}

void recursiveMethodA() {
	byte iteration = (byte) int64getArgument(0);

  finishIterationNumber(iteration, 'A');
  storeIterationNumber(iteration, 'A');
  incrementIterCounter();
  incrementBigIntCounter();

  byte functionNameB[] = "recursiveMethodB";
  if (iteration > 0) {
    arguments[0] = iteration - 1;
    int result = executeOnSameContext(
        10000,
        selfAddress,
        executeValue,
        functionNameB,
        16,
        1,
        (byte*)argumentsLengths,
        arguments
    );
    finishResult(result);
  }
}

void storeIterationNumber(byte iteration, byte prefix) {
  byte keyIter[] = "XkeyNNN.........................";
  byte valueIter[] = "XvalueNNN";
  intTo3String(iteration, keyIter, 4);
  intTo3String(iteration, valueIter, 6);
  keyIter[0] = prefix;
  valueIter[0] = prefix;
  storageStore(keyIter, valueIter, 9);
}

void finishIterationNumber(byte iteration, byte prefix) {
  byte finishIter[] = "XfinishNNN";
  intTo3String(iteration, finishIter, 7);
  finishIter[0] = prefix;
  finish(finishIter, 10);
}

void intTo3String(int value, byte *string, int startPos) {
  string[startPos + 2] = (byte)('0' + value % 10);
  string[startPos + 1] = (byte)('0' + (value / 10) % 10);
  string[startPos + 0] = (byte)('0' + (value / 100) % 10);
}

void incrementIterCounter() {
  byte counterValue;
  int len = storageGetValueLength(recursiveIterationCounterKey);
  if (len == 0) {
    counterValue = 1;
    storageStore(recursiveIterationCounterKey, &counterValue, 1);
  } else {
    storageLoad(recursiveIterationCounterKey, &counterValue); 
    counterValue = counterValue + 1;
    storageStore(recursiveIterationCounterKey, &counterValue, 1);
  }
}

void incrementBigIntCounter() {
  bigIntSetInt64(42, 1);
  bigIntAdd(bigIntIterationCounterID, bigIntIterationCounterID, 42);
}

void finishResult(int result) {
	if (result == 0) {
		byte message[] = "succ";
		finish(message, 4);
	}
	if (result == 1) {
		byte message[] = "fail";
		finish(message, 4);
	}
	if (result != 0 && result != 1) {
		byte message[] = "unkn";
		finish(message, 4);
	}
}
