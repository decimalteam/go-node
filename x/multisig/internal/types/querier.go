package types

import "strings"

// QueryWallets specifies type containing set of multisig wallets.
type QueryWallets []Wallet

// String implements fmt.Stringer interface.
func (n QueryWallets) String() string {
	wallets := make([]string, len(n))
	for i, wallet := range n {
		wallets[i] = wallet.String()
	}
	return strings.Join(wallets[:], "\n")
}

// QueryTransactions specifies type containing set of multisig transactions.
type QueryTransactions []Transaction

// String implements fmt.Stringer interface.
func (n QueryTransactions) String() string {
	transactions := make([]string, len(n))
	for i, transaction := range n {
		transactions[i] = transaction.String()
	}
	return strings.Join(transactions[:], "\n")
}
