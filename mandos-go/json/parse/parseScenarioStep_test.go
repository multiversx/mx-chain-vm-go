package mandosjsonparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseScenario(t *testing.T) {
	snippet := `
	{
		"step": "scCall",
		"txId": "1",
		"comment": "just an example",
		"tx": {
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000",
			"to": "0x1000000000000000000000000000000000000000000000000000000000000000",
			"value": "0x00",
			"function": "someFunctionName",
			"arguments": [
				"0x1234123400000000000000000000000000000000000000000000000000000004",
				"0x00",
				"",
				"''a message (as bytes)"
			],
			"gasLimit": "0x100000",
			"gasPrice": "0x01"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [
				{
					"address": "''smart_contract_address________s1",
					"identifier": "0xf099cd8bde557814842a3121e8ddfd433a539b8c9f14bf31ebf108d12e6196e9",
					"topics": [
						"0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b000000000000000000000000",
						"0x1234123400000000000000000000000000000000000000000000000000000004"
					],
					"data": "0x00"
				}
			],
			"gas": "0x1234",
			"refund": "*"
		}
	}`

	p := Parser{}
	step, parseErr := p.ParseScenarioStep(snippet)
	require.Nil(t, parseErr)
	require.NotNil(t, step)
	require.Equal(t, "scCall", step.StepTypeName())
}
