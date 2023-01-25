#include "../../mxvm/context.h"

void betaMethod() {
  byte arg[4] = {0};
  getArgument(0, arg);
  finish(arg, 4);
}
