package keys

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
)

const (
	flagUnarmoredHex = "unarmored-hex"
	flagUnsafe       = "unsafe"
)

// ExportKeyCommand exports private keys from the key store.
func ExportKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export <name>",
		Short: "Export private keys",
		Long:  `Export a private key from the local keybase in ASCII-armored encrypted format.`,
		Args:  cobra.ExactArgs(1),
		RunE:  runExportCmd,
	}
	cmd.Flags().Bool(flagUnarmoredHex, false, "Export unarmored hex privkey. Requires --unsafe.")
	cmd.Flags().Bool(flagUnsafe, false, "Enable unsafe operations. This flag must be switched on along with all unsafe operation-specific options.")

	return cmd
}

func runExportCmd(cmd *cobra.Command, args []string) error {
	buf := bufio.NewReader(cmd.InOrStdin())
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}
	kb := clientCtx.Keyring

	unarmored, _ := cmd.Flags().GetBool(flagUnarmoredHex)
	unsafe, _ := cmd.Flags().GetBool(flagUnsafe)

	if unarmored && unsafe {
		return exportUnsafeUnarmored(cmd, args[0], buf, clientCtx.Keyring)
	} else if unarmored || unsafe {
		return fmt.Errorf("the flags %s and %s must be used together", flagUnsafe, flagUnarmoredHex)
	}

	encryptPassword, err := input.GetPassword("Enter passphrase to encrypt the exported key:", buf)
	if err != nil {
		return err
	}

	armored, err := kb.ExportPrivKeyArmor(args[0], encryptPassword)
	if err != nil {
		return err
	}

	cmd.Println(armored)
	return nil
}



func exportUnsafeUnarmored(cmd *cobra.Command, uid string, buf *bufio.Reader, kr keyring.Keyring) error {
	// confirm deletion, unless -y is passed
	if yes, err := input.GetConfirmation("WARNING: The private key will be exported as an unarmored hexadecimal string. USE AT YOUR OWN RISK. Continue?", buf, cmd.ErrOrStderr()); err != nil {
		return err
	} else if !yes {
		return nil
	}

	hexPrivKey, err := keyring.NewUnsafe(kr).UnsafeExportPrivKeyHex(uid)
	if err != nil {
		return err
	}

	cmd.Println(hexPrivKey)

	return nil
}