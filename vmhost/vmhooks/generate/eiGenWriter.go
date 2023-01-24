package vmhooksgenerate

import (
	"os"
	"path/filepath"
)

type eiGenWriter struct {
	outFile *os.File
}

func NewEIGenWriter(pathToElrondApiPackage string, relativePath string) *eiGenWriter {
	outFile, err := os.Create(filepath.Join(pathToElrondApiPackage, relativePath))
	if err != nil {
		panic(err)
	}
	return &eiGenWriter{
		outFile: outFile,
	}
}

func (writer *eiGenWriter) WriteString(s string) {
	_, _ = writer.outFile.WriteString(s)
}

func (writer *eiGenWriter) Close() {
	_ = writer.outFile.Close()
}
