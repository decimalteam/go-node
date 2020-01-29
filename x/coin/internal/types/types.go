package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type Coin struct {
	Title                string  `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
	ConstantReserveRatio uint    `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
	Symbol               string  `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
	InitialAmount        sdk.Int `json:"initial_amount" yaml:"initial_amount"`
	InitialReserve       sdk.Int `json:"initial_reserve" yaml:"initial_reserve"`
	LimitAmount          sdk.Int `json:"limit_amount" yaml:"limit_amount"` // How many coins can be issued
}

func (c Coin) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Title: %s
		CRR: %d
		Symbol: %s
		InitAmount: %s
		InitReserve: %s
		LimitAmount: %s
	`, c.Title, c.ConstantReserveRatio, c.Symbol, c.InitialAmount.String(), c.InitialReserve.String(), c.LimitAmount.String()))
}
