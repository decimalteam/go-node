package keeper_test

/*
func TestSetCollection(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// create a new nft with id = "id" and owner = "address"
	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(
		id,
		address,
		address,
		tokenURI,
		sdk.NewInt(1),
		sdk.NewInt(101),
		false,
	)

	err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// collection should exist
	collection, exists := NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)

	// create a new NFT and add it to the collection created with the NFT mint
	nft2 := types.NewBaseNFT(
		id2,
		address,
		address,
		tokenURI,
		sdk.NewInt(1),
		sdk.NewInt(101),
		true,
	)
	collection2, err2 := collection.AddNFT(nft2)
	require.NoError(t, err2)
	NFTKeeper.SetCollection(ctx, denom, collection2)

	collection2, exists = NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)
	require.Len(t, collection2.NFTs, 2)

	// reset collection for invariant sanity
	NFTKeeper.SetCollection(ctx, denom, collection)

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
func TestGetCollection(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// collection shouldn't exist
	collection, exists := NFTKeeper.GetCollection(ctx, denom)
	require.Empty(t, collection)
	require.False(t, exists)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(2), sdk.NewInt(101), true)
	err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// collection should exist
	collection, exists = NFTKeeper.GetCollection(ctx, denom)
	require.True(t, exists)
	require.NotEmpty(t, collection)

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
func TestGetCollections(t *testing.T) {
	ctx, _, NFTKeeper := createTestApp(t, false)

	// collections should be empty
	collections := NFTKeeper.GetCollections(ctx)
	require.Empty(t, collections)

	// MintNFT shouldn't fail when collection does not exist
	nft := types.NewBaseNFT(id, address, address, tokenURI, sdk.NewInt(1), sdk.NewInt(101), true)
	err := NFTKeeper.MintNFT(ctx, denom, nft)
	require.NoError(t, err)

	// collections should equal 1
	collections = NFTKeeper.GetCollections(ctx)
	require.NotEmpty(t, collections)
	require.Equal(t, len(collections), 1)

	msg, fail := keeper.SupplyInvariant(NFTKeeper)(ctx)
	require.False(t, fail, msg)
}
*/
