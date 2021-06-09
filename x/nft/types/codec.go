package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"

	"bitbucket.org/decimalteam/go-node/x/nft/exported"
)

// RegisterCodec concrete types on codec
func RegisterCodec(cdc *codec.LegacyAmino) {
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
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("NFT", (*exported.NFT)(nil))
	registry.RegisterInterface("TokenOwners", (*exported.TokenOwners)(nil))
	registry.RegisterInterface("TokenOwner", (*exported.TokenOwner)(nil))
	registry.RegisterImplementations(&BaseNFT{})
	registry.RegisterImplementations(&IDCollection{})
	registry.RegisterImplementations(&Collection{})
	registry.RegisterImplementations(&Owner{})
	registry.RegisterImplementations(&TokenOwner{})
	registry.RegisterImplementations(&TokenOwners{})
	registry.RegisterImplementations(MsgTransferNFT{})
	registry.RegisterImplementations(MsgEditNFTMetadata{})
	registry.RegisterImplementations(MsgMintNFT{})
	registry.RegisterImplementations(MsgBurnNFT{})
}

// ModuleCdc generic sealed codec to be used throughout this module
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	codec2.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
