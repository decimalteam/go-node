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
	GetReserve() sdk.Int
	GetAllowMint() bool
	String() string
}

type TokenOwner interface {
	GetAddress() sdk.AccAddress
	GetSubTokenIDs() []int64
	SetSubTokenID(id int64) TokenOwner
	RemoveSubTokenID(id int64) TokenOwner
	String() string
}

type TokenOwners interface {
	GetOwners() []TokenOwner
	SetOwner(owner TokenOwner) TokenOwners
	GetOwner(owner sdk.AccAddress) TokenOwner
	String() string
}
