#include "../elrond/context.h"

byte counterKey[] = "counter";

void init() {
    int64storageStore(counterKey, 0);
}

void incrementCounter() {
    i64 counter = int64storageLoad(counterKey);
    if (isStorageLocked(counterKey) == 0) {
        counter++;
        int64storageStore(counterKey, counter);
    }
    int64finish(counter);
}

void lockCounter() {
    long long lockTimestamp = getBlockTimestamp();
    lockTimestamp += 3600*24;
    setStorageLock(counterKey, lockTimestamp);
}

void releaseCounter() {
    clearStorageLock(counterKey);
}
