package main

import (
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"
)

func main() {
	facade := &arwendebug.DebugFacade{}
	app := arwendebug.Initialize(facade)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
