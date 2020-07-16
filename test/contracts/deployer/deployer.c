#include "../elrond/context.h"
#include "../elrond/test_utils.h"

byte contractCode[5000] = {};
byte contractID = 0;
byte newAddress[32] = {};

byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};

void preverifyDeployment() {
	getArgument(0, &contractID);
	int codeLength = storageLoadLength(&contractID, 1);
	int loadedLength = storageLoad(&contractID, 1, contractCode);
	int64finish(loadedLength);
}

void deployChildContract() {
	getArgument(0, &contractID);
	int loadedLength = storageLoad(&contractID, 1, contractCode);
	int64finish(loadedLength);

	byte initArgLengths[] = {1};
	byte initArgs[] = {0};
	int result = createContract(
			deploymentValue,
			contractCode,
			loadedLength,
			newAddress,
			1,
			initArgLengths,
			initArgs);

	finishResult(result);
}
