package main

import (
	"log"
	"os"
)

func main() {
	server, err := NewServer(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
