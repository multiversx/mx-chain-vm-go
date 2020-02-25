#include "../elrond/context.h"

byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte childKey[] =  "childKey........................";
byte childData[] = "childData";
byte childFinish[] = "childFinish";

byte recipient[32]     = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,96};
byte dataA[14] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte dataB[14] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void not_ok() {
	byte msg[] = "not ok";
	finish(msg, 6);
}

void childFunction() {
	int numArgs = getNumArguments();
	if (numArgs != 2) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

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
	storageLoad(parentKeyA, dataA);
	storageLoad(parentKeyB, dataB);

	finish(dataA, 11);
	finish(dataB, 11);

	dataA[5] = 'D';

	for (int i = 0; i < 6; i++) {
		finish(&dataB[i], 1);
	}
	
	// finish(&dataB[6], 1);

	/*
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
	*/
}
