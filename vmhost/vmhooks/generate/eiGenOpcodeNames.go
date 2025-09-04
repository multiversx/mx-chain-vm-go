package vmhooksgenerate

import (
	"bufio"
	"os"
)

type OpcodeNames struct {
	AllowedV1     []string
	AllowedV2     []string
	RelevantCodes []string
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
	}
}
