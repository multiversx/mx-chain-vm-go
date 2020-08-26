package secp256k1

type Secp256k1 interface {
	VerifySecp256k1(key []byte,  msg []byte, sig []byte) error
	Ecrecover(hash []byte, recoveryID []byte, r []byte, s []byte) ([]byte, error)
}
