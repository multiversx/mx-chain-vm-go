package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	ARGUMENTS_MANDOS_FIELD_NAME = "arguments"
	OUT_MANDOS_FIELD_NAME       = "out"

	SUCCESS_STATUS_CODE = 0
	FAIL_STATUS_CODE    = 4
)

func (fe *fuzzExecutor) createAccount(address string, balance *big.Int) error {
	err := fe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"%s": {
				"nonce": "0",
				"balance": "%s",
				"storage": {},
				"code": ""
			}
		}
	}`,
		address,
		balance.String(),
	))
	if err != nil {
		return err
	}

	return nil
}

func (fe *fuzzExecutor) deployContract(deployerAddress string, scAddress string,
	contractCodeFileName string, initArguments ...string) error {

	err := fe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"newAddresses": [
                {
                    "creatorAddress": "%s",
                    "creatorNonce": "%d",
                    "newAddress": "%s"
                }
            ]
	}`,
		deployerAddress,
		fe.getNonce(deployerAddress),
		scAddress,
	))
	if err != nil {
		return err
	}

	deployArgs := constructArrayMandosField(ARGUMENTS_MANDOS_FIELD_NAME, initArguments...)
	deployMandosSnippet := `
	{
		"step": "scDeploy",
		"txId": "deploy",
		"tx": {
			"from": "%s",
			"contractCode": "file:%s",
			"value": "0",`
	deployMandosSnippet += deployArgs
	deployMandosSnippet += `,
			"gasLimit": "500,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"status": "0",
			"message": "",
			"gas": "*",
			"refund": "*"
		}
	}`

	err = fe.executeStep(fmt.Sprintf(deployMandosSnippet,
		deployerAddress,
		contractCodeFileName,
	))
	if err != nil {
		return err
	}

	return nil
}

func (fe *fuzzExecutor) performSmartContractCall(caller string, scAddress string,
	value *big.Int, scFunction string, arguments []string,
	expectedSuccess bool, expectedMessage string, expectedOutData []string) ([][]byte, error) {

	scCallArgs := constructArrayMandosField(ARGUMENTS_MANDOS_FIELD_NAME, arguments...)
	scCallExpectedOut := constructArrayMandosField(OUT_MANDOS_FIELD_NAME, expectedOutData...)

	var expectedStatusCode int
	if expectedSuccess {
		expectedStatusCode = SUCCESS_STATUS_CODE
	} else {
		expectedStatusCode = FAIL_STATUS_CODE
	}

	scCallMandosSnippet := `
	{
		"step": "scCall",
		"txId": "%05d",
		"tx": {
			"from": "%s",
			"to": "%s",
			"value": "%s",
			"function": "%s",`
	scCallMandosSnippet += scCallArgs
	scCallMandosSnippet += `,
			"gasLimit": "500,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"status": "%d",
			"message": "%s",`
	scCallMandosSnippet += scCallExpectedOut
	scCallMandosSnippet += `,
			"gas": "*",
			"refund": "*"
		}
	}`

	output, err := fe.executeTxStep(fmt.Sprintf(scCallMandosSnippet,
		fe.nextTxIndex(),
		caller,
		scAddress,
		value.String(),
		scFunction,
		expectedStatusCode,
		expectedMessage,
	))
	if err != nil {
		return [][]byte{}, err
	}

	return output.ReturnData, nil
}

func (fe *fuzzExecutor) createChildContractAddresses() error {
	err := fe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"newAddresses": [
			{
				"creatorAddress": "%s",
				"creatorNonce": "0",
				"newAddress": "%s"
			},
			{
				"creatorAddress": "%s",
				"creatorNonce": "1",
				"newAddress": "%s"
			},
			{
				"creatorAddress": "%s",
				"creatorNonce": "2",
				"newAddress": "%s"
			},
			{
				"creatorAddress": "%s",
				"creatorNonce": "3",
				"newAddress": "%s"
			}
		]
	}`,
		fe.data.actorAddresses.multisig,
		fe.data.actorAddresses.egldEsdtSwap,
		fe.data.actorAddresses.multisig,
		fe.data.actorAddresses.multiTransferEsdt,
		fe.data.actorAddresses.multisig,
		fe.data.actorAddresses.ethereumFeePrepay,
		fe.data.actorAddresses.multisig,
		fe.data.actorAddresses.esdtSafe,
	))
	if err != nil {
		return err
	}

	return nil
}

func constructArrayMandosField(mandosFieldName string, arguments ...string) string {
	nrArguments := len(arguments)
	if nrArguments == 0 {
		return fmt.Sprintf(`"%s": []`, mandosFieldName)
	}

	// no comma after last one, so we do it separately
	repeatedStringFormatSpecifier := strings.Repeat(`"%s", `, nrArguments-1)
	repeatedStringFormatSpecifier += `"%s"`

	mandosArgumentsSnippet := `"%s": [ ` + repeatedStringFormatSpecifier + ` ]`

	// first arg in the snippet is the field name
	argsAsInterface := make([]interface{}, 0, nrArguments+1)
	argsAsInterface = append(argsAsInterface, mandosFieldName)
	for _, arg := range arguments {
		argsAsInterface = append(argsAsInterface, arg)
	}

	return fmt.Sprintf(mandosArgumentsSnippet, argsAsInterface...)
}

func (fe *fuzzExecutor) setEsdtLocalRoles(scAddress string, tokenId string, scOwner string, scCodePath string) error {
	err := fe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"%s": {
				"nonce": "0",
				"balance": "0",
				"esdt": {
					"%s": {
						"balance": "0",
						"roles": [
							"ESDTRoleLocalMint",
							"ESDTRoleLocalBurn"
						]
					}
				},
				"storage": {},
				"owner": "%s",
				"code": "%s"
			}
		}
	}`,
		scAddress,
		tokenId,
		scOwner,
		scCodePath,
	))
	if err != nil {
		return err
	}

	return nil
}
