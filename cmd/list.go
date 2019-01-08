package cmd

import (
	"fmt"

	"github.com/dollarshaveclub/vault-shared-users/lib/vsu"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available credentials",
	Long:  `List available credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := vsu.List(config)
		if err != nil {
			cliError("error trying to list avaialble credentials", err)
		}

		if len(creds) > 0 {
			fmt.Println("Available Credentials")
			fmt.Println("=====================")
			fmt.Println()

			for _, name := range creds {
				fmt.Printf("- %s\n", name)
			}
		} else {
			fmt.Println("No avaiable secrets")
		}
	},
}
