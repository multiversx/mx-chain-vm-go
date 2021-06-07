#include "../elrond/context.h"
#include "../elrond/test_utils.h"

byte alphaAddress[32] = {};
byte contractCode[5000] = {};

// byte contractMetadata[2] = {1, 0};
// byte contractID = 0;
// byte newAddress[32] = {};

// byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};

// byte arguments[1] = {0};
// int argumentsLengths[1] = {1};

void init() {
	// byte arg = 0;
	// getArgument(0, &contractCode);

	byte arg = 0;
	getArgument(0, &alphaAddress);


	int isSelfContract = isSmartContract(alphaAddress);
	if (isSelfContract == 0) {
		byte message[] = "alpha not a contract";
		signalError(message, sizeof(message));
	}

	// int result = createContract(
	// 		2000,
	// 		deploymentValue,
	// 		contractCode,
	// 		contractMetadata,
	// 		loadedLength,
	// 		newAddress,
	// 		1,
	// 		(byte*)argumentsLengths,
	// 		arguments);

	// finishResult(result);
}
