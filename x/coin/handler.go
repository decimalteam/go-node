package coin

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

// NewHandler creates an sdk.Handler for all the coin type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		//ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
		case types.MsgBuyCoin:
			return handleMsgBuyCoin(ctx, k, msg)
		case types.MsgSellCoin:
			return handleMsgSellCoin(ctx, k, msg)
		case types.MsgSendCoin:
			return handleMsgSendCoin(ctx, k, msg)
		case types.MsgMultiSendCoin:
			return handleMsgMultiSendCoin(ctx, k, msg)
		case types.MsgSellAllCoin:
			return handleMsgSellAllCoin(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg types.MsgCreateCoin) (*sdk.Result, error) {
	var coin = types.Coin{
		Title:       msg.Title,
		CRR:         msg.ConstantReserveRatio,
		Symbol:      msg.Symbol,
		Reserve:     msg.InitialReserve,
		LimitVolume: msg.LimitVolume,
		Volume:      msg.InitialVolume,
	}
	// TODO: take reserve from creator and give it initial volume
	acc := k.AccountKeeper.GetAccount(ctx, msg.Creator)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(msg.InitialReserve) {
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "")
	}

	err := k.UpdateBalance(ctx, strings.ToLower(cliUtils.GetBaseCoin()), msg.InitialReserve.Neg(), msg.Creator)
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, "")
	}

	k.SetCoin(ctx, coin)
	err = k.UpdateBalance(ctx, strings.ToLower(coin.Symbol), msg.InitialVolume, msg.Creator)
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, "")
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateCoin),
			sdk.NewAttribute(types.AttributeSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeTitle, msg.Title),
			sdk.NewAttribute(types.AttributeInitVolume, msg.InitialVolume.String()),
			sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
			sdk.NewAttribute(types.AttributeCRR, strconv.FormatUint(uint64(msg.ConstantReserveRatio), 10)),
			sdk.NewAttribute(types.AttributeLimitVolume, msg.LimitVolume.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgBuyCoin(ctx sdk.Context, k Keeper, msg types.MsgBuyCoin) (*sdk.Result, error) {
	// TODO: Commission

	coinToBuy, _ := k.GetCoin(ctx, msg.CoinToBuy)
	if coinToBuy.Symbol != msg.CoinToBuy {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	// Check if coin to sell exists
	coinToSell, _ := k.GetCoin(ctx, msg.CoinToSell)
	if coinToSell.Symbol != msg.CoinToSell {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}

	amountBuy, amountSell, calcErr := cliUtils.BuyCoinCalculateAmounts(coinToBuy, coinToSell, msg.AmountToBuy, msg.AmountToSell)

	if calcErr != nil {
		return nil, sdkerrors.New(calcErr.Codespace(), calcErr.ABCICode(), calcErr.Error())
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Buyer)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(msg.CoinToSell)).LT(amountSell) {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, errMsg)
	}

	err := k.UpdateBalance(ctx, msg.CoinToBuy, amountBuy, msg.Buyer)
	if err != nil {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountSell.Neg(), msg.Buyer)
	if err != nil {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	amountBuyInBaseCoin := formulas.CalculateSaleReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountBuy)
	amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)

	k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountSellInBaseCoin), coinToSell.Volume.Sub(amountSell))
	k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountBuyInBaseCoin), coinToBuy.Volume.Add(amountBuy))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeBuyCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountBuy.String()),
			sdk.NewAttribute(types.AttributeAmountToSell, amountSell.String()),
			sdk.NewAttribute(types.AttributeAmountToBuyInBaseCoin, amountBuyInBaseCoin.String()),
			sdk.NewAttribute(types.AttributeAmountToSellInBaseCoin, amountSellInBaseCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSellCoin(ctx sdk.Context, k Keeper, msg types.MsgSellCoin) (*sdk.Result, error) {
	coinToBuy, _ := k.GetCoin(ctx, msg.CoinToBuy)
	if coinToBuy.Symbol != msg.CoinToBuy {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	// Check if coin to sell exists
	coinToSell, _ := k.GetCoin(ctx, msg.CoinToSell)
	if coinToSell.Symbol != msg.CoinToSell {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}

	amountBuy, amountSell, calcErr := cliUtils.SellCoinCalculateAmounts(coinToBuy, coinToSell, msg.AmountToBuy, msg.AmountToSell)

	if calcErr != nil {
		return nil, sdkerrors.New(calcErr.Codespace(), calcErr.ABCICode(), calcErr.Error())
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Seller)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(msg.CoinToSell)).LT(amountSell) {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, errMsg)
	}

	err := k.UpdateBalance(ctx, msg.CoinToBuy, amountBuy, msg.Seller)
	if err != nil {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountSell.Neg(), msg.Seller)
	if err != nil {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	amountBuyInBaseCoin := formulas.CalculateSaleReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountBuy)
	amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
	k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountBuyInBaseCoin), coinToSell.Volume.Sub(amountSell))
	k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountSellInBaseCoin), coinToBuy.Volume.Add(amountBuy))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSellCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountBuy.String()),
			sdk.NewAttribute(types.AttributeAmountToSell, amountSell.String()),
			sdk.NewAttribute(types.AttributeAmountToBuyInBaseCoin, amountBuyInBaseCoin.String()),
			sdk.NewAttribute(types.AttributeAmountToSellInBaseCoin, amountSellInBaseCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Seller.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSendCoin(ctx sdk.Context, k Keeper, msg types.MsgSendCoin) (*sdk.Result, error) {
	// TODO: commission
	log.Println("Sequence value: ", ctx.Value("sequence"))
	err := k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Receiver, sdk.Coins{sdk.NewCoin(strings.ToLower(msg.Coin), msg.Amount)})
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, 6, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSendCoin),
			sdk.NewAttribute(types.AttributeCoin, msg.Coin),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeReceiver, msg.Receiver.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// TODO: add it (don't know how to add it to cli)
func handleMsgMultiSendCoin(ctx sdk.Context, k Keeper, msg types.MsgMultiSendCoin) (*sdk.Result, error) {
	for i := range msg.Coins {
		// TODO: Commission
		_ = k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Coins[i].Receiver, sdk.Coins{sdk.NewCoin(strings.ToLower(msg.Coins[i].Coin), msg.Coins[i].Amount)})
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeMultiSendCoin),
				sdk.NewAttribute(types.AttributeCoin, msg.Coins[i].Coin),
				sdk.NewAttribute(types.AttributeAmount, msg.Coins[i].Amount.String()),
				sdk.NewAttribute(types.AttributeReceiver, msg.Coins[i].Receiver.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
		})
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSellAllCoin(ctx sdk.Context, k Keeper, msg types.MsgSellAllCoin) (*sdk.Result, error) {
	coinToBuy, _ := k.GetCoin(ctx, msg.CoinToBuy)
	if coinToBuy.Symbol != msg.CoinToBuy {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	// Check if coin to sell exists
	coinToSell, _ := k.GetCoin(ctx, msg.CoinToSell)
	if coinToSell.Symbol != msg.CoinToSell {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}
	acc := k.AccountKeeper.GetAccount(ctx, msg.Seller)
	balance := acc.GetCoins()

	amountBuy, amountSell, calcErr := cliUtils.SellCoinCalculateAmounts(coinToBuy, coinToSell, msg.AmountToBuy, balance.AmountOf(strings.ToLower(msg.CoinToSell)))

	if calcErr != nil {
		return nil, sdkerrors.New(calcErr.Codespace(), calcErr.ABCICode(), calcErr.Error())
	}

	err := k.UpdateBalance(ctx, msg.CoinToBuy, amountBuy, msg.Seller)
	if err != nil {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountSell.Neg(), msg.Seller)
	if err != nil {
		// TODO: Add proper error message
		errMsg := ""
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	amountBuyInBaseCoin := formulas.CalculateSaleReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountBuy)
	amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
	k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountBuyInBaseCoin), coinToSell.Volume.Sub(amountSell))
	k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountSellInBaseCoin), coinToBuy.Volume.Add(amountBuy))

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSellCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountBuy.String()),
			sdk.NewAttribute(types.AttributeAmountToSell, amountSell.String()),
			sdk.NewAttribute(types.AttributeAmountToBuyInBaseCoin, amountBuyInBaseCoin.String()),
			sdk.NewAttribute(types.AttributeAmountToSellInBaseCoin, amountSellInBaseCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Seller.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
