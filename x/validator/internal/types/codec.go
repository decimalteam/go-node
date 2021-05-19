package types

import (
	"bitbucket.org/decimalteam/go-node/x/validator/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
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

// ModuleCdc defines the module codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterLegacyAminoCodec(ModuleCdc)
	codec2.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
