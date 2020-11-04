package cli

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/spf13/cobra"
	"strings"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	swapQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	swapQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQuerySwap(queryRoute, cdc),
		GetCmdQueryActiveSwap(queryRoute, cdc),
		GetCmdQueryPool(queryRoute, cdc),
	)...)

	return swapQueryCmd
}

func GetCmdQuerySwap(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "swap [hashed_secret]",
		Short: "Query a swap",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			hashRaw, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}
			var hash types.Hash
			copy(hash[:], hashRaw)

			res, _, err := cliCtx.QueryStore(types.GetSwapKey(hash), storeName)

			var swap types.Swap
			cdc.MustUnmarshalJSON(res, &swap)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(swap)
		},
	}
}

func GetCmdQueryActiveSwap(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "active-swap",
		Args:  cobra.NoArgs,
		Short: "Query all active swaps",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", storeName, types.QueryActiveSwaps), nil)
			if err != nil {
				return err
			}

			var swaps types.Swaps
			if err := cdc.UnmarshalJSON(bz, &swaps); err != nil {
				return err
			}

			return cliCtx.PrintOutput(swaps)
		},
	}
}

func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "Query the swap pool",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values for amounts stored in the swap pool.

Example:
$ %s query swap pool
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", storeName), nil)
			if err != nil {
				return err
			}

			var pool exported.ModuleAccountI
			if err := cdc.UnmarshalJSON(bz, &pool); err != nil {
				return err
			}

			return cliCtx.PrintOutput(pool)
		},
	}
}
