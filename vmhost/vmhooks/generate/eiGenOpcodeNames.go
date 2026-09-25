package vmhooksgenerate

import (
	"bufio"
	"os"
)

type OpcodeNames struct {
	AllowedV1 []string
	AllowedV2 []string

	// PerByteCodes are the opcodes that, on top of their flat cost, are also charged
	// for each byte they process (the bulk memory operators).
	PerByteCodes []string

	// AdditionalCosts are executor costs that are not charged per opcode, so they cannot
	// be found in any of the allowed opcode lists.
	AdditionalCosts []string
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
		AllowedV1:       loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_allowed_opcodes_v1.txt"),
		AllowedV2:       loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_allowed_opcodes_v2.txt"),
		PerByteCodes:    loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_per_byte_opcodes.txt"),
		AdditionalCosts: loadOpcodeNamesForVersion("generate/cmd/input/wasmer2_additional_costs.txt"),
	}
}

// AllCostNames yields every cost that the VM hands over to the executor: the opcodes allowed
// in any version, the per-byte costs of the bulk memory operators, and the additional costs,
// which are not charged per opcode.
//
// The result is the exact layout of the OpcodeCost structs on both sides of the FFI boundary,
// so changing it requires rebuilding the executor.
func (on *OpcodeNames) AllCostNames() []string {
	var costNames []string
	alreadyAdded := make(map[string]bool)
	addCost := func(costName string) {
		if alreadyAdded[costName] {
			return
		}
		alreadyAdded[costName] = true
		costNames = append(costNames, costName)
	}

	for _, opcodeName := range on.AllowedV1 {
		addCost(opcodeName)
	}
	for _, opcodeName := range on.AllowedV2 {
		addCost(opcodeName)
	}
	for _, opcodeName := range on.PerByteCodes {
		addCost(PerByteCostName(opcodeName))
	}
	for _, costName := range on.AdditionalCosts {
		addCost(costName)
	}

	return costNames
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

// maxCostNameLength helps with aligning to the right
func maxCostNameLength(costNames []string) int {
	var max int
	for _, name := range costNames {
		if len(name) > max {
			max = len(name)
		}
	}
	return max
}
