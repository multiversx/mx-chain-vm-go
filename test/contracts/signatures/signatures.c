#include "../mxvm/context.h"

void panic() {
	byte msg[] = "don't call";
	signalError(msg, 10);
}

void goodFunction() {
	panic();
}

byte wrongReturn() {
	panic();
	return 0;
}

void wrongParams(int param) {
	panic();
	int64finish(param);
}

void* wrongParamsAndReturn(int q, byte *p) {
	panic();
	getArgument(q, p);
	return 0;
}
