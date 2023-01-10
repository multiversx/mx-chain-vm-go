package mandoscontroller

import (
	fr "github.com/multiversx/mx-chain-vm-v1_4-go/mandos-go/fileresolver"
)

// NewDefaultFileResolver yields a new DefaultFileResolver instance.
// Reexported here to avoid having all external packages importing the parser.
// DefaultFileResolver is in parse for local tests only.
func NewDefaultFileResolver() *fr.DefaultFileResolver {
	return fr.NewDefaultFileResolver()
}
