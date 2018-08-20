package validator

import (
	"fmt"
	"log"
	"path"

	"github.com/kyokan/plasma/db"
	"github.com/kyokan/plasma/eth"
	"github.com/kyokan/plasma/config"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
		)

func Start(config *config.GlobalConfig, privateKey *ecdsa.PrivateKey, rootHost string) error {
	log.Println("Validator Starting")

	userAddress := crypto.PubkeyToAddress((privateKey.Public()).(ecdsa.PublicKey)).Hex()
	rootUrl := fmt.Sprintf("http://%s/rpc", rootHost)
	plasma, err := eth.NewClient(config.NodeURL, config.ContractAddr, privateKey)
	if err != nil {
		return err
	}

	ldb, storage, err := db.CreateStorage(path.Join(config.DBPath, "validator", userAddress), plasma)
	if err != nil {
		return err
	}
	defer ldb.Close()

	go RootNodeListener(rootUrl, storage, plasma, userAddress)
	go ExitStartedListener(rootUrl, storage, plasma)
	go Run(config.RPCPort)

	select {}
	return nil
}
