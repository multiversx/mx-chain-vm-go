#include "../mxvm/context.h"

// Forward declaration of the C functions, assuming they are linked externally
void mx_activate_unsafe_mode(void);
void mx_deactivate_unsafe_mode(void);
void bigIntTDiv(int, int, int);

void activateUnsafeMode(void) {
    mx_activate_unsafe_mode();
}

void deactivateUnsafeMode(void) {
    mx_deactivate_unsafe_mode();
}

void testDivByZero(void) {
    int a = bigIntNew(1);
    int b = bigIntNew(0);
    int c = bigIntNew(0);
    bigIntTDiv(c, a, b);
}

// Dummy function to satisfy the build system
void _start() {}