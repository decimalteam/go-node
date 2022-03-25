package swap

import (
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// Migrate state to updated prefixes if necessary
	if !k.IsMigratedToUpdatedPrefixes(ctx) {
		err := k.MigrateToUpdatedPrefixes(ctx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed migrate to updated prefixes: %v", err))
			os.Exit(4)
		}
	}
}
