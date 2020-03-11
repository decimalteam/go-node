package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	k.set(ctx, types.GetValidatorKey(validator.ValAddress), validator)
}
