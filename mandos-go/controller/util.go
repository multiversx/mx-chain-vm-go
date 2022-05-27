package mandoscontroller

import (
	fr "github.com/ElrondNetwork/arwen-wasm-vm/v1_5/mandos-go/fileresolver"
)

// NewDefaultFileResolver yields a new DefaultFileResolver instance.
// Reexported here to avoid having all external packages importing the parser.
// DefaultFileResolver is in parse for local tests only.
func NewDefaultFileResolver() *fr.DefaultFileResolver {
	return fr.NewDefaultFileResolver()
}
