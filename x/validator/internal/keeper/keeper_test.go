package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 0)
	expParams := types.DefaultParams()

	//check that the empty keeper loads the default
	resParams := keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))

	//modify a params, save, and retrieve
	expParams.MaxValidators = 777
	keeper.SetParams(ctx, expParams)
	resParams = keeper.GetParams(ctx)
	require.True(t, expParams.Equal(resParams))
}

//func TestChangeCodec(t *testing.T) {
//	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 0)
//	validatorAddr := sdk.ValAddress(Addrs[0])
//	delegatorAddr := sdk.AccAddress(validatorAddr)
//
//	var cdc = codec.New()
//
//	// Register AppAccount
//	cdc.RegisterInterface((*authexported.Account)(nil), nil)
//	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/validator/base_account", nil)
//	supply.RegisterCodec(cdc)
//	nft.RegisterCodec(cdc)
//	codec.RegisterCrypto(cdc)
//
//	cdc.Seal()
//
//	keeper.cdc = cdc
//
//	// create nft
//	const denom = "denom1"
//	const tokenID = "token1"
//	quantity := sdk.NewInt(100)
//	reserve := sdk.NewInt(100)
//
//	competitionTime := ctx.BlockTime().Add(time.Second * 5)
//
//	keeper.SetUnbondingDelegation(ctx, types.NewUnbondingDelegation(
//		delegatorAddr,
//		validatorAddr,
//		types.NewUnbondingDelegationEntry(ctx.BlockHeight(), competitionTime, sdk.NewCoin(denom, quantity)),
//	))
//
//	delegations := keeper.GetAllUnbondingDelegations(ctx)
//
//	var cdc2 = codec.New()
//
//	cdc2.RegisterInterface((*exported.UnbondingDelegationEntryI)(nil), nil)
//	cdc2.RegisterConcrete(types.UnbondingDelegationEntry{}, "validator/unbonding_delegation_entry", nil)
//	cdc2.RegisterConcrete(types.UnbondingDelegationNFTEntry{}, "validator/unbonding_delegation_nft_entry", nil)
//
//	// Register AppAccount
//	cdc2.RegisterInterface((*authexported.Account)(nil), nil)
//	cdc2.RegisterConcrete(&auth.BaseAccount{}, "test/validator/base_account", nil)
//	supply.RegisterCodec(cdc2)
//	nft.RegisterCodec(cdc2)
//	codec.RegisterCrypto(cdc2)
//	cdc2.Seal()
//
//	keeper.cdc = cdc2
//
//
//}
