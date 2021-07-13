# Go-node REST API

> go-node/x/{api}/client/rest/{query|tx}

## Base structures
```go
type BaseReq struct {
    From          string       `json:"from"`
    Memo          string       `json:"memo"`
    ChainID       string       `json:"chain_id"`
    AccountNumber uint64       `json:"account_number"`
    Sequence      uint64       `json:"sequence"`
    Fees          sdk.Coins    `json:"fees"`
    GasPrices     sdk.DecCoins `json:"gas_prices"`
    Gas           string       `json:"gas"`
    GasAdjustment string       `json:"gas_adjustment"`
    Simulate      bool         `json:"simulate"`
}
    type Coins []Coin
        type Coin struct {
            Denom string `json:"denom"`
            Amount Int   `json:"amount"`
        }
    type DecCoins []DecCoin
        type DecCoin struct {
            Denom  string `json:"denom"`
            Amount Dec    `json:"amount"`
        }
            type Dec struct {
                *big.Int `json:"int"`
            }
type AccAddress []byte
type ValAddress []byte
```

### Coins

```
GET: /coins
GET: /coin/{symbol}
```

```
POST: /coin/create
```
```go
type CoinCreateReq struct {
	BaseReq              rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title                string       `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
	ConstantReserveRatio string       `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
	Symbol               string       `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
	InitialVolume        string       `json:"initial_volume" yaml:"initial_volume"`
	InitialReserve       string       `json:"initial_reserve" yaml:"initial_reserve"`
	LimitVolume          string       `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
	Identity             string       `json:"identity" yaml:"identity"`
}
```

```
POST: /coin/send
```
```go
type CoinSendReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Coin     string       `json:"coin" yaml:"coin"`
	Amount   string       `json:"amount" yaml:"amount"`
	Receiver string       `json:"receiver" yaml:"receiver"`
}
```

```
POST: /coin/sell
```
```go
type CoinSellReq struct {
	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
	CoinToSell   string       `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToSell string       `json:"amount_to_sell" yaml:"amount_to_sell"`
	CoinToBuy    string       `json:"coin_to_buy" yaml:"coin_to_buy"`
}
```

```
POST: /coin/buy
```
```go
type CoinBuyReq struct {
	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
	CoinToSell   string       `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToSell string       `json:"amount_to_sell" yaml:"amount_to_sell"`
	CoinToBuy    string       `json:"coin_to_buy" yaml:"coin_to_buy"`
}
```

### Gov

```
GET: /gov/parameters/type
GET: /gov/proposals
GET: /gov/proposals/proposal-id
GET: /gov/proposals/proposal-id/tally
GET: /gov/proposals/proposal-id/votes
GET: /gov/proposals/proposal-id/votes/voter
```

```
POST: /gov/proposals
```
```go
type PostProposalReq struct {
	BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Title            string         `json:"title" yaml:"title"`             // Title of the proposal
	Description      string         `json:"description" yaml:"description"` // Description of the proposal
	Proposer         sdk.AccAddress `json:"proposer" yaml:"proposer"`       // Address of the proposer
	VotingStartBlock string         `json:"voting_start_block" yaml:"voting_start_block"`
	VotingEndBlock   string         `json:"voting_end_block" yaml:"voting_end_block"`
}
```

```
POST: /gov/proposals/proposal-id/votes
```
```go
type VoteReq struct {
	BaseReq rest.BaseReq   `json:"base_req" yaml:"base_req"`
	Voter   sdk.ValAddress `json:"voter" yaml:"voter"`   // address of the voter
	Option  string         `json:"option" yaml:"option"` // option from OptionSet chosen by the voter
}
```

### Multisig

```
GET: /multisig/parameters
```

### NFT

```
GET: /nft/supply/{denom}
GET: /nft/owner/{delegatorAddr}
GET: /nft/owner/{delegatorAddr}/collection/{denom}
GET: /nft/denoms
GET: /nft/collection/{denom}/nft/{id}
```

```
POST: /nfts/transfer
```
```go
type transferNFTReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Denom       string       `json:"denom"`
	ID          string       `json:"id"`
	Recipient   string       `json:"recipient"`
	SubTokenIDs []string     `json:"subTokenIDs"`
}
```

```
POST: /nfts/collection/{denom}/nft/{id}/metadata
```
```go
type editNFTMetadataReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Denom    string       `json:"denom"`
	ID       string       `json:"id"`
	TokenURI string       `json:"tokenURI"`
}
```

```
POST: /nfts/mint
```
```go
type mintNFTReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Recipient sdk.AccAddress `json:"recipient"`
	Denom     string         `json:"denom"`
	ID        string         `json:"id"`
	TokenURI  string         `json:"tokenURI"`
	Quantity  string         `json:"quantity"`
}
```

```
PUT: /nfts/collection/{denom}/nft/{id}/burn
```
```go
type burnNFTReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Denom       string       `json:"denom"`
	ID          string       `json:"id"`
	SubTokenIDs []string     `json:"subTokenIDs"`
}
```

### Validator

```
GET: /validator/delegators/{delegatorAddr}/delegations
GET: /validator/delegators/{delegatorAddr}/unbonding_delegations
GET: /validator/delegators/{delegatorAddr}/txs
GET: /validator/delegators/{delegatorAddr}/validators
GET: /validator/delegators/{delegatorAddr}/validators/{validatorAddr}
GET: /validator/delegators/{delegatorAddr}/delegations/{validatorAddr}
GET: /validator/delegators/{delegatorAddr}/unbonding_delegations/{validatorAddr}
GET: /validator/validators
GET: /validator/validators/{validatorAddr}
GET: /validator/validators/{validatorAddr}/delegations
GET: /validator/validators/{validatorAddr}/unbonding_delegations
GET: /validator/pool
GET: /validator/parameters
```

```
POST: /validator/delegators/{delegatorAddr}/delegations
```
```go
type DelegateRequest struct {
    BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
    DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"` // in bech32
    ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"` // in bech32
    Amount           sdk.Coin       `json:"amount" yaml:"amount"`
}
```

```
POST: /validator/delegators/{delegatorAddr}/unbonding_delegations
```
```go
type UndelegateRequest struct {
    BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
    DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"` // in bech32
    ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"` // in bech32
    Amount           sdk.Coin       `json:"amount" yaml:"amount"`
}
```
