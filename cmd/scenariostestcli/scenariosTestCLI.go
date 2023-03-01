package scenariostestcli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mc "github.com/multiversx/mx-chain-scenario-go/controller"
	am "github.com/multiversx/mx-chain-vm-go/scenarioexec"
)

func resolveArgument(exeDir string, arg string) (string, bool, error) {
	fi, err := os.Stat(arg)
	if os.IsNotExist(err) {
		arg = filepath.Join(exeDir, arg)
		fmt.Println(arg)
		fi, err = os.Stat(arg)
	}
	if err != nil {
		return "", false, err
	}
	return arg, fi.IsDir(), nil
}

func parseOptionFlags() *mc.RunScenarioOptions {
	forceTraceGas := flag.Bool("force-trace-gas", false, "overrides the traceGas option in the scenarios")
	flag.Parse()

	return &mc.RunScenarioOptions{
		ForceTraceGas: *forceTraceGas,
	}
}

// ScenariosTestCLI provides the functionality for any scenarios test executor.
func ScenariosTestCLI() {
	options := parseOptionFlags()

	// directory of this executable
	exeDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// argument
	args := flag.Args()
	if len(args) != 1 {
		panic("One argument expected - the path to the json test or directory.")
	}
	jsonFilePath, isDir, err := resolveArgument(exeDir, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// init
	executor, err := am.NewVMTestExecutor()
	if err != nil {
		panic("Could not instantiate VM VM")
	}

	// execute
	switch {
	case isDir:
		runner := mc.NewScenarioController(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunAllJSONScenariosInDirectory(
			jsonFilePath,
			"",
			".scen.json",
			[]string{},
			options)
	case strings.HasSuffix(jsonFilePath, ".scen.json"):
		runner := mc.NewScenarioController(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONScenario(jsonFilePath, options)
	default:
		runner := mc.NewTestRunner(
			executor,
			mc.NewDefaultFileResolver(),
		)
		err = runner.RunSingleJSONTest(jsonFilePath)
	}

	// print result
	if err == nil {
		fmt.Println("SUCCESS")
	} else {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
