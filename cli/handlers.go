package cli

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kyokan/plasma/chain"
	"github.com/kyokan/plasma/eth"
	"context"
	"github.com/kyokan/plasma/db"
	"github.com/kyokan/plasma/config"
	"crypto/ecdsa"
	"path"
	"google.golang.org/grpc"
	"github.com/kyokan/plasma/rpc/pb"
	"github.com/kyokan/plasma/rpc"
	"github.com/pkg/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

func initHandler(config *config.GlobalConfig, privateKey *ecdsa.PrivateKey, pCtx context.Context) (context.CancelFunc, eth.Client, db.PlasmaStorage, error) {
	ctx, cancel := context.WithCancel(pCtx)
	client, err := eth.NewClient(config.NodeURL, config.ContractAddr, privateKey)
	if err != nil {
		return cancel, nil, nil, err
	}

	ldb, storage, err := db.CreateStorage(path.Join(config.DBPath, "root"), client)
	if err != nil {
		return cancel, nil, nil, err
	}

	go func() {
		<-ctx.Done()
		ldb.Close()
	}()

	return cancel, client, storage, nil
}

func Finalize(config *config.GlobalConfig, privateKey *ecdsa.PrivateKey) error {
	cancel, client, _, err := initHandler(config, privateKey, context.Background())
	defer cancel()
	if err != nil {
		return err
	}

	return client.Finalize()
}

func Exit(config *config.GlobalConfig, privateKey *ecdsa.PrivateKey, rootHost string, blockNum *big.Int, txIndex uint, oIndex uint) error {
	ctx := context.Background()
	cancel, client, _, err := initHandler(config, privateKey, ctx)
	defer cancel()
	if err != nil {
		return err
	}

	fmt.Printf("Exit starting for blocknum: %d, txIndex: %d, oIndex: %d\n", blockNum, txIndex, oIndex)

	conn, err := grpc.Dial(fmt.Sprintf("http://%s/rpc", rootHost))
	if err != nil {
		return err
	}
	defer conn.Close()

	rc := pb.NewRootClient(conn)
	res, err := rc.GetBlock(ctx, &pb.GetBlockRequest{
		Number: rpc.SerializeBig(blockNum),
	})
	if err != nil {
		return err
	}
	if res == nil {
		return errors.New("block does not exist")
	}

	return client.StartExit(&eth.StartExitOpts{
		Block:    rpc.DeserializeBlock(res.Block),
		Txs:      rpc.DeserializeTxs(res.Transactions),
		BlockNum: blockNum,
		TxIndex:  txIndex,
		OutIndex: oIndex,
	})
}

func Deposit(config *config.GlobalConfig, privateKey *ecdsa.PrivateKey, amount *big.Int) error {
	ctx := context.Background()
	cancel, client, _, err := initHandler(config, privateKey, ctx)
	defer cancel()
	if err != nil {
		return err
	}

	fmt.Printf("Deposit starting for amount: %s\n", amount.Text(10))
	userAddress := crypto.PubkeyToAddress((privateKey.Public()).(ecdsa.PublicKey))
	t := createDepositTx(userAddress, amount)
	client.Deposit(amount, &t)
	time.Sleep(3 * time.Second)
	curr, err := client.CurrentChildBlock()
	if err != nil {
		return err
	}

	fmt.Printf("Last child block: %v\n", curr)
	return nil
}

// TODO: Use same code as transaction sink.
func createDepositTx(userAddress common.Address, value *big.Int) chain.Transaction {
	return chain.Transaction{
		Input0: chain.ZeroInput(),
		Input1: chain.ZeroInput(),
		Output0: &chain.Output{
			NewOwner: userAddress,
			Amount:   value,
		},
		Output1: chain.ZeroOutput(),
		Fee:     big.NewInt(0),
	}
}
