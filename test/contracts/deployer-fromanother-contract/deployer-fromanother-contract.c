#include "../elrond/context.h"
#include "../elrond/test_utils.h"

byte sourceContractAddress[32] = {};
byte newAddress[32] = {};

byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};

byte arguments[] = {0, 0};
int argumentsLengths[] = {1, 1};

byte contractMetadata[2] = {3, 0};


void deployCodeFromAnotherContract() {
	getArgument(0, sourceContractAddress);

	int result = deployFromSourceContract(
			2000,
			deploymentValue,
			sourceContractAddress,
			contractMetadata,
			newAddress,
			2,
			(byte*)argumentsLengths,
			arguments);

	finishResult(result);
}