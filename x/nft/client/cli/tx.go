package cli

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

// Edit metadata flags
const (
	flagTokenURI = "tokenURI"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.LegacyAmino) *cobra.Command {
	nftTxCmd := &cobra.Command{
		Use:   types.ModuleName,
		Short: "NFT transactions subcommands",
		RunE:  client.ValidateCmd,
	}

	nftTxCmd.AddCommand(flags.PostCommands(
		GetCmdTransferNFT(cdc),
		GetCmdEditNFTMetadata(cdc),
		GetCmdMintNFT(cdc),
		GetCmdBurnNFT(cdc),
	)...)

	return nftTxCmd
}

// GetCmdTransferNFT is the CLI command for sending a TransferNFT transaction
func GetCmdTransferNFT(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [sender] [recipient] [denom] [tokenID] [quantity]",
		Short: "transfer a NFT to a recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Transfer a NFT from a given collection that has a 
			specific id (SHA-256 hex hash) to a specific recipient.

Example:
$ %s tx %s transfer 
dx1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p dx1l2rsakp388kuv9k8qzq6lrm9taddae7fpx59wm \
crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
--from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sender, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denom := args[2]
			tokenID := args[3]

			quantity, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid quantity")
			}

			msg := types.NewMsgTransferNFT(sender, recipient, denom, tokenID, quantity)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdEditNFTMetadata is the CLI command for sending an EditMetadata transaction
func GetCmdEditNFTMetadata(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-metadata [denom] [tokenID]",
		Short: "edit the metadata of an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Edit the metadata of an NFT from a given collection that has a 
			specific id (SHA-256 hex hash).

Example:
$ %s tx %s edit-metadata crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
--tokenURI path_to_token_URI_JSON --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenID := args[1]
			tokenURI := viper.GetString(flagTokenURI)

			msg := types.NewMsgEditNFTMetadata(cliCtx.GetFromAddress(), tokenID, denom, tokenURI)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTokenURI, "", "Extra properties available for querying")
	return cmd
}

// GetCmdMintNFT is the CLI command for a MintNFT transaction
func GetCmdMintNFT(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint [denom] [tokenID] [recipient] [quantity] [reserve] [allow_mint]",
		Short: "mint an NFT and set the owner to the recipient",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Mint an NFT from a given collection that has a 
			specific id (SHA-256 hex hash) and set the ownership to a specific address.

Example:
$ %s tx %s mint crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
dx1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenID := args[1]

			recipient, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return err
			}

			tokenURI := viper.GetString(flagTokenURI)

			quantity, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return errors.New("invalid quantity")
			}

			reserve, ok := sdk.NewIntFromString(args[4])
			if !ok {
				return errors.New("invalid quantity")
			}

			var allowMint bool
			if args[5] == "t" {
				allowMint = true
			}

			msg := types.NewMsgMintNFT(cliCtx.GetFromAddress(), recipient, tokenID, denom, tokenURI, quantity, reserve, allowMint)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagTokenURI, "", "URI for supplemental off-chain metadata (should return a JSON object)")

	return cmd
}

// GetCmdBurnNFT is the CLI command for sending a BurnNFT transaction
func GetCmdBurnNFT(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "burn [denom] [tokenID] [quantity]",
		Short: "burn an NFT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn (i.e permanently delete) an NFT from a given collection that has a 
			specific id (SHA-256 hex hash).

Example:
$ %s tx %s burn crypto-kitties d04b98f48e8f8bcc15c6ae5ac050801cd6dcfd428fb5f9e65c4e16e7807340fa \
--from mykey
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := authtypes.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			denom := args[0]
			tokenID := args[1]

			quantity, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return errors.New("invalid quantity")
			}

			msg := types.NewMsgBurnNFT(cliCtx.GetFromAddress(), tokenID, denom, quantity)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
