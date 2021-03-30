#include "../../elrond/context.h"
#include "../../elrond/test_utils.h"
#include "../../elrond/args.h"

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte self[32] = "\0\0\0\0\0\0\0\0\x0f\x0f" "parentSC..............";
byte vaultSC[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "vaultSC...............";
byte ESDTTransfer[] = "ESDTTransfer";

void basic_transfer() {
	byte tokenName[265] = {0};
	int tokenNameLen = getESDTTokenName(tokenName);
	finish(tokenName, tokenNameLen);

	byte callValue[32] = {0};
	int callValueLen = getCallValue(callValue);

	BinaryArgs args = NewBinaryArgs();
	AddBinaryArg(&args, tokenName, tokenNameLen);
	AddBinaryArg(&args, callValue, callValueLen);

	byte arguments[100];
	int argsLen = SerializeBinaryArgs(&args, arguments);
	finish(arguments, argsLen);

	int result = executeOnDestContext(
			1000000,
			self,
			callValue,
			ESDTTransfer,
			sizeof ESDTTransfer - 1,
			args.numArgs,
			args.lengths,
			args.serialized
	);
}
