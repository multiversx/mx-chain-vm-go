void memGrowWrongIndex() {
	asm (
		"i32.const 10\n"
		"memory.grow 1\n"
		"drop\n"
	);
}
