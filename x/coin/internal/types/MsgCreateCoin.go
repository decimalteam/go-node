package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"regexp"
)

var _ sdk.Msg = &MsgCreateCoin{}

type MsgCreateCoin struct {
	Creator              sdk.AccAddress `json:"creator" yams:"creator"`
	Title                string         `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
	ConstantReserveRatio uint           `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
	Symbol               string         `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
	InitialAmount        sdk.Int        `json:"initial_amount" yaml:"initial_amount"`
	InitialReserve       sdk.Int        `json:"initial_reserve" yaml:"initial_reserve"`
	LimitAmount          sdk.Int        `json:"limit_amount" yaml:"limit_amount"` // How many coins can be issued
}

func NewMsgCreateCoin(title string, crr uint, symbol string, initAmount sdk.Int, initReserve sdk.Int, limitAmount sdk.Int, creator sdk.AccAddress) MsgCreateCoin {
	return MsgCreateCoin{
		Creator:              creator,
		Title:                title,
		ConstantReserveRatio: crr,
		Symbol:               symbol,
		InitialAmount:        initAmount,
		InitialReserve:       initReserve,
		LimitAmount:          limitAmount,
	}
}

const CreateCoinConst = "CreateCoin"
const maxCoinNameBytes = 64
const allowedCoinSymbols = "^[A-Z0-9]{3,10}$"

var minCoinSupply = sdk.NewInt(1)
var maxCoinSupply = sdk.NewInt(1000000000000)

var minCoinReserve = sdk.NewInt(10)

func (msg MsgCreateCoin) Route() string { return RouterKey }
func (msg MsgCreateCoin) Type() string  { return CreateCoinConst }
func (msg MsgCreateCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

func (msg MsgCreateCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateCoin) ValidateBasic() sdk.Error {
	// Check coin CRR validity
	if msg.ConstantReserveRatio < 10 || msg.ConstantReserveRatio > 100 {
		return sdk.NewError(DefaultCodespace, InvalidCRR, "Coin CRR must be between 10 and 100")
	}
	// Check coin title maximum length
	if len(msg.Title) > maxCoinNameBytes {
		return sdk.NewError(DefaultCodespace, InvalidCoinTitle, fmt.Sprintf("Coin name is invalid. Allowed up to %d bytes.", maxCoinNameBytes))
	}
	// Check coin symbol for correct regexp
	if match, _ := regexp.MatchString(allowedCoinSymbols, msg.Symbol); !match {
		return sdk.NewError(DefaultCodespace, InvalidCoinSymbol, fmt.Sprintf("Invalid coin symbol. Should be %s", allowedCoinSymbols))
	}
	// Check coin initial amount to be correct
	if msg.InitialAmount.LT(minCoinSupply) || msg.InitialAmount.GT(maxCoinSupply) {
		return sdk.NewError(DefaultCodespace, InvalidCoinInitAmount, fmt.Sprintf("Coin initial amount should be between %s and %s. Given %s", minCoinSupply.String(), maxCoinSupply.String(), msg.InitialAmount.String()))
	}
	// Check coin initial reserve to be correct
	if msg.InitialReserve.LT(minCoinReserve) {
		return sdk.NewError(DefaultCodespace, InvalidCoinInitReserve, fmt.Sprintf("Coin reserve should be greater than or equal to %s", minCoinReserve.String()))
	}
	return nil
}
