package scenfileresolver

// FileResolver resolves scenario values starting with "file:"
type FileResolver interface {
	// Clone creates new instance of the same type.
	Clone() FileResolver

	// SetContext sets directory where the test runs, to help resolve relative paths.
	SetContext(contextPath string)

	// ResolveAbsolutePath yields absolute value based on context.
	ResolveAbsolutePath(value string) string

	// ResolveFileValue converts a value prefixed with "file:" and replaces it with the file contents.
	ResolveFileValue(value string) ([]byte, error)
}
