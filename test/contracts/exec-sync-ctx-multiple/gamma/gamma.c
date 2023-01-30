#include "../../mxvm/context.h"

void gammaMethod() {
  byte arg[4] = {0};
  getArgument(0, arg);
  finish(arg, 4);
}
