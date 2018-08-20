package cmd

import "github.com/spf13/cobra"

var depositCmd = &cobra.Command{
	Use: "deposit",
	Short: "performs a deposit",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(depositCmd)
	depositCmd.Flags().String(FlagAmount, "", "the amount to deposit")
	depositCmd.MarkFlagRequired(FlagAmount)
}