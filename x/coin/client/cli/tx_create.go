package cli

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

func GetCmdCreateCoin(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [title] [symbol] [crr] [initReserve] [initVolume] [limitVolume]",
		Short: "Creates new coin",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			// Parsing parameters to variables
			var title = args[0]
			var symbol = args[1]
			var crr, err = strconv.ParseUint(args[2], 10, 8)
			// If error when convert crr
			if err != nil {
				return sdkerrors.New(types.DefaultCodespace, types.DecodeError, "Failed to convert CRR to uint")
			}
			var initReserve, _ = sdk.NewIntFromString(args[3])
			var initVolume, _ = sdk.NewIntFromString(args[4])
			var limitVolume, _ = sdk.NewIntFromString(args[5])
			// TODO: take reserve from creator and give it initial volume
			price := formulas.CalculateSaleReturn(initVolume, initReserve, uint(crr), sdk.NewIntWithDecimal(1, 18))
			_price := big.NewFloat(0).SetInt(price.BigInt())
			_price = _price.Quo(_price, big.NewFloat(1000000000000000000))
			fmt.Printf("Цена: (%v) tCDL \n", _price)

			msg := types.NewMsgCreateCoin(title, uint(crr), symbol, initVolume, initReserve, limitVolume, cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			// Check if coin does not exist yet
			coinExists, _ := cliUtils.ExistsCoin(cliCtx, symbol)
			if coinExists {
				return sdkerrors.New(types.DefaultCodespace, types.CoinAlreadyExists, fmt.Sprintf("Coin with symbol %s already exists", symbol))
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
