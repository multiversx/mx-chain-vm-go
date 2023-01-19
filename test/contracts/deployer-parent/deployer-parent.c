#include "../mxvm/context.h"
#include "../mxvm/test_utils.h"
#include "../mxvm/args.h"

byte parentAddress[32] = {};

byte childGeneratedAddress[32] = {};
byte childCode[5000] = {};
byte childMetadata[2] = {1, 0};

byte deploymentValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};

void init() {
	int parentAddressSize = getArgument(0, parentAddress);
	int codeSize = getArgument(1, childCode);	

	int isSelfContract = isSmartContract(parentAddress);
	if (isSelfContract == 0) {
		byte message[] = "parent not a contract";
		signalError(message, sizeof(message) - 1);
	}

	BinaryArgs args = NewBinaryArgs();

	int lastArg = 0;
	lastArg = AddBinaryArg(&args, parentAddress, parentAddressSize);

	byte arguments[100];
	int argsLen = SerializeBinaryArgs(&args, arguments);

	int result = createContract(
			50000,
			deploymentValue,
			childCode,
			childMetadata,
			codeSize,
			childGeneratedAddress,
			lastArg + 1,
			(byte*)args.lengthsAsI32,
			args.serialized);

	int isChildContract = isSmartContract(childGeneratedAddress);
	if (isChildContract == 0) {
		byte message[] = "child not a contract";
		signalError(message, sizeof(message) - 1);
	}

	finishResult(result);
}
