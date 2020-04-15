package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rlp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/config"
)

// TODO: Move to sdk.Int (now use big.Int, because NewSdkIntFromBigInt has limit for bits length)
type CheckData struct {
	Nonce    []byte
	ChainID  string
	DueBlock uint64
	Coin     [10]byte
	Value    *sdk.Int
	GasCoin  [10]byte
	Lock     *big.Int
	V        *big.Int
	R        *big.Int
	S        *big.Int
}

func rlpHash(x interface{}) (h common.Hash, err error) {
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, x)
	if err != nil {
		return common.Hash{}, err
	}
	hw.Sum(h[:0])
	return h, nil
}

func NewSdkIntFromBytes(bytes []byte) *sdk.Int {
	newInt := sdk.Int{}
	return &newInt
}

func (check *CheckData) Sender() (sdk.AccAddress, error) {
	pub, err := check.PublicKey()
	if err != nil {
		return sdk.AccAddress{}, err
	}

	return sdk.AccAddressFromBech32(pub)
}

func (check *CheckData) String() string {
	sender, err := check.Sender()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Check sender: %s nonce: %x, dueBlock: %d, value: %s %s", sender, check.Nonce,
		check.DueBlock, check.Value.String(), string(bytes.Trim(check.Coin[:], "\x00")))
}

func (check *CheckData) PublicKey() (string, error) {

	if check.V.BitLen() > 8 {
		return "", sdkerrors.New(DefaultCodespace, InvalidVRS, "Invalid V, R, S values")
	}

	v := byte(check.V.Uint64() - 27)
	if !crypto.ValidateSignatureValues(v, check.R, check.S, true) {
		return "", sdkerrors.New(DefaultCodespace, InvalidVRS, "Invalid V, R, S values")
	}

	r := check.R.Bytes()
	s := check.S.Bytes()

	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = v

	hash, err := rlpHash([]interface{}{
		check.Nonce,
		check.ChainID,
		check.DueBlock,
		check.Coin,
		check.Value,
		check.GasCoin,
		check.Lock,
	})
	if err != nil {
		return "", err
	}

	pub, err := secp256k1.RecoverPubkey(hash[:], sig)
	if err != nil {
		return "", err
	}

	if len(pub) == 0 || pub[0] != 4 {
		return "", sdkerrors.New(DefaultCodespace, InvalidPublicKey, "Invalid public key")
	}

	return fmt.Sprintf("%s%s", config.DecimalPrefixAccAddr, hex.EncodeToString(pub)), nil
}

type Signed interface {
	Encode() (string, error)
}

type CheckInterface interface {
	SetPassphrase(passphrase string) CheckInterface
	Sign(prKey string) (Signed, error)
}

type Check struct {
	*CheckData
	passphrase string
}

// Create Check
// Nonce - unique "id" of the check. Coin Symbol - symbol of coin. Value - amount of coins.
// Due Block - defines last block height in which the check can be used.
func NewCheck(nonce uint64, chainID string, dueBlock uint64, coin string, value *sdk.Int, gasCoin string) CheckInterface {
	check := &Check{
		CheckData: &CheckData{
			Nonce:    []byte(strconv.Itoa(int(nonce))),
			ChainID:  chainID,
			DueBlock: dueBlock,
			Value:    value,
		},
	}
	copy(check.Coin[:], coin)
	copy(check.GasCoin[:], gasCoin)
	return check
}

// Prepare check string and convert to data
func DecodeCheck(rawCheck string) (*CheckData, error) {
	decode, err := base64.StdEncoding.DecodeString(rawCheck)
	if err != nil {
		panic(err)
	}

	res := new(CheckData)
	if err := rlp.DecodeBytes(decode, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Set secret phrase which you will pass to receiver of the check
func (check *Check) SetPassphrase(passphrase string) CheckInterface {
	check.passphrase = passphrase
	return check
}

//
func (check *Check) Encode() (string, error) {
	src, err := rlp.EncodeToBytes(check.CheckData)
	if err != nil {
		return "", err
	}

	return config.DecimalCheckPrefix + hex.EncodeToString(src), err
}

// Sign Check
func (check *Check) Sign(prKey string) (Signed, error) {
	msgHash, err := rlpHash([]interface{}{
		check.Nonce,
		check.ChainID,
		check.DueBlock,
		check.Coin,
		check.Value,
		check.GasCoin,
	})
	if err != nil {
		return nil, err
	}

	passphraseSum256 := sha256.Sum256([]byte(check.passphrase))

	lock, err := secp256k1.Sign(msgHash[:], passphraseSum256[:])
	if err != nil {
		return nil, err
	}
	check.Lock = big.NewInt(0).SetBytes(lock)

	msgHashWithLock, err := rlpHash([]interface{}{
		check.Nonce,
		check.ChainID,
		check.DueBlock,
		check.Coin,
		check.Value,
		check.GasCoin,
		check.Lock,
	})
	if err != nil {
		return nil, err
	}
	privateKey := crypto.ToECDSAUnsafe([]byte(prKey))

	sig, err := crypto.Sign(msgHashWithLock[:], privateKey)
	if err != nil {
		return nil, err
	}

	check.R = new(big.Int).SetBytes(sig[:32])
	check.S = new(big.Int).SetBytes(sig[32:64])
	check.V = new(big.Int).SetBytes([]byte{sig[64] + 27})

	return check, nil
}

type CheckAddress struct {
	address    [20]byte
	passphrase string
}

//func NewCheckAddress(address string, passphrase string) (*CheckAddress, error) {
//	toHex, err := wallet.AddressToHex(address)
//	if err != nil {
//		return nil, err
//	}
//
//	check := &CheckAddress{passphrase: passphrase}
//	copy(check.address[:], toHex)
//
//	return check, nil
//}

// Proof Check
func (check *CheckAddress) Proof() (string, error) {

	passphraseSum256 := sha256.Sum256([]byte(check.passphrase))

	addressHash, err := rlpHash([]interface{}{
		check.address[:],
	})
	if err != nil {
		return "", err
	}

	lock, err := secp256k1.Sign(addressHash[:], passphraseSum256[:])
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(lock), nil
}
