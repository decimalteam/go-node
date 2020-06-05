package utils

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/validator"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"log"
)

// Ante
func NewAnteHandler(ak keeper.AccountKeeper, vk validator.Keeper, ck coin.Keeper, sk supply.Keeper, consumer ante.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		NewFeeDecorator(ck, sk, ak, vk),
		ante.NewSigGasConsumeDecorator(ak, consumer),
		ante.NewSigVerificationDecorator(ak),
		ante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}

var (
	_ GasTx = (*types.StdTx)(nil) // assert StdTx implements GasTx
)

// GasTx defines a Tx with a GetGas() method which is needed to use SetUpContextDecorator
type GasTx interface {
	sdk.Tx
	GetGas() uint64
}

// SetUpContextDecorator sets the GasMeter in the Context and wraps the next AnteHandler with a defer clause
// to recover from any downstream OutOfGas panics in the AnteHandler chain to return an error with information
// on gas provided and gas used.
// CONTRACT: Must be first decorator in the chain
// CONTRACT: Tx must implement GasTx interface
type SetUpContextDecorator struct{}

// NewSetUpContextDecorator creates new SetUpContextDecorator.
func NewSetUpContextDecorator() SetUpContextDecorator {
	return SetUpContextDecorator{}
}

// AnteHandle implements sdk.AnteHandler function.
func (sud SetUpContextDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// all transactions must implement GasTx
	gasTx, ok := tx.(GasTx)
	if !ok {
		// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
		// during runTx.
		newCtx = SetGasMeter(simulate, ctx, 0)
		return newCtx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be GasTx")
	}

	newCtx = SetGasMeter(simulate, ctx, gasTx.GetGas())

	// Decorator will catch an OutOfGasPanic caused in the next antehandler
	// AnteHandlers must have their own defer/recover in order for the BaseApp
	// to know how much gas was used! This is because the GasMeter is created in
	// the AnteHandler, but if it panics the context won't be set properly in
	// runTx's recover call.
	defer func() {
		if r := recover(); r != nil {
			switch rType := r.(type) {
			case sdk.ErrorOutOfGas:
				log := fmt.Sprintf(
					"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
					rType.Descriptor, gasTx.GetGas(), newCtx.GasMeter().GasConsumed())

				err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, log)
			default:
				panic(r)
			}
		}
	}()

	return next(newCtx, tx, simulate)
}

// SetGasMeter returns a new context with a gas meter set from a given context.
func SetGasMeter(simulate bool, ctx sdk.Context, gasLimit uint64) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(NewGasMeter(gasLimit))
}

type FeeDecorator struct {
	ck coin.Keeper
	sk supply.Keeper
	ak auth.AccountKeeper
	vk validator.Keeper
}

// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	FeePayer() sdk.AccAddress
}

// NewSetUpContextDecorator creates new SetUpContextDecorator.
func NewFeeDecorator(ck coin.Keeper, sk supply.Keeper, ak auth.AccountKeeper, vk validator.Keeper) FeeDecorator {
	return FeeDecorator{
		ck: ck,
		sk: sk,
		ak: ak,
		vk: vk,
	}
}

const (
	declareCandidateFee = 10000
	editCandidateFee    = 10000
	delegateFee         = 200
	unbondFee           = 200
	setOnlineFee        = 100
	setOfflineFee       = 100

	sendFee        = 10
	sellFee        = 100
	buyFee         = 100
	redeemCheckFee = 30
	createCoinFee  = 100

	createWalletFee      = 100
	createTransactionFee = 100
	signTransactionFee   = 100
)

// AnteHandle implements sdk.AnteHandler function.
func (fd FeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := fd.sk.GetModuleAddress(types.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
	}

	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(len(ctx.TxBytes()) * 2))

	log.Println(len(ctx.TxBytes()) * 2)

	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch msg.Type() {
		case validator.DeclareCandidateConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(declareCandidateFee)
		case validator.DelegateConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(delegateFee)
		case validator.SetOnlineConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(setOnlineFee)
		case validator.SetOfflineConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(setOfflineFee)
		case validator.UnbondConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(unbondFee)
		case validator.EditCandidateConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(editCandidateFee)
		case coin.SendCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sendFee)
		case coin.MultiSendCoinConst:
			multiSend := msg.(coin.MsgMultiSendCoin)
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sendFee + int64((len(multiSend.Sends)-1)*5))
		case coin.BuyCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(buyFee)
		case coin.SellCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sellFee)
		case coin.SellAllConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sellFee)
		case coin.RedeemCheckConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(redeemCheckFee)
		case multisig.CreateTransactionConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createTransactionFee)
		case multisig.CreateWalletConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createWalletFee)
		case multisig.SignTransactionConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(signTransactionFee)
		case coin.CreateCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createCoinFee)
			return next(ctx, tx, simulate)
		}
	}

	commissionInBaseCoin = helpers.UnitToPip(commissionInBaseCoin)

	ctx = ctx.WithValue("fee", feeTx.GetFee())

	feePayer := feeTx.FeePayer()
	feePayerAcc := fd.ak.GetAccount(ctx, feePayer)

	if feePayerAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if feeTx.GetFee().IsZero() {
		// deduct the fees
		err = DeductFees(fd.sk, ctx, feePayerAcc, sdk.NewCoins(sdk.NewCoin(fd.vk.BondDenom(ctx), commissionInBaseCoin)))
		if err != nil {
			return ctx, err
		}
		ctx.GasMeter().ConsumeGas(helpers.PipToUnit(commissionInBaseCoin).Uint64(), "commission")
		return next(ctx, tx, simulate)
	}

	feeInBaseCoin := sdk.ZeroInt()

	f := feeTx.GetFee()[0]

	if f.Denom != fd.vk.BondDenom(ctx) {
		coinInfo, err := fd.vk.GetCoin(ctx, f.Denom)
		if err != nil {
			return ctx, err
		}

		if coinInfo.Reserve.LT(commissionInBaseCoin) {
			return ctx, fmt.Errorf("coin reserve balance is not sufficient for transaction. Has: %s, required %s",
				coinInfo.Reserve.String(),
				commissionInBaseCoin.String())
		}

		feeInBaseCoin = formulas.CalculateSaleAmount(coinInfo.Volume, coinInfo.Reserve, coinInfo.CRR, f.Amount)
	} else {
		feeInBaseCoin = f.Amount
	}

	if feeInBaseCoin.LT(commissionInBaseCoin) {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", feeInBaseCoin, commissionInBaseCoin)
	}

	// deduct the fees
	err = DeductFees(fd.sk, ctx, feePayerAcc, feeTx.GetFee())
	if err != nil {
		return ctx, err
	}
	ctx.GasMeter().ConsumeGas(helpers.PipToUnit(feeInBaseCoin).Uint64(), "commission")

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the BankKeeper (in addition to the AccountKeeper, because
// the BankKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(supplyKeeper supply.Keeper, ctx sdk.Context, acc exported.Account, fees sdk.Coins) error {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	err := supplyKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	s := supplyKeeper.GetSupply(ctx)
	s = s.Inflate(fees)
	supplyKeeper.SetSupply(ctx, s)

	return nil
}
