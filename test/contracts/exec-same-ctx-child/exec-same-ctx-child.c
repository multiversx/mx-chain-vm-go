#include "../elrond/context.h"
#include "../elrond/bigInt.h"

byte dataA[20] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte dataB[20] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte childKey[] =  "childKey........................";
byte childData[] = "childData";
byte childFinish[] = "childFinish";

byte recipient[32]     = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,96};

void not_ok() {
	byte msg[] = "not ok";
	finish(msg, 6);
}

void didCallerPay() {
	bigInt bigInt_payment = bigIntNew(0);
	bigIntGetCallValue(bigInt_payment);

	long long payment = bigIntGetInt64(bigInt_payment);
	if (payment != 99) {
		byte message[] = "child execution requires tx value of 99";
		signalError(message, 39);
	}
}

void childFunction() {
	int numArgs = getNumArguments();
	if (numArgs != 2) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	didCallerPay();

	// This transfer will appear alongside the transfers made by the parent.
  getArgument(0, recipient);
	byte transferData[100];
	getArgument(1, transferData);
	int dataLength = getArgumentLength(1);
	transferValue(recipient, value, transferData, dataLength);

	// This storage update will appear alongside the storage updates made by the parent.
	storageStore(childKey, childData, 9);

	// This finish value will appear alongside the finish values set by the parent.
	finish(childFinish, 11);

	// The child has access to the storage of the parent.
	int lenA = storageGetValueLength(parentKeyA);
	if (lenA != 11) {
		byte err[] = "err lenA";
		finish(err, 8);
		int64finish(lenA);
		not_ok();
		return;
	}
	int lenB = storageGetValueLength(parentKeyB);
	if (lenB != 11) {
		byte err[] = "err lenB";
		finish(err, 8);
		not_ok();
		return;
	}
	u64 slLenA = storageLoad(parentKeyA, dataA);
	u64 slLenB = storageLoad(parentKeyB, dataB);

	finish(dataA, 11);

	for (int i = 0; i < 11; i++) {
		finish(&dataA[i], 1);
	}

	finish(dataB, 11);

	for (int i = 0; i < 11; i++) {
		finish(&dataB[i], 1);
	}
	
	int i;
	int status = 0;
	for (i = 0; i < 11; i++) {
		if (dataA[i] != parentDataA[i]) {
			status = 1;
			break;
		}
		if (dataB[i] != parentDataB[i]) {
			status = 2;
			break;
		}
	}

	if (status == 0) {
		byte msg[] = "child ok";
		finish(msg, 8);
	}
}

void childFunction_BigInts() {
	int numArgs = getNumArguments();
	if (numArgs != 3) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	didCallerPay();

	int status = 0;

	bigInt intA = int64getArgument(0);
	bigInt intB = int64getArgument(1);
	bigInt intC = int64getArgument(2);

	long long a = bigIntGetInt64(intA);
	long long b = bigIntGetInt64(intB);
	long long c = bigIntGetInt64(intC);

	if (a != 84) {
		not_ok();
		int64finish(a);
		status = 1;
	}
	if (b != 96) {
		not_ok();
		int64finish(b);
		status = 1;
	}
	if (c != 1024) {
		not_ok();
		int64finish(c);
		status = 1;
	}

	bigInt intX = bigIntNew(256);
	if (intX != 4) {
		not_ok();
		int64finish(intX);
		status = 1;
	}

	if (status == 0) {
		byte msg[] = "child ok";
		finish(msg, 8);
	}
}
