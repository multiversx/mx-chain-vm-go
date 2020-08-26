package bls

type BLS interface {
	VerifyBLS(key []byte,  msg []byte, sig []byte) error
}
