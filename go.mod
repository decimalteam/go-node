module bitbucket.org/decimalteam/go-node

go 1.15

require (
	github.com/99designs/keyring v1.1.6
	github.com/armon/go-metrics v0.3.6 // indirect
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/btcsuite/btcutil v1.0.2
	github.com/cosmos/cosmos-sdk v0.42.4
	github.com/ethereum/go-ethereum v1.9.11
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/regen-network/cosmos-proto v0.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/iavl v0.13.2 // indirect
	github.com/tendermint/tendermint v0.34.9
	github.com/tendermint/tm-db v0.6.4
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	google.golang.org/genproto v0.0.0-20210114201628-6edceaf6022f // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
