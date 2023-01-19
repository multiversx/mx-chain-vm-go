package scenjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/model"
	oj "github.com/multiversx/mx-chain-vm-v1_4-go/scenarios/orderedjson"
)

func (p *Parser) processTx(txType mj.TransactionType, blrRaw oj.OJsonObject) (*mj.Transaction, error) {
	bltMap, isMap := blrRaw.(*oj.OJsonMap)
	if !isMap {
		return nil, errors.New("unmarshalled transaction is not a map")
	}

	blt := mj.Transaction{
		Type:      txType,
		EGLDValue: mj.JSONBigIntZero(),
		ESDTValue: nil,
	}

	var err error
	for _, kvp := range bltMap.OrderedKV {

		switch kvp.Key {
		case "nonce":
			blt.Nonce, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction nonce: %w", err)
			}
		case "from":
			if !txType.HasSender() {
				return nil, errors.New("`from` not allowed in transaction, it is always the zero address")
			}
			fromStr, err := p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction from: %w", err)
			}
			var fromErr error
			blt.From, fromErr = p.parseAccountAddress(fromStr)
			if fromErr != nil {
				return nil, fromErr
			}
		case "to":
			toStr, err := p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction to: %w", err)
			}

			if txType == mj.ScDeploy {
				if len(toStr) > 0 {
					return nil, errors.New("transaction to field not allowed for scDeploy transactions")
				}
			} else {
				blt.To, err = p.parseAccountAddress(toStr)
				if err != nil {
					return nil, err
				}
			}
		case "function":
			blt.Function, err = p.parseString(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction function: %w", err)
			}
			if !txType.HasFunction() && len(blt.Function) > 0 {
				return nil, errors.New("transaction function field not allowed in this context")
			}
		case "value":
			// backwards compatibility
			fallthrough
		case "egldValue":
			if !txType.HasValue() {
				return nil, errors.New("`egldValue` not allowed in this context")
			}
			blt.EGLDValue, err = p.processBigInt(kvp.Value, bigIntUnsignedBytes)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction egldValue: %w", err)
			}
		case "esdt":
			// backwards compatibility
			fallthrough
		case "esdtValue":
			if !txType.HasESDT() {
				return nil, errors.New("`esdtValue` not allowed in this context")
			}
			blt.ESDTValue, err = p.processTxESDT(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction esdtValue: %w", err)
			}
		case "arguments":
			blt.Arguments, err = p.parseSubTreeList(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction arguments: %w", err)
			}
			if txType == mj.Transfer && len(blt.Arguments) > 0 {
				return nil, errors.New("function arguments not allowed for transfer transactions")
			}
		case "contractCode":
			blt.Code, err = p.processStringAsByteArray(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction contract code: %w", err)
			}
			if txType != mj.ScDeploy && len(blt.Code.Value) > 0 {
				return nil, errors.New("transaction contractCode field only allowed int scDeploy transactions")
			}
		case "gasLimit":
			if !txType.HasGasLimit() {
				return nil, errors.New("`gasLimit` not allowed in this context")
			}
			blt.GasLimit, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction gasLimit: %w", err)
			}
		case "gasPrice":
			if !txType.HasGasPrice() {
				return nil, errors.New("`gasPrice` not allowed in this context")
			}
			blt.GasPrice, err = p.processUint64(kvp.Value)
			if err != nil {
				return nil, fmt.Errorf("invalid transaction gasPrice: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown field in transaction: %s", kvp.Key)
		}
	}

	return &blt, nil
}
