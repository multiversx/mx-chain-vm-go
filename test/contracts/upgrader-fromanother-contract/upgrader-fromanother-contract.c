#include "../mxvm/context.h"
#include "../mxvm/test_utils.h"

byte initialContractAddress[32] = {};
byte sourceContractAddress[32] = {};
byte newAddress[32] = {};

byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};

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
