package cli

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"crypto/sha256"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	"math/big"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/btcsuite/btcutil/base58"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

func GetCmdIssueCheck(cdc *codec.LegacyAmino) *cobra.Command {
	return &cobra.Command{
		Use:   "issue-check [coin] [amount] [nonce] [dueBlock] [passphrase]",
		Short: "Issue check",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd).WithLegacyAmino(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinSymbol = args[0]
			var amount, _ = sdk.NewIntFromString(args[1])
			var nonce, _ = sdk.NewIntFromString(args[2])
			var dueBlock, _ = strconv.ParseUint(args[3], 10, 64)
			var passphrase = args[4] // TODO: Read passphrase by request to avoid saving it in terminal history

			// Check if coin exists
			coin, _ := cliUtils.GetCoin(cliCtx, coinSymbol)
			if coin.Symbol != coinSymbol {
				return types2.ErrCoinDoesNotExist(coinSymbol)
			}

			// TODO: Check amount

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, _ := crypto.ToECDSA(passphraseHash[:])

			// Prepare check without lock
			check := &types2.Check{
				ChainID:  cliCtx.ChainID,
				Coin:     coin.Symbol,
				Amount:   amount.BigInt(),
				Nonce:    nonce.BigInt().Bytes(),
				DueBlock: dueBlock,
			}

			// Prepare check lock
			checkHash := check.HashWithoutLock()
			lock, _ := crypto.Sign(checkHash[:], passphrasePrivKey)

			// Fill check with prepared lock
			check.Lock = big.NewInt(0).SetBytes(lock)

			// Retrieve private key from the keybase account
			ExportPrivKeyArmor
			privKeyArmored, err := txBldr.Keybase().ExportPrivKey(cliCtx.FromName, "", "")
			keyring.Exporter().ExportPrivKeyArmor()
			if err != nil {
				msgError := fmt.Sprintf("unable to retrieve armored private key for account %s: %s", cliCtx.FromName, err.Error())
				return sdkerrors.New(types2.DefaultCodespace, sdkerrors.ErrInvalidRequest.ABCICode(), msgError)
			}
			privKey, algo, err := mintkey.UnarmorDecryptPrivKey(privKeyArmored, "")
			if err != nil {
				msgError := fmt.Sprintf("unable to retrieve private key for account %s: %s", cliCtx.FromName, err.Error())
				return sdkerrors.New(types2.DefaultCodespace, sdkerrors.ErrInvalidRequest.ABCICode(), msgError)
			}
			if algo != "secp256k1" {
				msgError := fmt.Sprintf("unable to retrieve secp256k1 private key for account %s: %s private key retrieved instead", cliCtx.FromName, algo)
				return sdkerrors.New(types2.DefaultCodespace, sdkerrors.ErrInvalidRequest.ABCICode(), msgError)
			}
			privKeySecp256k1, ok := privKey.(secp256k1.PrivKeySecp256k1)
			if !ok {
				panic("invalid private key")
			}
			privKeyECDSA, err := crypto.ToECDSA(privKeySecp256k1[:])
			if err != nil {
				panic(err)
			}
			// address := sdk.AccAddress(privKey.PubKey().Address())

			// Sign check by check issuer
			checkHash = check.Hash()
			signature, err := crypto.Sign(checkHash[:], privKeyECDSA)
			if err != nil {
				panic(err)
			}
			check.SetSignature(signature)

			// Return issued raw check encoded to base64 format to the issuer
			checkBytes, err := rlp.EncodeToBytes(check)
			if err != nil {
				panic(err)
			}
			return clientCtx.PrintObjectLegacy(struct {
				Check string
			}{
				Check: base58.Encode(checkBytes),
			})
		},
	}
}
