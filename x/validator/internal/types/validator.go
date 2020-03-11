package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

type Validator struct {
	ValAddress sdk.ValAddress `json:"val_address"`
	PubKey     crypto.PubKey  `json:"pub_key"`
	StakeCoins sdk.Coins      `json:"stake_coins"`
	Status     uint8          `json:"status"`
}
