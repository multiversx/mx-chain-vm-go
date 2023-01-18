//nolint:all
package hosttest

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/multiversx/wasm-vm/testcommon"
	test "github.com/multiversx/wasm-vm/testcommon"
)

func TestGraph_Generated(t *testing.T) {
	t.Skip("needs trace input")
	path := os.Args[len(os.Args)-1]

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()

	lineNo := -1
	var line string
	defer func() {
		fmt.Printf("Line no %d\nLast test %s\n", lineNo, line)
	}()

	var traceStepsMap map[int][]traceStep
	var traceStepsIndex map[int]traceStep
	var usedGasMap map[string]map[int]int
	var accumulatedGas, finalRemainingGas int

	includeLogLine := 1
	lineSetCount := 3 + includeLogLine

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo++
		line = scanner.Text()
		if line == "" {
			continue
		}
		if lineNo%100000 == 0 {
			fmt.Println(lineNo)
		}

		switch {
		case lineNo%lineSetCount == 0+includeLogLine:
			traceStepsMap, traceStepsIndex = parseTraceLine(line)
		case lineNo%lineSetCount == 1+includeLogLine:
			usedGasMap = parseUsedGasPerContractAndCallLine(line)
		case lineNo%lineSetCount == 2+includeLogLine:
			accumulatedGas, finalRemainingGas = parseFinalGasLine(line)
		}

		if lineNo%lineSetCount == 2+includeLogLine {
			callGraph :=
				createGraphFromScenario(traceStepsMap, traceStepsIndex, usedGasMap, accumulatedGas, finalRemainingGas)
			test.MakeGraphAndImage(callGraph)
			RunGraphCallTestTemplate(t, callGraph)
		}
	}
}

func parseUsedGasPerContractAndCallLine(line string) map[string]map[int]int {
	usedGasMap := make(map[string]map[int]int)
	usedGasElements := strings.Split(line[1:len(line)-1], "sc")
	for _, usedGasElement := range usedGasElements[1:] {
		usedGasElement = strings.Trim(usedGasElement, " @")
		separatorForElement := strings.Index(usedGasElement, ":>")
		contract := "sc" + strings.Trim(usedGasElement[:separatorForElement], " ")
		usedGasElements := strings.Trim(usedGasElement[separatorForElement+2:], " ()")

		if strings.HasPrefix(usedGasElements, "<<") {
			// no used gas values for contract
			continue
		}

		usedGasForContractMap := make(map[int]int)
		usedGasForContractElements := strings.Split(usedGasElements, "@@")
		for _, usedGasForContractElement := range usedGasForContractElements {
			usedGasForContractElement = strings.Trim(usedGasForContractElement, " ")
			separatorForElement := strings.Index(usedGasElement, ":>")
			usedGasForContractElementName := strings.Trim(usedGasForContractElement[:separatorForElement], " ")
			usedGasForContractElementValue := strings.Trim(usedGasForContractElement[separatorForElement+3:], " ")

			callId, _ := strconv.Atoi(usedGasForContractElementName)
			usedGasValue, _ := strconv.Atoi(usedGasForContractElementValue)
			usedGasForContractMap[callId*(-1)] = usedGasValue
		}

		usedGasMap[contract] = usedGasForContractMap
	}

	return usedGasMap
}

type traceStep struct {
	callType       string
	acc            string
	gas            int
	lockedGas      int
	callerCallId   int
	callId         int
	failed         bool
	propagatedFail bool
}

func parseTraceLine(line string) (map[int][]traceStep, map[int]traceStep) {

	line = eliminateCallPathElements(line)

	line = line[2 : len(line)-2]

	traceSteps := strings.Split(line, "],")
	traceStepsMap := make(map[int][]traceStep)
	traceStepsIndex := make(map[int]traceStep)
	for _, fullStep := range traceSteps {
		fullStep = strings.Trim(fullStep, " ")[1:]

		stepElements := strings.Split(fullStep, ",")
		step := traceStep{}
		for _, stepElement := range stepElements {
			separatorForStep := strings.Index(stepElement, "|")
			stepElementName := strings.Trim(stepElement[:separatorForStep], " ")
			stepElementValue := strings.Trim(stepElement[separatorForStep+3:], " ")

			switch stepElementName {
			case "type":
				step.callType = stepElementValue[1 : len(stepElementValue)-1]
			case "acc":
				step.acc = stepElementValue
			case "gas":
				gas, _ := strconv.Atoi(stepElementValue)
				step.gas = gas
			case "lockedGas":
				lockedGas, _ := strconv.Atoi(stepElementValue)
				step.lockedGas = lockedGas
			case "callerCallId":
				callerCallId, _ := strconv.Atoi(stepElementValue)
				step.callerCallId = callerCallId
			case "callId":
				callId, _ := strconv.Atoi(stepElementValue)
				step.callId = callId
			case "failed":
				if stepElementValue == "TRUE" {
					step.failed = true
				}
			case "propagatedFail":
				if stepElementValue == "TRUE" {
					step.propagatedFail = true
				}
			case "callPath":
				continue
			default:
				panic("unknown element found " + stepElementName)
			}
		}

		traceStepsIndex[step.callId] = step

		var stepsArray []traceStep
		if value, ok := traceStepsMap[step.callerCallId]; ok {
			stepsArray = append(value, step)
		} else {
			stepsArray = []traceStep{step}
		}
		traceStepsMap[step.callerCallId] = stepsArray
	}

	return traceStepsMap, traceStepsIndex
}

func eliminateCallPathElements(line string) string {
	callPath := "callPath |-> "
	for {
		indexOfCallPath := strings.Index(line, callPath)
		if indexOfCallPath == -1 {
			break
		}
		line = line[:indexOfCallPath] + line[strings.Index(line, ">>")+3:]
	}
	return line
}

func parseFinalGasLine(line string) (int, int) {
	line = line[2 : len(line)-2]
	parts := strings.Split(line, ",")
	accumulatedGas, _ := strconv.Atoi(strings.Trim(parts[0], " "))
	finalRemainingGas, _ := strconv.Atoi(strings.Trim(parts[1], " "))
	return accumulatedGas, finalRemainingGas
}

func createGraphFromScenario(
	traceStepsMap map[int][]traceStep,
	traceStepsIndex map[int]traceStep,
	usedGasMap map[string]map[int]int,
	_ int,
	_ int,
) *testcommon.TestCallGraph {

	callGraph := testcommon.CreateTestCallGraph()

	callIdToNode := make(map[int]*testcommon.TestCallNode)
	asyncCallIdToCallbackNode := make(map[int]*testcommon.TestCallNode)
	asyncCallIdToEdge := make(map[int]*testcommon.TestCallEdge)

	rootStep := traceStepsMap[0][0]
	traceStepsMap[0] = traceStepsMap[0][1:]

	root := callGraph.AddStartNode(rootStep.acc, "f1", uint64(rootStep.gas), getGasUsedForStep(usedGasMap, rootStep))
	callIdToNode[0] = root

	contractToFunctionCounter := make(map[string]int)
	contractToCallbackCounter := make(map[string]int)

	keys := make([]int, 0)
	for k := range traceStepsMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, parentId := range keys {
		parentNode := callIdToNode[parentId]
		childSteps := traceStepsMap[parentId]
		for _, childStep := range childSteps {
			newFunctionCounter := getUpdatedCounter(contractToFunctionCounter, childStep.acc)

			if childStep.callType != "Callback" && childStep.callType != "CallbackCross" {
				childNode := callGraph.AddNode(childStep.acc, fmt.Sprintf("f%d", newFunctionCounter))
				callIdToNode[childStep.callId] = childNode

				switch childStep.callType {

				case "Sync":
					edge := callGraph.AddSyncEdge(parentNode, childNode).
						SetGasLimit(uint64(childStep.gas)).
						SetGasUsed(getGasUsedForStep(usedGasMap, childStep))

					if childStep.failed && !childStep.propagatedFail {
						edge.SetFail()
					}

				case "AsyncLocal", "AsyncCross":
					newCallbackCounter := getUpdatedCounter(contractToCallbackCounter, traceStepsIndex[parentId].acc)
					callbackFunction := fmt.Sprintf("cb%d", newCallbackCounter)

					var edge *testcommon.TestCallEdge
					if childStep.callType == "AsyncLocal" {
						edge = callGraph.AddAsyncEdge(parentNode, childNode, callbackFunction, "")
					} else {
						edge = callGraph.AddAsyncCrossShardEdge(parentNode, childNode, callbackFunction, "")
					}

					minimumGasLocked := 150
					extraGasLocked := childStep.lockedGas - minimumGasLocked
					edge.
						SetGasLimit(uint64(childStep.gas)).
						SetGasLocked(uint64(extraGasLocked)).
						SetGasUsed(getGasUsedForStep(usedGasMap, childStep))

					if childStep.failed && !childStep.propagatedFail {
						edge.SetFail()
					}

					callbackNode := callGraph.AddNode(traceStepsIndex[childStep.callerCallId].acc, callbackFunction)
					asyncCallIdToCallbackNode[childStep.callId] = callbackNode
					asyncCallIdToEdge[childStep.callId] = edge
				}
			} else {
				callerEdge := asyncCallIdToEdge[childStep.callerCallId]
				callerEdge.SetGasUsedByCallback(getGasUsedForStep(usedGasMap, childStep))

				if childStep.failed && !childStep.propagatedFail {
					callerEdge.SetCallbackFail()
				}

				callbackNode := asyncCallIdToCallbackNode[childStep.callerCallId]
				callIdToNode[childStep.callId] = callbackNode
			}
		}
	}

	return callGraph
}

func getGasUsedForStep(usedGasMap map[string]map[int]int, childStep traceStep) uint64 {
	return uint64(usedGasMap[childStep.acc][childStep.callId])
}

func getUpdatedCounter(contractToCounter map[string]int, account string) int {
	var newCounterValue int
	if counterValue, ok := contractToCounter[account]; ok {
		newCounterValue = counterValue + 1
	} else {
		newCounterValue = 1
	}
	contractToCounter[account] = newCounterValue
	return newCounterValue
}
