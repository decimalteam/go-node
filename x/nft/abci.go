package nft

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func BeginBlocker(ctx sdk.Context, k Keeper) {
	if ctx.BlockHeight() == 838400 {
		nft , _ := k.GetNFT(ctx , "Colibri" , "56c9fa969e77b35e39d0f0042eac0249077fe079")
		addr , _  := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
		senderOwner := nft.GetOwners().GetOwner(addr)
		senderOwner = senderOwner.RemoveSubTokenID(144)
		senderOwner = senderOwner.SetSubTokenID(144)
		nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))
	}
}
