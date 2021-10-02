package swap

import (
	"encoding/hex"
	"fmt"
	"reflect"
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
