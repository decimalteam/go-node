package keeper

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/decimalteam/go-node/x/validator/exported"

	ncfg "bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
// of it, updating unbonding delegations & redelegations appropriately
//
// CONTRACT:
//    slashFactor is non-negative
// CONTRACT:
//    Infraction was committed equal to or less than an unbonding period in the past,
//    so all unbonding delegations and redelegations from that height are stored
// CONTRACT:
//    Slash will not slash unbonded validators (for the above reason)
// CONTRACT:
//    Infraction was committed at the current height or at a past height,
//    not at a height in the future
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, slashFactor sdk.Dec) sdk.Int {
	logger := k.Logger(ctx)

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	// ref https://github.com/cosmos/cosmos-sdk/issues/1348

	validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(fmt.Sprintf(
			"WARNING: Ignored attempt to slash a nonexistent validator with address %s, we recommend you investigate immediately",
			consAddr))
		return sdk.ZeroInt()
	}

	// Coin of slashing = slash slashFactor * power at time of infraction
	amount := validator.Tokens
	slashAmountDec := amount.ToDec().Mul(slashFactor)
	slashAmount := slashAmountDec.TruncateInt()

	// should not be slashing an unbonded validator
	if validator.IsUnbonded() {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.ValAddress))
	}
	// call the before-modification hook
	k.BeforeValidatorModified(ctx, validator.ValAddress)

	k.DeleteValidatorByPowerIndex(ctx, validator)

	delegations := k.GetValidatorDelegations(ctx, validator.ValAddress)
	amountSlashed := k.slashBondedDelegations(ctx, delegations, slashFactor)

	switch {
	case infractionHeight > ctx.BlockHeight():

		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():

		// Special-case slash at current height for efficiency - we don't need to look through unbonding delegations or redelegations
		logger.Info(fmt.Sprintf(
			"slashing at current height %d, not scanning unbonding delegations & redelegations",
			infractionHeight))

	case infractionHeight < ctx.BlockHeight():
		// Iterate through unbonding delegations from slashed validator
		unbondingDelegations := k.GetUnbondingDelegationsFromValidator(ctx, validator.ValAddress)
		for _, unbondingDelegation := range unbondingDelegations {
			amountSlashed = amountSlashed.Add(k.slashUnbondingDelegation(ctx, unbondingDelegation, infractionHeight, slashFactor)...)
		}
	}

	// Log that a slash occurred!
	logger.Info(fmt.Sprintf(
		"validator %s slashed by slash factor of %s; burned %v tokens",
		validator.GetOperator(), slashFactor.String(), amountSlashed))

	validator, err = k.GetValidator(ctx, validator.ValAddress)
	if err != nil {
		panic(err)
	}
	k.SetValidatorByPowerIndexWithoutCalc(ctx, validator)

	return slashAmount
}

// jail a validator
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		panic(err)
	}
	k.jailValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("validator %s jailed", consAddr))
}

// unjail a validator
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		panic(err)
	}
	err = k.unjailValidator(ctx, validator)
	if err != nil {
		panic(err)
	}
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("validator %s unjailed", consAddr))
}

// slash an unbonding delegation and update the pool
// return the amount that would have been slashed assuming
// the unbonding delegation had enough stake to slash
// (the amount actually slashed may be less if there's
// insufficient stake remaining)
func (k Keeper) slashUnbondingDelegation(ctx sdk.Context, unbondingDelegation types.UnbondingDelegation,
	infractionHeight int64, slashFactor sdk.Dec) sdk.Coins {

	now := ctx.BlockHeader().Time
	totalSlashAmount := sdk.NewCoins()
	burnedAmount := sdk.NewCoins()

	// perform slashing on all entries within the unbonding delegation
	for i, entry := range unbondingDelegation.Entries {
		if _, ok := entry.(types.UnbondingDelegationNFTEntry); ok {
			continue
		}
		entry := entry.(types.UnbondingDelegationEntry)

		// If unbonding started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Unbonding delegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance.Amount)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(sdk.NewCoin(entry.Balance.Denom, slashAmount))

		// Don't slash more tokens than held
		// Possible since the unbonding delegation may already
		// have been slashed, and slash amounts are calculated
		// according to stake held at time of infraction
		unbondingSlashAmount := sdk.MinInt(slashAmount, entry.Balance.Amount)

		// Update unbonding delegation if necessary
		if unbondingSlashAmount.IsZero() {
			continue
		}

		burnedAmount = burnedAmount.Add(sdk.NewCoin(entry.Balance.Denom, unbondingSlashAmount))
		entry.Balance.Amount = entry.Balance.Amount.Sub(unbondingSlashAmount)
		unbondingDelegation.Entries[i] = entry
		k.SetUnbondingDelegation(ctx, unbondingDelegation)

		if entry.Balance.Denom != k.BondDenom(ctx) {
			coin, err := k.GetCoin(ctx, entry.Balance.Denom)
			if err != nil {
				panic(err)
			}
			ret := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, unbondingSlashAmount)
			k.CoinKeeper.UpdateCoin(ctx, coin, coin.Reserve.Sub(ret), coin.Volume.Sub(unbondingSlashAmount))
			k.SubtractDelegatedCoin(ctx, sdk.NewCoin(entry.Balance.Denom, unbondingSlashAmount))
		}
	}

	tokensToBurn := sdk.NewCoins()
	for _, coin := range burnedAmount {
		tokensToBurn = tokensToBurn.Add(sdk.NewCoin(coin.Denom, sdk.MaxInt(coin.Amount, sdk.ZeroInt()))) // defensive.
	}
	if err := k.burnNotBondedTokens(ctx, tokensToBurn); err != nil {
		panic(err)
	}

	return totalSlashAmount
}

// return total slashed coins
func (k Keeper) slashBondedDelegations(ctx sdk.Context, delegations []exported.DelegationI, slashFactor sdk.Dec) sdk.Coins {
	totalSlashAmount := sdk.ZeroInt()
	burnedAmount := sdk.NewCoins()

	for _, delegation := range delegations {
		switch delegation := delegation.(type) {
		case types.Delegation:
			// Calculate slash amount proportional to stake contributing to infraction
			slashAmountDec := slashFactor.MulInt(delegation.GetCoin().Amount)
			slashAmount := slashAmountDec.TruncateInt()
			totalSlashAmount = totalSlashAmount.Add(slashAmount)

			bondSlashAmount := sdk.MinInt(slashAmount, delegation.GetCoin().Amount)
			bondSlashAmount = sdk.MaxInt(bondSlashAmount, sdk.ZeroInt())

			if bondSlashAmount.IsZero() {
				continue
			}

			burnedAmount = burnedAmount.Add(sdk.NewCoin(delegation.GetCoin().Denom, bondSlashAmount))
			delegation.Coin.Amount = delegation.GetCoin().Amount.Sub(bondSlashAmount)
			k.SetDelegation(ctx, delegation)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeLiveness,
					sdk.NewAttribute(types.AttributeKeyValidator, delegation.GetValidatorAddr().String()),
					sdk.NewAttribute(types.AttributeKeyDelegator, delegation.GetDelegatorAddr().String()),
					sdk.NewAttribute(types.AttributeKeySlashAmount, sdk.NewCoin(delegation.GetCoin().Denom, bondSlashAmount).String()),
					sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
				),
			)

			validator, err := k.GetValidator(ctx, delegation.GetValidatorAddr())
			if err != nil {
				panic(err)
			}
			if delegation.GetCoin().Denom != k.BondDenom(ctx) {
				coin, err := k.GetCoin(ctx, delegation.GetCoin().Denom)
				if err != nil {
					panic(err)
				}
				ret := formulas.CalculateSaleReturn(coin.Volume, coin.Reserve, coin.CRR, bondSlashAmount)
				k.CoinKeeper.UpdateCoin(ctx, coin, coin.Reserve.Sub(ret), coin.Volume.Sub(bondSlashAmount))
				validator.Tokens = validator.Tokens.Sub(ret)

				k.SubtractDelegatedCoin(ctx, sdk.NewCoin(delegation.GetCoin().Denom, bondSlashAmount))
			} else {
				validator.Tokens = validator.Tokens.Sub(bondSlashAmount)
			}
			err = k.SetValidator(ctx, validator)
			if err != nil {
				panic(err)
			}
		case types.DelegationNFT:
			validator, err := k.GetValidator(ctx, delegation.GetValidatorAddr())
			if err != nil {
				panic(err)
			}

			for _, subTokenID := range delegation.SubTokenIDs {
				reserve, found := k.nftKeeper.GetSubToken(ctx, delegation.Denom, delegation.TokenID, subTokenID)
				if !found {
					panic(fmt.Errorf("subToken with ID = %d not found", subTokenID))
				}
				// Calculate slash amount proportional to stake contributing to infraction
				slashAmountDec := slashFactor.MulInt(reserve)
				slashAmount := slashAmountDec.TruncateInt()
				totalSlashAmount = totalSlashAmount.Add(slashAmount)

				bondSlashAmount := sdk.MinInt(slashAmount, delegation.GetCoin().Amount)
				bondSlashAmount = sdk.MaxInt(bondSlashAmount, sdk.ZeroInt())

				if bondSlashAmount.IsZero() {
					continue
				}

				burnedAmount = burnedAmount.Add(sdk.NewCoin(delegation.GetCoin().Denom, bondSlashAmount))
				delegation.Coin.Amount = delegation.GetCoin().Amount.Sub(bondSlashAmount)
				validator.Tokens = validator.Tokens.Sub(bondSlashAmount)

				if !reserve.IsZero() {
					reserve = reserve.Sub(bondSlashAmount)
					k.nftKeeper.SetSubToken(ctx, delegation.Denom, delegation.TokenID, subTokenID, reserve)
				}

				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeLivenessNFT,
						sdk.NewAttribute(types.AttributeKeyValidator, delegation.GetValidatorAddr().String()),
						sdk.NewAttribute(types.AttributeKeyDelegator, delegation.GetDelegatorAddr().String()),
						sdk.NewAttribute(types.AttributeKeySlashReserve, reserve.String()),
						sdk.NewAttribute(types.AttributeKeySlashSubTokenID, strconv.FormatInt(subTokenID, 10)),
						sdk.NewAttribute(types.AttributeKeySlashAmount, sdk.NewCoin(delegation.GetCoin().Denom, bondSlashAmount).String()),
						sdk.NewAttribute(types.AttributeKeyDenom, delegation.Denom),
						sdk.NewAttribute(types.AttributeKeyID, delegation.TokenID),
						sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
					),
				)
			}

			k.SetDelegationNFT(ctx, delegation)

			err = k.SetValidator(ctx, validator)
			if err != nil {
				panic(err)
			}
		}
	}

	tokensToBurn := sdk.NewCoins()
	for _, coin := range burnedAmount {
		tokensToBurn = tokensToBurn.Add(sdk.NewCoin(coin.Denom, sdk.MaxInt(coin.Amount, sdk.ZeroInt()))) // defensive.
	}
	if err := k.burnBondedTokens(ctx, tokensToBurn); err != nil {
		panic(err)
	}

	return tokensToBurn
}

// handle a validator signature, must be called once per validator per block
func (k Keeper) HandleValidatorSignature(ctx sdk.Context, addr crypto.Address, power int64, signed bool) {
	var (
		logger   = k.Logger(ctx)
		height   = ctx.BlockHeight()
		consAddr = sdk.ConsAddress(addr)
	)

	pubkey, err := k.getPubkey(ctx, addr)
	if err != nil {
		panic(fmt.Sprintf("Validator consensus-address %s not found", consAddr))
	}

	validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		panic(err)
	}

	if validator.Jailed {
		return
	}

	// fetch signing info
	signInfo, found := k.getValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}

	const blockToClearMessedBlocks = 477514
	if !validator.Online || height == blockToClearMessedBlocks {
		if signInfo.MissedBlocksCounter > 0 {
			k.clearValidatorSigningInfo(ctx, consAddr, &signInfo)
			k.setValidatorSigningInfo(ctx, consAddr, signInfo)
		}
		return
	}

	index := signInfo.IndexOffset
	signInfo.IndexOffset = (signInfo.IndexOffset + 1) % types.SignedBlocksWindow

	missedInWindow := k.getValidatorMissedBlockBitArray(ctx, consAddr, index)
	missed := !signed

	switch missed {
	case true:
		// If in grace period then pass missing block
		if inGracePeriod(ctx) {
			// log.Println(consAddr.String())
			ctx.Logger().Info(
				fmt.Sprintf("Missed block in grace period (%s)", validator.ValAddress))
			return
		}
		// If missed < 24 then missed = missed + 1
		if signInfo.MissedBlocksCounter < types.SignedBlocksWindow && !missedInWindow {
			k.setValidatorMissedBlockBitArray(ctx, consAddr, index, true)
			signInfo.MissedBlocksCounter++
		}
	case false:
		// If in grace perid and missed > 0 then missed = missed - 1
		// If missed in bit array and missed > 0 then missed = missed - 1
		grMissedBlocks := signInfo.MissedBlocksCounter > 0
		if (inGracePeriod(ctx) && grMissedBlocks) || (missedInWindow && grMissedBlocks) {
			k.setValidatorMissedBlockBitArray(ctx, consAddr, index, false)
			signInfo.MissedBlocksCounter--
		}
	}

	if missed {
		// log.Println(fmt.Sprintf("Missed blocks: %d", signInfo.MissedBlocksCounter), signInfo.Address.String())
		ctx.Logger().Info(
			fmt.Sprintf("Missed blocks %d in slash period (%s)", signInfo.MissedBlocksCounter, validator.ValAddress))

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeLiveness,
				sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
				sdk.NewAttribute(types.AttributeKeyMissedBlocks, fmt.Sprintf("%d", signInfo.MissedBlocksCounter)),
			),
		)

		logger.Info(
			fmt.Sprintf("Absent validator %s (%s) at height %d, %d missed, threshold %d", consAddr, pubkey, height, signInfo.MissedBlocksCounter, types.MinSignedPerWindow))
	}

	maxMissed := types.SignedBlocksWindow - types.MinSignedPerWindow

	// if we are past the minimum height and the validator has missed too many blocks, punish them
	if signInfo.MissedBlocksCounter > maxMissed {
		if !validator.IsJailed() {

			// Downtime confirmed: slash and jail the validator
			logger.Info(fmt.Sprintf("Validator %s past and below signed blocks threshold of %d",
				consAddr, types.MinSignedPerWindow))

			// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height,
			// and subtract an additional 1 since this is the LastCommit.
			// Note that this *can* result in a negative "distributionHeight" up to -ValidatorUpdateDelay-1,
			// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
			// That's fine since this is just used to filter unbonding delegations & redelegations.
			distributionHeight := height - sdk.ValidatorUpdateDelay - 1

			slashAmount := k.Slash(ctx, consAddr, distributionHeight, types.SlashFractionDowntime)
			k.Jail(ctx, consAddr)

			ctx.EventManager().EmitEvent(sdk.NewEvent(
				types.EventTypeSlash,
				sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
				sdk.NewAttribute(types.AttributeKeySlashAmount, slashAmount.String()),
				sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueMissingSignature),
			))

			// We need to reset the counter & array so that the validator won't be immediately slashed for downtime upon rebonding.
			k.clearValidatorSigningInfo(ctx, consAddr, &signInfo)
		} else {
			// Validator was (a) not found or (b) already jailed, don't slash
			logger.Info(
				fmt.Sprintf("Validator %s would have been slashed for downtime, but was either not found in store or already jailed", consAddr),
			)
		}
	}

	// Set the updated signing infoconsAddr
	k.setValidatorSigningInfo(ctx, consAddr, signInfo)
}

// Set index, counter = 0, bit array = nil
func (k Keeper) clearValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress, signInfo *types.ValidatorSigningInfo) {
	signInfo.MissedBlocksCounter = 0
	signInfo.IndexOffset = 0
	k.clearValidatorMissedBlockBitArray(ctx, address)
}

// Stored by *validator* address (not operator address)
func (k Keeper) clearValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorMissedBlockBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// Stored by *validator* address (not operator address)
func (k Keeper) getValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorMissedBlockBitArrayKey(address, index))
	if bz == nil {
		// lazy: treat empty key as not missed
		missed = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) setValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missed)
	store.Set(types.GetValidatorMissedBlockBitArrayKey(address, index), bz)
}

// Stored by *validator* address (not operator address)
func (k Keeper) getValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress) (info types.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorSigningInfoKey(address))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	found = true
	return
}

// Stored by *validator* address (not operator address)
func (k Keeper) setValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress, info types.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(info)
	store.Set(types.GetValidatorSigningInfoKey(address), bz)
}

// handle a validator signing two blocks at the same height
// power: power of the double-signing validator at the height of infraction
func (k Keeper) HandleDoubleSign(ctx sdk.Context, addr crypto.Address, infractionHeight int64, timestamp time.Time, power int64) {
	logger := k.Logger(ctx)

	// calculate the age of the evidence
	t := ctx.BlockHeader().Time
	age := t.Sub(timestamp)

	// fetch the validator public key
	consAddr := sdk.ConsAddress(addr)
	validator, err := k.GetValidatorByConsAddr(ctx, consAddr)
	if err != nil {
		panic(fmt.Sprintf("Validator %s not found", consAddr))
	}

	pubkey := validator.PubKey

	if validator.IsUnbonded() {
		// Defensive.
		// Simulation doesn't take unbonding periods into account, and
		// Tendermint might break this assumption at some point.
		return
	}

	// fetch the validator signing info
	signInfo, found := k.getValidatorSigningInfo(ctx, consAddr)
	if !found {
		panic(fmt.Sprintf("Expected signing info for validator %s but not found", consAddr))
	}

	// validator is already tombstoned
	if signInfo.Tombstoned {
		logger.Info(fmt.Sprintf("Ignored double sign from %s at height %d, validator already tombstoned", sdk.ConsAddress(pubkey.Address()), infractionHeight))
		return
	}

	// double sign confirmed
	logger.Info(fmt.Sprintf("Confirmed double sign from %s at height %d, age of %d", sdk.ConsAddress(pubkey.Address()), infractionHeight, age))

	// We need to retrieve the stake distribution which signed the block, so we subtract ValidatorUpdateDelay from the evidence height.
	// Note that this *can* result in a negative "distributionHeight", up to -ValidatorUpdateDelay,
	// i.e. at the end of the pre-genesis block (none) = at the beginning of the genesis block.
	// That's fine since this is just used to filter unbonding delegations & redelegations.
	distributionHeight := infractionHeight - sdk.ValidatorUpdateDelay

	// get the percentage slash penalty fraction
	fraction := types.SlashFractionDoubleSign

	// Slash validator
	// `power` is the int64 power of the validator as provided to/by
	// Tendermint. This value is validator.Tokens as sent to Tendermint via
	// ABCI, and now received as evidence.
	// The fraction is passed in to separately to slash unbonding and rebonding delegations.
	slashAmount := k.Slash(ctx, consAddr, distributionHeight, fraction)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSlash,
		sdk.NewAttribute(types.AttributeKeyAddress, consAddr.String()),
		sdk.NewAttribute(types.AttributeKeySlashAmount, slashAmount.String()),
		sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueDoubleSign),
	))

	// Jail validator if not already jailed
	// begin unbonding validator if not already unbonding (tombstone)
	if !validator.IsJailed() {
		k.Jail(ctx, consAddr)
	}

	// Set tombstoned to be true
	signInfo.Tombstoned = true

	// Set validator signing info
	k.setValidatorSigningInfo(ctx, consAddr, signInfo)
}

func (k Keeper) getPubkey(ctx sdk.Context, address crypto.Address) (crypto.PubKey, error) {
	store := ctx.KVStore(k.storeKey)
	var pubkey crypto.PubKey
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(types.GetAddrPubkeyRelationKey(address)), &pubkey)
	if err != nil {
		return nil, fmt.Errorf("address %s not found", sdk.ConsAddress(address))
	}
	return pubkey, nil
}

func (k Keeper) addPubkey(ctx sdk.Context, pubkey crypto.PubKey) {
	addr := pubkey.Address()
	k.setAddrPubkeyRelation(ctx, addr, pubkey)
}

func (k Keeper) setAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address, pubkey crypto.PubKey) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(pubkey)
	store.Set(types.GetAddrPubkeyRelationKey(addr), bz)
}

func inGracePeriod(ctx sdk.Context) bool {
	var (
		height           = ctx.BlockHeight()
		gracePeriodStart = ncfg.UpdatesInfo.LastBlock
		gracePeriodEnd   = gracePeriodStart + (ncfg.OneHour * 24 * 4)
	)
	return height >= gracePeriodStart && height <= gracePeriodEnd
}
