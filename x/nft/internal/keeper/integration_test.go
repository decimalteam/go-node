package keeper

import (
	"testing"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	Addrs = types.CreateTestAddrs(4)

	Denom1    = "test_denom1"
	Denom2    = "test_denom2"
	Denom3    = "test_denom3"
	ID1       = "1"
	ID2       = "2"
	ID3       = "3"
	TokenURI1 = "https://google.com/token-1.json"
	TokenURI2 = "https://google.com/token-2.json"
)

func createTestApp(t *testing.T, isCheckTx bool) (sdk.Context, *codec.Codec, Keeper) {
	ctx, nftKeeper := CreateTestInput(t, isCheckTx, 10000000)

	return ctx, MakeTestCodec(), nftKeeper
}
