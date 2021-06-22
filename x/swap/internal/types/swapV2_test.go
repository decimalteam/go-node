package types

import (
	"bitbucket.org/decimalteam/go-node/config"
	"encoding/hex"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestEcrecover(t *testing.T) {
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	_r, err := hex.DecodeString("b8b3eb4980e649a65b7e136fbcafda4d12e3b11a40d8aaa7d951e13fbe483579")
	require.NoError(t, err)

	var r Hash
	copy(r[:], _r)

	_s, err := hex.DecodeString("74de77f4a9f4045992cf6f220cff9be67a2c0332124e60af0a6791c9b0a64c36")
	require.NoError(t, err)

	var s Hash
	copy(s[:], _s)

	sender, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	recipient, err := sdk.AccAddressFromBech32("dx1twj64nphm8zl98uxv7gnt6xg4tpkk4gyr3tux9")
	require.NoError(t, err)

	//wantAddress := SwapServiceAddress()

	msg := NewMsgRedeemV2(
		sender,
		recipient,
		"0x856F08B12cB844fa05CDF1eBfFd303B091D34d09",
		sdk.NewInt(2000000000000000000),
		"muh coin",
		"coin",
		"qqqqqqqq",
		2,
		1,
		28,
		r,
		s)

	hash, err := GetHash(msg.TransactionNumber, msg.TokenName, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.FromChain, msg.DestChain)
	require.NoError(t, err)

	require.Equal(t, "d90ed147ca8100c8329314b74466e1b2f154eeeb26bdfcd9af84f68901f9bf4c", hex.EncodeToString(hash[:]))

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
