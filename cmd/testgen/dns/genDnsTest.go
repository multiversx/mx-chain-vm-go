// Generates dns deploy scenario step, with 256 dns contracts, 1 per shard.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	scenio "github.com/multiversx/mx-chain-scenario-go/scenario/io"
	scenjsonparse "github.com/multiversx/mx-chain-scenario-go/scenario/json/parse"
	scenjsonwrite "github.com/multiversx/mx-chain-scenario-go/scenario/json/write"
	scenmodel "github.com/multiversx/mx-chain-scenario-go/scenario/model"
)

// DefaultVMType helps us set up the scenario generator.
var DefaultVMType = []byte{5, 0}

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	vmTestRoot := filepath.Join(exePath, "../../../test")
	return vmTestRoot
}

type testGenerator struct {
	parser            scenjsonparse.Parser
	generatedScenario *scenmodel.Scenario
}

func (tg *testGenerator) addStep(stepSnippet string) {
	step, err := tg.parser.ParseScenarioStep(stepSnippet)
	if err != nil {
		panic(err)
	}
	tg.generatedScenario.Steps = append(tg.generatedScenario.Steps, step)
}

func main() {
	fileResolver := scenio.NewDefaultFileResolver().
		ReplacePath(
			"dns.wasm",
			filepath.Join(getTestRoot(), "dns/dns.wasm"))
	tg := &testGenerator{
		parser: scenjsonparse.NewParser(fileResolver, DefaultVMType),
		generatedScenario: &scenmodel.Scenario{
			Name: "dns test",
		},
	}

	newAddressesSnippets := ""
	for shard := 0; shard < 256; shard++ {
		if shard > 0 {
			newAddressesSnippets += ","
		}
		newAddressesSnippets += fmt.Sprintf(`{
				"creatorAddress": "''dns_owner_______________________",
				"creatorNonce": "0x%02x",
				"newAddress": "''dns____________________________|0x%02x"
			}`,
			shard,
			shard)
	}
	tg.addStep(fmt.Sprintf(`
			{
				"step": "setState",
				"accounts": {
					"''dns_owner_______________________": {
						"nonce": "0",
						"balance": "0",
						"storage": {},
						"code": ""
					}
				},
				"newAddresses": [
					%s
				]
			}`,
		newAddressesSnippets))

	for shard := 0; shard < 256; shard++ {
		tg.addStep(fmt.Sprintf(`
			{
				"step": "scDeploy",
				"txId": "deploy-0x%02x",
				"tx": {
					"from": "''dns_owner_______________________",
					"value": "0",
					"contractCode": "file:dns.wasm",
					"arguments": [ "123,000" ],
					"gasLimit": "100,000",
					"gasPrice": "0"
				},
				"expect": {
					"out": [],
					"status": "",
					"logs": "*",
					"gas": "*",
					"refund": "*"
				}
			}`,
			shard))
	}

	for shard := 0; shard < 256; shard++ {
		tg.addStep(fmt.Sprintf(`
			{
				"step": "scCall",
				"txId": "feature-register-0x%02x",
				"tx": {
					"from": "''dns_owner_______________________",
					"to": "''dns____________________________|0x%02x",
					"value": "0",
					"function": "setFeatureFlag",
					"arguments": [ "''register", "true" ],
					"gasLimit": "100,000",
					"gasPrice": "0"
				},
				"expect": {
					"out": [],
					"status": "",
					"logs": "*",
					"gas": "*",
					"refund": "*"
				}
			}`,
			shard, shard))
	}

	// save
	serialized := scenjsonwrite.ScenarioToJSONString(tg.generatedScenario)
	err := os.WriteFile(
		filepath.Join(getTestRoot(), "dns/dns_init.steps.json"),
		[]byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
