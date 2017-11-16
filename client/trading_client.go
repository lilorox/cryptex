package client

import (
	"math/big"
)

type TradingClient interface {
	GetName() string
	Balance() (map[string]*big.Float, error)
}
