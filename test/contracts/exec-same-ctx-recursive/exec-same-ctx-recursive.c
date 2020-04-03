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

void callRecursive() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	byte iteration = (byte) int64getArgument(0);

  // Add iteration number to finish()
  byte finishIter[] = "finishNNN";
  intTo3String(iteration, finishIter, 6);
  finish(finishIter, 9);

  // Add iteration number to storage
  byte keyIter[] = "keyNNN..........................";
  byte valueIter[] = "valueNNN";
  intTo3String(iteration, keyIter, 3);
  intTo3String(iteration, valueIter, 5);
  storageStore(keyIter, valueIter, 8);

  // Increment a stored counter
  incrementIterCounter();

  // Increment a bigInt counter (possible because of executeOnSameContext())
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
