#include "../elrond/context.h"

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,99};

void parentFunctionChildCall() {
	byte childAddress[] = "secondSC........................";
	byte functionName[] = "childFunction";

	i64 result = executeOnSameContext(
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
