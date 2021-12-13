package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	k.SetBaseDenom()

	if ctx.BlockHeight() == 7_348_751 {
		nft, _ := k.GetNFT(ctx, "Signs_of_Zodiac", "7206c987e670ad90b7e7c9ffba2ba90bb061c533")
		addr, _ := sdk.AccAddressFromBech32("dx1hn6n8rwgtc53mmjmzn7d6e5sg44ztq0c7kpl90")
		senderOwner := nft.GetOwners().GetOwner(addr)

		senderOwner = senderOwner.SortSubTokensFix()

		nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))

		collection, found := k.GetCollection(ctx, "Signs_of_Zodiac")
		if !found {
			return
		}

		collection.NFTs, _ = collection.NFTs.Update("7206c987e670ad90b7e7c9ffba2ba90bb061c533", nft)
		k.SetCollection(ctx, "Signs_of_Zodiac", collection)
	}
}
