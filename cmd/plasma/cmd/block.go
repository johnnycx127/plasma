package cmd

import "github.com/spf13/cobra"

var blockCmd = &cobra.Command{
	Use: "block",
	Short: "outputs transactions in a block",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(blockCmd)
	blockCmd.Flags().String(FlagRootHost, "", "hostname and port of the root node")
	blockCmd.Flags().Uint(FlagBlockNum, 0, "the block number to show")
	blockCmd.MarkFlagRequired(FlagRootHost)
	blockCmd.MarkFlagRequired(FlagBlockNum)
}