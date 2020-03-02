package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
)

func main() {
	//mainWithFiles()
	//mainWithStdPipes()
	mainWithPipes()
}

func mainWithFiles() {
	user, _ := user.Current()
	home := user.HomeDir
	folder := path.Join(home, "Arwen")
	os.MkdirAll(folder, os.ModePerm)

	nodeToArwen := filepath.Join(folder, fmt.Sprintf("node-to-arwen.bin"))
	arwenToNode := filepath.Join(folder, fmt.Sprintf("arwen-to-node.bin"))

	// Create the communication files
	nodeToArwenFile, err := os.Create(nodeToArwen)
	if err != nil {
		log.Fatal("Cannot create file")
	}

	arwenToNodeFile, err := os.Create(arwenToNode)
	if err != nil {
		log.Fatal("Cannot create file")
	}

	nodeToArwenFile.Close()
	arwenToNodeFile.Close()

	// Open the files as required
	nodeToArwenFile, err = os.Open(nodeToArwen)
	if err != nil {
		log.Fatal("Cannot open file [nodeToArwen]")
	}

	arwenToNodeFile, err = os.OpenFile(arwenToNode, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot open file [arwenToNode]")
	}
	part, err := arwenpart.NewArwenPart(nodeToArwenFile, arwenToNodeFile)
	if err != nil {
		log.Fatal(err)
	}

	err = part.StartLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func mainWithStdPipes() {
	part, err := arwenpart.NewArwenPart(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	err = part.StartLoop()
	if err != nil {
		log.Fatal(err)
	}
}

func mainWithPipes() {
	// NewFile returns a new File with the given file descriptor and name.
	// On Unix systems, if the file descriptor is in non-blocking mode, NewFile will attempt to return a pollable File (one for which the SetDeadline methods work).

	// A deadline is an absolute time after which I/O operations fail with an error instead of blocking.
	// The deadline applies to all future and pending I/O, not just the immediately following call to Read or Write.
	// After a deadline has been exceeded, the connection can be refreshed by setting a deadline in the future.

	nodeToArwenFile := os.NewFile(3, "/proc/self/fd/3")
	if nodeToArwenFile == nil {
		log.Fatal("Cannot create file")
	}

	arwenToNodeFile := os.NewFile(4, "/proc/self/fd/4")
	if arwenToNodeFile == nil {
		log.Fatal("Cannot create file")
	}

	part, err := arwenpart.NewArwenPart(nodeToArwenFile, arwenToNodeFile)
	if err != nil {
		log.Fatal(err)
	}

	err = part.StartLoop()
	if err != nil {
		log.Fatal(err)
	}
}
