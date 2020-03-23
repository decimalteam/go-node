package cli

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/check/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"strconv"
)

func GetCmdIssueCheck(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "issue [coin] [amount] [gasCoin] [untilBlock] [nonce] [passphrase]",
		Short: "Returns coin information by symbol",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			coinSymbol := args[0]
			amount, _ := sdk.NewIntFromString(args[1])
			gasCoinSymbol := args[2]
			untilBlock, _ := sdk.NewIntFromString(args[3])
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			passphrase := args[5]
			// Можно ли выпустить чек с несуществующей валютой?

			_, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/coin/get/%s", coinSymbol), nil)
			if err != nil {
				fmt.Printf("Could not resolve coin %s \n%s\n", coinSymbol, err.Error())

				return nil
			}
			_, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/coin/get/%s", gasCoinSymbol), nil)
			if err != nil {
				fmt.Printf("Could not resolve coin %s \n%s\n", gasCoinSymbol, err.Error())

				return nil
			}
			if sdk.NewInt(cliCtx.Height).GT(untilBlock) {
				fmt.Printf("Until block less than current block. Current block: %d", cliCtx.Height)

				return nil
			}
			if amount.LTE(sdk.NewInt(0)) {
				fmt.Printf("Invalid amount. Should be more than 0.")

				return nil
			}
			check := types.NewCheck(nonce, config.ChainID, untilBlock.BigInt().Uint64(), coinSymbol, &amount, gasCoinSymbol)
			signed, err := check.Sign(passphrase)
			if err != nil {
				fmt.Printf("Signing check error: %s", err.Error())
				return nil
			}
			hash, err := signed.Encode()
			if err != nil {
				fmt.Printf("Encoding check error: %s", err.Error())
			}
			fmt.Printf("Check: %s", hash)
			return nil
		},
	}
}
