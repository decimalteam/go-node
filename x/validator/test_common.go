package validator

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/auth"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

// nolint: deadcode unused
var (
	priv1 = ed25519.GenPrivKey()
	addr1 = decsdk.AccAddress(priv1.PubKey().Address())
	priv2 = ed25519.GenPrivKey()
	addr2 = decsdk.AccAddress(priv2.PubKey().Address())
	addr3 = decsdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	priv4 = ed25519.GenPrivKey()
	addr4 = decsdk.AccAddress(priv4.PubKey().Address())
	coins = sdk.Coins{sdk.NewCoin("foocoin", sdk.NewInt(10))}
	fee   = auth.NewStdFee(
		100000,
		sdk.Coins{sdk.NewCoin("foocoin", sdk.NewInt(0))},
	)

	commissionRates = sdk.ZeroDec()
)

func NewTestMsgDeclareCandidate(address decsdk.ValAddress, pubKey crypto.PubKey, amt sdk.Int) MsgDeclareCandidate {
	return types.NewMsgDeclareCandidate(
		address, pubKey, commissionRates, sdk.NewCoin(types.DefaultBondDenom, amt), types.Description{}, decsdk.AccAddress(address),
	)
}

func NewTestMsgDeclareCandidateWithCommission(address decsdk.ValAddress, pubKey crypto.PubKey,
	amt sdk.Int, commissionRate sdk.Dec) MsgDeclareCandidate {

	return types.NewMsgDeclareCandidate(
		address, pubKey, commissionRate, sdk.NewCoin(types.DefaultBondDenom, amt), types.Description{}, decsdk.AccAddress(address),
	)
}

func NewTestMsgDelegate(delAddr decsdk.AccAddress, valAddr decsdk.ValAddress, amt sdk.Int) MsgDelegate {
	amount := sdk.NewCoin(types.DefaultBondDenom, amt)
	return types.NewMsgDelegate(valAddr, delAddr, amount)
}
