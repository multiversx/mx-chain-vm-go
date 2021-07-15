package dex

import (
	"math/big"

	worldmock "github.com/ElrondNetwork/arwen-wasm-vm/v1_4/mock/world"
)

func (pfe *fuzzDexExecutor) interpretExpr(expression string) []byte {
	bytes, err := pfe.mandosParser.ExprInterpreter.InterpretString(expression)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (pfe *fuzzDexExecutor) getTokensWithNonce(address string, toktik string, nonce int) (*big.Int, error) {
	token := worldmock.MakeTokenKey([]byte(toktik), uint64(nonce))
	return pfe.world.BuiltinFuncs.GetTokenBalance(pfe.interpretExpr(address), token)
}

func (pfe *fuzzDexExecutor) getTokens(address string, toktik string) (*big.Int, error) {
	token := worldmock.MakeTokenKey([]byte(toktik), 0)
	return pfe.world.BuiltinFuncs.GetTokenBalance(pfe.interpretExpr(address), token)
}
