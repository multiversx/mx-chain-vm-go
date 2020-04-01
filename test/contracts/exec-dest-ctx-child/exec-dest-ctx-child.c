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
	if (payment != 99) {
		byte message[] = "child execution requires tx value of 99";
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

	if (a != 0) {
		not_ok();
		int64finish(a);
	}
	if (b != 0) {
		not_ok();
		int64finish(b);
	}
	if (c != 0) {
		not_ok();
		int64finish(c);
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
