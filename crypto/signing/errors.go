package signing

import (
	"errors"
)

// ErrInvalidPublicKey is raised when an invalid public key is used
var ErrInvalidPublicKey = errors.New("public key is invalid")

// ErrInvalidSignature will be returned when ed25519 signature verification fails
var ErrInvalidSignature = errors.New("invalid signature")

// ErrHasherNotSupported will be returned when a provided hasher type is not supported by the signature scheme
var ErrHasherNotSupported = errors.New("hasher not supported")
