#ifndef _ARGS_H_
#define _ARGS_H_

#include "types.h"

const byte Base16Figures[] = {
	'0', '1', '2', '3',
	'4', '5', '6', '7',
	'8', '9', 'A', 'B',
	'C', 'D', 'E', 'F'
};

typedef struct binaryArgs {
	byte *arguments[10];
	int lengths[10];
	int numArgs;
} BinaryArgs;

BinaryArgs NewBinaryArgs();
void AddBinaryArg(BinaryArgs *args, byte *arg, int arglen);
int SerializeBinaryArgs(BinaryArgs *args, byte *result);
void* memset(void *str, int c, unsigned long n);

BinaryArgs NewBinaryArgs() {
	BinaryArgs args;
	args.numArgs = 0;

	return args;
}

void AddBinaryArg(BinaryArgs *args, byte *arg, int arglen) {
	int n = args->numArgs;
	args->arguments[n] = arg;
	args->lengths[n] = arglen;
	args->numArgs = n + 1;
}

int SerializeBinaryArgs(BinaryArgs *args, byte *result) {
	int cursor = 0;
	for (int i = 0; i < args->numArgs; i++) {
		for (int j = 0; j < args->lengths[i]; j++) {
			byte b = args->arguments[i][j];

			result[cursor] = Base16Figures[b / 16];
			cursor += 1;
			result[cursor] = Base16Figures[b % 16];
			cursor += 1;
		}

		if (i < args->numArgs - 1) {
			result[cursor] = '@';
			cursor += 1;
		}
	}

	return cursor;
}

#endif
