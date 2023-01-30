#include "../mxvm/context.h"

byte msg_ok[] = "ok";
byte msg_not_ok[] = "not ok";
byte msg_unexpected[] = "unexpected";

byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void test_getCallValue_1byte() {
	int length = getCallValue(value);
	if (length != 32) {
		signalError(msg_unexpected, 10);
	}
	if (value[31] == 64) {
		finish(msg_ok, 2);
	} else {
		finish(msg_not_ok, 6);
	}

	finish((byte*)&length, 4);
	finish(value, length);
}

void test_getCallValue_4bytes() {
	int length = getCallValue(value);
	if (length != 32) {
		signalError(msg_unexpected, 10);
	}
	int ok = 0;
	ok = ok + (value[28] == 64);
	ok = ok + (value[29] == 12);
	ok = ok + (value[30] == 16);
	ok = ok + (value[31] == 99);

	if (ok == 4) {
		finish(msg_ok, 2);
	} else {
		finish(msg_not_ok, 6);
	}

	finish((byte*)&length, 4);
	finish(value, length);
}

void test_getCallValue_bigInt_to_Bytes() {
	int length = getCallValue(value);
	if (length != 32) {
		signalError(msg_unexpected, 10);
	}

	int ok = 0;
	ok = ok + (value[30] == 19);
	ok = ok + (value[31] == 233);

	if (ok == 2) {
		finish(msg_ok, 2);
	} else {
		finish(msg_not_ok, 6);
	}

	finish((byte*)&length, 4);
	finish(value, length);

	// Construct the bigInt 12345, on 4 bytes.
	value[28] = 0;
	value[29] = 0;
	value[30] = 48;
	value[31] = 57;
	finish(value, 32);
}

void test_int64getArgument() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		signalError(msg_unexpected, 10);
	}

	i64 argument = int64getArgument(0);

	if (argument == 12345) {
		finish(msg_ok, 2);
	} else {
		finish(msg_not_ok, 6);
	}

	finish((byte*)&argument, 4);
	int64finish(argument);
}
