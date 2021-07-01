package utils

import (
	"bitbucket.org/decimalteam/go-node/x/swap"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"strings"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/validator"
)

// Ante
func NewAnteHandler(ak keeper.AccountKeeper, bankKeeper bankKeeper.Keeper, vk validator.Keeper, ck coin.Keeper, consumer ante.SignatureVerificationGasConsumer) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewCountMsgDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		NewPreCreateAccountDecorator(ak), // should be before SetPubKeyDecorator
		ante.NewSetPubKeyDecorator(ak),   // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		NewFeeDecorator(ck, ak, bankKeeper, vk),
		ante.NewSigGasConsumeDecorator(ak, consumer),
		ante.NewSigVerificationDecorator(ak),
		NewPostCreateAccountDecorator(ak),      // should be after SigVerificationDecorator
		ante.NewIncrementSequenceDecorator(ak), // innermost AnteDecorator
	)
}

var (
	_ GasTx = (*legacytx.StdTx)(nil) // assert StdTx implements GasTx
)

// GasTx defines a Tx with a GetGas() method which is needed to use SetUpContextDecorator
type GasTx interface {
	sdk.Tx
	GetGas() uint64
}

// PreCreateAccountDecorator creates account in case of check redeeming from account unknown in the blockchain.
// Such accounts sign their first transaction with account number equal to 0. This is the reason why
// creating account is divided in two parts (PreCreateAccountDecorator and PostCreateAccountDecorator):
// it is necessary to create account in the beginning of the Ante chain with account number 0, but just
// before the end of the Ante chain it is necessary to restore correct account number.
type PreCreateAccountDecorator struct {
	ak keeper.AccountKeeper
}

// NewPreCreateAccountDecorator creates new PreCreateAccountDecorator.
func NewPreCreateAccountDecorator(ak keeper.AccountKeeper) PreCreateAccountDecorator {
	return PreCreateAccountDecorator{ak: ak}
}

// AnteHandle implements sdk.AnteHandler function.
func (cad PreCreateAccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	if len(msgs) > 0 {
		switch msgs[0].String() {
		case coin.RedeemCheckConst:
			signers := msgs[0].GetSigners()
			if len(signers) == 1 {
				acc := cad.ak.GetAccount(ctx, signers[0])
				if acc == nil {
					acc = cad.ak.NewAccountWithAddress(ctx, signers[0])
					ctx = ctx.WithValue("created_account_address", signers[0].String())
					ctx = ctx.WithValue("created_account_number", acc.GetAccountNumber())
					acc.SetAccountNumber(0) // necessary to validate signature
					cad.ak.SetAccount(ctx, acc)
				}
			}
		}
	}

	return next(ctx, tx, simulate)
}

// PostCreateAccountDecorator restores account number in case of check redeeming from account unknown for the blockchain.
type PostCreateAccountDecorator struct {
	ak keeper.AccountKeeper
}

// NewPostCreateAccountDecorator creates new PostCreateAccountDecorator.
func NewPostCreateAccountDecorator(ak keeper.AccountKeeper) PostCreateAccountDecorator {
	return PostCreateAccountDecorator{ak: ak}
}

// AnteHandle implements sdk.AnteHandler function.
func (cad PostCreateAccountDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	accAddress := ctx.Value("created_account_address")
	accNumber := ctx.Value("created_account_number")
	if accAddress != nil && accNumber != nil {
		ctx = ctx.WithValue("created_account_address", nil)
		ctx = ctx.WithValue("created_account_number", nil)
		accAddr, err := sdk.AccAddressFromBech32(accAddress.(string))
		if err != nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "invalid address of created accout")
		}
		acc := cad.ak.GetAccount(ctx, accAddr)
		if acc == nil {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownAddress, "unable to find created accout")
		}
		acc.SetAccountNumber(accNumber.(uint64))
		cad.ak.SetAccount(ctx, acc)
	}

	return next(ctx, tx, simulate)
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
				logStr := fmt.Sprintf(
					"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
					rType.Descriptor, gasTx.GetGas(), newCtx.GasMeter().GasConsumed())

				err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, logStr)
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
	bk bankKeeper.Keeper
	ak keeper.AccountKeeper
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
func NewFeeDecorator(ck coin.Keeper, ak keeper.AccountKeeper, bankKeeper bankKeeper.Keeper, vk validator.Keeper) FeeDecorator {
	return FeeDecorator{
		ck: ck,
		bk: bankKeeper,
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

	htltFee = 33000
)

// AnteHandle implements sdk.AnteHandler function.
func (fd FeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if ctx.BlockHeight() == 0 {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	if addr := fd.ak.GetModuleAddress(types.FeeCollectorName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
	}

	// all transactions must implement GasTx
	stdTx, ok := tx.(legacytx.StdTx)
	if !ok {
		return newCtx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be StdTx")
	}

	commissionInBaseCoin := sdk.ZeroInt()
	commissionInBaseCoin = commissionInBaseCoin.AddRaw(int64(len(ctx.TxBytes()) * 2))

	ctx = ctx.WithValue("fee", feeTx.GetFee())

	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		switch msg.String() {
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
			multiSend := msg.(*coin.MsgMultiSendCoin)
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sendFee + int64((len(multiSend.Sends)-1)*5))
		case coin.BuyCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(buyFee)
		case coin.SellCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sellFee)
		case coin.SellAllConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(sellFee)
		case coin.RedeemCheckConst:
			commissionInBaseCoin = sdk.ZeroInt()
		case multisig.CreateTransactionConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createTransactionFee)
		case multisig.CreateWalletConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createWalletFee)
		case multisig.SignTransactionConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(signTransactionFee)
		case coin.CreateCoinConst:
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(createCoinFee)
		case swap.MsgHTLTConst:
			if msg.(*swap.MsgHTLT).From.Equals(swap.SwapServiceAddress) {
				return next(ctx, tx, simulate)
			}
			commissionInBaseCoin = commissionInBaseCoin.AddRaw(htltFee)
		}
	}

	commissionInBaseCoin = helpers.UnitToPip(commissionInBaseCoin)

	feePayer := feeTx.FeePayer()
	feePayerAcc := fd.ak.GetAccount(ctx, feePayer)

	if feePayerAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if feeTx.GetFee().IsZero() {
		// deduct the fees
		if commissionInBaseCoin.IsZero() {
			return next(ctx, tx, simulate)
		}
		err = DeductFees(fd.bk, fd.ak, ctx, feePayerAcc, fd.ck, sdk.NewCoin(fd.vk.BondDenom(ctx), commissionInBaseCoin), commissionInBaseCoin)
		if err != nil {
			return ctx, err
		}
		if len(msgs) == 1 {
			if msgs[0].String() == validator.DelegateConst {
				stdTx.Fee.Gas = helpers.PipToUnit(commissionInBaseCoin).Uint64() * 10
				ctx = SetGasMeter(simulate, ctx, stdTx.GetGas())
				ctx.GasMeter().ConsumeGas(helpers.PipToUnit(commissionInBaseCoin).Uint64()*10, "commission")
			} else {
				stdTx.Fee.Gas = helpers.PipToUnit(commissionInBaseCoin).Uint64()
				ctx = SetGasMeter(simulate, ctx, stdTx.GetGas())
				ctx.GasMeter().ConsumeGas(helpers.PipToUnit(commissionInBaseCoin).Uint64(), "commission")
			}
		}
		return next(ctx, tx, simulate)
	}

	feeInBaseCoin := sdk.ZeroInt()

	f := feeTx.GetFee()[0]

	if f.Denom != fd.vk.BondDenom(ctx) {
		coinInfo, err := fd.vk.GetCoin(ctx, f.Denom)
		if err != nil {
			return ctx, err
		}

		feeInBaseCoin = formulas.CalculateSaleReturn(coinInfo.Volume, coinInfo.Reserve, coinInfo.CRR, f.Amount)

		if coinInfo.Reserve.Sub(feeInBaseCoin).LT(coin.MinCoinReserve(ctx)) {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("coin reserve balance is not sufficient for transaction. Has: %s, fee: %s",
				coinInfo.Reserve.String(),
				feeInBaseCoin.String()))
		}

	} else {
		feeInBaseCoin = f.Amount
	}

	if feeInBaseCoin.LT(commissionInBaseCoin) {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", feeInBaseCoin, commissionInBaseCoin)
	}

	// deduct the fees
	err = DeductFees(fd.bk, fd.ak, ctx, feePayerAcc, fd.ck, f, feeInBaseCoin)
	if err != nil {
		return ctx, err
	}
	if msgs[0].String() == validator.DelegateConst {
		stdTx.Fee.Gas = helpers.PipToUnit(feeInBaseCoin).Uint64() * 10
		ctx = SetGasMeter(simulate, ctx, stdTx.GetGas())
		ctx.GasMeter().ConsumeGas(helpers.PipToUnit(feeInBaseCoin).Uint64()*10, "commission")
	} else {
		stdTx.Fee.Gas = helpers.PipToUnit(feeInBaseCoin).Uint64()
		ctx = SetGasMeter(simulate, ctx, stdTx.GetGas())
		ctx.GasMeter().ConsumeGas(helpers.PipToUnit(feeInBaseCoin).Uint64(), "commission")
	}

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the BankKeeper (in addition to the AccountKeeper, because
// the BankKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(bankKeeper bankKeeper.Keeper, accKeeper keeper.AccountKeeper, ctx sdk.Context, acc client.Account, coinKeeper coin.Keeper, fee sdk.Coin, feeInBaseCoin sdk.Int) error {
	//blockTime := ctx.BlockHeader().Time
	var coins sdk.Coins

	bankKeeper.IterateTotalSupply(ctx, func (c sdk.Coin) bool {
		coins = coins.Add(c)

		return true
	})

	feeCoin, err := coinKeeper.GetCoin(ctx, strings.ToLower(fee.Denom))
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "coin not exist: %s", fee.Denom)
	}

	if !coinKeeper.IsCoinBase(fee.Denom) {
		if feeCoin.Reserve.Sub(feeInBaseCoin).LT(coin.MinCoinReserve(ctx)) {
			return coin.ErrTxBreaksMinReserveRule(ctx, feeCoin.Reserve.Sub(fee.Amount).String())
		}
	}

	if !fee.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fee)
	}

	// verify the account has enough funds to pay for fee
	_, hasNeg := coins.SafeSub(sdk.NewCoins(fee))
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fee; %s < %s", coins, fee)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	var spendableCoins sdk.Coins
	bankKeeper.IterateAccountBalances(ctx, acc.GetAddress(), func (c sdk.Coin) bool {
		spendableCoins = spendableCoins.Add(c)

		return false
	})

	if _, hasNeg := spendableCoins.SafeSub(sdk.NewCoins(fee)); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fee; %s < %s", spendableCoins, fee)
	}

	err = bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, sdk.NewCoins(fee))
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	s := bankKeeper.GetSupply(ctx, "del").Add(fee)
	newcoins := sdk.NewCoins(s)
	bankKeeper.MintCoins(ctx, minttypes.ModuleName, newcoins)

	bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, acc.GetAddress(), newcoins)

	// update coin: decrease reserve and volume
	if !coinKeeper.IsCoinBase(fee.Denom) {
		coinKeeper.UpdateCoin(ctx, feeCoin, feeCoin.Reserve.Sub(feeInBaseCoin), feeCoin.Volume.Sub(fee.Amount))
	} else {
		coinKeeper.UpdateCoin(ctx, feeCoin, feeCoin.Reserve, feeCoin.Volume.Sub(fee.Amount))
	}

	return nil
}

type CountMsgDecorator struct {
}

func NewCountMsgDecorator() CountMsgDecorator {
	return CountMsgDecorator{}
}

func (cd CountMsgDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if len(tx.GetMsgs()) > 1 {
		return ctx, sdkerrors.New("ante", 100, "Too many msgs in the transaction. Max = 1")
	}
	return next(ctx, tx, simulate)
}
