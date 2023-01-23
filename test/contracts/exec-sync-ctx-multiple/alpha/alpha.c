#include "../../mxvm/context.h"
#include "../../mxvm/bigInt.h"
#include "../../mxvm/test_utils.h"

byte betaSC[]  = "\0\0\0\0\0\0\0\0\x0F\x0F" "betaSC................";
byte gammaSC[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "gammaSC...............";
byte deltaSC[] = "\0\0\0\0\0\0\0\0\x0F\x0F" "deltaSC...............";

byte betaMethod[] = "betaMethod";
byte gammaMethod[] = "gammaMethod";
byte deltaMethod[] = "deltaMethod";

byte callValue[] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,12};

void callChildrenDirectly_DestCtx() {
  u32 argumentLengths[] = {4};
  byte argumentData[] = "argx";
	int result = 0;

  argumentData[3] = '1';
  result = executeOnDestContext(
			10000,
			betaSC,
			callValue,
			betaMethod,
			10,
			1,
			(byte*)argumentLengths,
			argumentData
	);
	finishResult(result);

  argumentData[3] = '2';
  result = executeOnDestContext(
			10000,
			gammaSC,
			callValue,
			gammaMethod,
			11,
			1,
			(byte*)argumentLengths,
			argumentData
	);
	finishResult(result);

  argumentData[3] = '3';
  result = executeOnDestContext(
			10000,
			deltaSC,
			callValue,
			deltaMethod,
			11,
			1,
			(byte*)argumentLengths,
			argumentData
	);
	finishResult(result);
}

void callChildrenDirectly_SameCtx() {
  u32 argumentLengths[] = {4};
  byte argumentData[] = "argx";
	int result = 0;

  argumentData[3] = '1';
  result = executeOnSameContext(
			10000,
			betaSC,
			callValue,
			betaMethod,
			10,
			1,
			(byte*)argumentLengths,
			argumentData
	);
	finishResult(result);

  argumentData[3] = '2';
  result = executeOnSameContext(
			10000,
			gammaSC,
			callValue,
			gammaMethod,
			11,
			1,
			(byte*)argumentLengths,
			argumentData
	);
	finishResult(result);

  argumentData[3] = '3';
  result = executeOnSameContext(
			10000,
			deltaSC,
			callValue,
			deltaMethod,
			11,
			1,
			(byte*)argumentLengths,
			argumentData
	);
	finishResult(result);
}
