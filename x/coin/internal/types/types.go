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
	Reserve              sdk.Int `json:"reserve" yaml:"reserve"`
	LimitVolume          sdk.Int `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
	Volume               sdk.Int `json:"volume" yaml:"volume"`
}

func (c Coin) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Title: %s
		CRR: %d
		Symbol: %s
		Reserve: %s
		LimitVolume: %s
		Volume: %s
	`, c.Title, c.ConstantReserveRatio, c.Symbol, c.Reserve.String(), c.LimitVolume.String(), c.Volume.String()))
}

func (c Coin) IsBase() bool {
	return c.Symbol == "DEU"
}

//func (c Coin) UnmarshalJSON(bz []byte) error {
//	var alias Coin
//
//	err := json.Unmarshal(bz, &alias)
//	if err != nil {
//		return err
//	}
//
//	c.Volume = alias.Volume
//	c.Symbol = alias.Symbol
//	c.Title = alias.Title
//	c.ConstantReserveRatio = alias.ConstantReserveRatio
//	c.Reserve = alias.Reserve
//	c.LimitVolume = alias.LimitVolume
//
//	return nil
//}
