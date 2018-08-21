package cmd

import (
	"github.com/kyokan/plasma/config"
	"github.com/spf13/viper"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewGlobalConfig() *config.GlobalConfig {
	return &config.GlobalConfig{
		DBPath:       viper.GetString(FlagDB),
		NodeURL:      viper.GetString(FlagNodeURL),
		RPCPort:      viper.GetInt(FlagRPCPort),
		ContractAddr: viper.GetString(FlagContractAddr),
	}
}

func ParsePrivateKey() (*ecdsa.PrivateKey, error) {
	privateKeyStr := viper.GetString(FlagPrivateKey)
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}