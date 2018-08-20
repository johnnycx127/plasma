package cmd

import "github.com/spf13/cobra"

const (
	FlagTxIndex  = "txindex"
	FlagOIndex   = "oindex"
)

var exitCmd = &cobra.Command{
	Use:   "exit",
	Short: "performs an exit",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exitCmd)
	exitCmd.Flags().String(FlagRootHost, "", "hostname and port of the root node")
	exitCmd.Flags().Uint(FlagBlockNum, 0, "block number to exit")
	exitCmd.Flags().Uint(FlagTxIndex, 0, "transaction to exit")
	exitCmd.Flags().Uint(FlagOIndex, 0, "output to exit")
	exitCmd.MarkFlagRequired(FlagRootHost)
	exitCmd.MarkFlagRequired(FlagBlockNum)
	exitCmd.MarkFlagRequired(FlagTxIndex)
	exitCmd.MarkFlagRequired(FlagOIndex)
}
