#include "../elrond/context.h"
#include "../elrond/test_utils.h"

byte initialContractAddress[32] = {};
byte sourceContractAddress[32] = {};
byte newAddress[32] = {};

byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};
byte zeroValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

byte arguments[0] = {};
int argumentsLengths[0] = {};

byte contractMetadata[2] = {3, 0};

void upgradeCodeFromAnotherContract() {
	getArgument(0, initialContractAddress);
	getArgument(1, sourceContractAddress);

	upgradeFromSourceContract(
			initialContractAddress,
			500000,
			deploymentValue,
			sourceContractAddress,
			contractMetadata,
			0,
			(byte*)argumentsLengths,
			arguments);			
}

byte childContractAddress[32] = {};
byte childCode[5000] = {};

void upgradeChildContract() {
	getArgument(0, childContractAddress);
	int codeLen = getArgument(1, childCode);

	upgradeContract(
			childContractAddress,
			500000,
			zeroValue,
			childCode,
			contractMetadata,
			codeLen,
			0,
			(byte*)argumentsLengths,
			arguments);			
}

void dummy() {
	byte msg[] = "dummy text";
	finish(msg, 10);
}
