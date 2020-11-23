package mock

type CryptoHookMock struct {
	Result []byte
	Err    error
}

func (c *CryptoHookMock) Sha256(data []byte) ([]byte, error) {
	return c.Result, c.Err
}

func (c *CryptoHookMock) Keccak256(data []byte) ([]byte, error) {
	return c.Result, c.Err
}

func (c *CryptoHookMock) Ripemd160(data []byte) ([]byte, error) {
	return c.Result, c.Err
}

func (c *CryptoHookMock) VerifyBLS(key []byte,  msg []byte, sig []byte) error {
	return c.Err
}

func (c *CryptoHookMock) VerifyEd25519(key []byte,  msg []byte, sig []byte) error {
	return c.Err
}

func (c *CryptoHookMock) VerifySecp256k1(key []byte,  msg []byte, sig []byte) error {
	return c.Err
}

func (c *CryptoHookMock) Ecrecover(hash []byte, recoveryID []byte, r []byte, s []byte) ([]byte, error) {
	return c.Result, c.Err
}
