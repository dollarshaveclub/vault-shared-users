package cmd

import (
	"fmt"

	"github.com/dollarshaveclub/vault-shared-users/lib/vsu"
	"github.com/spf13/cobra"
)

var (
	otopOnly      = false
	recoveryCodes = false
)

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&otopOnly, "2fa-only", false, "only return a 2fa code for the user")
	getCmd.Flags().BoolVar(&recoveryCodes, "recovery-codes", false, "return the stored recovery codes")
}

var getCmd = &cobra.Command{
	Use:   "get [flags] <name>",
	Short: "Get a set of credentials for the provided user",
	Long:  `Get a set of credentials for the provided user`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		secrets, err := vsu.Get(config, name)
		if err != nil {
			cliError("unable to get secret"+name, err)
		}

		output(secrets)
	},
}

func output(secrets map[string]string) {
	otherKeys := []string{}
	// get user defined key
	for name := range secrets {
		if !vsu.IsDefaultKey(name) {
			otherKeys = append(otherKeys, name)
		}
	}

	if !otopOnly {
		fmt.Printf("Username: %s\n", secrets[vsu.UsernameKey])
		fmt.Printf("Password: %s\n", secrets[vsu.PasswordKey])
	}

	fmt.Printf("TOTP Code: %s\n", secrets[vsu.TotpKey])

	if !otopOnly {
		for i := range otherKeys {
			key := otherKeys[i]
			fmt.Printf("%s: %s\n", key, secrets[key])
		}
	}

	if recoveryCodes {
		fmt.Printf("Recovery Codes: %s\n", secrets[vsu.RecoveryCodesKey])
	}
}
