package cli

import flag "github.com/spf13/pflag"

const (
	FlagHash = "hash"
)

var (
	FsHash = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsHash.String(FlagHash, "", "HashedSecret of secret. If not specified, it will be random")
}
