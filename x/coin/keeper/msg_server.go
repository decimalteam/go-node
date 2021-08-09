package keeper

import (
	"context"

	"bitbucket.org/decimalteam/go-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) BuyCoin(ctx context.Context, msg *types.MsgBuyCoin) (*types.MsgBuyCoinResponse, error) {
	panic("implement me")
}

func (k msgServer) CreateCoin(ctx context.Context, msg *types.MsgCreateCoin) (*types.MsgCreateCoinResponse, error) {
	panic("implement me")
}

func (k msgServer) MultiSendCoin(ctx context.Context, msg *types.MsgMultiSendCoin) (*types.MsgMultisendCoinResponse, error) {
	panic("implement me")
}

func (k msgServer) RedeemCheck(ctx context.Context, msg *types.MsgRedeemCheck) (*types.MsgRedeemCheckResponse, error) {
	panic("implement me")
}

func (k msgServer) SellAllCoin(ctx context.Context, msg *types.MsgSellAllCoin) (*types.MsgSellAllCoinResponse, error) {
	panic("implement me")
}

func (k msgServer) SellCoin(ctx context.Context, msg *types.MsgSellCoin) (*types.MsgSellCoinReponse, error) {
	panic("implement me")
}

func (k msgServer) SendCoin(goCtx context.Context, msg *types.MsgSendCoin) (*types.MsgSendCoinResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	from, _ := sdk.AccAddressFromBech32(msg.Sender)
	to, _ := sdk.AccAddressFromBech32(msg.Receiver)

	err := k.BankKeeper.SendCoins(ctx, from, to, sdk.Coins{msg.Coin})
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types.AttributeCoin, msg.Coin.String()),
		sdk.NewAttribute(types.AttributeReceiver, msg.Receiver),
	))

	return &types.MsgSendCoinResponse{}, err
}

func (k msgServer) UpdateCoin(ctx context.Context, msg *types.MsgUpdateCoin) (*types.MsgUpdateCoinResponse, error) {
	panic("implement me")
}
