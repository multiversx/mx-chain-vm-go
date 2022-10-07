package elrondapigenerate

import (
	"fmt"
	"os"
)

func WriteEIInterface(eiMetadata *EIMetadata, out *os.File) {
	out.WriteString("package elrondapi \n\n")
	out.WriteString("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
	out.WriteString("// !!!!!!!!!!!!!!!!!!!!!! AUTO-GENERATED FILE !!!!!!!!!!!!!!!!!!!!!!\n")
	out.WriteString("// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
	out.WriteString("\n")
	out.WriteString("type ImportsInterface interface {\n")

	for _, funcMetadata := range eiMetadata.AllFunctions {
		out.WriteString("\t")
		out.WriteString(funcMetadata.PublicName)
		out.WriteString("(")
		for argIndex, arg := range funcMetadata.Arguments {
			if argIndex > 0 {
				out.WriteString(", ")
			}
			out.WriteString(fmt.Sprintf("%s %s", arg.Name, arg.Type))
		}
		out.WriteString(")")
		if funcMetadata.Result != nil {
			out.WriteString(fmt.Sprintf(" %s", funcMetadata.Result.Type))
		}

		out.WriteString("\n")
	}

	out.WriteString("}\n")
}
