package validator

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin"
	vtypes "bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

var (
	// simulation signature values used to estimate gas consumption
	simSecp256k1Pubkey secp256k1.PubKeySecp256k1
	simSecp256k1Sig    [64]byte
)

func init() {
	// This decodes a valid hex string into a sepc256k1Pubkey for use in transaction simulation
	bz, _ := hex.DecodeString("035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A")
	copy(simSecp256k1Pubkey[:], bz)
}

// SignatureVerificationGasConsumer is the type of function that is used to both consume gas when verifying signatures
// and also to accept or reject different types of PubKey's. This is where apps can define their own PubKey
type SignatureVerificationGasConsumer = func(meter sdk.GasMeter, sig []byte, pubkey crypto.PubKey, params auth.Params) error

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(ak auth.AccountKeeper, supplyKeeper types.SupplyKeeper, coinKeeper coin.Keeper, sigGasConsumer SignatureVerificationGasConsumer) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (newCtx sdk.Context, err error) {

		if addr := supplyKeeper.GetModuleAddress(types.FeeCollectorName); addr == nil {
			panic(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
		}

		// all transactions must be of type auth.StdTx
		stdTx, ok := tx.(auth.StdTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = SetGasMeter(simulate, ctx, 0)
			err = errors.New("tx must be StdTx")
			return
		}

		params := ak.GetParams(ctx)

		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && !simulate {
			_, err = EnsureSufficientMempoolFees(ctx, stdTx.Fee)
			if err != nil {
				return
			}
		}

		newCtx = SetGasMeter(simulate, ctx, stdTx.Fee.Gas)

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
						rType.Descriptor, stdTx.Fee.Gas, newCtx.GasMeter().GasConsumed(),
					)
					err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, log)

					// TODO: How to report these?
					// res := sdk.Result{}
					// res.GasWanted = stdTx.Fee.Gas
					// res.GasUsed = newCtx.GasMeter().GasConsumed()
				default:
					panic(r)
				}
			}
		}()

		if _, err = ValidateSigCount(stdTx, params); err != nil {
			return
		}

		if err = tx.ValidateBasic(); err != nil {
			return
		}

		newCtx.GasMeter().ConsumeGas(params.TxSizeCostPerByte*sdk.Gas(len(newCtx.TxBytes())), "txSize")

		if _, err = ValidateMemo(stdTx, params); err != nil {
			return
		}

		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		signerAddrs := stdTx.GetSigners()
		signerAccs := make([]exported.Account, len(signerAddrs))
		isGenesis := ctx.BlockHeight() == 0

		// fetch first signer, who's going to pay the fees
		signerAccs[0], _, err = GetSignerAcc(newCtx, ak, signerAddrs[0])
		if err != nil {
			return
		}

		// deduct the fees
		if !stdTx.Fee.Amount.IsZero() {
			_, err = DeductFees(supplyKeeper, coinKeeper, newCtx, signerAccs[0], stdTx.Fee.Amount)
			if err != nil {
				return
			}

			// reload the account as fees have been deducted
			signerAccs[0] = ak.GetAccount(newCtx, signerAccs[0].GetAddress())
		}

		// stdSigs contains the sequence number, account number, and signatures.
		// When simulating, this would just be a 0-length slice.
		stdSigs := stdTx.GetSignatures()

		for i := 0; i < len(stdSigs); i++ {
			// skip the fee payer, account is cached and fees were deducted already
			if i != 0 {
				signerAccs[i], _, err = GetSignerAcc(newCtx, ak, signerAddrs[i])
				if err != nil {
					return
				}
			}

			// check signature, return account with incremented nonce
			signBytes := GetSignBytes(newCtx.ChainID(), stdTx, signerAccs[i], isGenesis)
			signerAccs[i], _, err = processSig(newCtx, signerAccs[i], stdTx.Signatures[i], signBytes, simulate, params, sigGasConsumer)
			if err != nil {
				return
			}

			ak.SetAccount(newCtx, signerAccs[i])
		}

		// TODO: tx tags (?)
		// TODO: GasWanted?
		// return newCtx, sdk.Result{GasWanted: stdTx.Fee.Gas}, false // continue...
		return
	}
}

// GetSignerAcc returns an account for a given address that is expected to sign
// a transaction.
func GetSignerAcc(ctx sdk.Context, ak auth.AccountKeeper, addr sdk.AccAddress) (exported.Account, *sdk.Result, error) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, &sdk.Result{}, nil
	}
	return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
}

// ValidateSigCount validates that the transaction has a valid cumulative total
// amount of signatures.
func ValidateSigCount(stdTx auth.StdTx, params auth.Params) (*sdk.Result, error) {
	pubKeys := stdTx.GetPubKeys()

	sigCount := 0
	for _, pk := range pubKeys {
		sigCount += auth.CountSubKeys(pk)
		if uint64(sigCount) > params.TxSigLimit {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrTooManySignatures, "signatures: %d, limit: %d", sigCount, params.TxSigLimit)
		}
	}

	return &sdk.Result{}, nil
}

// ValidateMemo validates the memo size.
func ValidateMemo(stdTx auth.StdTx, params auth.Params) (*sdk.Result, error) {
	memoLength := len(stdTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrMemoTooLarge,
			"maximum number of characters is %d but received %d characters",
			params.MaxMemoCharacters, memoLength,
		)
	}

	return &sdk.Result{}, nil
}

// verify the signature and increment the sequence. If the account doesn't have
// a pubkey, set it.
func processSig(
	ctx sdk.Context, acc exported.Account, sig auth.StdSignature, signBytes []byte, simulate bool, params auth.Params,
	sigGasConsumer SignatureVerificationGasConsumer,
) (updatedAcc exported.Account, res *sdk.Result, err error) {

	pubKey, res, err := ProcessPubKey(acc, sig, simulate)
	if err != nil {
		return nil, nil, err
	}

	err = acc.SetPubKey(pubKey)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "setting PubKey on signer's account")
	}

	if simulate {
		// Simulated txs should not contain a signature and are not required to
		// contain a pubkey, so we must account for tx size of including a
		// StdSignature (Amino encoding) and simulate gas consumption
		// (assuming a SECP256k1 simulation key).
		consumeSimSigGas(ctx.GasMeter(), pubKey, sig, params)
	}

	if err = sigGasConsumer(ctx.GasMeter(), sig.Signature, pubKey, params); err != nil {
		return nil, nil, err
	}

	if !simulate && !pubKey.VerifyBytes(signBytes, sig.Signature) {
		return nil, nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "signature verification failed; verify correct account sequence and chain-id")
	}

	if err = acc.SetSequence(acc.GetSequence() + 1); err != nil {
		panic(err)
	}

	return acc, res, nil
}

func consumeSimSigGas(gasmeter sdk.GasMeter, pubkey crypto.PubKey, sig auth.StdSignature, params auth.Params) {
	simSig := auth.StdSignature{PubKey: pubkey}
	if len(sig.Signature) == 0 {
		simSig.Signature = simSecp256k1Sig[:]
	}

	sigBz := ModuleCdc.MustMarshalBinaryLengthPrefixed(simSig)
	cost := sdk.Gas(len(sigBz) + 6)

	// If the pubkey is a multi-signature pubkey, then we estimate for the maximum
	// number of signers.
	if _, ok := pubkey.(multisig.PubKeyMultisigThreshold); ok {
		cost *= params.TxSigLimit
	}

	gasmeter.ConsumeGas(params.TxSizeCostPerByte*cost, "txSize")
}

// ProcessPubKey verifies that the given account address matches that of the
// StdSignature. In addition, it will set the public key of the account if it
// has not been set.
func ProcessPubKey(acc exported.Account, sig auth.StdSignature, simulate bool) (crypto.PubKey, *sdk.Result, error) {
	// If pubkey is not known for account, set it from the StdSignature.
	pubKey := acc.GetPubKey()
	if simulate {
		// In simulate mode the transaction comes with no signatures, thus if the
		// account's pubkey is nil, both signature verification and gasKVStore.Set()
		// shall consume the largest amount, i.e. it takes more gas to verify
		// secp256k1 keys than ed25519 ones.
		if pubKey == nil {
			return simSecp256k1Pubkey, &sdk.Result{}, nil
		}

		return pubKey, &sdk.Result{}, nil
	}

	if pubKey == nil {
		pubKey = sig.PubKey
		if pubKey == nil {
			return nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "PubKey not found")
		}

		if !bytes.Equal(pubKey.Address(), acc.GetAddress()) {
			return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "PubKey does not match Signer address %s", acc.GetAddress())
		}
	}

	return pubKey, &sdk.Result{}, nil
}

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter sdk.GasMeter, sig []byte, pubkey crypto.PubKey, params auth.Params,
) (*sdk.Result, error) {
	switch pubkey := pubkey.(type) {
	case ed25519.PubKeyEd25519:
		meter.ConsumeGas(params.SigVerifyCostED25519, "ante verify: ed25519")
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "ed25519 public keys are unsupported")

	case secp256k1.PubKeySecp256k1:
		meter.ConsumeGas(params.SigVerifyCostSecp256k1, "ante verify: secp256k1")
		return &sdk.Result{}, nil

	case multisig.PubKeyMultisigThreshold:
		var multisignature multisig.Multisignature
		codec.Cdc.MustUnmarshalBinaryBare(sig, &multisignature)

		consumeMultisignatureVerificationGas(meter, multisignature, pubkey, params)
		return &sdk.Result{}, nil

	default:
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "unrecognized public key type: %T", pubkey)
	}
}

func consumeMultisignatureVerificationGas(meter sdk.GasMeter,
	sig multisig.Multisignature, pubkey multisig.PubKeyMultisigThreshold,
	params auth.Params) {

	size := sig.BitArray.Size()
	sigIndex := 0
	for i := 0; i < size; i++ {
		if sig.BitArray.GetIndex(i) {
			DefaultSigVerificationGasConsumer(meter, sig.Sigs[sigIndex], pubkey.PubKeys[i], params)
			sigIndex++
		}
	}
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the CoinKeeper (in addition to the AccountKeeper, because
// the CoinKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(supplyKeeper types.SupplyKeeper, coinKeeper coin.Keeper, ctx sdk.Context, acc exported.Account, fees sdk.Coins) (*sdk.Result, error) {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	fee := sdk.Coin{}
	if len(coins) >= 1 {
		fee = fees[0]
	} else {
		return &sdk.Result{}, nil
	}

	commission := sdk.NewIntFromBigInt(fee.Amount.BigInt())

	// TODO вопрос с регистрозависимостью
	var denom string
	if fee.Denom == "tdcl" {
		denom = "tDCL"
	}
	feeCoin, errFee := coinKeeper.GetCoin(ctx, strings.ToUpper(denom))
	if errFee != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "fees coins is not exist: %s", errFee)
	}
	if fee.Denom != DefaultBondDenom {
		if feeCoin.Reserve.LT(fee.Amount) {
			return nil, vtypes.ErrCoinReserveIsNotSufficient(DefaultCodespace, feeCoin.Reserve.String(), fee.Amount.String())
		}

		commission = formulas.CalculateSaleAmount(feeCoin.Volume, feeCoin.Reserve, feeCoin.CRR, fee.Amount)
	}

	coinKeeper.UpdateCoin(ctx, feeCoin, feeCoin.Reserve.Sub(fee.Amount), feeCoin.Volume.Sub(commission))

	err := supplyKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{}, nil
}

// EnsureSufficientMempoolFees verifies that the given transaction has supplied
// enough fees to cover a proposer's minimum fees. A result object is returned
// indicating success or failure.
//
// Contract: This should only be called during CheckTx as it cannot be part of
// consensus.
func EnsureSufficientMempoolFees(ctx sdk.Context, stdFee auth.StdFee) (*sdk.Result, error) {
	minGasPrices := ctx.MinGasPrices()
	if !minGasPrices.IsZero() {
		requiredFees := make(sdk.Coins, len(minGasPrices))

		// Determine the required fees by multiplying each required minimum gas
		// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
		glDec := sdk.NewDec(int64(stdFee.Gas))
		for i, gp := range minGasPrices {
			fee := gp.Amount.Mul(glDec)
			requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
		}

		if !stdFee.Amount.IsAnyGTE(requiredFees) {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %q required: %q", stdFee.Amount, requiredFees)
		}
	}

	return &sdk.Result{}, nil
}

// SetGasMeter returns a new context with a gas meter set from a given context.
func SetGasMeter(simulate bool, ctx sdk.Context, gasLimit uint64) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
}

// GetSignBytes returns a slice of bytes to sign over for a given transaction
// and an account.
func GetSignBytes(chainID string, stdTx auth.StdTx, acc exported.Account, genesis bool) []byte {
	var accNum uint64
	if !genesis {
		accNum = acc.GetAccountNumber()
	}

	return auth.StdSignBytes(
		chainID, accNum, acc.GetSequence(), stdTx.Fee, stdTx.Msgs, stdTx.Memo,
	)
}
