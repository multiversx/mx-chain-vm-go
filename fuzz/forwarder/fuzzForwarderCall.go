package fuzzForwarder

import (
	"fmt"
	"strings"
)

func (pfe *fuzzExecutor) programCall(fromIndex, toIndex int, token string, nonce int, amount string) error {
	pfe.data.programmedCalls[fromIndex] = append(pfe.data.programmedCalls[fromIndex], &programmedCall{
		fromIndex: fromIndex,
		toIndex:   toIndex,
		token:     token,
		nonce:     nonce,
		amount:    amount,
	})

	_, err := pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "add_programmed_call",
			"arguments": [
				"%s",
				"str:%s",
				"%d",
				"%s"
			],
			"gasLimit": "50,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": ""
		}
	}`,
		pfe.nextTxIndex(),
		pfe.data.mainCallerAddress,
		pfe.forwarderAddress(fromIndex),
		pfe.forwarderAddress(toIndex),
		token,
		nonce,
		amount,
	))
	if err != nil {
		return err
	}

	return nil
}

func (pfe *fuzzExecutor) executeCall(toIndex int) error {
	var sb strings.Builder
	initialCall := &programmedCall{
		fromIndex: -1,
		toIndex:   toIndex,
		token:     "EGLD",
		nonce:     0,
		amount:    "0",
	}
	err := pfe.getExpectedLogs(initialCall, &sb)
	if err != nil {
		return err
	}

	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scCall",
		"txId": "%d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "0",
			"function": "forward_programmed_calls",
			"arguments": [],
			"gasLimit": "500,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [
				%s
			]
		}
	}`,
		pfe.nextTxIndex(),
		pfe.data.mainCallerAddress,
		pfe.forwarderAddress(toIndex),
		sb.String(),
	))
	if err != nil {
		return err
	}
	// for _, log := range output.Logs {
	// 	fmt.Println(pfe.arwenTestExecutor.LogToString(log))
	// }

	return nil
}

func (pfe *fuzzExecutor) popCall(forwarderIndex int) *programmedCall {
	if len(pfe.data.programmedCalls[forwarderIndex]) == 0 {
		return nil
	}
	result := pfe.data.programmedCalls[forwarderIndex][0]
	pfe.data.programmedCalls[forwarderIndex] = pfe.data.programmedCalls[forwarderIndex][1:]
	return result
}

func (pfe *fuzzExecutor) getExpectedLogs(call *programmedCall, sb *strings.Builder) error {
	if sb.Len() > 0 {
		sb.WriteString(",")
	}
	sb.WriteString(fmt.Sprintf(`
	{
		"address": "%s",
		"identifier": "str:forward_programmed_calls",
		"topics": [
			"str:%s",
			"%d",
			"%s"
		],
		"data": ""
	}`,
		pfe.forwarderAddress(call.toIndex),
		call.token,
		call.nonce,
		call.amount,
	))

	nextCall := pfe.popCall(call.toIndex)
	for nextCall != nil {
		pfe.getExpectedLogs(nextCall, sb)
		nextCall = pfe.popCall(call.toIndex)
	}

	return nil
}
