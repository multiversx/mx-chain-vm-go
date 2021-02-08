package worldmock

import "errors"

// Errors the mimic the ones from elrond-go.

// ErrInsufficientFunds signals the funds are insufficient for the move balance operation but the
// transaction fee is covered by the current balance
var ErrInsufficientFunds = errors.New("insufficient funds")
