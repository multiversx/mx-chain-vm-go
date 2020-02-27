package main

import (
	"log"
	"os"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
)

func main() {
	server, err := arwenpart.NewChildServer(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
