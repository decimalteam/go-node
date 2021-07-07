package keeper

/*import (
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/validator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestKeeper_PayRewardsWithMultisigAddress(t *testing.T) {
	ctx, _, keeper, _, _, _ := CreateTestInput(t, false, 1000)
	validatorAddr1 := sdk.ValAddress(Addrs[0])

	wallet, err := multisig.NewWallet([]sdk.AccAddress{addrAcc1, addrAcc2}, []uint{1, 1}, 2, []byte{})
	require.NoError(t, err)
	require.NotNil(t, wallet)

	err = keeper.SetValidator(ctx, types.Validator{
		ValAddress:              validatorAddr1,
		PubKey:                  pk1,
		Tokens:                  sdk.NewInt(100000),
		Status:                  types.Bonded,
		Commission:              sdk.Dec{},
		Jailed:                  false,
		UnbondingCompletionTime: time.Time{},
		UnbondingHeight:         0,
		Description:             types.Description{},
		AccumRewards:            sdk.NewInt(100000),
		RewardAddress:           wallet.Address,
		Online:                  true,
	})
	require.NoError(t, err)

	err = keeper.PayRewards(ctx)
	require.NoError(t, err)
}
*/