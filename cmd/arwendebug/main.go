package main

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwendebug")

const (
	// ErrCodeSuccess signals success
	ErrCodeSuccess = iota
	// ErrCodeCriticalError signals a critical error
	ErrCodeCriticalError
)

func main() {
	logger.ToggleLoggerName(true)
	logger.SetLogLevel("*:TRACE")

	facade := &arwendebug.DebugFacade{}
	app := initializeCLI(facade)

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(ErrCodeCriticalError)
	}

	os.Exit(ErrCodeSuccess)
}
