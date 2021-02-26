package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NFT non fungible token interface
type NFT interface {
	GetID() string
	GetOwners() TokenOwners
	SetOwner(owner TokenOwner)
	GetCreator() sdk.AccAddress
	GetTokenURI() string
	EditMetadata(tokenURI string)
	String() string
}

type TokenOwner interface {
	GetAddress() sdk.AccAddress
	GetQuantity() sdk.Int
	SetQuantity(quantity sdk.Int)
}

type TokenOwners interface {
	GetOwners() []TokenOwner
	SetOwner(owner TokenOwner)
	GetOwner(owner sdk.AccAddress) TokenOwner
}
