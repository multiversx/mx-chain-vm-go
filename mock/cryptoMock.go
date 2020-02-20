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

func (c *CryptoHookMock) Ecrecover(hash []byte, recoveryID []byte, r []byte, s []byte) ([]byte, error) {
	return c.Result, c.Err
}
