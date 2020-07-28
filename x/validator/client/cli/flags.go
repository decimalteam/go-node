package cli

import (
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	flag "github.com/spf13/pflag"
)

const (
	FlagAddressValidator    = "validator"
	FlagAddressValidatorSrc = "addr-validator-source"
	FlagAddressValidatorDst = "addr-validator-dest"
	FlagPubKey              = "pubkey"
	FlagAmount              = "amount"
	FlagSharesAmount        = "shares-amount"
	FlagSharesFraction      = "shares-fraction"
	FlagRewardAddress       = "reward-addr"

	FlagMoniker         = "moniker"
	FlagIdentity        = "identity"
	FlagWebsite         = "website"
	FlagSecurityContact = "security-contact"
	FlagDetails         = "details"

	FlagCommissionRate = "commission-rate"

	FlagGenesisFormat = "genesis-format"
	FlagNodeID        = "node-id"
	FlagIP            = "ip"
)

// common flagsets to add to various functions
var (
	FsPk                = flag.NewFlagSet("", flag.ContinueOnError)
	FsAmount            = flag.NewFlagSet("", flag.ContinueOnError)
	FsShares            = flag.NewFlagSet("", flag.ContinueOnError)
	FsDescriptionCreate = flag.NewFlagSet("", flag.ContinueOnError)
	FsCommissionCreate  = flag.NewFlagSet("", flag.ContinueOnError)
	FsCommissionUpdate  = flag.NewFlagSet("", flag.ContinueOnError)
	FsDescriptionEdit   = flag.NewFlagSet("", flag.ContinueOnError)
	FsValidator         = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsPk.String(FlagPubKey, "", "The Bech32 encoded PubKey of the validator")
	FsAmount.String(FlagAmount, "", "Amount of coins to bond")
	FsShares.String(FlagSharesAmount, "", "Amount of source-shares to either unbond or redelegate as a positive integer or decimal")
	FsShares.String(FlagSharesFraction, "", "Fraction of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1")
	FsDescriptionCreate.String(FlagMoniker, "", "The validator's name")
	FsDescriptionCreate.String(FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	FsDescriptionCreate.String(FlagWebsite, "", "The validator's (optional) website")
	FsDescriptionCreate.String(FlagSecurityContact, "", "The validator's (optional) security contact email")
	FsDescriptionCreate.String(FlagDetails, "", "The validator's (optional) details")
	FsCommissionCreate.String(FlagCommissionRate, "", "The commission rate percentage")
	FsCommissionUpdate.String(FlagCommissionRate, "", "The new commission rate percentage")
	FsDescriptionEdit.String(FlagMoniker, types.DoNotModifyDesc, "The validator's name")
	FsDescriptionEdit.String(FlagIdentity, types.DoNotModifyDesc, "The (optional) identity signature (ex. UPort or Keybase)")
	FsDescriptionEdit.String(FlagWebsite, types.DoNotModifyDesc, "The validator's (optional) website")
	FsDescriptionEdit.String(FlagSecurityContact, types.DoNotModifyDesc, "The validator's (optional) security contact email")
	FsDescriptionEdit.String(FlagDetails, types.DoNotModifyDesc, "The validator's (optional) details")
	FsValidator.String(FlagAddressValidator, "", "The Bech32 address of the validator")
}
