package cli

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdGetCoin(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get [symbol]",
		Short: "Returns coin information by symbol",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			symbol := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryGetCoin, symbol), nil)
			if err != nil {
				fmt.Printf("could not resolve coin %s \n%s\n", symbol, err.Error())

				return nil
			}

			var out types.Coin
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
