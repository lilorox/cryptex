package types

type Balance struct {
	FiatCurrency string

	ClientBalances []ClientBalance
}
