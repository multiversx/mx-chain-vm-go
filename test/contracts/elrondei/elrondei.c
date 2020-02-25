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
