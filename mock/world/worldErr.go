package worldmock

import "errors"

// ErrInsufficientFunds signals the funds are insufficient for the move balance operation but the
// transaction fee is covered by the current balance.
var ErrInsufficientFunds = errors.New("insufficient funds")

// ErrNilWorldMock signals that the WorldMock is nil but shouldn't be.
var ErrNilWorldMock = errors.New("nil worldmock")
