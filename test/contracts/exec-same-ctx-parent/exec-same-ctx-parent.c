#include "../elrond/context.h"
#include "../elrond/bigInt.h"

byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte parentFinishA[] = "parentFinishA";
byte parentFinishB[] = "parentFinishB";

byte parentTransferReceiver[] = "parentTransferReceiver..........";
byte parentTransferValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,42};
byte parentTransferData[] = "parentTransferData";

byte executeValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,99};
i32 executeArgumentsLengths[] = {32, 6};
byte executeArgumentsData[] = "asdfoottxxwlllllllllllwrraatttttqwerty";

void finishResult(int);

void parentFunctionPrepare() {
	storageStore(parentKeyA, parentDataA, 11);
	storageStore(parentKeyB, parentDataB, 11);
	finish(parentFinishA, 13);
	finish(parentFinishB, 13);
	transferValue(
			parentTransferReceiver,
			parentTransferValue,
			parentTransferData,
			18
	);
}

void parentFunctionWrongCall() {
	parentFunctionPrepare();
	byte childAddress[] = "wrongSC.........................";
	byte functionName[] = "childFunction";

	int result = executeOnSameContext(
			10000,
			childAddress,
			executeValue,
			functionName,
			13,
			2,
			(byte*)executeArgumentsLengths,
			executeArgumentsData
	);
	finishResult(result);
}

void parentFunctionChildCall() {
	parentFunctionPrepare();
	byte childAddress[] = "secondSC........................";
	byte functionName[] = "childFunction";
	int result = executeOnSameContext(
			20000,
			childAddress,
			executeValue,
			functionName,
			13,
			2,
			(byte*)executeArgumentsLengths,
			executeArgumentsData
	);
	finishResult(result);
}

void finishResult(int result) {
	if (result == 0) {
		byte message[] = "success";
		finish(message, 7);
	}
	if (result == 1) {
		byte message[] = "failed";
		finish(message, 6);
	}
	if (result != 0 && result != 1) {
		byte message[] = "unknown result";
		finish(message, 14);
	}
}
