package types

import (
	"bitbucket.org/decimalteam/go-node/x/nft/exported"
	"bytes"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: deadcode unused
var (
	addresses = CreateTestAddrs(500)

	denom     = "denom"
	denom2    = "test-denom2"
	denom3    = "test-denom3"
	id        = "1"
	id2       = "2"
	id3       = "3"
	address   = addresses[0]
	address2  = addresses[1]
	address3  = addresses[1]
	tokenURI  = "https://google.com/token-1.json"
	tokenURI2 = "https://google.com/token-2.json"
)

// MakeTestCodec creates a codec for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*exported.NFT)(nil), nil)
	cdc.RegisterInterface((*exported.TokenOwners)(nil), nil)
	cdc.RegisterInterface((*exported.TokenOwner)(nil), nil)
	cdc.RegisterConcrete(&BaseNFT{}, "nft/BaseNFT", nil)
	cdc.RegisterConcrete(&IDCollection{}, "nft/IDCollection", nil)
	cdc.RegisterConcrete(&Collection{}, "nft/Collection", nil)
	cdc.RegisterConcrete(&Owner{}, "nft/Owner", nil)
	cdc.RegisterConcrete(&TokenOwner{}, "nft/TokenOwner", nil)
	cdc.RegisterConcrete(&TokenOwners{}, "nft/TokenOwners", nil)
	cdc.RegisterConcrete(MsgTransferNFT{}, "nft/msg_transfer", nil)
	cdc.RegisterConcrete(MsgEditNFTMetadata{}, "nft/msg_edit_metadata", nil)
	cdc.RegisterConcrete(MsgMintNFT{}, "nft/msg_mint", nil)
	cdc.RegisterConcrete(MsgBurnNFT{}, "nft/msg_burn", nil)

	// Register AppAccount
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/coin/base_account", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

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
