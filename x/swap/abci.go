package swap

import (
	"fmt"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/utils/updates"
)

func BeginBlocker(ctx sdk.Context, k Keeper) {
	// Migrate state to updated prefixes if necessary
	if ctx.BlockHeight() == updates.Update14Block {
		err := k.MigrateToUpdatedPrefixes(ctx)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("failed migrate to updated prefixes: %v", err))
			os.Exit(4)
		}
	}
}
