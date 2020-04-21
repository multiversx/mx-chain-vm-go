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

	flagServerAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "",
		Value:       ":9091",
		Destination: &args.ServerAddress,
	}

	flagDatabase := cli.StringFlag{
		Name:        "database",
		Usage:       "",
		Value:       "./db",
		Destination: &args.Database,
	}

	flagSession := cli.StringFlag{
		Name:        "session",
		Usage:       "",
		Value:       "default",
		Destination: &args.Session,
	}

	flagImpersonated := cli.StringFlag{
		Required:    true,
		Name:        "impersonated",
		Usage:       "",
		Destination: &args.Impersonated,
	}

	flagAccountAddress := cli.StringFlag{
		Required:    true,
		Name:        "address",
		Usage:       "",
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
				facade.StartServer(args.ServerAddress)
				return nil
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
				flagSession,
				flagDatabase,
				flagImpersonated,
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
				flagSession,
				flagDatabase,
				flagImpersonated,
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
				flagSession,
				flagDatabase,
				flagAccountAddress,
				flagAccountBalance,
				flagAccountNonce,
			},
		},
	}

	return app
}
