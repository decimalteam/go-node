package keys

import (
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/cosmos/cosmos-sdk/client/flags"
)

func TestCommands(t *testing.T) {
	rootCommands := Commands("home")
	assert.NotNil(t, rootCommands)

	// Commands are registered
	assert.Equal(t, 11, len(rootCommands.Commands()))
}

func TestMain(m *testing.M) {
	viper.Set(flags.FlagKeyringBackend, keyring.BackendTest)
	os.Exit(m.Run())
}
