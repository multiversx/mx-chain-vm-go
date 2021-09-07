package mock

// CryptoHookMock is used in tests to check that interface methods were called
type CryptoHookMock struct {
	Result []byte
	Err    error
}

// Sha256 mocked method
func (c *CryptoHookMock) Sha256(data []byte) ([]byte, error) {
	return c.Result, c.Err
}

// Keccak256 mocked method
func (c *CryptoHookMock) Keccak256(data []byte) ([]byte, error) {
	return c.Result, c.Err
}

// Ripemd160 mocked method
func (c *CryptoHookMock) Ripemd160(data []byte) ([]byte, error) {
	return c.Result, c.Err
}

// VerifyBLS mocked method
func (c *CryptoHookMock) VerifyBLS(key []byte, msg []byte, sig []byte) error {
	return c.Err
}

// VerifyEd25519 mocked method
func (c *CryptoHookMock) VerifyEd25519(key []byte, msg []byte, sig []byte) error {
	return c.Err
}

// VerifySecp256k1 mocked method
func (c *CryptoHookMock) VerifySecp256k1(key []byte, msg []byte, sig []byte, hashType uint8) error {
	return c.Err
}

// EncodeSecp256k1DERSignature mocked method
func (c *CryptoHookMock) EncodeSecp256k1DERSignature(r, s []byte) []byte {
	return make([]byte, 0)
}

// Ecrecover mocked method
func (c *CryptoHookMock) Ecrecover(hash []byte, recoveryID []byte, r []byte, s []byte) ([]byte, error) {
	return c.Result, c.Err
}
