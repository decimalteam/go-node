package types

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgDeclareCandidate{}, "validator/declare_candidate", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "validator/delegate", nil)
	cdc.RegisterConcrete(MsgDelegateNFT{}, "validator/delegate_nft", nil)
	cdc.RegisterConcrete(MsgUnbond{}, "validator/unbond", nil)
	cdc.RegisterConcrete(MsgUnbondNFT{}, "validator/unbond_nft", nil)
	cdc.RegisterConcrete(MsgEditCandidate{}, "validator/edit_candidate", nil)
	cdc.RegisterConcrete(MsgSetOnline{}, "validator/set_online", nil)
	cdc.RegisterConcrete(MsgSetOffline{}, "validator/set_offline", nil)
	cdc.RegisterInterface((*exported.DelegationI)(nil), nil)
	cdc.RegisterConcrete(Delegation{}, "validator/delegation", nil)
	cdc.RegisterConcrete(DelegationNFT{}, "validator/delegation_nft", nil)
	cdc.RegisterInterface((*exported.UnbondingDelegationEntryI)(nil), nil)
	cdc.RegisterConcrete(UnbondingDelegationEntry{}, "validator/unbonding_delegation_entry", nil)
	cdc.RegisterConcrete(UnbondingDelegationNFTEntry{}, "validator/unbonding_delegation_nft_entry", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeclareCandidate{},
		&MsgDelegate{},
		&MsgDelegateNFT{},
		&MsgUnbond{},
		&MsgUnbondNFT{},
		&MsgEditCandidate{},
		&MsgSetOnline{},
		&MsgSetOffline{},
		&Delegation{},
		&DelegationNFT{},
		&UnbondingDelegationEntry{},
		&UnbondingDelegationNFTEntry{},
	)

	registry.RegisterInterface("DelegationI",
		(*exported.DelegationI)(nil),
	)
	registry.RegisterInterface("UnbondingDelegationEntryI",
		(*exported.UnbondingDelegationEntryI)(nil),
	)
}

// ModuleCdc defines the module codec
var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	codec2.RegisterCrypto(amino)
	amino.Seal()
}
