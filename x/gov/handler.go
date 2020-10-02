package gov

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

// NewHandler creates an sdk.Handler for all the gov type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {

		case types.MsgSubmitProposal:
			return handleMsgSubmitProposal(ctx, keeper, msg)

		case types.MsgVote:
			return handleMsgVote(ctx, keeper, msg)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleMsgSubmitProposal(ctx sdk.Context, keeper Keeper, msg types.MsgSubmitProposal) (*sdk.Result, error) {
	proposal, err := keeper.SubmitProposal(ctx, msg.Content, msg.VotingStartBlock, msg.VotingEndBlock)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSubmitProposal,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyProposalTitle, msg.Content.Title),
			sdk.NewAttribute(types.AttributeKeyProposalDescription, msg.Content.Description),
			sdk.NewAttribute(types.AttributeKeyProposalVotingStartBlock, strconv.FormatUint(msg.VotingStartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyProposalVotingEndBlock, strconv.FormatUint(msg.VotingEndBlock, 10)),
		),
	)

	return &sdk.Result{
		Data:   types.GetProposalIDBytes(proposal.ProposalID),
		Events: ctx.EventManager().Events(),
	}, nil
}

func handleMsgVote(ctx sdk.Context, keeper Keeper, msg types.MsgVote) (*sdk.Result, error) {
	err := keeper.AddVote(ctx, msg.ProposalID, msg.Voter, msg.Option)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalVote,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Voter.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, strconv.FormatUint(msg.ProposalID, 10)),
			sdk.NewAttribute(types.AttributeKeyOption, msg.Option.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
