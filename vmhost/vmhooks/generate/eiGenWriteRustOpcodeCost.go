package vmhooksgenerate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const useStatements = `
use serde::{Deserialize, Serialize};

`

// WriteRustOpcodeCost generates code for opcode_cost.rs
func WriteRustOpcodeCost(out *eiGenWriter) {
	out.WriteString(`// Code generated by vmhooks generator. DO NOT EDIT.

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
`)
	out.WriteString(useStatements)
	out.WriteString("#[derive(Clone, Debug, Default, Deserialize, Serialize, PartialEq)]\n")
	out.WriteString("#[serde(default)]\n")
	out.WriteString("pub struct OpcodeCost {\n")

	readFile, err := os.Open("generate/cmd/input/wasmer2_opcodes_short.txt")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		out.WriteString(fmt.Sprintf("    #[serde(rename = \"%s\", default)]\n", line))
		out.WriteString(fmt.Sprintf("    pub opcode_%s: u32,\n", strings.ToLower(line)))
	}
	out.WriteString("}\n")
}
