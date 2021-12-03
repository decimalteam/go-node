package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	/*if ctx.BlockHeight() == 889401 {
		problems := [][3]string{
			[3]string{"BravoMeme", "98bc347afb23484bc802f36083ca2672d7a72e2a", "dx13fgegk6aaymdzhjmx0cdt9g3y7a3se0azf675e"},
			[3]string{"ter", "8271ad9f2a18c5a6a130f73dbda43f74c114a684", "dx1ulat3e0s25g6amuxd6p059enkl5pe5rlu3u9tu"},
			[3]string{"KunRA", "b91b5822909db1a23c6ddce6f85a66d108b54012", "dx1xv02l6dq3jxcdernk9ur7s9pkewxclxcsza0sa"},
		}

		for _, arr := range problems {
			nft, _ := k.GetNFT(ctx, arr[0], arr[1])
			addr, _ := sdk.AccAddressFromBech32(arr[2])
			senderOwner := nft.GetOwners().GetOwner(addr)

			senderOwner = senderOwner.SortSubTokensFix()

			nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))

			collection, found := k.GetCollection(ctx, arr[0])
			if !found {
				fmt.Println("Error")
			}

			collection.NFTs, _ = collection.NFTs.Update(arr[1], nft)

			k.SetCollection(ctx, arr[0], collection)
		}

	}*/

}
