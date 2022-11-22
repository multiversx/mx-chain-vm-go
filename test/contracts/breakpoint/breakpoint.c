#include "../elrond/context.h"

byte array[32] = {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0};

void testFunc() {
  i64 arg = int64getArgument(0);

  if (arg == 1) {
    byte msg[] = "exit here";
    signalError(msg, 9);
    byte msg2[] = "exit later";
    signalError(msg2, 10);
  }

  if (arg == 2) {
    array[2147483647] = 42;
    int64finish(array[2147483647]);
  }

	int64finish(100);
}
