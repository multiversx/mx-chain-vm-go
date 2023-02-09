package wasmer2

// SetLogLevel sets the log level for the Executor.
func SetLogLevel(logLevel LogLevel) {
	cWasmerSetLogLevel(uint64(logLevel))
}
