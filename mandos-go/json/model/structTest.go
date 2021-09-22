package mandosjsonmodel

// Test is a json object representing a test.
type Test struct {
	TestName    string
	CheckGas    bool
	Pre         []*Account
	Blocks      []*Block
	Network     string
	BlockHashes []JSONBytesFromString
	PostState   *CheckAccounts
}

// Block is a json object representing a block.
type Block struct {
	Results      []*TransactionResult
	Transactions []*Transaction
	BlockHeader  *BlockHeader
}

// BlockHeader is a json object representing the block header.
type BlockHeader struct {
	Beneficiary JSONBigInt // "coinbase"
	Difficulty  JSONBigInt
	Number      JSONBigInt
	GasLimit    JSONBigInt
	Timestamp   JSONUint64
}
