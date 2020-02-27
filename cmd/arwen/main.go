package main

import (
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
)

func main() {
	part, err := arwenpart.NewArwenPart(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = part.StartLoop()
	if err != nil {
		log.Fatal(err)
	}
}
