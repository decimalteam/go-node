package types

import (
	"bytes"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func GetHash(transactionNumber, tokenName, tokenSymbol string, amount sdk.Int, recipient sdk.AccAddress, fromChain int, destChain int) (Hash, error) {
	var hash [32]byte

	encoded := encodePacked(
		encodeString(transactionNumber),
		encodeUint256(amount.BigInt()),
		encodeString(tokenName),
		encodeString(tokenSymbol),
		encodeString(recipient.String()),
		encodeUint8(uint8(fromChain)),
		encodeUint8(uint8(destChain)),
	)

	copy(hash[:], crypto.Keccak256(encoded))

	return hash, nil
}

func encodePacked(input ...[]byte) []byte {
	return bytes.Join(input, nil)
}

func encodeString(v string) []byte {
	return []byte(v)
}

func encodeUint256(v *big.Int) []byte {
	return abi.U256(v)
}

func encodeUint8(v uint8) []byte {
	return new(big.Int).SetUint64(uint64(v)).Bytes()
}

func Ecrecover(sighash [32]byte, R, S, Vb *big.Int) (ethcmn.Address, error) {
	if Vb.BitLen() > 8 {
		return ethcmn.Address{}, errors.New("invalid sig")
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, true) {
		return ethcmn.Address{}, errors.New("invalid sig")
	}
	// encode the signature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return ethcmn.Address{}, err
	}

	if len(pub) == 0 || pub[0] != 4 {
		return ethcmn.Address{}, errors.New("invalid public key")
	}

	var addr ethcmn.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])

	return addr, nil
}
