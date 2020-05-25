package coin

import (
	"reflect"
	"testing"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Test_handleMsgBuyCoin(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		msg types.MsgBuyCoin
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleMsgBuyCoin(tt.args.ctx, tt.args.k, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMsgBuyCoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgBuyCoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgCreateCoin(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		msg types.MsgCreateCoin
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.Result
		wantErr bool
	}{
		{
			name: "TEST1",
			args: args{
				ctx: sdk.Context{},
				k:   Keeper{},
				msg: types.MsgCreateCoin{},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleMsgCreateCoin(tt.args.ctx, tt.args.k, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMsgCreateCoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgCreateCoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgMultiSendCoin(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		msg types.MsgMultiSendCoin
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleMsgMultiSendCoin(tt.args.ctx, tt.args.k, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMsgMultiSendCoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgMultiSendCoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgSellCoin(t *testing.T) {
	type args struct {
		ctx     sdk.Context
		k       Keeper
		msg     types.MsgSellCoin
		sellAll bool
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleMsgSellCoin(tt.args.ctx, tt.args.k, tt.args.msg, tt.args.sellAll)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMsgSellCoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgSellCoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMsgSendCoin(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
		msg types.MsgSendCoin
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleMsgSendCoin(tt.args.ctx, tt.args.k, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMsgSendCoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMsgSendCoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}
