#include "../elrond/context.h"

byte msg_ok[] = "ok";
byte msg_not_ok[] = "not ok";
byte msg_unexpected[] = "unexpected";

byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void test_getCallValue_1byte() {
	int length = getCallValue(value);
	if (length != 1) {
		signalError(msg_unexpected, 10);
	}
	if (value[0] == 64) {
		finish(msg_ok, 2);
	} else {
		finish(msg_not_ok, 6);
	}

	finish((byte*)&length, 4);
	finish(value, length);
}

void test_getCallValue_4bytes() {
	int length = getCallValue(value);
	if (length != 4) {
		signalError(msg_unexpected, 10);
	}
	int ok = 0;
	ok = ok + (value[0] == 64);
	ok = ok + (value[1] == 12);
	ok = ok + (value[2] == 16);
	ok = ok + (value[3] == 99);

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
	if (length != 2) {
		signalError(msg_unexpected, 10);
	}

	int ok = 0;
	ok = ok + (value[0] == 19);
	ok = ok + (value[1] == 233);

	if (ok == 2) {
		finish(msg_ok, 2);
	} else {
		finish(msg_not_ok, 6);
	}

	finish((byte*)&length, 4);
	finish(value, length);

	// Construct the bigInt 12345, on 4 bytes.
	value[0] = 0;
	value[1] = 0;
	value[2] = 48;
	value[3] = 57;
	finish(value, 4);
}

void test_getCallValue_Int64Argument() {
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
}
