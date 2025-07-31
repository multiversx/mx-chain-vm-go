package mock

// CryptoHookMock is used in tests to check that interface methods were called
type CryptoHookMock struct {
	Result []byte
	Err    error
}

// Sha256 mocked method
func (c *CryptoHookMock) Sha256(_ []byte) ([]byte, error) {
	return c.Result, c.Err
}

// Keccak256 mocked method
func (c *CryptoHookMock) Keccak256(_ []byte) ([]byte, error) {
	return c.Result, c.Err
}

// Ripemd160 mocked method
func (c *CryptoHookMock) Ripemd160(_ []byte) ([]byte, error) {
	return c.Result, c.Err
}

// VerifyBLS mocked method
func (c *CryptoHookMock) VerifyBLS(_ []byte, _ []byte, _ []byte) error {
	return c.Err
}

// VerifyAggregatedSig mocked method
func (c *CryptoHookMock) VerifyAggregatedSig(_ [][]byte, _ []byte, _ []byte) error {
	return c.Err
}

// VerifySignatureShare mocked method
func (c *CryptoHookMock) VerifySignatureShare(_ []byte, _ []byte, _ []byte) error {
	return c.Err
}

// VerifyEd25519 mocked method
func (c *CryptoHookMock) VerifyEd25519(_ []byte, _ []byte, _ []byte) error {
	return c.Err
}

// VerifySecp256k1 mocked method
func (c *CryptoHookMock) VerifySecp256k1(_ []byte, _ []byte, _ []byte, _ uint8) error {
	return c.Err
}

// VerifySecp256r1 mocked method
func (c *CryptoHookMock) VerifySecp256r1(_ []byte, _ []byte, _ []byte) error {
	return c.Err
}

// EncodeSecp256k1DERSignature mocked method
func (c *CryptoHookMock) EncodeSecp256k1DERSignature(_, _ []byte) []byte {
	return make([]byte, 0)
}

// Ecrecover mocked method
func (c *CryptoHookMock) Ecrecover(_ []byte, _ []byte, _ []byte, _ []byte) ([]byte, error) {
	return c.Result, c.Err
}
