#include "../elrond/context.h"

// No imports provided on purpose.
// We are using it in a wasmer instance test, in isolation, with no access to the elrondapi package.
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
