package output

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mgutz/ansi"

	"github.com/lilorox/cryptex/types"
)

const (
	LightHorizontal = "\u2500"
	LightDownRight  = "\u250c"
	LightDownLeft   = "\u2510"
	LightUpRight    = "\u2514"
	LightUpLeft     = "\u2518"
)

var brightblue = ansi.ColorFunc("blue+h")

type PrettyFormatter struct{}

func (p *PrettyFormatter) FormatError(err error) string {
	return ansi.Color(err.Error(), "red+b")
}

func (p *PrettyFormatter) ShowBalance(b *types.Balance) {
	first := true
	for _, cb := range b.ClientBalances {
		if first {
			first = false
		} else {
			fmt.Println()
		}

		title := fmt.Sprintf("Balance on %s :", cb.Origin)
		underline := strings.Repeat(LightHorizontal, utf8.RuneCountInString(title)+1)
		fmt.Printf(
			"%s %s\n%s\n",
			brightblue(LightDownRight),
			title,
			brightblue(LightUpRight+underline),
		)
		for cur, value := range cb.Currencies {
			fmt.Printf("  %- 10s %f\n", cur, value)
		}
	}
}
