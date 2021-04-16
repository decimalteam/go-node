package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NFT non fungible token interface
type NFT interface {
	GetID() string
	GetOwners() TokenOwners
	SetOwners(owners TokenOwners) NFT
	GetCreator() sdk.AccAddress
	GetTokenURI() string
	EditMetadata(tokenURI string) NFT
	String() string
}

type TokenOwner interface {
	GetAddress() sdk.AccAddress
	GetQuantity() sdk.Int
	SetQuantity(quantity sdk.Int) TokenOwner
}

type TokenOwners interface {
	GetOwners() []TokenOwner
	SetOwner(owner TokenOwner) TokenOwners
	GetOwner(owner sdk.AccAddress) TokenOwner
}
