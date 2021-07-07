package types

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

// CreateTestAddrs creates test addresses
func CreateTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, testAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

// for incode address generation
func testAddr(addr string, bech string) sdk.AccAddress {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}

func TestMsgChainActivateDeactivateValidateBasicMethod(t *testing.T) {
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	Addrs := CreateTestAddrs(100)

	newMsgChainActivate := NewMsgChainActivate(Addrs[0], 1, "del")
	err := newMsgChainActivate.ValidateBasic()
	require.Error(t, err)

	swapServiceAddress, err := sdk.AccAddressFromBech32(ChainActivatorAddress)
	require.NoError(t, err)

	newMsgChainActivate = NewMsgChainActivate(swapServiceAddress, 1, "del")
	err = newMsgChainActivate.ValidateBasic()
	require.NoError(t, err)

	newMsgChainDeactivate := NewMsgChainDeactivate(Addrs[0], 1)
	err = newMsgChainDeactivate.ValidateBasic()
	require.Error(t, err)

	newMsgChainDeactivate = NewMsgChainDeactivate(swapServiceAddress, 1)
	err = newMsgChainDeactivate.ValidateBasic()
	require.NoError(t, err)
}
