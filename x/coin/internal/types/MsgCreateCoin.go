package types

import (
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateCoin{}

type MsgCreateCoin struct {
	Sender               sdk.AccAddress `json:"sender" yaml:"sender"`
	Title                string         `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
	Symbol               string         `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
	ConstantReserveRatio uint           `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
	InitialVolume        sdk.Int        `json:"initial_volume" yaml:"initial_volume"`
	InitialReserve       sdk.Int        `json:"initial_reserve" yaml:"initial_reserve"`
	LimitVolume          sdk.Int        `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
}

func NewMsgCreateCoin(sender sdk.AccAddress, title string, symbol string, crr uint, initVolume sdk.Int, initReserve sdk.Int, limitVolume sdk.Int) MsgCreateCoin {
	return MsgCreateCoin{
		Sender:               sender,
		Title:                title,
		Symbol:               symbol,
		ConstantReserveRatio: crr,
		InitialVolume:        initVolume,
		InitialReserve:       initReserve,
		LimitVolume:          limitVolume,
	}
}

const CreateCoinConst = "create_coin"
const maxCoinNameBytes = 64
const allowedCoinSymbols = "^[a-zA-Z][a-zA-Z0-9]{2,9}$"

var minCoinSupply = sdk.NewInt(1)
var maxCoinSupply = helpers.BipToPip(sdk.NewInt(1000000000000000))

var MinCoinReserve = helpers.BipToPip(sdk.NewInt(10000))

func (msg MsgCreateCoin) Route() string { return RouterKey }
func (msg MsgCreateCoin) Type() string  { return CreateCoinConst }
func (msg MsgCreateCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgCreateCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateCoin) ValidateBasic() error {
	// Validate coin title
	if len(msg.Title) > maxCoinNameBytes {
		return ErrInvalidCoinTitle()
	}
	// Validate coin symbol
	if match, _ := regexp.MatchString(allowedCoinSymbols, msg.Symbol); !match {
		return ErrInvalidCoinSymbol(msg.Symbol)
	}
	// Validate coin CRR
	if msg.ConstantReserveRatio < 10 || msg.ConstantReserveRatio > 100 {
		return ErrInvalidCRR()
	}
	// Check coin initial volume to be correct
	if msg.InitialVolume.LT(minCoinSupply) || msg.InitialVolume.GT(maxCoinSupply) {
		return ErrInvalidCoinInitialVolume(msg.InitialVolume.String())
	}
	// Check coin initial reserve to be correct
	if msg.InitialReserve.LT(MinCoinReserve) {
		return ErrInvalidCoinInitialReserve()
	}

	if msg.InitialVolume.GT(msg.LimitVolume) {
		return ErrLimitVolumeBroken(msg.InitialVolume.String(), msg.LimitVolume.String())
	}
	return nil
}
