package dex

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/esdt"
)

func (pfe *fuzzDexExecutor) interpretExpr(expression string) []byte {
	bytes, err := pfe.mandosParser.ExprInterpreter.InterpretString(expression)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (pfe *fuzzDexExecutor) getTokensWithNonce(address string, toktik string, nonce int) (*big.Int, error) {
	return pfe.world.BuiltinFuncs.GetTokenBalance(pfe.interpretExpr(address), []byte(toktik), uint64(nonce))
}

func (pfe *fuzzDexExecutor) getTokens(address string, toktik string) (*big.Int, error) {
	return pfe.world.BuiltinFuncs.GetTokenBalance(pfe.interpretExpr(address), []byte(toktik), 0)
}

func (pfe *fuzzDexExecutor) getTokenData(address string, toktik string, nonce int) (*esdt.ESDigitalToken, error) {
	return pfe.world.BuiltinFuncs.GetTokenData(pfe.interpretExpr(address), []byte(toktik), uint64(nonce))
}
