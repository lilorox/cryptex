package client

import (
	"errors"
	"fmt"

	"github.com/lilorox/cryptex/types"
)

type CryptexClient struct {
	FiatReference string
	clients       []TradingClient
}

func New(fiat string) *CryptexClient {
	return &CryptexClient{
		FiatReference: fiat,
	}
}

func (c *CryptexClient) RegisterTradingClient(t TradingClient) {
	c.clients = append(c.clients, t)
}

func (c *CryptexClient) Balances() (*types.Balance, error) {
	if len(c.clients) == 0 {
		return nil, errors.New("No registered clients")
	}

	b := types.Balance{
		FiatCurrency:   c.FiatReference,
		ClientBalances: make([]types.ClientBalance, len(c.clients)),
	}

	for i, client := range c.clients {
		balance, err := client.Balance()
		if err != nil {
			return nil, fmt.Errorf("Client %s return error %s\n", client.GetName(), err.Error())
		}

		cb := types.ClientBalance{
			Origin:     client.GetName(),
			Currencies: balance,
		}
		b.ClientBalances[i] = cb
	}

	return &b, nil
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
