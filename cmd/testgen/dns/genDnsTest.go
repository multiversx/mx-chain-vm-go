// Generates dns deploy scenario step, with 256 dns contracts, 1 per shard.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	mc "github.com/multiversx/mx-chain-vm-go/mandos-go/controller"
	mjparse "github.com/multiversx/mx-chain-vm-go/mandos-go/json/parse"
	mjwrite "github.com/multiversx/mx-chain-vm-go/mandos-go/json/write"
	mj "github.com/multiversx/mx-chain-vm-go/mandos-go/model"
)

func getTestRoot() string {
	exePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	arwenTestRoot := filepath.Join(exePath, "../../../test")
	return arwenTestRoot
}

type testGenerator struct {
	mandosParser      mjparse.Parser
	generatedScenario *mj.Scenario
}

func (tg *testGenerator) addStep(stepSnippet string) {
	step, err := tg.mandosParser.ParseScenarioStep(stepSnippet)
	if err != nil {
		panic(err)
	}
	tg.generatedScenario.Steps = append(tg.generatedScenario.Steps, step)
}

func main() {
	fileResolver := mc.NewDefaultFileResolver().
		ReplacePath(
			"dns.wasm",
			filepath.Join(getTestRoot(), "dns/dns.wasm"))
	tg := &testGenerator{
		mandosParser: mjparse.NewParser(fileResolver),
		generatedScenario: &mj.Scenario{
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
	serialized := mjwrite.ScenarioToJSONString(tg.generatedScenario)
	err := ioutil.WriteFile(
		filepath.Join(getTestRoot(), "dns/dns_init.steps.json"),
		[]byte(serialized), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
