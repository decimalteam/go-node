package types

import (
	"fmt"

	"golang.org/x/crypto/sha3"

	"github.com/google/uuid"

	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MultisigTransactionIDPrefix = "dxmstx"

////////////////////////////////////////////////////////////////
// Multisig Wallet
////////////////////////////////////////////////////////////////

// Wallet is a struct that contains all the metadata of a multi-signature wallet.
type Wallet struct {
	Address   sdk.AccAddress   `json:"address" yaml:"address"`
	Owners    []sdk.AccAddress `json:"owners" yaml:"owners"`
	Weights   []uint           `json:"weights" yaml:"weights"`
	Threshold uint             `json:"threshold" yaml:"threshold"`
}

// NewWallet returns a new Wallet.
func NewWallet(owners []sdk.AccAddress, weights []uint, threshold uint) (*Wallet, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	bz, err := uid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	hz := sha3.Sum256(bz)
	address := sdk.AccAddress(hz[12:])

	return &Wallet{
		Address:   address,
		Owners:    owners,
		Weights:   weights,
		Threshold: threshold,
	}, nil
}

// String implements fmt.Stringer interface.
func (w *Wallet) String() string {
	weightsSum := uint(0)
	for _, weight := range w.Weights {
		weightsSum += weight
	}
	return fmt.Sprintf(`Wallet %s: (%d of %d)`, w.Address, w.Threshold, weightsSum)
}

////////////////////////////////////////////////////////////////
// Multisig Transaction
////////////////////////////////////////////////////////////////

// Transaction is a struct that contains all the metadata of a multi-signature wallet transaction.
type Transaction struct {
	ID        string           `json:"id" yaml:"id"`
	Wallet    sdk.AccAddress   `json:"wallet" yaml:"wallet"`
	Receiver  sdk.AccAddress   `json:"receiver" yaml:"receiver"`
	Coins     sdk.Coins        `json:"coins" yaml:"coins"`
	Signers   []sdk.AccAddress `json:"signers" yaml:"signers"`
	CreatedAt int64            `json:"created_at" yaml:"created_at"` // block height
}

// NewTransaction returns a new Transaction.
func NewTransaction(wallet, receiver sdk.AccAddress, coins sdk.Coins, signers []sdk.AccAddress, height int64) (*Transaction, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	bz, err := uid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	hz := sha3.Sum256(bz)

	id, err := bech32.ConvertAndEncode(MultisigTransactionIDPrefix, hz[12:])
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
