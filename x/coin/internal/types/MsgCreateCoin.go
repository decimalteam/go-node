package types

import (
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

var minCoinReserve = helpers.BipToPip(sdk.NewInt(10000))

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
		return sdkerrors.New(DefaultCodespace, InvalidCoinTitle, fmt.Sprintf("Coin name is invalid. Allowed up to %d bytes.", maxCoinNameBytes))
	}
	// Validate coin symbol
	if match, _ := regexp.MatchString(allowedCoinSymbols, msg.Symbol); !match {
		return sdkerrors.New(DefaultCodespace, InvalidCoinSymbol, fmt.Sprintf("Invalid coin symbol. Should be %s", allowedCoinSymbols))
	}
	// Validate coin CRR
	if msg.ConstantReserveRatio < 10 || msg.ConstantReserveRatio > 100 {
		return sdkerrors.New(DefaultCodespace, InvalidCRR, "Coin CRR must be between 10 and 100")
	}
	// Check coin initial volume to be correct
	if msg.InitialVolume.LT(minCoinSupply) || msg.InitialVolume.GT(maxCoinSupply) {
		return sdkerrors.New(DefaultCodespace, InvalidCoinInitVolume, fmt.Sprintf("Coin initial volume should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), msg.InitialVolume.String()))
	}
	// Check coin initial reserve to be correct
	if msg.InitialReserve.LT(minCoinReserve) {
		return sdkerrors.New(DefaultCodespace, InvalidCoinInitReserve, fmt.Sprintf("Coin initial reserve should be greater than or equal to %s", minCoinReserve.String()))
	}
	return nil
}
