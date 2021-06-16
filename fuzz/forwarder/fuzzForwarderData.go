package fuzzForwarder

type programmedCallType int

const (
	syncCall programmedCallType = iota
	asyncCall
	transferExecute
)

type programmedCall struct {
	callType  programmedCallType
	fromIndex int
	toIndex   int
	token     string
	nonce     int
	amount    string
}

type fuzzData struct {
	mainCallerAddress     string
	numForwarders         int
	programmedCalls       map[int][]*programmedCall
	numFungibleTokens     int
	numSemiFungibleTokens int
}
