package mandosjsonmodel

// ESDTData models an account holding an ESDT token or a transaction sending and ESDT token
type ESDTData struct {
	TokenIdentifier JSONBytesFromString
	Nonce           JSONUint64
	Value           JSONBigInt
	Frozen          JSONUint64
}

// CheckESDTData checks the ESDT tokens held by an account
type CheckESDTData struct {
	TokenIdentifier JSONBytesFromString
	Nonce           JSONCheckUint64
	Value           JSONCheckBigInt
	Frozen          JSONCheckUint64
}

// ESDTRoles specifies token role initializations
type ESDTRoles struct {
	TokenIdentifier JSONBytesFromString
	Roles           []string
}
