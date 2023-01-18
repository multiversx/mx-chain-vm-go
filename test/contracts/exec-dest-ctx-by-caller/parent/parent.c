#include "../../mxvm/context.h"
#include "../../mxvm/test_utils.h"

byte childSC[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "childSC...............";

void call_child() {
	int numArgs = getNumArguments();
	if (numArgs != 0) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}


  byte caller[32] = {0};
  getCaller(caller);

	byte value[32] = {0};
	byte function[] = "give";
	i32 argLengths[1] = {1};

  byte value_to_give = 42;
  executeOnDestContextByCaller(
      1000,
      childSC,
			value,
			function,
			4,
			1,
			(byte*)argLengths,
			&value_to_give
  );

	byte msg[] = "child called";
	finish(msg, 12);
}
