#include "../elrond/context.h"

void testFunc() {
  i64 arg = int64getArgument(0);

  if (arg == 1) {
    byte msg[] = "exit here";
    signalError(msg, 9);
    byte msg2[] = "exit later";
    signalError(msg2, 10);
  } else {
    int64finish(100);
  }
}

void init() {
}

void _main() {
}
