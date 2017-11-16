package kraken

import (
	"math/big"
	//"encoding/json"
)

// GenericResponse wraps the Kraken API JSON response
type GenericResponse struct {
	Error  []string    `json:"error"`
	Result interface{} `json:"result"`
}

type BalanceResponse map[string]*big.Float
