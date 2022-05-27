package mandosjsonparse

import (
	ei "github.com/ElrondNetwork/arwen-wasm-vm/v1_5/mandos-go/expression/interpreter"
	fr "github.com/ElrondNetwork/arwen-wasm-vm/v1_5/mandos-go/fileresolver"
)

// Parser performs parsing of both json tests (older) and scenarios (new).
type Parser struct {
	ExprInterpreter            ei.ExprInterpreter
	AllowEsdtTxLegacySyntax    bool
	AllowEsdtLegacySetSyntax   bool
	AllowEsdtLegacyCheckSyntax bool
}

// NewParser provides a new Parser instance.
func NewParser(fileResolver fr.FileResolver) Parser {
	return Parser{
		ExprInterpreter: ei.ExprInterpreter{
			FileResolver: fileResolver,
		},
		AllowEsdtTxLegacySyntax:    true,
		AllowEsdtLegacySetSyntax:   true,
		AllowEsdtLegacyCheckSyntax: true,
	}
}
