package arwendebug

import (
	"github.com/urfave/cli"
)

// CLIArguments -
type CLIArguments struct {
	ServerAddress   string
	DatabasePath    string
	Session         string
	ContractAddress string
	Action          string
	Function        string
	Arguments       []string
	Code            string
	CodePath        string
	CodeMetadata    string
}

// Initialize -
func Initialize(facade *DebugFacade) *cli.App {
	app := cli.NewApp()
	app.Name = "Arwen Debug"
	app.Usage = ""

	args := &CLIArguments{}

	flagServerAddress := cli.StringFlag{
		Name:        "address",
		Usage:       "",
		Value:       "localhost:9091",
		Destination: &args.ServerAddress,
	}

	flagSession := cli.StringFlag{
		Name:        "session",
		Usage:       "",
		Value:       "default",
		Destination: &args.Session,
	}

	app.Flags = []cli.Flag{
		flagSession,
	}

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
				facade.DeploySmartContract()
				return nil
			},
		},
		{
			Name:  "upgrade",
			Usage: "upgrade smart contract",
			Action: func(context *cli.Context) error {
				facade.UpgradeSmartContract()
				return nil
			},
		},
	}

	return app
}
