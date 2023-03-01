#include "../mxvm/context.h"
#include "../mxvm/bigInt.h"
#include "../mxvm/test_utils.h"

byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte parentFinishA[] = "parentFinishA";
byte parentFinishB[] = "parentFinishB";

byte childAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "childSC...............";
byte vaultAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "vaultAddress..........";
byte thirdPartyAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "thirdPartyAddress.....";

byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void handleBehaviorArgument();
void handleTransferToVault();
int mustTransferToVault();
int isVault();

void parentPerformAsyncCall() {

	int numArgs = getNumArguments();
	if (numArgs < 3) {
		byte msg[] = "wrong num of arguments";
		signalError(msg, 22);
	}

	storageStore(parentKeyA, 32, parentDataA, 11);
	storageStore(parentKeyB, 32, parentDataB, 11);
	finish(parentFinishA, 13);
	finish(parentFinishB, 13);

	value[31] = 3;
	byte transferData[] = "hello";
	transferValue(thirdPartyAddress, value, transferData, 5);
	
	// 207468657265 is the word 'there', hex-encoded ASCII
	byte callData[] = "transferToThirdParty@03@207468657265@00";
	callData[38] = int64getArgument(0) + '0';

	byte successCallback[] = "callBackSucc";
	byte errorCallback[] = "callBackErr";

	long long gas = int64getArgument(1);
	long long extraGasForCallback = int64getArgument(2);

	value[31] = 7;
	createAsyncCall(childAddress, value, callData, 39, successCallback, 12, errorCallback, 11, gas, extraGasForCallback);
}

void callBackSucc() {
	int numArgs = getNumArguments();
	if (numArgs < 2) {
		byte msg[] = "wrong num of arguments";
		signalError(msg, 22);
	}

	byte loadedData[11];
	storageLoad(parentKeyB, 32, loadedData);

	int status = 0;
	for (int i = 0; i < 11; i++) {
		if (loadedData[i] != parentDataB[i]) {
			status = 1;
			break;
		}
	}

	handleBehaviorArgument();
	handleTransferToVault();

	finishResult(status);
}

void callBackErr() {
	int numArgs = getNumArguments();
	if (numArgs < 2) {
		byte msg[] = "wrong num of arguments";
		signalError(msg, 22);
	}

	byte loadedData[11];
	storageLoad(parentKeyB, 32, loadedData);

	int status = 0;
	for (int i = 0; i < 11; i++) {
		if (loadedData[i] != parentDataB[i]) {
			status = 1;
			break;
		}
	}

	handleBehaviorArgument();
	handleTransferToVault();

	byte message[] = "succCallbackErr";
	finish(message, 15);
}

void handleTransferToVault() {
	if (mustTransferToVault()) {
		value[31] = 4;

		transferValue(vaultAddress, value, 0, 0);
	}
}

int mustTransferToVault() {
	int numArgs = getNumArguments();
	byte childArgument[10];

	if (numArgs == 3) {
		getArgument(2, childArgument);
		if (isVault(childArgument)) {
			return 0;
		}
	}

	if (numArgs == 4) {
		getArgument(3, childArgument);
		if (isVault(childArgument)) {
			return 0;
		}
	}

	return 1;
}

int isVault(byte *string) {
	byte vault[] = "vault";
	for (int i = 0; i < 5; i++) {
		if (vault[i] != string[i]) {
			return 0;
		}
	}

	return 1;
}

void handleBehaviorArgument() {
	int numArgs = getNumArguments();
	if (numArgs < 4) {
		return;
	}

	byte behavior = int64getArgument(1);

	if (behavior == 3) {
		byte msg[] = "callBack error";
		signalError(msg, 14);
	}
	if (behavior == 4) {
		byte msg[] = "loop";
		while (1) {
			finish(msg, 4);
		}
	}

	finish(&behavior, 1);
}

