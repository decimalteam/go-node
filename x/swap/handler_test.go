package swap

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"
)

func Test_getHash(t *testing.T) {
	type args struct {
		secret [32]byte
	}

	var secret [32]byte
	d, _ := hex.DecodeString("927c1ac33100bdbb001de19c626a05a7c3c11304fc825f5eabb22e741507711b")
	copy(secret[:], d)

	var want [32]byte
	w, _ := hex.DecodeString("5efc1cee17257e9f9e34827bce9f827fa305bc39e80fd77081e4b7780f2b0ca7")
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
