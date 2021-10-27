package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"bufio"
	"fmt"


	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"

	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	govutils "bitbucket.org/decimalteam/go-node/x/gov/client/utils"
)

// Proposal flags
const (
	FlagTitle            = "title"
	FlagDescription      = "description"
	FlagVotingStartBlock = "voting-start-block"
	FlagVotingEndBlock   = "voting-end-block"
	flagVoter            = "voter"
	flagStatus           = "status"
	FlagProposal         = "proposal"
)

type proposal struct {
	Title            string
	Description      string
	VotingStartBlock uint64
	VotingEndBlock   uint64
}

// ProposalFlags defines the core required fields of a proposal. It is used to
// verify that these values are not provided in conjunction with a JSON proposal
// file.
var ProposalFlags = []string{
	FlagTitle,
	FlagDescription,
	FlagVotingStartBlock,
	FlagVotingEndBlock,
}

// GetTxCmd returns the transaction commands for this module
// governance ModuleClient is slightly different from other ModuleClients in that
// it contains a slice of "proposal" child commands. These commands are respective
// to proposal type handlers that are implemented in other modules but are mounted
// under the governance CLI (eg. parameter change proposals).
func GetTxCmd(storeKey string, cdc *codec.LegacyAmino) *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:                        types2.ModuleName,
		Short:                      "Governance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govTxCmd.AddCommand(
		GetCmdSubmitProposal(cdc),
		GetCmdVote(cdc),
		GetCmdSubmitUpgradeProposal(cdc),
	)

	return govTxCmd
}

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit a proposal along with an initial deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal along with an initial deposit.
Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Example:
$ %s tx gov submit-proposal --proposal="path/to/proposal.json" --from mykey

Where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "voting_start_block": 10000,
  "voting_end_block": 20000,
}

Which is equivalent to:

$ %s tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --voting_start_block 10000 --voting_end_block 20000 --from mykey
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc).WithInput(inBuf)

			proposal, err := parseSubmitProposalFlags(cmd.Flags())
			if err != nil {
				return err
			}

			msg := types2.NewMsgSubmitProposal(types2.Content{
				Title:       proposal.Title,
				Description: proposal.Description,
			}, clientCtx.GetFromAddress(), proposal.VotingStartBlock, proposal.VotingEndBlock)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "title of proposal")
	cmd.Flags().String(FlagDescription, "", "description of proposal")
	cmd.Flags().String(FlagVotingStartBlock, "", "start block of voting")
	cmd.Flags().String(FlagVotingEndBlock, "", "end block of voting")
	cmd.Flags().String(FlagProposal, "", "proposal file path (if this path is given, other proposal flags are ignored)")

	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%s query gov proposals".


Example:
$ %s tx gov vote 1 yes --from mykey
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc).WithInput(inBuf)

			// Get voting address
			from := clientCtx.GetFromAddress()

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// Find out which vote option user chose
			byteVoteOption, err := types2.VoteOptionFromString(govutils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := types2.NewMsgVote(sdk.ValAddress(from), proposalID, byteVoteOption)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

const (
	// TimeFormat specifies ISO UTC format for submitting the time for a new upgrade proposal
	TimeFormat = "2006-01-02T15:04:05Z"

	FlagUpgradeHeight = "upgrade-height"
	FlagUpgradeTime   = "time"
	FlagUpgradeInfo   = "upgrade-info"
	FlagToDownload    = "upgrade-to-download"
)

const (
	flagProposalType = "type"
)

func parseArgsToContent(cmd *cobra.Command, name string, proposer sdk.AccAddress) (types2.MsgSoftwareUpgradeProposal, error) {
	title, err := cmd.Flags().GetString(FlagTitle)
	if err != nil {
		return types2.MsgSoftwareUpgradeProposal{}, err
	}

	description, err := cmd.Flags().GetString(FlagDescription)
	if err != nil {
		return types2.MsgSoftwareUpgradeProposal{}, err
	}

	height, err := cmd.Flags().GetInt64(FlagUpgradeHeight)
	if err != nil {
		return types2.MsgSoftwareUpgradeProposal{}, err
	}

	timeStr, err := cmd.Flags().GetString(FlagUpgradeTime)
	if err != nil {
		return types2.MsgSoftwareUpgradeProposal{}, err
	}

	if height != 0 && len(timeStr) != 0 {
		return types2.MsgSoftwareUpgradeProposal{}, fmt.Errorf("only one of --upgrade-time or --upgrade-height should be specified")
	}

	var upgradeTime time.Time
	if len(timeStr) != 0 {
		upgradeTime, err = time.Parse(TimeFormat, timeStr)
		if err != nil {
			return types2.MsgSoftwareUpgradeProposal{}, err
		}
	}

	info, err := cmd.Flags().GetString(FlagUpgradeInfo)
	if err != nil {
		return types2.MsgSoftwareUpgradeProposal{}, err
	}

	toDownload, err := cmd.Flags().GetInt64(FlagToDownload)
	if err != nil {
		return types2.MsgSoftwareUpgradeProposal{}, err
	}

	plan := types2.Plan{Name: name, Time: upgradeTime, Height: height, Info: info, ToDownload: toDownload}
	msg := types2.NewSoftwareUpgradeProposal(title, description, plan, proposer)
	return msg, nil
}

// GetCmdSubmitUpgradeProposal implements a command handler for submitting a software upgrade proposal transaction.
func GetCmdSubmitUpgradeProposal(cdc *codec.LegacyAmino) *cobra.Command  {
	cmd := &cobra.Command{
		Use:   "software-upgrade [name] (--upgrade-height [height] | --upgrade-time [time]) (--upgrade-info [info]) [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a software upgrade proposal",
		Long: "Submit a software upgrade along with an initial deposit.\n" +
			"Please specify a unique name and height OR time for the upgrade to take effect.\n" +
			"You may include info to reference a binary download link, in a format compatible with: https://github.com/regen-network/cosmosd",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			inBuf := bufio.NewReader(cmd.InOrStdin())

			cliCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc).WithInput(inBuf)
			from := cliCtx.GetFromAddress()

			msg, err := parseArgsToContent(cmd, name, from)
			if err != nil {
				return err
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().Int64(FlagUpgradeHeight, 0, "The height at which the upgrade must happen (not to be used together with --upgrade-time)")
	cmd.Flags().String(FlagUpgradeTime, "", fmt.Sprintf("The time at which the upgrade must happen (ex. %s) (not to be used together with --upgrade-height)", TimeFormat))
	cmd.Flags().String(FlagUpgradeInfo, "", "Optional info for the planned upgrade such as commit hash, etc.")
	cmd.Flags().Int64(FlagToDownload, 0, "How many blocks before the update you need to start downloading the new version")

	return cmd
}

// DONTCOVER
