package mandosjsonmodel

// ESDTTransfer models the transfer of tokens in a tx
type ESDTTxData struct {
	TokenIdentifier JSONBytesFromString
	Nonce           JSONUint64
	Value           JSONBigInt
}

// ESDTInstance models an instance of an NFT/SFT, with its own nonce
type ESDTInstance struct {
	Nonce   JSONUint64
	Balance JSONBigInt
	Uri     JSONBytesFromTree
}

// ESDTData models an account holding an ESDT token
type ESDTData struct {
	TokenIdentifier JSONBytesFromString
	Instances       []*ESDTInstance
	LastNonce       JSONUint64
	Roles           []string
	Frozen          JSONUint64
}

// CheckESDTInstance checks an instance of an NFT/SFT, with its own nonce
type CheckESDTInstance struct {
	Nonce   JSONCheckUint64
	Balance JSONCheckBigInt
	Uri     JSONCheckBytes
}

// CheckESDTData checks the ESDT tokens held by an account
type CheckESDTData struct {
	TokenIdentifier JSONBytesFromString
	Instances       []*CheckESDTInstance
	LastNonce       JSONCheckUint64
	Roles           []string
	Frozen          JSONCheckUint64
}
