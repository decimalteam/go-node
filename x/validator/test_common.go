package validator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// nolint: deadcode unused
var (
	priv1 = ed25519.GenPrivKey()
	addr1 = sdk.AccAddress(priv1.PubKey().Address())
	priv2 = ed25519.GenPrivKey()
	addr2 = sdk.AccAddress(priv2.PubKey().Address())
	addr3 = sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	priv4 = ed25519.GenPrivKey()
	addr4 = sdk.AccAddress(priv4.PubKey().Address())
	coins = sdk.Coins{sdk.NewCoin("foocoin", sdk.NewInt(10))}
	fee   = auth.NewStdFee(
		100000,
		sdk.Coins{sdk.NewCoin("foocoin", sdk.NewInt(0))},
	)

	commissionRates = sdk.ZeroDec()
)

func NewTestMsgDeclareCandidate(address sdk.ValAddress, pubKey crypto.PubKey, amt sdk.Int) MsgDeclareCandidate {
	return types.NewMsgDeclareCandidate(
		address, pubKey, commissionRates, sdk.NewCoin(DefaultBondDenom, amt), types.Description{}, sdk.AccAddress(address),
	)
}

func NewTestMsgDeclareCandidateWithCommission(address sdk.ValAddress, pubKey crypto.PubKey,
	amt sdk.Int, commissionRate sdk.Dec) MsgDeclareCandidate {

	return types.NewMsgDeclareCandidate(
		address, pubKey, commissionRate, sdk.NewCoin(DefaultBondDenom, amt), types.Description{}, sdk.AccAddress(address),
	)
}

func NewTestMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amt sdk.Int) MsgDelegate {
	amount := sdk.NewCoin(DefaultBondDenom, amt)
	return types.NewMsgDelegate(valAddr, delAddr, amount)
}
