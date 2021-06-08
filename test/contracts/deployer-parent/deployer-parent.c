#include "../elrond/context.h"
#include "../elrond/test_utils.h"
#include "../elrond/args.h"

byte parentAddress[32] = {};

byte childAddress[32] = {};
byte childCode[5000] = {};
byte childMetadata[2] = {1, 0};

byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};

void init() {
	getArgument(0, parentAddress);
	getArgument(1, childAddress);
	getArgument(2, childCode);
	bigInt codeSize = bigIntNew(0);
  	bigIntGetUnsignedArgument(3, codeSize);

	int isSelfContract = isSmartContract(parentAddress);
	if (isSelfContract == 0) {
		byte message[] = "parent not a contract";
		signalError(message, sizeof(message) - 1);
	}

	BinaryArgs args = NewBinaryArgs();

	int lastArg = 0;
	lastArg = AddBinaryArg(&args, parentAddress, 32);
	TrimLeftZeros(&args, lastArg);

	byte arguments[100];
	int argsLen = SerializeBinaryArgs(&args, arguments);

	int result = createContract(
			1000,
			deploymentValue,
			childCode,
			childMetadata,
			codeSize,
			childAddress,
			lastArg,
			args.lengths,
			args.serialized);

	finishResult(result);
}
