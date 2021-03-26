#include "../../elrond/context.h"
#include "../../elrond/test_utils.h"
#include "../../elrond/args.h"
#include "../../elrond/cbuiltins.h"

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};
byte vaultSC[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "vaultSC...............";


void basic_transfer() {
	byte tokenName[265] = {0};
	int tokenNameLen = getESDTTokenName(tokenName);
	finish(tokenName, tokenNameLen);

	bigInt callValue = bigIntNew(0);
	bigIntGetESDTCallValue(callValue);
	bigIntFinishUnsigned(callValue);

	BinaryArgs args = NewBinaryArgs();
	byte arg1[] = {1, 2, 3, 4};
	byte arg2[] = "hello";
	AddBinaryArg(&args, arg1, 4);
	AddBinaryArg(&args, arg2, 5);

	byte arguments[100];
	int argsLen = SerializeBinaryArgs(&args, arguments);
	finish(arguments, argsLen);
}
