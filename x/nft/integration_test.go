package nft

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

// nolint: deadcode unused
var (
	Addrs = types.CreateTestAddrs(100)

	Denom1    = "test-denom1"
	Denom2    = "test-denom2"
	Denom3    = "test-denom3"
	ID1       = "1"
	ID2       = "2"
	ID3       = "3"
	TokenURI1 = "https://google.com/token-1.json"
	TokenURI2 = "https://google.com/token-2.json"
)

func createTestApp(t *testing.T, isCheckTx bool) (sdk.Context, *codec.Codec, keeper.Keeper) {
	ctx, nftKeeper := keeper.CreateTestInput(t, isCheckTx, 1000)

	return ctx, keeper.MakeTestCodec(), nftKeeper
}

// CheckInvariants checks the invariants
func CheckInvariants(k Keeper, ctx sdk.Context) bool {
	collectionsSupply := make(map[string]int)
	ownersCollectionsSupply := make(map[string]int)

	k.IterateCollections(ctx, func(collection types.Collection) bool {
		collectionsSupply[collection.Denom] = collection.Supply()
		return false
	})

	owners := k.GetOwners(ctx)
	for _, owner := range owners {
		for _, idCollection := range owner.IDCollections {
			ownersCollectionsSupply[idCollection.Denom] += idCollection.Supply()
		}
	}

	for denom, supply := range collectionsSupply {
		if supply != ownersCollectionsSupply[denom] {
			fmt.Printf("denom is %s, supply is %d, ownerSupply is %d", denom, supply, ownersCollectionsSupply[denom])
			return false
		}
	}
	return true
}
