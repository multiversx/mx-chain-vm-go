package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const suffix = ".toml"

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	fs, _ := ioutil.ReadDir(".")
	out, _ := os.Create("gasScheduleEmbedGenerated.go")
	out.Write([]byte("package gasschedules \n\n"))
	out.Write([]byte("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n"))
	out.Write([]byte("// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!\n"))
	out.Write([]byte("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n"))
	out.Write([]byte("\n"))
	out.Write([]byte("// Please do not edit manually!\n"))
	out.Write([]byte("// Call `go generate` in `arwen-wasm-vm/arwenmandos/gasSchedules` to update it.\n"))
	out.Write([]byte("\n"))
	out.Write([]byte("const (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), suffix) {
			out.Write([]byte("\t" + strings.TrimSuffix(f.Name(), suffix) + " = `"))
			f, _ := os.Open(f.Name())
			io.Copy(out, f)
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}