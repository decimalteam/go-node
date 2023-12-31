package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/decimalteam/go-node/x/nft/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SortedIntArray []int64

// IDCollection defines a set of nft ids that belong to a specific
// collection
type IDCollection struct {
	Denom string            `json:"denom" yaml:"denom"`
	IDs   SortedStringArray `json:"ids" yaml:"ids"`
}

// SortedStringArray is an array of strings whose sole purpose is to help with find
type SortedStringArray []string

// String is the string representation
func (sa SortedStringArray) String() string { return strings.Join(sa[:], ",") }

// NewIDCollection creates a new IDCollection instance
func NewIDCollection(denom string, ids []string) IDCollection {
	return IDCollection{
		Denom: strings.TrimSpace(denom),
		IDs:   SortedStringArray(ids).Sort(),
	}
}

// Exists determines whether an ID is in the IDCollection
func (idCollection IDCollection) Exists(id string) (exists bool) {
	index := idCollection.IDs.find(id)
	return index != -1
}

// AddID adds an ID to the idCollection
func (idCollection IDCollection) AddID(id string) IDCollection {
	idCollection.IDs = append(idCollection.IDs, id).Sort()
	return idCollection
}

// DeleteID deletes an ID from an ID Collection
func (idCollection IDCollection) DeleteID(id string) (IDCollection, error) {
	index := idCollection.IDs.find(id)
	if index == -1 {
		return idCollection, ErrUnknownNFT(idCollection.Denom, id)
	}

	idCollection.IDs = append(idCollection.IDs[:index], idCollection.IDs[index+1:]...)

	return idCollection, nil
}

// Supply gets the total supply of NFTIDs of a balance
func (idCollection IDCollection) Supply() int {
	return len(idCollection.IDs)
}

// String follows stringer interface
func (idCollection IDCollection) String() string {
	return fmt.Sprintf(`Denom: 			%s
IDs:        	%s`,
		idCollection.Denom,
		strings.Join(idCollection.IDs, ","),
	)
}

// ----------------------------------------------------------------------------
// Owners

// IDCollections is an array of ID Collections whose sole purpose is for find
type IDCollections []IDCollection

// String follows stringer interface
func (idCollections IDCollections) String() string {
	if len(idCollections) == 0 {
		return ""
	}

	out := ""
	for _, idCollection := range idCollections {
		out += fmt.Sprintf("%v\n", idCollection.String())
	}
	return out[:len(out)-1]
}

// Append appends IDCollections to IDCollections
func (idCollections IDCollections) Append(idCollections2 ...IDCollection) IDCollections {
	return append(idCollections, idCollections2...).Sort()
}
func (idCollections IDCollections) find(denom string) int {
	return FindUtil(idCollections, denom)
}

// Owner of non fungible tokens
type Owner struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	IDCollections IDCollections  `json:"idCollections" yaml:"idCollections"`
}

// NewOwner creates a new Owner
func NewOwner(owner sdk.AccAddress, idCollections ...IDCollection) Owner {
	return Owner{
		Address:       owner,
		IDCollections: idCollections,
	}
}

// Supply gets the total supply of an Owner
func (owner Owner) Supply() int {
	total := 0
	for _, idCollection := range owner.IDCollections {
		total += idCollection.Supply()
	}
	return total
}

// GetIDCollection gets the IDCollection from the owner
func (owner Owner) GetIDCollection(denom string) (IDCollection, bool) {
	index := owner.IDCollections.find(denom)
	if index == -1 {
		return IDCollection{}, false
	}
	return owner.IDCollections[index], true
}

// UpdateIDCollection updates the ID Collection of an owner
func (owner Owner) UpdateIDCollection(idCollection IDCollection) (Owner, error) {
	index := owner.IDCollections.find(idCollection.Denom)
	if index == -1 {
		return owner, ErrUnknownCollection(idCollection.Denom)
	}

	owner.IDCollections = append(append(owner.IDCollections[:index], idCollection), owner.IDCollections[index+1:]...)

	return owner, nil
}

// DeleteID deletes an ID from an owners ID Collection
func (owner Owner) DeleteID(denom string, id string) (Owner, error) {
	idCollection, found := owner.GetIDCollection(denom)
	if !found {
		return owner, ErrUnknownNFT(denom, id)
	}
	idCollection, err := idCollection.DeleteID(id)
	if err != nil {
		return owner, err
	}
	owner, err = owner.UpdateIDCollection(idCollection)
	if err != nil {
		return owner, err
	}
	return owner, nil
}

// String follows stringer interface
func (owner Owner) String() string {
	return fmt.Sprintf(`
	Address: 				%s
	IDCollections:        	%s`,
		owner.Address,
		owner.IDCollections.String(),
	)
}
func (sa SortedStringArray) find(el string) (idx int) {
	return FindUtil(sa, el)
}

// ----------------------------------------------------------------------------
// TokenOwner

type TokenOwner struct {
	Address     sdk.AccAddress `json:"address"`
	SubTokenIDs SortedIntArray `json:"sub_token_ids"`
}

func NewTokenOwner(address sdk.AccAddress, subTokenIDs []int64) TokenOwner {
	return TokenOwner{
		Address:     address,
		SubTokenIDs: subTokenIDs,
	}
}

func (t TokenOwner) GetAddress() sdk.AccAddress {
	return t.Address
}

func (t TokenOwner) GetSubTokenIDs() []int64 {
	return t.SubTokenIDs
}

func (t TokenOwner) SetSubTokenID(subTokenID int64) exported.TokenOwner {
	index := t.SubTokenIDs.Find(subTokenID)
	if index == -1 {
		t.SubTokenIDs = append(t.SubTokenIDs, subTokenID).Sort()
	} else {
		t.SubTokenIDs[index] = subTokenID
	}
	return t
}

func (t TokenOwner) RemoveSubTokenID(subTokenID int64) exported.TokenOwner {
	index := t.SubTokenIDs.Find(subTokenID)
	if index != -1 {
		t.SubTokenIDs = append(t.SubTokenIDs[:index], t.SubTokenIDs[index+1:]...)
	}
	return t
}

func (t TokenOwner) SortSubTokensFix() exported.TokenOwner {
	t.SubTokenIDs = t.SubTokenIDs.Sort()
	return t
}

func (t TokenOwner) String() string {
	return fmt.Sprintf("%s %s", t.Address, t.SubTokenIDs.String())
}

// ----------------------------------------------------------------------------
// TokenOwners

type TokenOwners struct {
	Owners []exported.TokenOwner `json:"owners"`
}

func (t TokenOwners) GetOwners() []exported.TokenOwner {
	return t.Owners
}

func (t TokenOwners) SetOwner(owner exported.TokenOwner) exported.TokenOwners {
	for i, o := range t.Owners {
		if o.GetAddress().Equals(owner.GetAddress()) {
			t.Owners[i] = owner
			return t
		}
	}

	t.Owners = append(t.Owners, TokenOwner{
		Address:     owner.GetAddress(),
		SubTokenIDs: owner.GetSubTokenIDs(),
	})

	return t
}

func (t TokenOwners) GetOwner(address sdk.AccAddress) exported.TokenOwner {
	for _, owner := range t.Owners {
		if owner.GetAddress().Equals(address) {
			return owner
		}
	}
	return nil
}

func (t TokenOwners) String() string {
	if len(t.Owners) == 0 {
		return ""
	}

	out := ""
	for _, owner := range t.Owners {
		out += fmt.Sprintf("%v\n", owner)
	}
	return out[:len(out)-1]
}

type TokenOwnersJSON struct {
	Owners []TokenOwner `json:"owners"`
}

func (t *TokenOwners) UnmarshalJSON(b []byte) error {
	var owners TokenOwnersJSON
	err := json.Unmarshal(b, &owners)
	if err != nil {
		return err
	}

	for _, owner := range owners.Owners {
		t.Owners = append(t.Owners, owner)
	}
	return nil
}

//-----------------------------------------------------------------------------
// Sort and Findable interface for SortedStringArray

func (sa SortedStringArray) ElAtIndex(index int) string { return sa[index] }
func (sa SortedStringArray) Len() int                   { return len(sa) }
func (sa SortedStringArray) Less(i, j int) bool {
	return strings.Compare(sa[i], sa[j]) == -1
}
func (sa SortedStringArray) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

var _ sort.Interface = SortedStringArray{}

// Sort is a helper function to sort the set of strings in place
func (sa SortedStringArray) Sort() SortedStringArray {
	sort.Sort(sa)
	return sa
}

//-----------------------------------------------------------------------------
// Sort and Findable interface for SortedIntArray

func (sa SortedIntArray) ElAtIndex(index int) int64 { return sa[index] }
func (sa SortedIntArray) Len() int                  { return len(sa) }
func (sa SortedIntArray) Less(i, j int) bool {
	return sa[i] < sa[j]
}
func (sa SortedIntArray) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

var _ sort.Interface = SortedStringArray{}

// Sort is a helper function to sort the set of strings in place
func (sa SortedIntArray) Sort() SortedIntArray {
	sort.Sort(sa)
	return sa
}

func (sa SortedIntArray) Find(el int64) (idx int) {
	return FindUtilInt64(sa, el)
}

// String is the string representation
func (sa SortedIntArray) String() string {
	str := make([]string, sa.Len())
	for i, v := range sa {
		str[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(str[:], ",")
}

//-----------------------------------------------------------------------------
// Sort and Findable interface for IDCollections

func (idCollections IDCollections) ElAtIndex(index int) string { return idCollections[index].Denom }
func (idCollections IDCollections) Len() int                   { return len(idCollections) }
func (idCollections IDCollections) Less(i, j int) bool {
	return strings.Compare(idCollections[i].Denom, idCollections[j].Denom) == -1
}
func (idCollections IDCollections) Swap(i, j int) {
	idCollections[i].Denom, idCollections[j].Denom = idCollections[j].Denom, idCollections[i].Denom
}

var _ sort.Interface = IDCollections{}

// Sort is a helper function to sort the set of strings in place
func (idCollections IDCollections) Sort() IDCollections {
	sort.Sort(idCollections)
	return idCollections
}
