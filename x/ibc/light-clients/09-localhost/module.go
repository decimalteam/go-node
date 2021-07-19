package localhost

import (
	"bitbucket.org/decimalteam/go-node/x/ibc/light-clients/09-localhost/types"
)

// Name returns the IBC client name
func Name() string {
	return types.SubModuleName
}
