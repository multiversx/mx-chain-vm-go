package main

import (
	"github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"
	"github.com/urfave/cli"
)

func initializeCLI(facade *arwendebug.DebugFacade) *cli.App {
	app := cli.NewApp()
	app.Name = "Arwen Debug"
	app.Usage = ""

	args := &cliArguments{}

	// For server
	flagServerAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "",
		Value:       ":9091",
		Destination: &args.ServerAddress,
	}

	// Common for all actions
	flagDatabase := cli.StringFlag{
		Name:        "database",
		Usage:       "",
		Value:       "./db",
		Destination: &args.Database,
	}

	flagWorld := cli.StringFlag{
		Name:        "world",
		Usage:       "",
		Value:       "default",
		Destination: &args.World,
	}

	flagOutcome := cli.StringFlag{
		Required:    true,
		Name:        "outcome",
		Usage:       "",
		Destination: &args.Outcome,
	}

	// Common for contract actions
	flagContract := cli.StringFlag{
		Required:    true,
		Name:        "contract",
		Destination: &args.ContractAddress,
	}

	flagImpersonated := cli.StringFlag{
		Required:    true,
		Name:        "impersonated",
		Usage:       "",
		Destination: &args.Impersonated,
	}

	flagFunction := cli.StringFlag{
		Required:    true,
		Name:        "function",
		Destination: &args.Function,
	}

	flagArguments := cli.StringSliceFlag{
		Required: false,
		Name:     "arguments",
		Value:    &args.Arguments,
	}

	flagValue := cli.StringFlag{
		Name:        "value",
		Destination: &args.Value,
	}

	flagGasLimit := cli.Uint64Flag{
		Name:        "gas-limit",
		Destination: &args.GasLimit,
	}

	flagGasPrice := cli.Uint64Flag{
		Name:        "gas-price",
		Destination: &args.GasPrice,
	}

	// For deploy / upgrade
	flagCode := cli.StringFlag{
		Name:        "code",
		Destination: &args.Code,
	}

	flagCodePath := cli.StringFlag{
		Name:        "code-path",
		Destination: &args.CodePath,
	}

	flagCodeMetadata := cli.StringFlag{
		Name:        "code-metadata",
		Destination: &args.CodeMetadata,
	}

	// For create-account
	flagAccountAddress := cli.StringFlag{
		Required:    true,
		Name:        "address",
		Destination: &args.AccountAddress,
	}

	flagAccountBalance := cli.StringFlag{
		Required:    true,
		Name:        "balance",
		Usage:       "",
		Destination: &args.AccountBalance,
	}

	flagAccountNonce := cli.Uint64Flag{
		Required:    true,
		Name:        "nonce",
		Usage:       "",
		Destination: &args.AccountNonce,
	}

	app.Flags = []cli.Flag{}

	app.Authors = []cli.Author{
		{
			Name:  "The Elrond Team",
			Email: "contact@elrond.com",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "start debug server",
			Action: func(context *cli.Context) error {
				return arwendebug.StartServer(facade, args.ServerAddress)
			},
			Flags: []cli.Flag{
				flagServerAddress,
			},
		},
		{
			Name:  "deploy",
			Usage: "deploy a smart contract",
			Action: func(context *cli.Context) error {
				_, err := facade.DeploySmartContract(args.toDeployRequest())
				return err
			},
			Flags: []cli.Flag{
				flagOutcome,
				flagWorld,
				flagDatabase,
				flagImpersonated,
				flagCode,
				flagCodePath,
				flagCodeMetadata,
				flagArguments,
				flagValue,
				flagGasLimit,
				flagGasPrice,
			},
		},
		{
			Name:  "upgrade",
			Usage: "upgrade smart contract",
			Action: func(context *cli.Context) error {
				_, err := facade.UpgradeSmartContract(args.toUpgradeRequest())
				return err
			},
			Flags: []cli.Flag{
				flagOutcome,
				flagWorld,
				flagDatabase,
				flagContract,
				flagImpersonated,
				flagCode,
				flagCodePath,
				flagCodeMetadata,
				flagArguments,
				flagValue,
				flagGasLimit,
				flagGasPrice,
			},
		},
		{
			Name:  "run",
			Usage: "run smart contract",
			Action: func(context *cli.Context) error {
				_, err := facade.RunSmartContract(args.toRunRequest())
				return err
			},
			Flags: []cli.Flag{
				flagOutcome,
				flagWorld,
				flagDatabase,
				flagContract,
				flagImpersonated,
				flagFunction,
				flagArguments,
				flagValue,
				flagGasLimit,
				flagGasPrice,
			},
		},
		{
			Name:  "query",
			Usage: "query smart contract",
			Action: func(context *cli.Context) error {
				_, err := facade.QuerySmartContract(args.toQueryRequest())
				return err
			},
			Flags: []cli.Flag{
				flagOutcome,
				flagWorld,
				flagDatabase,
				flagContract,
				flagImpersonated,
				flagFunction,
				flagArguments,
			},
		},
		{
			Name:  "create-account",
			Usage: "create account",
			Action: func(context *cli.Context) error {
				_, err := facade.CreateAccount(args.toCreateAccountRequest())
				return err
			},
			Flags: []cli.Flag{
				flagOutcome,
				flagWorld,
				flagDatabase,
				flagAccountAddress,
				flagAccountBalance,
				flagAccountNonce,
			},
		},
	}

	return app
}
