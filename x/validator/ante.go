package validator

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin"
	vtypes "bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"log"
	"strconv"
)

// Ante
func NewAnteHandler(ak keeper.AccountKeeper, vk Keeper, ck coin.Keeper, sk supply.Keeper, consumer ante.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		return NewFeeCoinUpdateDecorator(vk, ck).AnteHandle(ctx, tx, simulate, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			return NewSequenceEventDecorator(ak).AnteHandle(ctx, tx, simulate, auth.NewAnteHandler(ak, sk, consumer))
		})
	}
}

type SequenceEventDecorator struct {
	ak keeper.AccountKeeper
}

func NewSequenceEventDecorator(ak keeper.AccountKeeper) SequenceEventDecorator {
	return SequenceEventDecorator{
		ak: ak,
	}
}

func (sed SequenceEventDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// no need to increment sequence on CheckTx or RecheckTx
	if ctx.IsCheckTx() && !simulate {
		return next(ctx, tx, simulate)
	}

	sigTx, ok := tx.(auth.StdTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
	}

	// increment sequence of all signers
	for _, addr := range sigTx.GetSigners() {
		acc := sed.ak.GetAccount(ctx, addr)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(vtypes.AttributeKeySequence, strconv.FormatUint(acc.GetSequence(), 10)),
		))
	}

	log.Println("next")
	return next(ctx, tx, simulate)
}

type FeeCoinUpdateDecorator struct {
	vk Keeper
	ck coin.Keeper
}

func NewFeeCoinUpdateDecorator(vk Keeper, ck coin.Keeper) FeeCoinUpdateDecorator {
	return FeeCoinUpdateDecorator{
		vk: vk,
		ck: ck,
	}
}

func (d FeeCoinUpdateDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// all transactions must be of type auth.StdTx
	stdTx, ok := tx.(auth.StdTx)
	if !ok {
		// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
		// during runTx.
		return ctx, errors.New("tx must be StdTx")
	}

	if !stdTx.Fee.Amount.IsZero() {
		for _, fee := range stdTx.GetFee() {
			commission := sdk.NewIntFromBigInt(fee.Amount.BigInt())

			feeCoin, errFee := d.vk.GetCoin(ctx, fee.Denom)
			if errFee != nil {
				return ctx, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "fees coins is not exist: %s", errFee)
			}
			if fee.Denom != d.vk.BondDenom(ctx) {
				if feeCoin.Reserve.LT(fee.Amount) {
					return ctx, vtypes.ErrCoinReserveIsNotSufficient(DefaultCodespace, feeCoin.Reserve.String(), fee.Amount.String())
				}

				commission = formulas.CalculateSaleAmount(feeCoin.Volume, feeCoin.Reserve, feeCoin.CRR, fee.Amount)
			}

			d.ck.UpdateCoin(ctx, feeCoin, feeCoin.Reserve.Sub(fee.Amount), feeCoin.Volume.Sub(commission))
		}
	}

	log.Println("next")
	return next(ctx, tx, simulate)
}
