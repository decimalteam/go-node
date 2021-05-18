package cli

import (
	"crypto/sha256"
	"math/big"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/btcsuite/btcutil/base58"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

const (
	ALGO_SECP256K1 string = "secp256k1"
)

func GetCmdIssueCheck(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "issue-check [coin] [amount] [nonce] [dueBlock] [passphrase]",
		Short: "Issue check",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(cliCtx.Input).WithTxEncoder(utils.GetTxEncoder(cdc))

			var coinSymbol = args[0]
			var amount, _ = sdk.NewIntFromString(args[1])
			var nonce, _ = sdk.NewIntFromString(args[2])
			var dueBlock, _ = strconv.ParseUint(args[3], 10, 64)
			var passphrase = args[4] // TODO: Read passphrase by request to avoid saving it in terminal history

			// Check if coin exists
			coin, _ := cliUtils.GetCoin(cliCtx, coinSymbol)
			if coin.Symbol != coinSymbol {
				return types.ErrCoinDoesNotExist(coinSymbol)
			}

			// TODO: Check amount

			// Prepare private key from passphrase
			passphraseHash := sha256.Sum256([]byte(passphrase))
			passphrasePrivKey, _ := crypto.ToECDSA(passphraseHash[:])

			// Prepare check without lock
			check := &types.Check{
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
			privKeyArmored, err := txBldr.Keybase().ExportPrivKey(cliCtx.FromName, "", "")
			if err != nil {
				return types.ErrUnableRetriveArmoredPkey(cliCtx.FromName, err.Error())
			}
			privKey, algo, err := mintkey.UnarmorDecryptPrivKey(privKeyArmored, "")
			if err != nil {
				return types.ErrUnableRetrivePkey(cliCtx.FromName, err.Error())
			}
			if algo != ALGO_SECP256K1 {
				return types.ErrUnableRetriveSECPPkey(cliCtx.FromName, algo)
			}
			privKeySecp256k1, ok := privKey.(secp256k1.PrivKeySecp256k1)
			if !ok {
				return types.ErrInvalidPkey()
			}
			privKeyECDSA, err := crypto.ToECDSA(privKeySecp256k1[:])
			if err != nil {
				return types.ErrInternal(err.Error())
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
			return cliCtx.PrintOutput(struct {
				Check string
			}{
				Check: base58.Encode(checkBytes),
			})
		},
	}
}
