#include "../mxvm/context.h"

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,99};
byte childSC[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "childSC...............";

void parentFunctionChildCall() {
	byte* childAddress = childSC;
	byte functionName[] = "childFunction";

	u64 result = executeOnSameContext(
			200000,
			childAddress,
			executeValue,
			functionName,
			13,
			0,
			0,
			0
	);
	int64finish(result);

	result = executeOnSameContext(
			200000,
			childAddress,
			executeValue,
			functionName,
			13,
			0,
			0,
			0
	);
	int64finish(result);
	
	byte msg[] = "parent";
	finish(msg, 6);
}
