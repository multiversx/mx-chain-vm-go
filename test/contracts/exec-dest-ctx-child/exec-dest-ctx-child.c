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

byte recipient[32]     = "childTransferReceiver...........";
byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,12};

void not_ok() {
	byte msg[] = "not ok";
	finish(msg, 6);
}

void didCallerPay() {
	bigInt bigInt_payment;
	bigIntGetCallValue(bigInt_payment);

	long long payment = bigIntGetInt64(bigInt_payment);
	if (payment != 33) {
		byte message[] = "child execution requires tx value of 33";
		signalError(message, 39);
	}
}

void childFunction() {
	int numArgs = getNumArguments();
	if (numArgs != 3) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

	didCallerPay();

	// This transfer will appear alongside the transfers made by the parent.
	byte transferData[100];
	getArgument(1, transferData);
	int dataLength = getArgumentLength(1);
	int result = transferValue(recipient, value, transferData, dataLength);
	if (result != 0) {
		not_ok();
	}
	
	// This storage update will appear separate from the storage updates made by the parent.
	storageStore(childKey, childData, 9);

	// This finish value will appear alongside the finish values set by the parent.
	finish(childFinish, 11);
}
