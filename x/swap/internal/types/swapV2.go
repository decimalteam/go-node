package types

import (
	"crypto/sha256"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/crypto/ripemd160"
	"math/big"
)

type SwapV2 struct {
	TransferType      TransferType   `json:"transfer_type"`
	From              sdk.AccAddress `json:"from"`
	Recipient         string         `json:"recipient"`
	Amount            sdk.Coin       `json:"amount"`
	TokenName         string         `json:"token_name"`
	Timestamp         uint64         `json:"timestamp"`
	DestChain         string         `json:"dest_chain"`
	TransactionNumber string         `json:"transaction_number"`
	V                 *big.Int       `json:"v"`
	R                 *big.Int       `json:"r"`
	S                 *big.Int       `json:"s"`
}

func (s SwapV2) Hash() Hash {
	var hash [32]byte
	copy(hash[:], crypto.Keccak256(
		[]byte(s.TransactionNumber),
		abi.U256(s.Amount.Amount.BigInt()),
		[]byte(s.TokenName),
		[]byte(s.Amount.Denom),
		[]byte(s.Recipient),
		[]byte(s.DestChain),
	))
	return hash
}

func GetHash(transactionNumber, tokenName, tokenSymbol string, amount sdk.Int, recipient sdk.AccAddress, destChain int) (Hash, error) {
	var hash [32]byte

	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	stringTy, _ := abi.NewType("string", "string", nil)

	arguments := abi.Arguments{
		{
			Type: stringTy, // transactionNumber
		},
		{
			Type: uint256Ty, // amount
		},
		{
			Type: stringTy, // tokenName
		},
		{
			Type: stringTy, // tokenSymbol
		},
		{
			Type: stringTy, // recipient
		},
		{
			Type: uint256Ty, // destChain
		},
	}

	bytes, err := arguments.Pack(
		transactionNumber,
		amount,
		tokenName,
		tokenSymbol,
		recipient,
		sdk.NewInt(int64(destChain)),
	)
	if err != nil {
		return [32]byte{}, err
	}

	copy(hash[:], crypto.Keccak256(bytes))

	return hash, nil
}

func Ecrecover(sighash [32]byte, R, S, Vb *big.Int) (sdk.AccAddress, error) {
	if Vb.BitLen() > 8 {
		return sdk.AccAddress{}, errors.New("invalid sig")
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, true) {
		return sdk.AccAddress{}, errors.New("invalid sig")
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return sdk.AccAddress{}, errors.New("invalid public key")
	}
	pub2, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	pub3 := crypto.CompressPubkey(pub2)
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(pub3)
	sha := hasherSHA256.Sum(nil)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha)
	return hasherRIPEMD160.Sum(nil), nil
}
