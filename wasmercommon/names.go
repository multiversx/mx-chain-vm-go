package wasmercommon

import "strings"

// ReadNames drives the two-step string-out convention: ask for the length, allocate, fill.
//
// Two calls rather than one because the executor does not allocate for the caller and the
// caller cannot guess the size. A zero length means there are no names at all, which is not
// an error.
//
// The two steps arrive as closures so that this stays free of cgo; each backend wraps its own
// prefixed C calls.
func ReadNames(length func() int32, fill func([]byte) int32) ([]string, error) {
	size := length()
	if size < 0 {
		return nil, ErrNullPointer
	}
	if size == 0 {
		return nil, nil
	}

	dest := make([]byte, int(size))
	written := fill(dest)
	if written <= 0 {
		return nil, ErrBufferTooSmall
	}

	// `written` counts the NUL terminator the executor wrote; the names stop before it.
	return strings.Split(string(dest[:written-1]), "|"), nil
}
