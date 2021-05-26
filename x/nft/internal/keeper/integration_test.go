package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

// nolint: deadcode unused
var (
	addresses = types.CreateTestAddrs(4)

	denom     = "test-denom"
	denom2    = "test-denom2"
	denom3    = "test-denom3"
	id        = "1"
	id2       = "2"
	id3       = "3"
	address   = addresses[0]
	address2  = addresses[1]
	address3  = addresses[2]
	tokenURI  = "https://google.com/token-1.json"
	tokenURI2 = "https://google.com/token-2.json"
)

func createTestApp(t *testing.T, isCheckTx bool) (sdk.Context, *codec.Codec, keeper.Keeper) {
	ctx, nftKeeper := keeper.CreateTestInput(t, isCheckTx, 1000)

	return ctx, types.MakeTestCodec(), nftKeeper
}
