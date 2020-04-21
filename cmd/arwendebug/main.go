package main

import (
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwendebug")

func main() {
	logger.ToggleLoggerName(true)
	logger.SetLogLevel("*:DEBUG")

	facade := &arwendebug.DebugFacade{}
	app := initializeCLI(facade)

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
}
