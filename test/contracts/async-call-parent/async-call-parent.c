#include "../elrond/context.h"
#include "../elrond/bigInt.h"
#include "../elrond/test_utils.h"

byte parentKeyA[] =  "parentKeyA......................";
byte parentDataA[] = "parentDataA";
byte parentKeyB[] =  "parentKeyB......................";
byte parentDataB[] = "parentDataB";
byte parentFinishA[] = "parentFinishA";
byte parentFinishB[] = "parentFinishB";

byte childAddress[] = "childSC.........................";
byte vaultAddress[] = "vaultAddress....................";
byte thirdPartyAddress[] = "thirdPartyAddress...............";

byte value[32] = {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0};

void parentPerformAsyncCall() {
	storageStore(parentKeyA, parentDataA, 11);
	storageStore(parentKeyB, parentDataB, 11);
	finish(parentFinishA, 13);
	finish(parentFinishB, 13);

	value[31] = 3;
	byte transferData[] = "hello";
	transferValue(thirdPartyAddress, value, transferData, 5);
	
	byte callData[] = "transferToThirdParty@3@207468657265";
	value[31] = 7;
	asyncCall(childAddress, value, callData, 35);
}

void callBack() {
}
