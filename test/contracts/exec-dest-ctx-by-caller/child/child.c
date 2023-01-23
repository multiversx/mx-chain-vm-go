#include "../../mxvm/context.h"
#include "../../mxvm/test_utils.h"


void give() {
	int numArgs = getNumArguments();
	if (numArgs != 1) {
		byte message[] = "wrong number of arguments";
		signalError(message, 25);
	}

  byte value_to_give = 0;
  getArgument(0, &value_to_give);

  byte caller[32] = {0};
  getCaller(caller);

	byte value[32] = {0};
	value[31] = value_to_give;

  transferValue(caller, value, 0, 0);

	byte msg[] = "sent";
	finish(msg, 4);
}
