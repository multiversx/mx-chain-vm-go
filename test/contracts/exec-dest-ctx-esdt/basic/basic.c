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

	byte esdtValue[32] = {0};
	int esdtValueLen = getESDTValue(esdtValue);

	BinaryArgs args = NewBinaryArgs();

	int lastArg = 0;
	lastArg = AddBinaryArg(&args, tokenName, tokenNameLen);
	lastArg = AddBinaryArg(&args, esdtValue, esdtValueLen);
	TrimLeftZeros(&args, lastArg);

	byte arguments[100];
	int argsLen = SerializeBinaryArgs(&args, arguments);
	finish(arguments, argsLen);
	finish(args.lengths, 2);
	finish((byte*)args.lengthsAsI32, 8);

	int result = executeOnDestContext(
			1000000,
			self,
			executeValue,
			ESDTTransfer,
			sizeof ESDTTransfer - 1,
			args.numArgs,
		  (byte*)args.lengthsAsI32,
			args.serialized
	);
}
