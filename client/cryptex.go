package client

import (
	"fmt"
	"strings"
	"unicode/utf8"
	//"github.com/beldur/kraken-go-api-client"
)

type CryptexClient struct {
	DefaultFiat string
	clients     []TradingClient
}

func (c *CryptexClient) RegisterTradingClient(t TradingClient) {
	c.clients = append(c.clients, t)
}

func (c *CryptexClient) ShowBalances() error {
	for _, client := range c.clients {
		balance, err := client.Balance()
		if err != nil {
			fmt.Printf("Client %s return error %s\n", client.GetName(), err.Error())
			return err
		}
		name := client.GetName()
		underline := strings.Repeat("\u2500", utf8.RuneCountInString(name)+14)
		fmt.Printf(
			"\u250c Balance on %s :\n\u2514%s\n",
			client.GetName(),
			underline,
		)
		for cur, value := range balance {
			fmt.Printf("  %- 10s %f\n", cur, value)
		}
	}

	return nil
}

/*
func (c *CryptexClient) ListOpenOrders(ctx *cli.Context) error {
	result, err := c.API.Kraken.OpenOrders(nil)
	if err != nil {
		return err
	}

	fmt.Printf("Open orders: %+v\n", result)

	return nil
}
*/
