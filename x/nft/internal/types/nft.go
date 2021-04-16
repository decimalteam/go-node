package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"bitbucket.org/decimalteam/go-node/x/nft/exported"
)

var _ exported.NFT = (*BaseNFT)(nil)

// BaseNFT non fungible token definition
type BaseNFT struct {
	ID       string               `json:"id,omitempty" yaml:"id"` // id of the token; not exported to clients
	Owners   exported.TokenOwners `json:"owners" yaml:"owners"`   // account addresses that owns the NFT
	Creator  sdk.AccAddress       `json:"creator"`
	TokenURI string               `json:"token_uri" yaml:"token_uri"` // optional extra properties available for querying
}

// NewBaseNFT creates a new NFT instance
func NewBaseNFT(id string, creator, owner sdk.AccAddress, tokenURI string, quantity sdk.Int) BaseNFT {
	return BaseNFT{
		ID: id,
		Owners: &TokenOwners{Owners: []exported.TokenOwner{&TokenOwner{
			Address:  owner,
			Quantity: quantity,
		}}},
		TokenURI: strings.TrimSpace(tokenURI),
		Creator:  creator,
	}
}

// GetID returns the ID of the token
func (bnft BaseNFT) GetID() string { return bnft.ID }

// GetOwner returns the account address that owns the NFT
func (bnft BaseNFT) GetOwners() exported.TokenOwners { return bnft.Owners }

// GetTokenURI returns the path to optional extra properties
func (bnft BaseNFT) GetTokenURI() string { return bnft.TokenURI }

func (bnft BaseNFT) GetCreator() sdk.AccAddress { return bnft.Creator }

// EditMetadata edits metadata of an nft
func (bnft BaseNFT) EditMetadata(tokenURI string) exported.NFT {
	bnft.TokenURI = tokenURI
	return bnft
}

func (bnft BaseNFT) SetOwners(owners exported.TokenOwners) exported.NFT {
	bnft.Owners = owners
	return bnft
}

func (bnft BaseNFT) String() string {
	return fmt.Sprintf(`ID:				%s
Owners:			%s
TokenURI:		%s`,
		bnft.ID,
		bnft.Owners,
		bnft.TokenURI,
	)
}

func TransferNFT(nft exported.NFT, sender, recipient sdk.AccAddress, quantity sdk.Int) (exported.NFT, error) {
	senderOwner := nft.GetOwners().GetOwner(sender)
	if senderOwner.GetQuantity().LT(quantity) {
		return nil, ErrInvalidQuantity
	}

	senderOwner = senderOwner.SetQuantity(senderOwner.GetQuantity().Sub(quantity))

	recipientOwner := nft.GetOwners().GetOwner(recipient)
	if recipientOwner == nil {
		nft = nft.SetOwners(nft.GetOwners().SetOwner(NewTokenOwner(recipient, quantity)))
		nft = nft.SetOwners(nft.GetOwners().SetOwner(senderOwner))
		return nft, nil
	}

	recipientOwner = recipientOwner.SetQuantity(recipientOwner.GetQuantity().Add(quantity))

	return nft, nil
}

// ----------------------------------------------------------------------------
// NFT

// NFTs define a list of NFT
type NFTs []exported.NFT

// NewNFTs creates a new set of NFTs
func NewNFTs(nfts ...exported.NFT) NFTs {
	if len(nfts) == 0 {
		return NFTs{}
	}
	return NFTs(nfts).Sort()
}

// Append appends two sets of NFTs
func (nfts NFTs) Append(nftsB ...exported.NFT) NFTs {
	return append(nfts, nftsB...).Sort()
}

// Find returns the searched collection from the set
func (nfts NFTs) Find(id string) (nft exported.NFT, found bool) {
	index := nfts.find(id)
	if index == -1 {
		return nft, false
	}
	return nfts[index], true
}

// Update removes and replaces an NFT from the set
func (nfts NFTs) Update(id string, nft exported.NFT) (NFTs, bool) {
	index := nfts.find(id)
	if index == -1 {
		return nfts, false
	}

	return append(append(nfts[:index], nft), nfts[index+1:]...), true
}

// Remove removes an NFT from the set of NFTs
func (nfts NFTs) Remove(id string) (NFTs, bool) {
	index := nfts.find(id)
	if index == -1 {
		return nfts, false
	}

	return append(nfts[:index], nfts[index+1:]...), true
}

// String follows stringer interface
func (nfts NFTs) String() string {
	if len(nfts) == 0 {
		return ""
	}

	out := ""
	for _, nft := range nfts {
		out += fmt.Sprintf("%v\n", nft.String())
	}
	return out[:len(out)-1]
}

// Empty returns true if there are no NFTs and false otherwise.
func (nfts NFTs) Empty() bool {
	return len(nfts) == 0
}

func (nfts NFTs) find(id string) int {
	return FindUtil(nfts, id)
}

// ----------------------------------------------------------------------------
// Encoding

// NFTJSON is the exported NFT format for clients
type NFTJSON map[string]BaseNFT

// MarshalJSON for NFTs
func (nfts NFTs) MarshalJSON() ([]byte, error) {
	nftJSON := make(NFTJSON)
	for _, nft := range nfts {
		id := nft.GetID()
		nftJSON[id] = *nft.(*BaseNFT)
	}
	return json.Marshal(nftJSON)
}

// UnmarshalJSON for NFTs
func (nfts *NFTs) UnmarshalJSON(b []byte) error {
	nftJSON := make(NFTJSON)
	if err := json.Unmarshal(b, &nftJSON); err != nil {
		return err
	}

	for _, nft := range nftJSON {
		bnft := nft
		*nfts = append(*nfts, &bnft)
	}
	return nil
}

// Findable and Sort interfaces
func (nfts NFTs) ElAtIndex(index int) string { return nfts[index].GetID() }
func (nfts NFTs) Len() int                   { return len(nfts) }
func (nfts NFTs) Less(i, j int) bool         { return strings.Compare(nfts[i].GetID(), nfts[j].GetID()) == -1 }
func (nfts NFTs) Swap(i, j int)              { nfts[i], nfts[j] = nfts[j], nfts[i] }

var _ sort.Interface = NFTs{}

// Sort is a helper function to sort the set of coins in place
func (nfts NFTs) Sort() NFTs {
	sort.Sort(nfts)
	return nfts
}
