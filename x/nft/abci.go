package nft

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	k.SetBaseDenom()

	if ctx.BlockHeight() == 889401 {
		problems := [][3]string{
			[3]string{"Fur_and_Fury", "ba00925f5e66413a82277987d440b6bdd3226c94", "dx1nafxm7gn4kmyjtctya7cshj4nj956k5tq5p9wu"},
			[3]string{"Fur_and_Fury", "ba00925f5e66413a82277987d440b6bdd3226c94", "dx16xh0e5unylk28kwrctsuuazkwwgrk2hvsv5us3"},
			[3]string{"Spacebot_collection", "52d593f23eb27f2960f25d875be7583171939022", "dx1hn6n8rwgtc53mmjmzn7d6e5sg44ztq0c7kpl90"},
			[3]string{"Spacebot_collection", "59383cadbd72b1fac25bdafc355f4d19d4ee8ee8", "dx1hn6n8rwgtc53mmjmzn7d6e5sg44ztq0c7kpl90"},
			[3]string{"Spacebot_collection", "b33ab516a420c791056269cdad1aab312eaa24ba", "dx1hn6n8rwgtc53mmjmzn7d6e5sg44ztq0c7kpl90"},
			[3]string{"Signs_of_Zodiac", "d16b476d900589269948ab2c3325b090512577e1", "dx1xraedegpypve2ga5yzgt42nvhqq0qp6a7kck8s"},
			[3]string{"Signs_of_Zodiac", "fc43af59fbd228a2044f0887214ea9e331e7a73b", "dx1hn6n8rwgtc53mmjmzn7d6e5sg44ztq0c7kpl90"},
			[3]string{"QR", "034e38b4123bb36d6ba8e05baae94d928953be72", "dx1ehuyv5sn46pe54ntzanlpncpjpqt2lc7p3uahu"},
			[3]string{"BlindDevPortfoleo", "db711f3b8791ad211d70505f99c5eb3b78890f29", "dx1hn6n8rwgtc53mmjmzn7d6e5sg44ztq0c7kpl90"},
		}

		for _, arr := range problems {
			nft, _ := k.GetNFT(ctx, arr[0], arr[1])
			addr, _ := sdk.AccAddressFromBech32(arr[2])

			senderOwner := nft.GetOwners().GetOwner(addr)
			senderOwner = senderOwner.SortSubTokensFix()

			nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))
			collection, found := k.GetCollection(ctx, arr[0])
			if !found {
				fmt.Printf("collection undefined [%s, %s, %s]\n", arr[0], arr[1], arr[2])
			}

			collection.NFTs, _ = collection.NFTs.Update(arr[1], nft)
			k.SetCollection(ctx, arr[0], collection)
		}
	}
}
