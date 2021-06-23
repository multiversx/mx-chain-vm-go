package elrond_ethereum_bridge

import (
	"fmt"
	"math/big"
	"strings"
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

	deployArgs := constructArgumentsMandosField(initArguments...)
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
			"gasLimit": "200,000,000",
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

func constructArgumentsMandosField(arguments ...string) string {
	nrArguments := len(arguments)
	if nrArguments == 0 {
		return `"arguments": []`
	}

	// no comma after last one, so we do it separately
	repeatedStringFormatSpecifier := strings.Repeat(`"%s", `, nrArguments-1)
	repeatedStringFormatSpecifier += `"%s"`

	mandosArgumentsSnippet := `"arguments": [ ` + repeatedStringFormatSpecifier + ` ]`

	argsAsInterface := make([]interface{}, 0, nrArguments)
	for _, arg := range arguments {
		argsAsInterface = append(argsAsInterface, arg)
	}

	return fmt.Sprintf(mandosArgumentsSnippet, argsAsInterface...)
}
