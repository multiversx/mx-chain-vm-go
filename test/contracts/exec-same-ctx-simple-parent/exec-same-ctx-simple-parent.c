#include "../elrond/context.h"

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,99};

void parentFunctionChildCall() {
	byte childAddress[] = "\0\0\0\0\0\0\0\0\x0f\x0fsecondSC..............";
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
