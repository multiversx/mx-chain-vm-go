package evmhooks

import "errors"

var ErrInvalidEncodedData = errors.New("invalid encoded data")

var ErrInvalidReturnDataSize = errors.New("invalid return data size")
