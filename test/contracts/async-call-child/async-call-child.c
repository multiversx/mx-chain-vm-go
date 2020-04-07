#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/test_utils.h"

byte childKey[] =  "childKey........................";
byte childData[] = "childData";
byte childFinish[] = "childFinish";
byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

byte vaultAddress[] = "vaultAddress....................";
byte thirdPartyAddress[] = "thirdPartyAddress...............";

int sendToThirdParty();
int sendToVault();
byte getValueToSend();

void transferToThirdParty() {
	int numArgs = getNumArguments();
	if (numArgs != 2) {
		byte msg[] = "wrong num of arguments";
		signalError(msg, 22);
	}

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

	storageStore(childKey, childData, 9);
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
