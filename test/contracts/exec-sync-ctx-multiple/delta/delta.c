#include "../../mxvm/context.h"

void deltaMethod() {
  byte arg[4] = {0};
  getArgument(0, arg);
  finish(arg, 4);
}
