#include "args.h"

const byte Base16Figures[] = {
	'0', '1', '2', '3',
	'4', '5', '6', '7',
	'8', '9', 'A', 'B',
	'C', 'D', 'E', 'F'
};

BinaryArgs NewBinaryArgs() {
	BinaryArgs args;
	args.numArgs = 0;
	args.serialized = 0;
	args.lenSerialized = 0;

	return args;
}

void AddBinaryArg(BinaryArgs *args, byte *arg, int arglen) {
	int n = args->numArgs;
	args->arguments[n] = arg;
	args->lengths[n] = arglen;
	args->numArgs = n + 1;
}

int SerializeBinaryArgs(BinaryArgs *args, byte *serializedBuffer) {
	int cursor = 0;
	for (int i = 0; i < args->numArgs; i++) {
		for (int j = 0; j < args->lengths[i]; j++) {
			byte b = args->arguments[i][j];

			serializedBuffer[cursor] = Base16Figures[b / 16];
			cursor += 1;
			serializedBuffer[cursor] = Base16Figures[b % 16];
			cursor += 1;
		}

		if (i < args->numArgs - 1) {
			serializedBuffer[cursor] = '@';
			cursor += 1;
		}
	}

	args->serialized = serializedBuffer;
	args->lenSerialized = cursor;

	return cursor;
}
