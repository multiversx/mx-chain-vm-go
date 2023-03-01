#include "../mxvm/types.h"

// No imports provided on purpose.

// It is intended exclusively to test the arity checker.

void goodFunction() {
}

byte wrongReturn() {
	return 0;
}

void wrongParams(int param) {
}

void* wrongParamsAndReturn(int q, byte *p) {
	return 0;
}
