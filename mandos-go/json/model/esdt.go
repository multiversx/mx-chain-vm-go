package mandosjsonmodel

// ESDTData models an account holding an ESDT token or a transaction sending and ESDT token
type ESDTData struct {
	TokenIdentifier JSONBytesFromString
	Nonce           JSONUint64
	Value           JSONBigInt
	Frozen          JSONUint64
}
