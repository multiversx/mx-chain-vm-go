#ifndef _ARGS_H_
#define _ARGS_H_

#include "types.h"

typedef struct binaryArgs {
	byte *arguments[10];
	byte lengths[10];
	int numArgs;
	byte *serialized;
	int lenSerialized;
} BinaryArgs;

BinaryArgs NewBinaryArgs();
void AddBinaryArg(BinaryArgs *args, byte *arg, int arglen);
int SerializeBinaryArgs(BinaryArgs *args, byte *result);
void* memset(void *str, int c, unsigned long n);

#endif
