package cmd

import "github.com/spf13/cobra"

var balanceCmd = &cobra.Command{
	Use: "balance",
	Short: "shows an address's balance",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.Flags().String(FlagAddress, "", "the address to show balances for")
	depositCmd.MarkFlagRequired(FlagAddress)
}