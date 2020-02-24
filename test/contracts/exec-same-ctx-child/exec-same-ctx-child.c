#include "../elrond/context.h"

byte childKey[] =  "childKey........................";
byte childData[] = "childData";
byte childFinish[] = "childFinish";

byte recipient[32]     = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,96};

void childFunction() {
	int numArgs = getNumArguments();
	if (numArgs != 2) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

  getArgument(0, recipient);

	byte transferData[100];
	getArgument(1, transferData);
	int dataLength = getArgumentLength(1);
	transferValue(recipient, value, transferData, dataLength);

	storageStore(childKey, childData, 9);
	finish(childFinish, 11);
}
