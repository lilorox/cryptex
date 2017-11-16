package main

import (
	"fmt"
	"os"

	"github.com/beldur/kraken-go-api-client"
	"github.com/urfave/cli"
)

type CryptexClient struct {
	API struct {
		Kraken *krakenapi.KrakenApi
	}
}

func (c *CryptexClient) ConnectKraken(apiKey, privKey string) {
	c.API.Kraken = krakenapi.New(apiKey, privKey)
}

func (c *CryptexClient) ListOpenOrders(ctx *cli.Context) error {
	result, err := c.API.Kraken.OpenOrders(nil)
	if err != nil {
		return err
	}

	fmt.Printf("Open orders: %+v\n", result)

	return nil
}

func main() {
	client := &CryptexClient{}

	app := cli.NewApp()
	app.Name = "cryptex"
	app.Usage = "Manage your account on different cryptocurrency exchange"

	app.Flags = []cli.Flag{
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
		if ctx.String("kraken-api-key") != "" && ctx.String("kraken-private-key") != "" {
			client.ConnectKraken(ctx.String("kraken-api-key"), ctx.String("kraken-private-key"))
		}
		return nil
	}

	app.Commands = []cli.Command{
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
	}

	app.Run(os.Args)
}
