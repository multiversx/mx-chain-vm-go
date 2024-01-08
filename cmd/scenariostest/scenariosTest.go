package main

import (
	scenclibase "github.com/multiversx/mx-chain-scenario-go/clibase"
	mc "github.com/multiversx/mx-chain-scenario-go/controller"

	vmscenario "github.com/multiversx/mx-chain-vm-go/scenario"
	"github.com/multiversx/mx-chain-vm-go/wasmer"
	"github.com/multiversx/mx-chain-vm-go/wasmer2"
	cli "github.com/urfave/cli/v2"
)

var _ scenclibase.CLIRunConfig = (*vm15Flags)(nil)

func main() {
	scenclibase.ScenariosCLI("VM 1.5 internal", &vm15Flags{})
}

type vm15Flags struct{}

func (*vm15Flags) GetFlags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "force-trace-gas",
			Aliases: []string{"g"},
			Usage:   "overrides the traceGas option in the scenarios`",
		},
		&cli.BoolFlag{
			Name:  "wasmer1",
			Usage: "use the wasmer1 executor`",
		},
		&cli.BoolFlag{
			Name:  "wasmer2",
			Usage: "use the wasmer2 executor`",
		},
	}
}

func (*vm15Flags) ParseFlags(cCtx *cli.Context) scenclibase.CLIRunOptions {
	runOptions := &mc.RunScenarioOptions{
		ForceTraceGas: cCtx.Bool("force-trace-gas"),
	}

	vmBuilder := vmscenario.NewScenarioVMHostBuilder()
	if cCtx.Bool("wasmer1") {
		vmBuilder.OverrideVMExecutor = wasmer.ExecutorFactory()
	}
	if cCtx.Bool("wasmer2") {
		vmBuilder.OverrideVMExecutor = wasmer2.ExecutorFactory()
	}

	return scenclibase.CLIRunOptions{
		RunOptions: runOptions,
		VMBuilder:  vmBuilder,
	}
}
