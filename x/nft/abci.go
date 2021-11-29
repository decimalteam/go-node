package nft

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func BeginBlocker(ctx sdk.Context, k Keeper) {
	if ctx.BlockHeight() == 838750 {
		nft , _ := k.GetNFT(ctx , "Colibri" , "56c9fa969e77b35e39d0f0042eac0249077fe079")
		addr , _  := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
		senderOwner := nft.GetOwners().GetOwner(addr)
		fmt.Println("Before : ",senderOwner.GetSubTokenIDs())
		senderOwner = senderOwner.SortSubTokensFix()
		fmt.Println("After : ",senderOwner.GetSubTokenIDs())
		nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))
	}
}
