package ed25519

type Ed25519 interface {
	Ed25519Verify(key []byte,  msg []byte, sig []byte) error
}
