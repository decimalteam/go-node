package cli

import (
	gcutils "bitbucket.org/decimalteam/go-node/x/gov/client/utils"
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        types2.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(
		GetCmdQueryProposal(queryRoute, cdc),
		GetCmdQueryProposals(queryRoute, cdc),
		GetCmdQueryVote(queryRoute, cdc),
		GetCmdQueryVotes(queryRoute, cdc),
		GetCmdQueryParam(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryProposer(queryRoute, cdc),
		GetCmdQueryTally(queryRoute, cdc),
	)

	return govQueryCmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a proposal. You can find the
proposal-id by running "%s query gov proposals".

Example:
$ %s query gov proposal 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// Query the proposal
			res, err := gcutils.QueryProposalByID(proposalID, clientCtx, queryRoute)
			if err != nil {
				return err
			}

			var proposal types2.Proposal
			cdc.MustUnmarshalJSON(res, &proposal)
			return clientCtx.PrintProto(&proposal) // nolint:errcheck
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryProposals implements a query proposals command.
func GetCmdQueryProposals(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated proposals that match optional filters:

Example:
$ %s query gov proposals --depositor cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query gov proposals --voter cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
$ %s query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
$ %s query gov proposals --page=2 --limit=100
`,
				version.AppName, version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			bechVoterAddr, _ := cmd.Flags().GetString(flagVoter)
			strProposalStatus, _ := cmd.Flags().GetString(flagStatus)
			page := viper.GetInt(flags.FlagPage)
			limit := viper.GetInt(flags.FlagLimit)

			var depositorAddr sdk.AccAddress
			var voterAddr sdk.AccAddress
			var proposalStatus types2.ProposalStatus

			params := types2.NewQueryProposalsParams(page, limit, proposalStatus, voterAddr, depositorAddr)

			if len(bechVoterAddr) != 0 {
				voterAddr, err := sdk.AccAddressFromBech32(bechVoterAddr)
				if err != nil {
					return err
				}
				params.Voter = voterAddr
			}

			if len(strProposalStatus) != 0 {
				proposalStatus, err := types2.ProposalStatusFromString(gcutils.NormalizeProposalStatus(strProposalStatus))
				if err != nil {
					return err
				}
				params.ProposalStatus = proposalStatus
			}

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/proposals", queryRoute), bz)
			if err != nil {
				return err
			}

			var matchingProposals types2.Proposals
			err = cdc.UnmarshalJSON(res, &matchingProposals)
			if err != nil {
				return err
			}

			if len(matchingProposals) == 0 {
				return fmt.Errorf("no matching proposals found")
			}

			return clientCtx.PrintObjectLegacy(matchingProposals) // nolint:errcheck
		},
	}

	cmd.Flags().Int(flags.FlagPage, 1, "pagination page of proposals to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of proposals to query for")
	cmd.Flags().String(flagVoter, "", "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// Command to Get a Proposal Information
// GetCmdQueryVote implements the query proposal vote command.
func GetCmdQueryVote(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.

Example:
$ %s query gov vote 1 cosmos1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = gcutils.QueryProposalByID(proposalID, clientCtx, queryRoute)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			voterAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			params := types2.NewQueryVoteParams(proposalID, voterAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/vote", queryRoute), bz)
			if err != nil {
				return err
			}

			var vote types2.Vote

			// XXX: Allow the decoding to potentially fail as the vote may have been
			// pruned from state. If so, decoding will fail and so we need to check the
			// Empty() case. Consider updating Vote JSON decoding to not fail when empty.
			_ = cdc.UnmarshalJSON(res, &vote)

			if vote.Empty() {
				res, err = gcutils.QueryVoteByTxQuery(clientCtx, params)
				if err != nil {
					return err
				}

				if err := cdc.UnmarshalJSON(res, &vote); err != nil {
					return err
				}
			}

			return clientCtx.PrintString(vote.String())
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vote details for a single proposal by its identifier.

Example:
$ %[1]s query gov votes 1
$ %[1]s query gov votes 1 --page=2 --limit=100
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			page := viper.GetInt(flags.FlagPage)
			limit := viper.GetInt(flags.FlagLimit)

			params := types2.NewQueryProposalVotesParams(proposalID, page, limit)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// check to see if the proposal is in the store
			res, err := gcutils.QueryProposalByID(proposalID, clientCtx, queryRoute)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			var proposal types2.Proposal
			cdc.MustUnmarshalJSON(res, &proposal)

			propStatus := proposal.Status
			if !(propStatus == types2.StatusVotingPeriod) {
				res, err = gcutils.QueryVotesByTxQuery(clientCtx, params)
			} else {
				res, _, err = clientCtx.QueryWithData(fmt.Sprintf("custom/%s/votes", queryRoute), bz)
			}

			if err != nil {
				return err
			}

			var votes types2.Votes
			cdc.MustUnmarshalJSON(res, &votes)
			return clientCtx.PrintObjectLegacy(votes)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Int(flags.FlagPage, 1, "pagination page of votes to to query for")
	cmd.Flags().Int(flags.FlagLimit, 100, "pagination limit of votes to query for")
	return cmd
}

// GetCmdQueryTally implements the command to query for proposal tally result.
func GetCmdQueryTally(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a proposal vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query tally of votes on a proposal. You can find
the proposal-id by running "%s query gov proposals".

Example:
$ %s query gov tally 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = gcutils.QueryProposalByID(proposalID, clientCtx, queryRoute)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// Construct query
			params := types2.NewQueryProposalParams(proposalID)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			// Query store
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/tally", queryRoute), bz)
			if err != nil {
				return err
			}

			var tally types2.TallyResult
			cdc.MustUnmarshalJSON(res, &tally)
			return clientCtx.PrintObjectLegacy(tally)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryParams(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %s query gov params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)
			tp, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/params/tallying", queryRoute), nil)
			if err != nil {
				return err
			}

			var tallyParams types2.TallyParams
			cdc.MustUnmarshalJSON(tp, &tallyParams)

			return clientCtx.PrintString(types2.NewParams(tallyParams).String())
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryParam(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param [param-type]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the parameters (voting|tallying|deposit) of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.

Example:
$ %s query gov param voting
$ %s query gov param tallying
$ %s query gov param deposit
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			// Query store
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/params/%s", queryRoute, args[0]), nil)
			if err != nil {
				return err
			}
			var out fmt.Stringer
			switch args[0] {
			case "tallying":
				var param types2.TallyParams
				cdc.MustUnmarshalJSON(res, &param)
				out = param
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
			}

			return clientCtx.PrintObjectLegacy(out)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryProposer implements the query proposer command.
func GetCmdQueryProposer(queryRoute string, cdc *codec.LegacyAmino) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposer [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the proposer of a governance proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query which address proposed a proposal with a given ID.

Example:
$ %s query gov proposer 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			// validate that the proposalID is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint", args[0])
			}

			prop, err := gcutils.QueryProposerByTxQuery(clientCtx, proposalID)
			if err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(prop)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
