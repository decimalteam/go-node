package types

import (
	"encoding/hex"
	"math/big"
	"testing"

	"bitbucket.org/decimalteam/go-node/config"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestEcrecover(t *testing.T) {
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	_r, err := hex.DecodeString("8e2f625a4b8a149a08efe71848c55031a1b7d1e8f625d60ae284336001f448ab")
	require.NoError(t, err)

	var r Hash
	copy(r[:], _r)

	_s, err := hex.DecodeString("6f6ecd87173c2d567ed9b7a4a5687b773e6a924f4122276b063186001f69e9c4")
	require.NoError(t, err)

	var s Hash
	copy(s[:], _s)

	sender, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	recipient, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	//wantAddress := SwapServiceAddress()

	amount, ok := sdk.NewIntFromString("1000000000000000000")
	require.True(t, ok)

	msg := NewMsgRedeemV2(
		sender,
		recipient,
		"0x45376AD024c767577714C7B92882578aE8B7f98C",
		amount,
		"del",
		"1625633838875",
		2,
		1,
		28,
		r,
		s)

	transactionNumber, ok := sdk.NewIntFromString(msg.TransactionNumber)
	require.True(t, ok)

	hash, err := GetHash(transactionNumber, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.FromChain, msg.DestChain)
	require.NoError(t, err)

	require.Equal(t, "a1eb252d25bb4e1ea472f8f9671789a1418ad710fe7481418eb75828fcbf5b29", hex.EncodeToString(hash[:]))

	R := big.NewInt(0)
	R.SetBytes(msg.R[:])

	S := big.NewInt(0)
	S.SetBytes(msg.S[:])

	type args struct {
		sighash [32]byte
		R       *big.Int
		S       *big.Int
		Vb      *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    ethcmn.Address
		wantErr bool
	}{
		{
			"Test1",
			args{
				sighash: hash,
				R:       R,
				S:       S,
				Vb:      sdk.NewInt(int64(msg.V)).BigInt(),
			},
			ethcmn.HexToAddress(CheckingAddress),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Ecrecover(tt.args.sighash, tt.args.R, tt.args.S, tt.args.Vb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecrecover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if hex.EncodeToString(got.Bytes()) != CheckingAddress {
				t.Errorf("Ecrecover() got = %v, want %v", hex.EncodeToString(got.Bytes()), CheckingAddress)
			}
		})
	}
}
