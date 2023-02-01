#include "../mxvm/context.h"
#include "../mxvm/bigInt.h"

void managedAsyncCall(int addressBuffer, int valueBuffer, int functionBuffer, int argumentsBuffer);

void foo()
{
   int contractB = mBufferNew();
   bigInt value = bigIntNew(0);
   int function = mBufferNewFromBytes("bar", 3);
   int noArguments = mBufferNew();

   mBufferGetArgument(0, contractB);
   managedAsyncCall(contractB, value, function, noArguments);
}

void callBack()
{
   int numArguments = getNumArguments();

   int64finish(0xCA11BAC3);

   for (int i = 0; i < numArguments; i++)
   {
      int dump = mBufferNew();
      mBufferGetArgument(i, dump);
      mBufferFinish(dump);
   }

   int64finish(0xCA11BAC3);
}
