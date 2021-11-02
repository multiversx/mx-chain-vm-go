#include "../../elrond/context.h"

void memGrow() {
	i64 count = int64getArgument(0);
	for (i64 i = 0; i < count; i++) {
		int64finish(i);
	}
}
