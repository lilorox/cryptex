package output

import (
	"github.com/lilorox/cryptex/types"
)

type OutputFormatter interface {
	FormatError(error) string
	ShowBalance(*types.Balance)
}
