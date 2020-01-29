package types

import "strings"

// Query endpoints supported by the coin querier
const (
	QueryListCoins = "list"
	QueryGetCoin   = "get"
)

type QueryResCoins []string

func (n QueryResCoins) String() string {
	return strings.Join(n[:], "\n")
}
