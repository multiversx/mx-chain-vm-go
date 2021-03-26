#ifndef _CBUILTINS_H_
#define _CBUILTINS_H_

#include "types.h"

void* memset(void *dest, int c, unsigned long n) {
	byte v = (byte)c;
	for (unsigned long i = 0; i < n; i++) {
		((byte*)dest)[i] = v;
	}
	return dest;
}

#endif
