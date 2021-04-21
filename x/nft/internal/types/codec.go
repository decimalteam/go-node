package types

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"bitbucket.org/decimalteam/go-node/x/nft/exported"
)

// RegisterCodec concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
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

// ModuleCdc generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
