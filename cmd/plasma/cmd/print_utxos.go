package cmd

import "github.com/spf13/cobra"

var printUTXOsCmd = &cobra.Command{
	Use: "print-utxos",
	Short: "prints UTXOs for a given address",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printUTXOsCmd)
	printUTXOsCmd.Flags().String(FlagAddress, "", "the address to print UTXOs for")
	printUTXOsCmd.MarkFlagRequired(FlagAddress)
}