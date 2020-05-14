package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/bech32"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Ensure that different address types implement the interface
var _ sdk.Address = AccAddress{}
var _ sdk.Address = ValAddress{}
var _ sdk.Address = ConsAddress{}

var _ yaml.Marshaler = AccAddress{}
var _ yaml.Marshaler = ValAddress{}
var _ yaml.Marshaler = ConsAddress{}

////////////////////////////////////////////////////////////////
// Account address
////////////////////////////////////////////////////////////////

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses hex format.
type AccAddress []byte

// AccAddressFromPrefixedHex creates an AccAddress from a prefixed hex string.
func AccAddressFromPrefixedHex(address string) (AccAddress, error) {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()
	if !strings.HasPrefix(address, bech32PrefixAccAddr) {
		return nil, fmt.Errorf("expected prefix %q in prefixed account address in hex format", bech32PrefixAccAddr)
	}
	addr, err := sdk.AccAddressFromHex(address[len(bech32PrefixAccAddr):])
	return AccAddress(addr), err
}

// AccAddressFromHex creates an AccAddress from a hex string.
func AccAddressFromHex(address string) (AccAddress, error) {
	addr, err := sdk.AccAddressFromHex(address)
	return AccAddress(addr), err
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromBech32(address string) (AccAddress, error) {
	addr, err := sdk.AccAddressFromBech32(address)
	return AccAddress(addr), err
}

// VerifyAddressFormat verifies that the provided bytes form a valid address
// according to the default address rules or a custom address verifier set by
// GetConfig().SetAddressVerifier().
func VerifyAddressFormat(bz []byte) error {
	return sdk.VerifyAddressFormat(bz)
}

// Equals returns boolean for whether two AccAddresses are Equal.
func (aa AccAddress) Equals(aa2 sdk.Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Empty returns boolean for whether an AccAddress is empty.
func (aa AccAddress) Empty() bool {
	if aa == nil {
		return true
	}

	aa2 := AccAddress{}
	return bytes.Equal(aa.Bytes(), aa2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf compatibility.
func (aa AccAddress) Marshal() ([]byte, error) {
	return aa, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf compatibility.
func (aa *AccAddress) Unmarshal(data []byte) error {
	*aa = data
	return nil
}

// MarshalJSON marshals to JSON using hex format.
func (aa AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using hex format.
func (aa AccAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming hex format encoding.
func (aa *AccAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AccAddressFromPrefixedHex(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming hex format encoding.
func (aa *AccAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	aa2, err := AccAddressFromPrefixedHex(s)
	if err != nil {
		return err
	}

	*aa = aa2
	return nil
}

// Bytes returns the raw address bytes.
func (aa AccAddress) Bytes() []byte {
	return aa
}

// String implements the Stringer interface.
func (aa AccAddress) String() string {
	if aa.Empty() {
		return ""
	}

	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()
	hexAddr := bech32PrefixAccAddr + hex.EncodeToString(aa.Bytes())
	return hexAddr
}

// Format implements the fmt.Formatter interface.
func (aa AccAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

////////////////////////////////////////////////////////////////
// Validator operator address
////////////////////////////////////////////////////////////////

// ValAddress defines a wrapper around bytes meant to present a validator's
// operator. When marshaled to a string or JSON, it uses hex format.
type ValAddress []byte

// ValAddressFromPrefixedHex creates an ValAddress from a prefixed hex string.
func ValAddressFromPrefixedHex(address string) (ValAddress, error) {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()
	if !strings.HasPrefix(address, bech32PrefixValAddr) {
		return nil, fmt.Errorf("expected prefix %q in prefixed validator address in hex format", bech32PrefixValAddr)
	}
	addr, err := sdk.ValAddressFromHex(address[len(bech32PrefixValAddr):])
	return ValAddress(addr), err
}

// ValAddressFromHex creates a ValAddress from a hex string.
func ValAddressFromHex(address string) (ValAddress, error) {
	addr, err := sdk.ValAddressFromHex(address)
	return ValAddress(addr), err
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address string) (ValAddress, error) {
	addr, err := sdk.ValAddressFromBech32(address)
	return ValAddress(addr), err
}

// Equals returns boolean for whether two ValAddresses are Equal.
func (va ValAddress) Equals(va2 sdk.Address) bool {
	if va.Empty() && va2.Empty() {
		return true
	}

	return bytes.Equal(va.Bytes(), va2.Bytes())
}

// Empty returns boolean for whether an AccAddress is empty.
func (va ValAddress) Empty() bool {
	if va == nil {
		return true
	}

	va2 := ValAddress{}
	return bytes.Equal(va.Bytes(), va2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf compatibility.
func (va ValAddress) Marshal() ([]byte, error) {
	return va, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf compatibility.
func (va *ValAddress) Unmarshal(data []byte) error {
	*va = data
	return nil
}

// MarshalJSON marshals to JSON using hex format.
func (va ValAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(va.String())
}

// MarshalYAML marshals to YAML using hex format.
func (va ValAddress) MarshalYAML() (interface{}, error) {
	return va.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming hex format encoding.
func (va *ValAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	va2, err := ValAddressFromPrefixedHex(s)
	if err != nil {
		return err
	}

	*va = va2
	return nil
}

// UnmarshalYAML unmarshals from YAML assuming hex format encoding.
func (va *ValAddress) UnmarshalYAML(data []byte) error {
	var s string

	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	va2, err := ValAddressFromPrefixedHex(s)
	if err != nil {
		return err
	}

	*va = va2
	return nil
}

// Bytes returns the raw address bytes.
func (va ValAddress) Bytes() []byte {
	return va
}

// String implements the Stringer interface.
func (va ValAddress) String() string {
	if va.Empty() {
		return ""
	}

	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()
	hexAddr := bech32PrefixValAddr + hex.EncodeToString(va.Bytes())
	return hexAddr
}

// Format implements the fmt.Formatter interface.
func (va ValAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(va.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", va)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(va))))
	}
}

////////////////////////////////////////////////////////////////
// Consensus node address
////////////////////////////////////////////////////////////////

// ConsAddress defines a wrapper around bytes meant to present a consensus node.
// When marshaled to a string or JSON, it uses hex format.
type ConsAddress []byte

// ConsAddressFromPrefixedHex creates an ConsAddress from a prefixed hex string.
func ConsAddressFromPrefixedHex(address string) (ConsAddress, error) {
	bech32PrefixConsAddr := sdk.GetConfig().GetBech32ConsensusAddrPrefix()
	if !strings.HasPrefix(address, bech32PrefixConsAddr) {
		return nil, fmt.Errorf("expected prefix %q in prefixed consensus address in hex format", bech32PrefixConsAddr)
	}
	addr, err := sdk.ConsAddressFromHex(address[len(bech32PrefixConsAddr):])
	return ConsAddress(addr), err
}

// ConsAddressFromHex creates a ConsAddress from a hex string.
func ConsAddressFromHex(address string) (ConsAddress, error) {
	addr, err := sdk.ConsAddressFromHex(address)
	return ConsAddress(addr), err
}

// ConsAddressFromBech32 creates a ConsAddress from a Bech32 string.
func ConsAddressFromBech32(address string) (ConsAddress, error) {
	addr, err := sdk.ConsAddressFromBech32(address)
	return ConsAddress(addr), err
}

// GetConsAddress gets ConsAddress from pubkey.
func GetConsAddress(pubkey crypto.PubKey) ConsAddress {
	return ConsAddress(pubkey.Address())
}

// Equals returns boolean for whether two ConsAddress are Equal.
func (ca ConsAddress) Equals(ca2 sdk.Address) bool {
	if ca.Empty() && ca2.Empty() {
		return true
	}

	return bytes.Equal(ca.Bytes(), ca2.Bytes())
}

// Empty returns boolean for whether an ConsAddress is empty.
func (ca ConsAddress) Empty() bool {
	if ca == nil {
		return true
	}

	ca2 := ConsAddress{}
	return bytes.Equal(ca.Bytes(), ca2.Bytes())
}

// Marshal returns the raw address bytes. It is needed for protobuf compatibility.
func (ca ConsAddress) Marshal() ([]byte, error) {
	return ca, nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf compatibility.
func (ca *ConsAddress) Unmarshal(data []byte) error {
	*ca = data
	return nil
}

// MarshalJSON marshals to JSON using hex format.
func (ca ConsAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.String())
}

// MarshalYAML marshals to YAML using hex format.
func (ca ConsAddress) MarshalYAML() (interface{}, error) {
	return ca.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming hex format encoding.
func (ca *ConsAddress) UnmarshalJSON(data []byte) error {
	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	ca2, err := ConsAddressFromPrefixedHex(s)
	if err != nil {
		return err
	}

	*ca = ca2
	return nil
}

// UnmarshalYAML unmarshals from YAML assuming hex format encoding.
func (ca *ConsAddress) UnmarshalYAML(data []byte) error {
	var s string

	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	ca2, err := ConsAddressFromPrefixedHex(s)
	if err != nil {
		return err
	}

	*ca = ca2
	return nil
}

// Bytes returns the raw address bytes.
func (ca ConsAddress) Bytes() []byte {
	return ca
}

// String implements the Stringer interface.
func (ca ConsAddress) String() string {
	if ca.Empty() {
		return ""
	}

	bech32PrefixConsAddr := sdk.GetConfig().GetBech32ConsensusAddrPrefix()
	bech32Addr := bech32PrefixConsAddr + hex.EncodeToString(ca.Bytes())
	return bech32Addr
}

// Format implements the fmt.Formatter interface.
func (ca ConsAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(ca.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", ca)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(ca))))
	}
}

////////////////////////////////////////////////////////////////
// Public keys
////////////////////////////////////////////////////////////////

// Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

// Bech32 conversion constants
const (
	Bech32PubKeyTypeAccPub  Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub  Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub Bech32PubKeyType = "conspub"
)

// HexifyPubKey returns a prefixed hex format string containing the appropriate
// prefix based on the key type provided for a given PublicKey.
func HexifyPubKey(pkt Bech32PubKeyType, pubkey crypto.PubKey) (string, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = sdk.GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = sdk.GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = sdk.GetConfig().GetBech32ConsensusPubPrefix()

	}

	return bech32Prefix + hex.EncodeToString(pubkey.Bytes()), nil
}

// MustHexifyPubKey calls HexifyPubKey except it panics on error.
func MustHexifyPubKey(pkt Bech32PubKeyType, pubkey crypto.PubKey) string {
	res, err := HexifyPubKey(pkt, pubkey)
	if err != nil {
		panic(err)
	}

	return res
}

// GetPubKeyFromPrefixedHex returns a PublicKey from a prefixed hex format PublicKey with a given key type.
func GetPubKeyFromPrefixedHex(pkt Bech32PubKeyType, pubkeyStr string) (crypto.PubKey, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = sdk.GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = sdk.GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = sdk.GetConfig().GetBech32ConsensusPubPrefix()

	}

	if !strings.HasPrefix(pubkeyStr, bech32Prefix) {
		return nil, fmt.Errorf("expected prefix %q in prefixed public key in hex format", bech32Prefix)
	}

	bz, err := hex.DecodeString(pubkeyStr[len(bech32Prefix):])
	if err != nil {
		return nil, err
	}

	pk, err := tmamino.PubKeyFromBytes(bz)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// MustGetPubKeyFromPrefixedHex calls GetPubKeyFromPrefixedHex except it panics on error.
func MustGetPubKeyFromPrefixedHex(pkt Bech32PubKeyType, pubkeyStr string) crypto.PubKey {
	res, err := GetPubKeyFromPrefixedHex(pkt, pubkeyStr)
	if err != nil {
		panic(err)
	}

	return res
}

// Bech32ifyPubKey returns a Bech32 encoded string containing the appropriate
// prefix based on the key type provided for a given PublicKey.
func Bech32ifyPubKey(pkt Bech32PubKeyType, pubkey crypto.PubKey) (string, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = sdk.GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = sdk.GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = sdk.GetConfig().GetBech32ConsensusPubPrefix()

	}

	return bech32.ConvertAndEncode(bech32Prefix, pubkey.Bytes())
}

// MustBech32ifyPubKey calls Bech32ifyPubKey except it panics on error.
func MustBech32ifyPubKey(pkt Bech32PubKeyType, pubkey crypto.PubKey) string {
	res, err := Bech32ifyPubKey(pkt, pubkey)
	if err != nil {
		panic(err)
	}

	return res
}

// GetPubKeyFromBech32 returns a PublicKey from a bech32-encoded PublicKey with a given key type.
func GetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) (crypto.PubKey, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = sdk.GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = sdk.GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = sdk.GetConfig().GetBech32ConsensusPubPrefix()

	}

	bz, err := GetFromBech32(pubkeyStr, bech32Prefix)
	if err != nil {
		return nil, err
	}

	pk, err := tmamino.PubKeyFromBytes(bz)
	if err != nil {
		return nil, err
	}

	return pk, nil
}

// MustGetPubKeyFromBech32 calls GetPubKeyFromBech32 except it panics on error.
func MustGetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) crypto.PubKey {
	res, err := GetPubKeyFromBech32(pkt, pubkeyStr)
	if err != nil {
		panic(err)
	}

	return res
}

// GetFromBech32 decodes a bytestring from a Bech32 encoded string.
func GetFromBech32(bech32str, prefix string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
	}

	return bz, nil
}
