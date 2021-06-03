package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"golang.org/x/crypto/sha3"

	"github.com/spf13/cobra"

	"github.com/btcsuite/btcutil/base58"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func GetCmdRedeemCheck(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "redeem-check [check] [passphrase]",
		Short: "Redeem check",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)

			var checkBase58 = args[0]
			var passphrase = args[1] // TODO: Read passphrase by request to avoid saving it in terminal history

			// Decode provided check from base58 format to raw bytes
			checkBytes := base58.Decode(checkBase58)
			if len(checkBytes) == 0 {
				msgError := "unable to decode check from base58"
				return sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, msgError)
			}

			// Parse provided check from raw bytes to ensure it is valid
			_, err := types2.ParseCheck(checkBytes)
			if err != nil {
				msgError := fmt.Sprintf("unable to parse check: %s", err.Error())
				return sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, msgError)
			}

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, err := crypto.ToECDSA(passphraseHash[:])
			if err != nil {
				msgError := fmt.Sprintf("unable to create private key from passphrase: %s", err.Error())
				return sdkerrors.New(types2.DefaultCodespace, types2.InvalidPassphrase, msgError)
			}

			// Prepare bytes to sign by private key generated from passphrase
			receiverAddressHash := make([]byte, 32)
			hw := sha3.NewLegacyKeccak256()
			err = rlp.Encode(hw, []interface{}{
				clientCtx.FromAddress,
			})
			if err != nil {
				msgError := fmt.Sprintf("unable to RLP encode check receiver address: %s", err.Error())
				return sdkerrors.New(types2.DefaultCodespace, types2.InvalidPassphrase, msgError)
			}
			hw.Sum(receiverAddressHash[:0])

			// Sign receiver address by private key generated from passphrase
			signature, err := crypto.Sign(receiverAddressHash[:], passphrasePrivKey)
			if err != nil {
				msgError := fmt.Sprintf("unable to sign check receiver address by private key generated from passphrase: %s", err.Error())
				return sdkerrors.New(types2.DefaultCodespace, types2.InvalidPassphrase, msgError)
			}
			proofBase64 := base64.StdEncoding.EncodeToString(signature)

			// Prepare redeem check message
			msg := types2.NewMsgRedeemCheck(clientCtx.FromAddress, checkBase58, proofBase64)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), []sdk.Msg{&msg}...)
		},
	}
}
