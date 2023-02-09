#include "../mxvm/context.h"
#include "../mxvm/bigInt.h"
#include "../mxvm/test_utils.h"

byte childKey[] =  "childKey........................";
byte childData[] = "childData";
byte childFinish[] = "childFinish";
byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

byte vaultAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0FvaultAddress..........";
byte thirdPartyAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0FthirdPartyAddress.....";

int sendToThirdParty();
int sendToVault();
byte getValueToSend();
void handleBehaviorArgument();

void transferToThirdParty() {
	int numArgs = getNumArguments();
	if (numArgs != 3) {
		byte msg[] = "wrong num of arguments";
		signalError(msg, 22);
	}

	handleBehaviorArgument();

	int result = 0;

	result = sendToThirdParty();
	if (result == 0) {
		byte msg[] = "thirdparty";
		finish(msg, 10);
	}

	result = sendToVault();
	if (result == 0) {
		byte msg[] = "vault";
		finish(msg, 5);
	}

	storageStore(childKey, 32, childData, 9);
}

void handleBehaviorArgument() {
	int numArgs = getNumArguments();
	if (numArgs < 3) {
		return;
	}

	byte behavior = int64getArgument(2);

	if (behavior == 1) {
		byte msg[] = "child error";
		signalError(msg, 11);
	}
	if (behavior == 2) {
		byte msg[] = "loop";
		while (1) {
			finish(msg, 4);
		}
	}

	finish(&behavior, 1);
}

int sendToThirdParty() {
	value[31] = getValueToSend();

	byte data[100];
	int len = getArgument(1, data);

	return transferValue(thirdPartyAddress, value, data, len);
}

int sendToVault() {
	value[31] = 4;

	return transferValue(vaultAddress, value, 0, 0);
}

byte getValueToSend() {
	int len;

	len = getArgumentLength(0);
	if (len != 1) {
		byte msg[] = "wrong argument size";
		signalError(msg, 19);
	}

	byte valueToSend = 0;
	getArgument(0, &valueToSend);

	return valueToSend;
}
