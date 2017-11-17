package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/lilorox/cryptex/client"
	"github.com/lilorox/cryptex/client/kraken"
	"github.com/lilorox/cryptex/output"
)

var formatter output.OutputFormatter

const (
	BadArgument = iota
	BalanceError
)

func handleError(err error, statusCode int) *cli.ExitError {
	return cli.NewExitError(formatter.FormatError(err), statusCode)
}

func main() {
	cryptex := client.New("EUR")

	app := cli.NewApp()
	app.Name = "cryptex"
	app.Usage = "Manage your account on different cryptocurrency exchange"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "output",
			Usage:  "Specify the output type: pretty (default), json",
			Value:  "pretty",
			EnvVar: "OUTPUT",
		},
		cli.StringFlag{
			Name:   "cex-userid",
			Usage:  "CEX.io UserID",
			EnvVar: "CEX_USERID",
		},
		cli.StringFlag{
			Name:   "cex-api-key",
			Usage:  "CEX.io API key",
			EnvVar: "CEX_API_KEY",
		},
		cli.StringFlag{
			Name:   "cex-secret-key",
			Usage:  "CEX.io Secret key",
			EnvVar: "CEX_SECRET_KEY",
		},
		cli.StringFlag{
			Name:   "kraken-api-key",
			Usage:  "Kraken API key",
			EnvVar: "KRAKEN_API_KEY",
		},
		cli.StringFlag{
			Name:   "kraken-private-key",
			Usage:  "Kraken API private key",
			EnvVar: "KRAKEN_PRIV_KEY",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		switch ctx.String("output") {
		case "pretty":
			formatter = &output.PrettyFormatter{}
		default:
			return cli.NewExitError(fmt.Sprintf("Unknown output type '%s'", ctx.String("output")), BadArgument)
		}

		if ctx.String("kraken-api-key") != "" && ctx.String("kraken-private-key") != "" {
			k := kraken.New(
				ctx.String("kraken-api-key"),
				ctx.String("kraken-private-key"),
			)
			cryptex.RegisterTradingClient(k)
		}

		/*
			if ctx.String("cex-userid") != "" && ctx.String("cex-api-key") != "" && ctx.String("cex-secret-key") != "" {
				client.ConnectKraken(ctx.String("kraken-api-key"), ctx.String("kraken-private-key"))
			}
		*/
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "balance",
			Aliases: []string{"b"},
			Usage:   "Show current balances",
			Action: func(ctx *cli.Context) error {
				b, err := cryptex.Balances()
				if err != nil {
					return handleError(err, BalanceError)
				}
				formatter.ShowBalance(b)
				return nil
			},
		},
		/*
			cli.Command{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List open orders",
				Subcommands: []cli.Command{
					cli.Command{
						Name:    "orders",
						Aliases: []string{"o"},
						Usage:   "Manipulate orders",
						Action:  client.ListOpenOrders,
					},
				},
			},
		*/
	}

	app.Run(os.Args)
}
