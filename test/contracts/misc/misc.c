#include "../mxvm/context.h"

byte array[] = "this is some random string of bytes";
i64 arrayLength = 35;

void iterate_over_byte_array() {
	finish(array, arrayLength);
	finish((byte*)&arrayLength, 1);

	for (int i = 0; i < arrayLength; i++) {
		finish(&array[i], 1);
	}
}
