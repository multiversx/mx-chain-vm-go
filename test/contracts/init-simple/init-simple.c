#include "../elrond/context.h"

void init() {
	int64finish(42);
}

byte finishMsg[10] = "finish0000";

void dummy() {
	finish(finishMsg, 10);
}
