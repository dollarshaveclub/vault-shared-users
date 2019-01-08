package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/dollarshaveclub/vault-shared-users/lib/vsu"
	"github.com/spf13/cobra"
)

var (
	debug         bool          // flag
	login         bool          // flag
	retries       int           // flag
	testTokenPath string        // flag
	timeout       time.Duration // flag
	tokenPath     string        // flag
	vaultAddr     string        // flag

	config *vsu.Config
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debugging output")
	rootCmd.PersistentFlags().BoolVar(&login, "login", false, "attempt to login to Vault for the user")
	rootCmd.PersistentFlags().IntVar(&retries, "retries", -1, "number of retries")
	rootCmd.PersistentFlags().StringVar(&testTokenPath, "test-token-path", "", "path to find a GitHub token for login")
	rootCmd.PersistentFlags().DurationVar(&timeout, "vault-timeout", -1, "timeout for vault requests in milliseconds")
	rootCmd.PersistentFlags().StringVar(&tokenPath, "token-path", "", "path to find a valid Vault token")
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", os.Getenv("VAULT_ADDR"), "address to access vault and should be a full URL")
}

var rootCmd = &cobra.Command{
	Use:   "vault-shared-access",
	Short: "Vault shared access allows us to store credentials in Vault for third-party sites",
	Long:  `Vault shared access allows us to store credentials in Vault for third-party sites, including 2FA codes`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config = &vsu.Config{
			Login:         login,
			Retries:       retries,
			TestTokenPath: testTokenPath,
			Timeout:       timeout,
			TokenPath:     tokenPath,
			VaultAddr:     vaultAddr,
		}
	},
}

// Execute is the entrypoint for running the different commands of psst
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cliError(message string, err error) {
	fmt.Fprintln(os.Stderr, message)
	if debug {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
