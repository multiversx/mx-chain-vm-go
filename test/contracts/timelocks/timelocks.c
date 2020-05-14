#include "../elrond/context.h"

byte counterKey[] = "counter";

void init() {
    int64storageStore(counterKey, sizeof(counterKey), 0);
}

void incrementCounter() {
    int counterKeySize = sizeof(counterKey);
    i64 counter = int64storageLoad(counterKey, counterKeySize);
    if (isStorageLocked(counterKey, counterKeySize) == 0) {
        counter++;
        int64storageStore(counterKey, counterKeySize, counter);
    }
    int64finish(counter);
}

void lockCounter() {
    long long lockTimestamp = getBlockTimestamp();
    lockTimestamp += 3600*24;
    setStorageLock(counterKey, sizeof(counterKey), lockTimestamp);
}

void releaseCounter() {
    setStorageLock(counterKey, sizeof(counterKey), 0);
}
