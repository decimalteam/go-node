package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------- IDCollection ---------------------------------------------------

func TestNewIDCollection(t *testing.T) {
	ids := []string{ID1, ID2, ID3}
	idCollection := NewIDCollection(Denom1, ids)
	require.Equal(t, idCollection.Denom, Denom1)
	require.Equal(t, len(idCollection.IDs), 3)
}

func TestIDCollectionExistsMethod(t *testing.T) {
	ids := []string{ID2, ID1}
	idCollection := NewIDCollection(Denom1, ids)
	require.True(t, idCollection.Exists(ID1))
	require.True(t, idCollection.Exists(ID2))
	require.False(t, idCollection.Exists(ID3))
}

func TestIDCollectionAddIDMethod(t *testing.T) {
	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	idCollection = idCollection.AddID(ID3)
	require.Equal(t, len(idCollection.IDs), 3)
}

func TestIDCollectionDeleteIDMethod(t *testing.T) {
	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	newIDCollection, err := idCollection.DeleteID(ID3)
	require.Error(t, err)
	require.Equal(t, idCollection.String(), newIDCollection.String())

	idCollection, err = idCollection.DeleteID(ID2)
	require.NoError(t, err)
	require.Equal(t, len(idCollection.IDs), 1)
}

func TestIDCollectionSupplyMethod(t *testing.T) {
	idCollectionEmpty := IDCollection{}
	require.Equal(t, 0, idCollectionEmpty.Supply())

	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	require.Equal(t, 2, idCollection.Supply())

	idCollection, err := idCollection.DeleteID(ID1)
	require.Nil(t, err)
	require.Equal(t, idCollection.Supply(), 1)

	idCollection, err = idCollection.DeleteID(ID2)
	require.Nil(t, err)
	require.Equal(t, idCollection.Supply(), 0)

	idCollection = idCollection.AddID(ID1)
	require.Nil(t, err)
	require.Equal(t, idCollection.Supply(), 1)
}

func TestIDCollectionStringMethod(t *testing.T) {
	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	require.Equal(t, idCollection.String(), fmt.Sprintf(`Denom: 			%s
IDs:        	%s,%s`, Denom1, ID1, ID2))
}

// ---------------------------------------- IDCollections ---------------------------------------------------

func TestIDCollectionsString(t *testing.T) {
	emptyCollections := IDCollections([]IDCollection{})
	require.Equal(t, emptyCollections.String(), "")

	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	idCollection2 := NewIDCollection(Denom2, ids)

	idCollections := IDCollections([]IDCollection{idCollection, idCollection2})
	require.Equal(t, idCollections.String(), fmt.Sprintf(`Denom: 			%s
IDs:        	%s,%s
Denom: 			%s
IDs:        	%s,%s`, Denom1, ID1, ID2, Denom2, ID1, ID2))
}

// ---------------------------------------- Owner ---------------------------------------------------

func TestNewOwner(t *testing.T) {
	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	idCollection2 := NewIDCollection(Denom2, ids)

	owner := NewOwner(Addrs[0], idCollection, idCollection2)
	require.Equal(t, owner.Address.String(), Addrs[0].String())
	require.Equal(t, len(owner.IDCollections), 2)
}

func TestOwnerSupplyMethod(t *testing.T) {
	owner := NewOwner(Addrs[0])
	require.Equal(t, owner.Supply(), 0)

	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	owner = NewOwner(Addrs[0], idCollection)
	require.Equal(t, owner.Supply(), 2)

	idCollection2 := NewIDCollection(Denom2, ids)
	owner = NewOwner(Addrs[0], idCollection, idCollection2)
	require.Equal(t, owner.Supply(), 4)
}

func TestOwnerGetIDCollectionMethod(t *testing.T) {
	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	owner := NewOwner(Addrs[0], idCollection)

	gotCollection, found := owner.GetIDCollection(Denom2)
	require.False(t, found)
	require.Equal(t, gotCollection.Denom, "")
	require.Equal(t, len(gotCollection.IDs), 0)
	require.Equal(t, gotCollection.String(), IDCollection{}.String())

	gotCollection, found = owner.GetIDCollection(Denom1)
	require.True(t, found)
	require.Equal(t, gotCollection.String(), idCollection.String())

	idCollection2 := NewIDCollection(Denom2, ids)
	owner = NewOwner(Addrs[0], idCollection, idCollection2)

	gotCollection, found = owner.GetIDCollection(Denom1)
	require.True(t, found)
	require.Equal(t, gotCollection.String(), idCollection.String())

	gotCollection, found = owner.GetIDCollection(Denom2)
	require.True(t, found)
	require.Equal(t, gotCollection.String(), idCollection2.String())
}

func TestOwnerUpdateIDCollectionMethod(t *testing.T) {
	ids := []string{ID1}
	idCollection := NewIDCollection(Denom1, ids)
	owner := NewOwner(Addrs[0], idCollection)
	require.Equal(t, owner.Supply(), 1)

	ids2 := []string{ID1, ID2}
	idCollection2 := NewIDCollection(Denom2, ids2)

	// UpdateIDCollection should fail if denom doesn't exist
	returnedOwner, err := owner.UpdateIDCollection(idCollection2)
	require.Error(t, err)

	idCollection3 := NewIDCollection(Denom1, ids2)
	returnedOwner, err = owner.UpdateIDCollection(idCollection3)
	require.NoError(t, err)
	require.Equal(t, returnedOwner.Supply(), 2)

	owner = returnedOwner

	returnedCollection, _ := owner.GetIDCollection(Denom1)
	require.Equal(t, len(returnedCollection.IDs), 2)

	owner = NewOwner(Addrs[0], idCollection, idCollection2)
	require.Equal(t, owner.Supply(), 3)

	returnedOwner, err = owner.UpdateIDCollection(idCollection3)
	require.NoError(t, err)
	require.Equal(t, returnedOwner.Supply(), 4)
}

func TestOwnerDeleteIDMethod(t *testing.T) {
	ids := []string{ID1, ID2}
	idCollection := NewIDCollection(Denom1, ids)
	owner := NewOwner(Addrs[0], idCollection)

	returnedOwner, err := owner.DeleteID(Denom2, ID1)
	require.Error(t, err)
	require.Equal(t, owner.String(), returnedOwner.String())

	returnedOwner, err = owner.DeleteID(Denom1, ID3)
	require.Error(t, err)
	require.Equal(t, owner.String(), returnedOwner.String())

	owner, err = owner.DeleteID(Denom1, ID1)
	require.NoError(t, err)

	returnedCollection, _ := owner.GetIDCollection(Denom1)
	require.Equal(t, len(returnedCollection.IDs), 1)
}
