package simulation

import (
	"fmt"

	clientsim "bitbucket.org/decimalteam/go-node/x/ibc/core/02-client/simulation"
	connectionsim "bitbucket.org/decimalteam/go-node/x/ibc/core/03-connection/simulation"
	channelsim "bitbucket.org/decimalteam/go-node/x/ibc/core/04-channel/simulation"
	host "bitbucket.org/decimalteam/go-node/x/ibc/core/24-host"
	"bitbucket.org/decimalteam/go-node/x/ibc/core/keeper"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding ibc type.
func NewDecodeStore(k keeper.Keeper) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		if res, found := clientsim.NewDecodeStore(k.ClientKeeper, kvA, kvB); found {
			return res
		}

		if res, found := connectionsim.NewDecodeStore(k.Codec(), kvA, kvB); found {
			return res
		}

		if res, found := channelsim.NewDecodeStore(k.Codec(), kvA, kvB); found {
			return res
		}

		panic(fmt.Sprintf("invalid %s key prefix: %s", host.ModuleName, string(kvA.Key)))
	}
}
