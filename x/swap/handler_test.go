package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func Test_getHash(t *testing.T) {
	type args struct {
		secret []byte
	}

	var secret []byte
	d, _ := hex.DecodeString("927c1ac33100bdbb001de19c626a05a7c3c11304fc825f5eabb22e741507711b")
	copy(secret[:], d)

	var want [32]byte
	w, _ := hex.DecodeString("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	copy(want[:], w)

	tests := []struct {
		name string
		args args
		want [32]byte
	}{
		{
			name: "1",
			args: args{
				secret: secret,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHash(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				fmt.Println(hex.EncodeToString(got[:]))
				t.Errorf("getHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleMsgHTLT(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	h := NewHandler(swapKeeper)

	msghtlt := types.NewMsgHTLT(
		keeper.Addrs[0],
		"",
		skd.NewInt(2),
		"",
		"",
		"",
		2,
	)

	res, err := h(ctx, msghtlt)
	require.NoError(t, err)

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case types.AttributeKeyDestChain:
				require.Equal(t, value, strconv.Itoa(msghtlt.DestChain))
			case types.AttributeKeyFrom:
				require.Equal(t, value, msghtlt.From.String())
			case types.AttributeKeyRecipient:
				require.Equal(t, value, msghtlt.Recipient)
			case types.AttributeKeyAmount:
				require.Equal(t, value, msghtlt.Amount.String())
			case types.AttributeKeyTransactionNumber:
				require.Equal(t, value, msghtlt.TransactionNumber)
			case types.AttributeKeyTokenName:
				require.Equal(t, value, msghtlt)
			case types.AttributeKeyTokenSymbol:
				require.Equal(t, value, msghtlt.TokenSymbol)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}

	swapKeeper.GetDestChainName()
}

func TestHandleMsgRedeemV2(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	h := NewHandler(swapKeeper)

	msgredeemv2 := types.NewMsgRedeemV2(
		keeper.Addrs[0],
		keeper.Addrs[1],
		keeper.Addrs[3].String(),
		sdk.NewInt(4),
		keeper.TestTokenName,
		keeper.TestTokenSymbol,
		"1",
		1,
		10,
		uint8(1),
		[32]byte{},
		[32]byte{},
	)

	// handleMsgRedeemV2 should not fail when swap is not redeemed
	res, err := h(ctx, msgredeemv2)
	require.NoError(t, err)

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case types.AttributeKeyDestChain:
				require.Equal(t, value, strconv.Itoa(msgredeemv2.DestChain))
			case types.AttributeKeyFrom:
				require.Equal(t, value, msgredeemv2.From.String())
			case types.AttributeKeyRecipient:
				require.Equal(t, value, msgredeemv2.Recipient)
			case types.AttributeKeyAmount:
				require.Equal(t, value, msgredeemv2.Amount.String())
			case types.AttributeKeyTransactionNumber:
				require.Equal(t, value, msgredeemv2.TransactionNumber)
			case types.AttributeKeyTokenName:
				require.Equal(t, value, msgredeemv2)
			case types.AttributeKeyTokenSymbol:
				require.Equal(t, value, msgredeemv2.TokenSymbol)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
	// handleMsgRedeemV2 should fail when swap is redeemed
	hash, err := types.GetHash(msgredeemv2.TransactionNumber, msgredeemv2.TokenName, msgredeemv2.TokenSymbol, msgredeemv2.Amount, msgredeemv2.Recipient, msgredeemv2.DestChain)
	swapKeeper.SetSwapV2(ctx, hash)

	res, err = h(ctx, msgredeemv2)
	require.Error(t, err)
}

func TestHandleMsgSwapInitialize(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	h := NewHandler(swapKeeper)

	fromChain := 1
	destChain := 10

	msgswapinitialize := types.NewMsgSwapInitialize(
		keeper.Addrs[0],
		keeper.Addrs[1].String(),
		sdk.NewInt(4),
		keeper.TestTokenName,
		keeper.TestTokenSymbol,
		"1",
		fromChain,
		destChain,
	)

	// TestHandleMsgSwapInitialize should fail if DestChain does not exist
	_, err := h(ctx, msgswapinitialize)
	require.Error(t, err)

	chain1 := types.NewChain("chain1", false)

	swapKeeper.SetChain(ctx, destChain, chain1)

	// TestHandleMsgSwapInitialize should fail if FromChain does not exist
	_, err = h(ctx, msgswapinitialize)
	require.Error(t, err)

	chain2 := types.NewChain("chain2", false)

	swapKeeper.SetChain(ctx, fromChain, chain2)

	res, err := h(ctx, msgswapinitialize)
	require.NoError(t, err)

	for _, event := range res.Events {
		for _, attribute := range event.Attributes {
			value := string(attribute.Value)
			switch key := string(attribute.Key); key {
			case types.AttributeKeyDestChain:
				require.Equal(t, value, strconv.Itoa(msgswapinitialize.DestChain))
			case types.AttributeKeyFrom:
				require.Equal(t, value, msgswapinitialize.From.String())
			case types.AttributeKeyRecipient:
				require.Equal(t, value, msgswapinitialize.Recipient)
			case types.AttributeKeyAmount:
				require.Equal(t, value, msgswapinitialize.Amount.String())
			case types.AttributeKeyTransactionNumber:
				require.Equal(t, value, msgswapinitialize.TransactionNumber)
			case types.AttributeKeyTokenName:
				require.Equal(t, value, msgswapinitialize)
			case types.AttributeKeyTokenSymbol:
				require.Equal(t, value, msgswapinitialize.TokenSymbol)
			default:
				require.Fail(t, fmt.Sprintf("unrecognized event %s", key))
			}
		}
	}
}

func TestMsgChainActivate(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	h := NewHandler(swapKeeper)

	chainNumber := 1
	msgchainactivate := types.NewMsgChainActivate(
		keeper.Addr[0],
		chainNumber,
		"chain1",
	)

	// handleMsgSwapInitialize should create chain if it does not exist
	_, err := h(ctx, msgchainactivate)
	require.NoError(t, err)

	chain := types.NewChain(msgchainactivate.ChainName, true)
	// Chain should exist ane be active
	ch, found := swapKeeper.GetChain(ctx, chainNumber)
	require.True(t, found)
	require.True(t, reflect.DeepEqual(ch, chain))

	// handleMsgSwapInitialize should set chain to active
	ch.Active = false
	swapKeeper.SetChain(ctx, chainNumber, ch)

	_ = h(ctx, msgchainactivate)

	ch, found = swapKeeper.GetChain(ctx, chainNumber)
	require.True(t, found)
	require.True(t, ch.Active)
}

func TestMsgChainDeactivate(t *testing.T) {
	ctx, swapKeeper := keeper.CreateTestInput(t, false)

	h := NewHandler(swapKeeper)

	chainNumber := 1
	msgchaindeactivate := types.NewMsgChainDeactivate(
		keeper.Addr[0],
		chainNumber,
	)

	// handleMsgChainDeactivate should fail if chain does not exist
	_, err := h(ctx, msgchaindeactivate)
	require.Error(t, err)

	chain := types.NewChain("chain1", true)
	swapKeeper.SetChain(ctx, chainNumber, chain)

	_, err = h(ctx, msgchaindeactivate)
	require.NoError(t, err)

	// chain should be deactivated
	ch, found := swapKeeper.GetChain(ctx, chainNumber)
	require.True(t, found)
	require.False(t, ch.Active)
}
