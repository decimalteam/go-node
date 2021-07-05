package types

import (
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/helpers"
)

var _ sdk.Msg = &MsgCreateCoin{}

//type MsgCreateCoin struct {
//	Sender               sdk.AccAddress `json:"sender" yaml:"sender"`
//	Title                string         `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
//	Symbol               string         `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
//	ConstantReserveRatio uint           `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
//	InitialVolume        sdk.Int        `json:"initial_volume" yaml:"initial_volume"`
//	InitialReserve       sdk.Int        `json:"initial_reserve" yaml:"initial_reserve"`
//	LimitVolume          sdk.Int        `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
//	Identity             string         `json:"identity" yaml:"identity"`
//}

func NewMsgCreateCoin(sender sdk.AccAddress, title string, symbol string, crr uint, initVolume sdk.Int, initReserve sdk.Int, limitVolume sdk.Int, identity string) MsgCreateCoin {
	return MsgCreateCoin{
		Sender:               sender.String(),
		Title:                title,
		Symbol:               symbol,
		ConstantReserveRatio: uint64(crr),
		InitialVolume:        initVolume,
		InitialReserve:       initReserve,
		LimitVolume:          limitVolume,
		Identity:             identity,
	}
}

const CreateCoinConst = "create_coin"
const maxCoinNameBytes = 64
const allowedCoinSymbols = "^[a-zA-Z][a-zA-Z0-9]{2,9}$"

var minCoinSupply = sdk.NewInt(1)
var maxCoinSupply = helpers.BipToPip(sdk.NewInt(1000000000000000))

func MinCoinReserve(ctx sdk.Context) sdk.Int {
	return helpers.BipToPip(sdk.NewInt(1000))
}

func (msg *MsgCreateCoin) Route() string { return RouterKey }
func (msg *MsgCreateCoin) Type() string  { return CreateCoinConst }
func (msg *MsgCreateCoin) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Sender)

	return []sdk.AccAddress{accAddr}
}

func (msg *MsgCreateCoin) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCoin) ValidateBasic() error {
	// Validate coin title
	if len(msg.Title) > maxCoinNameBytes {
		return ErrInvalidCoinTitle()
	}
	// Validate coin symbol
	if match, _ := regexp.MatchString(allowedCoinSymbols, msg.Symbol); !match {
		return ErrInvalidCoinSymbol(msg.Symbol)
	}
	// Forbid creating coin with symbol DEL in testnet
	if strings.HasPrefix(config.ChainID, "decimal-testnet") {
		if strings.ToLower(msg.Symbol) == config.SymbolBaseCoin {
			return ErrForbiddenCoinSymbol(msg.Symbol)
		}
	}
	// Validate coin CRR
	if msg.ConstantReserveRatio < 10 || msg.ConstantReserveRatio > 100 {
		return ErrInvalidCRR()
	}
	// Check coin initial volume to be correct
	if msg.InitialVolume.LT(minCoinSupply) || msg.InitialVolume.GT(maxCoinSupply) {
		return ErrInvalidCoinInitialVolume(msg.InitialVolume.String())
	}

	if msg.InitialVolume.GT(msg.LimitVolume) {
		return ErrLimitVolumeBroken(msg.InitialVolume.String(), msg.LimitVolume.String())
	}
	return nil
}
