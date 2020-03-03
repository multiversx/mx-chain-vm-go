package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/ElrondNetwork/arwen-wasm-vm/ipc/arwenpart"
)

func main() {
	err := syscall.SetNonblock(3, true)
	if err != nil {
		fmt.Println("SetNoblock error")
		fmt.Println(err)
		return
	}

	err = syscall.SetNonblock(4, true)
	if err != nil {
		fmt.Println("SetNoblock error")
		fmt.Println(err)
		return
	}

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
