#include "../elrond/context.h"

void memSize() {
	i64 count = int64getArgument(0);
	for (i64 i = 0; i < count; i++) {
		asm (
			"memory.size 0\n"
			"drop\n"
		);
		int64finish(i);
	}
}

void memGrowDelta() {
	i64 count = int64getArgument(0);
	i32 delta = int64getArgument(1);
	for (i64 i = 0; i < count; i++) {
		asm (
			"local.get %[delta]\n"
			"memory.grow 0\n"
			"drop\n"
			: /* No outputs, only inputs below */
			: [delta] "r" (delta)
		);
		int64finish(i);
	}
}

void memGrowDeltaOpReps() {
	i64 count = int64getArgument(0);
	i32 delta = int64getArgument(1);
	i64 opreps = int64getArgument(2);
	for (i64 i = 0; i < count; i++) {
		for (i64 j = 0; j < opreps; j++) {
			asm (
				"local.get %[delta]\n"
				"memory.grow 0\n"
				"drop\n"
				: /* No outputs, only inputs below */
				: [delta] "r" (delta)
			);
		}
	}
}
