package types

import (
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"testing"
)

func TestGetRewardForBlock(t *testing.T) {
	type args struct {
		blockHeight uint64
	}
	tests := []struct {
		name string
		args args
		want sdk.Int
	}{
		{
			"0 block",
			args{
				0,
			},
			helpers.BipToPip(sdk.NewInt(50)),
		},
		{
			"1 month",
			args{
				432000,
			},
			helpers.BipToPip(sdk.NewInt(55)),
		},
		{
			"11 month",
			args{
				432000 * 11,
			},
			helpers.BipToPip(sdk.NewInt(105)),
		},
		{
			"2th year",
			args{
				432000 * 12,
			},
			helpers.BipToPip(sdk.NewInt(122)),
		},
		{
			"last block of 12th month",
			args{
				432000*24 - 1,
			},
			helpers.BipToPip(sdk.NewInt(309)),
		},
		{
			"3th year",
			args{
				5184000 * 2,
			},
			helpers.BipToPip(sdk.NewInt(338)),
		},
		{
			"last month",
			args{
				5184000*8 + 432000*11,
			},
			helpers.BipToPip(sdk.NewInt(5769)),
		},
		{
			"10th year",
			args{
				5184000 * 9,
			},
			helpers.BipToPip(sdk.NewInt(0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRewardForBlock(tt.args.blockHeight); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRewardForBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
