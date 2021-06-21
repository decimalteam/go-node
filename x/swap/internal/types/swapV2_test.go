package types

import (
	"bitbucket.org/decimalteam/go-node/config"
	"encoding/hex"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
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

	_r, err := hex.DecodeString("d9c49d2418e3badd092685cb9e5088519b088531c790e1b32d4fac4a7b4a57dc")
	require.NoError(t, err)

	var r Hash
	copy(r[:], _r)

	_s, err := hex.DecodeString("4b3a51a828aa92890fbe79008b44f5034a4cab9604d1e51351b6146620e812e8")
	require.NoError(t, err)

	var s Hash
	copy(s[:], _s)

	sender, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	recipient, err := sdk.AccAddressFromBech32("dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v")
	require.NoError(t, err)

	//wantAddress := SwapServiceAddress()

	msg := NewMsgRedeemV2(
		sender,
		recipient,
		"0x45376AD024c767577714C7B92882578aE8B7f98C",
		sdk.NewInt(1000000000000000000),
		"decimal",
		"del",
		"lksdnd-asvkla-SDCds",
		2,
		1,
		27,
		r,
		s)

	hash, err := GetHash(msg.TransactionNumber, msg.TokenName, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.FromChain, msg.DestChain)
	require.NoError(t, err)

	require.Equal(t, "83d14ee00f224fa84aa885da8c25a81a6b20c87b8d316f2c3060f3a57d3c436c", hex.EncodeToString(hash[:]))

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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ecrecover() got = %v, want %v", got, tt.want)
			}
		})
	}
}
