#include "../elrond/context.h"

byte key[] = "test_key";

void delete_existing_storage() {
	storageStore(key, sizeof(key) - 1, 0, 0);
}
