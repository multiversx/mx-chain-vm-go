package crypto

import (
	"errors"
)

// ErrNilPrivateKey is raised when a private key was expected but received nil
var ErrNilPrivateKey = errors.New("private key is nil")

// ErrInvalidPrivateKey is raised when an invalid private key is used
var ErrInvalidPrivateKey = errors.New("private key is invalid")

// ErrNilPrivateKeyScalar is raised when a private key with nil scalar is used
var ErrNilPrivateKeyScalar = errors.New("private key holds a nil scalar")

// ErrInvalidScalar is raised when an invalid scalar is used
var ErrInvalidScalar = errors.New("scalar is invalid")

// ErrInvalidPoint is raised when an invalid point is used
var ErrInvalidPoint = errors.New("point is invalid")

// ErrNilPublicKey is raised when public key is expected but received nil
var ErrNilPublicKey = errors.New("public key is nil")

// ErrInvalidPublicKey is raised when an invalid public key is used
var ErrInvalidPublicKey = errors.New("public key is invalid")

// ErrNilPublicKeyPoint is raised when a public key with nil point is used
var ErrNilPublicKeyPoint = errors.New("public key holds a nil point")

// ErrNilParam is raised for nil parameters
var ErrNilParam = errors.New("nil parameter")

// ErrInvalidParam is raised for invalid parameters
var ErrInvalidParam = errors.New("parameter is invalid")

// ErrNilSignature is raised for a nil signature
var ErrNilSignature = errors.New("signature is nil")

// ErrNilMessage is raised when trying to verify a nil signed message or trying to sign a nil message
var ErrNilMessage = errors.New("message to be signed or to be verified is nil")

// ErrSigNotValid is raised when a signature verification fails due to invalid signature
var ErrSigNotValid = errors.New("signature is invalid")

// ErrBLSInvalidSignature will be returned when the provided BLS signature is invalid
var ErrBLSInvalidSignature = errors.New("bls12-381: invalid signature")

// ErrWrongPrivateKeySize signals that the length of the provided private key is not the expected one
var ErrWrongPrivateKeySize = errors.New("wrong private key size")

// ErrNilSuite is raised when a nil crypto suite is used
var ErrNilSuite = errors.New("crypto suite is nil")
