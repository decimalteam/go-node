package types

const (
	QuerySwap        = "swap"
	QueryActiveSwaps = "active_swaps"
	QueryPool        = "pool"
)

type QuerySwapParams struct {
	HashedSecret Hash `json:"hashed_secret"`
}

func NewQuerySwapParams(hashedSecret Hash) QuerySwapParams {
	return QuerySwapParams{HashedSecret: hashedSecret}
}
