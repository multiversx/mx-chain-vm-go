package scenariostestcli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	am "github.com/multiversx/mx-chain-vm-go/scenarioexec"
	mc "github.com/multiversx/mx-chain-vm-go/scenarios/controller"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
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
	useWasmer1 := flag.Bool("wasmer1", false, "use the wasmer1 executor")
	useWasmer2 := flag.Bool("wasmer2", false, "use the wasmer2 executor")
	flag.Parse()

	return &mc.RunScenarioOptions{
		ForceTraceGas: *forceTraceGas,
		UseWasmer1:    *useWasmer1,
		UseWasmer2:    *useWasmer2,
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
	if len(args) < 1 {
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
	if options.UseWasmer1 {
		executor.OverrideVMExecutor = wasmer.ExecutorFactory()
	}
	if options.UseWasmer2 {
		executor.OverrideVMExecutor = wasmer2.ExecutorFactory()
	}

	// execute
	switch {
	case isDir:
		runner := mc.NewScenarioRunner(
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
		runner := mc.NewScenarioRunner(
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
