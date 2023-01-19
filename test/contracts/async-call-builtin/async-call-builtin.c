#include "../mxvm/context.h"

byte scAddress[] = "\0\0\0\0\0\0\0\0\x0F\x0FparentSC..............";
byte value[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void performAsyncCallToBuiltin() {
	i64 arg = int64getArgument(0);

	byte msg[] = "hello";
	finish(msg, 5);

	if (arg == 1) {
		byte callData[] = "builtinFail";
		asyncCall(scAddress, value, callData, 11);
	}
}

void callBack() {
	i64 returnCode = int64getArgument(0);

	int64finish(returnCode);
}
