package vmhooksgenerate

import (
	"bufio"
	"os"
)

type OpcodeNames struct {
	AllowedV1     []string
	AllowedV2     []string
	RelevantCodes []string

	// PerByteCodes are the opcodes that, on top of their flat cost, are also charged
	// for each byte they process (the bulk memory operators).
	PerByteCodes []string
}

func loadOpcodeNamesForVersion(filePath string) []string {
	var names []string

	readFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		names = append(names, line)
	}

	return names
}

func LoadOpcodeNames() *OpcodeNames {
	return &OpcodeNames{
		AllowedV1:     loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_allowed_opcodes_v1.txt"),
		AllowedV2:     loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_allowed_opcodes_v2.txt"),
		RelevantCodes: loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_relevant_opcodes.txt"),
		PerByteCodes:  loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_per_byte_opcodes.txt"),
	}
}

// IsPerByteOpcode returns true for the opcodes that are charged per processed byte,
// on top of their flat cost.
func (on *OpcodeNames) IsPerByteOpcode(opcodeName string) bool {
	for _, perByteOpcode := range on.PerByteCodes {
		if perByteOpcode == opcodeName {
			return true
		}
	}
	return false
}

// PerByteCostName is the name of the cost of a single byte processed by an opcode.
//
// It complements the cost of the opcode itself, which is charged flat, no matter how
// many bytes end up being processed.
func PerByteCostName(opcodeName string) string {
	return opcodeName + "PerByte"
}

// MaxNameLength helps with aligning to the right
func (on *OpcodeNames) MaxNameLength() int {
	var max int
	for _, name := range on.RelevantCodes {
		if len(name) > max {
			max = len(name)
		}
	}
	return max
}
