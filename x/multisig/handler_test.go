package multisig

import (
	mul "bitbucket.org/decimalteam/go-node/x/multisig/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMsgCreateWallet(t *testing.T) {
	ctx, keeper, _, _, _ := mul.CreateTestInput(t, false)

	msgCreateWallet := types.NewMsgCreateWallet(mul.Addrs[0], []sdk.AccAddress{mul.Addrs[0], mul.Addrs[1]}, []uint{1, 1}, 2)
	res, err := handleMsgCreateWallet(ctx, keeper, msgCreateWallet)
	require.NoError(t, err)
	require.NotNil(t, res)

	wallet, err := types.NewWallet([]sdk.AccAddress{mul.Addrs[0], mul.Addrs[1]}, []uint{1, 1}, 2, ctx.TxBytes())
	require.NoError(t, err)
	require.NotNil(t, wallet)

	w := keeper.GetWallet(ctx, wallet.Address.String())
	require.NotNil(t, w)

	require.Equal(t, *wallet, w)
}

func TestMsgCreateWalletWithExistWallet(t *testing.T) {
	ctx, keeper, _, _, _ := mul.CreateTestInput(t, false)

	msgCreateWallet := types.NewMsgCreateWallet(mul.Addrs[0], []sdk.AccAddress{mul.Addrs[0], mul.Addrs[1]}, []uint{1, 1}, 2)
	res, err := handleMsgCreateWallet(ctx, keeper, msgCreateWallet)
	require.NoError(t, err)
	require.NotNil(t, res)

	wallet, err := types.NewWallet([]sdk.AccAddress{mul.Addrs[0], mul.Addrs[1]}, []uint{1, 1}, 2, ctx.TxBytes())
	require.NoError(t, err)
	require.NotNil(t, wallet)

	w := keeper.GetWallet(ctx, wallet.Address.String())
	require.NotNil(t, w)

	require.Equal(t, *wallet, w)

	res, err = handleMsgCreateWallet(ctx, keeper, msgCreateWallet)
	require.Errorf(t, err, fmt.Sprintf("Multi-signature wallet with address %s already exists", wallet.Address.String()))
}

func TestMsgCreateWalletWithExistAddress(t *testing.T) {
	ctx, keeper, _, accountKeeper, bankKeeper := mul.CreateTestInput(t, false)

	wallet, err := types.NewWallet([]sdk.AccAddress{mul.Addrs[0], mul.Addrs[1]}, []uint{1, 1}, 2, ctx.TxBytes())
	require.NoError(t, err)
	require.NotNil(t, wallet)

	accountKeeper.NewAccountWithAddress(ctx, wallet.Address)
	err = bankKeeper.SetCoins(ctx, wallet.Address, sdk.NewCoins(sdk.NewCoin("del", sdk.NewInt(1000))))
	require.NoError(t, err)

	msgCreateWallet := types.NewMsgCreateWallet(mul.Addrs[0], []sdk.AccAddress{mul.Addrs[0], mul.Addrs[1]}, []uint{1, 1}, 2)
	_, err = handleMsgCreateWallet(ctx, keeper, msgCreateWallet)
	require.Errorf(t, err, fmt.Sprintf("Account with address %s already exists", wallet.Address.String()))
}
