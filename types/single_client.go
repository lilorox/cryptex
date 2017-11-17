package types

import (
	"math/big"
)

type ClientBalance struct {
	Origin     string
	Currencies map[string]*big.Float
}
