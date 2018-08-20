package cmd

import (
	"github.com/spf13/viper"
	"github.com/kyokan/plasma/config"
	)

func NewGlobalConfig() *config.GlobalConfig {
	return &config.GlobalConfig{
		DBPath:       viper.GetString(FlagDB),
		NodeURL:      viper.GetString(FlagNodeURL),
		RPCPort:      viper.GetInt(FlagRPCPort),
		ContractAddr: viper.GetString(FlagContractAddr),
	}
}
