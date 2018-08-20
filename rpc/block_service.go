package rpc

import (
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"

	"github.com/kyokan/plasma/chain"
	"github.com/kyokan/plasma/db"
)

type GetBlocksArgs struct {
	Height uint64
}

type GetBlocksResponse struct {
	Block        *chain.Block        `json:"Block"`
	Transactions []chain.Transaction `json:"Transactions"`
}

type GetUTXOsArgs struct {
	UserAddress string
}

type GetUTXOsResponse struct {
	Transactions []chain.Transaction
}


type GetBalanceArgs struct {
	UserAddress string
}

type GetBalanceResponse struct {
	Amount *big.Int
}

type BlockService struct {
	Storage db.PlasmaStorage
}

func (t *BlockService) GetBlock(r *http.Request, args *GetBlocksArgs, reply *GetBlocksResponse) error {
	log.Println("Received Block.GetBlocks request.")

	height := args.Height

	block, err := t.Storage.BlockAtHeight(height)

	if err != nil {
		return err
	}

	txs, err := t.Storage.FindTransactionsByBlockNum(height)

	if err != nil {
		return err
	}

	*reply = GetBlocksResponse{
		Block:        block,
		Transactions: txs,
	}

	return nil
}

func (t *BlockService) GetUTXOs(r *http.Request, args *GetUTXOsArgs, reply *GetUTXOsResponse) error {
	log.Println("Received Block.GetUTXOs request.")

	userAddress := common.HexToAddress(args.UserAddress)

	txs, err := t.Storage.UTXOs(&userAddress)

	if err != nil {
		return err
	}

	*reply = GetUTXOsResponse{
		Transactions: txs,
	}

	return nil
}

func (t *BlockService) GetBalance(r *http.Request, args *GetBalanceArgs, reply *GetBalanceResponse) error {
	log.Println("Received Block.Balance request.")

	userAddress := common.HexToAddress(args.UserAddress)

	amount, err := t.Storage.Balance(&userAddress)

	if err != nil {
		return err
	}

	*reply = GetBalanceResponse{
		Amount: amount,
	}

	return nil
}
