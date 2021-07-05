package types

import (
	"fmt"

	"golang.org/x/crypto/sha3"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MultisigTransactionIDPrefix is prefix for multisig transaction ID.
const MultisigTransactionIDPrefix = "dxmstx"

////////////////////////////////////////////////////////////////
// Multisig Wallet
////////////////////////////////////////////////////////////////

// Wallet is a struct that contains all the metadata of a multi-signature wallet.
//type Wallet struct {
//	Address   sdk.AccAddress   `json:"address" yaml:"address"`
//	Owners    []sdk.AccAddress `json:"owners" yaml:"owners"`
//	Weights   []uint64           `json:"weights" yaml:"weights"`
//	Threshold uint64             `json:"threshold" yaml:"threshold"`
//}

// NewWallet returns a new Wallet.
func NewWallet(owners []string, weights []uint64, threshold uint64, salt []byte) (*Wallet, error) {
	walletMetadata := struct {
		Owners    []string `json:"owners" yaml:"owners"`
		Weights   []uint64         `json:"weights" yaml:"weights"`
		Threshold uint64           `json:"threshold" yaml:"threshold"`
		Salt      []byte           `json:"salt" yaml:"salt"`
	}{
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
		Salt:      salt,
	}
	bz := sha3.Sum256(sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(walletMetadata)))
	address := sdk.AccAddress(bz[12:])

	return &Wallet{
		Address:   address.String(),
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
	}, nil
}

// String implements fmt.Stringer interface.
func (w *Wallet) String() string {
	weightsSum := uint64(0)
	for _, weight := range w.Weights {
		weightsSum += weight
	}
	return fmt.Sprintf(`Wallet %s: (%d of %d)`, w.Address, w.Threshold, weightsSum)
}

////////////////////////////////////////////////////////////////
// Multisig Transaction
////////////////////////////////////////////////////////////////

// Transaction is a struct that contains all the metadata of a multi-signature wallet transaction.
//type Transaction struct {
//	ID        string           `json:"id" yaml:"id"`
//	Wallet    sdk.AccAddress   `json:"wallet" yaml:"wallet"`
//	Receiver  sdk.AccAddress   `json:"receiver" yaml:"receiver"`
//	Coins     sdk.Coins        `json:"coins" yaml:"coins"`
//	Signers   []sdk.AccAddress `json:"signers" yaml:"signers"`
//	CreatedAt int64            `json:"created_at" yaml:"created_at"` // block height
//}

// NewTransaction returns a new Transaction.
func NewTransaction(wallet, receiver string, coins sdk.Coins, signers []string, height int64, salt []byte) (*Transaction, error) {

	transactionMetadata := struct {
		Wallet    string   `json:"wallet" yaml:"wallet"`
		Receiver  string   `json:"receiver" yaml:"receiver"`
		Coins     sdk.Coins        `json:"coins" yaml:"coins"`
		Signers   []string `json:"signers" yaml:"signers"`
		CreatedAt int64            `json:"created_at" yaml:"created_at"` // block height
		Salt      []byte           `json:"salt" yaml:"salt"`
	}{
		Wallet:    wallet,
		Receiver:  receiver,
		Coins:     coins,
		Signers:   signers,
		CreatedAt: height,
		Salt:      salt,
	}
	bz := sha3.Sum256(sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(transactionMetadata)))
	id, err := bech322.ConvertAndEncode(MultisigTransactionIDPrefix, bz[12:])
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:        id,
		Wallet:    wallet,
		Receiver:  receiver,
		Coins:     coins,
		Signers:   signers,
		CreatedAt: height,
	}, nil
}

// String implements fmt.Stringer interface.
func (t *Transaction) String() string {
	return fmt.Sprintf("Transaction %s: %s --> %s %+v", t.ID, t.Wallet, t.Receiver, t.Coins)
}

func toStrAddrs(owners []sdk.AccAddress) []string {
	strowners := make([]string, len(owners))
	for _, o := range owners {
		strowners = append(strowners, o.String())
	}

	return strowners
}
