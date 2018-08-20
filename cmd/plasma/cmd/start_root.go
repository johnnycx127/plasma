package cmd

import (
	"github.com/spf13/cobra"
		"github.com/spf13/viper"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kyokan/plasma/root"
)

var startRootCmd = &cobra.Command{
	Use: "start-root",
	Short: "starts running a Plasma root node",
	RunE: func(cmd *cobra.Command, args []string) error {
		privateKeyStr := viper.GetString(FlagPrivateKey)
		privateKey, err := crypto.HexToECDSA(privateKeyStr)
		if err != nil {
			return err
		}

		return root.Start(NewGlobalConfig(), privateKey)
	},
}

func init() {
	rootCmd.AddCommand(startRootCmd)
	startRootCmd.Flags().Uint(FlagRPCPort, 8643, "port for the RPC server to listen on")
	viper.BindPFlag(FlagRPCPort, startRootCmd.Flags().Lookup(FlagRPCPort))
}