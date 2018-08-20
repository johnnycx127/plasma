package rpc

import (
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kyokan/plasma/chain"
	"github.com/kyokan/plasma/types"
)

type SendArgs struct {
	chain.Transaction
	From   string
	To     string
	Amount string
}

type SendResponse struct {
	Transaction *chain.Transaction
}

type TransactionService struct {
	TxChan chan<- chan types.TransactionRequest
}

func (t *TransactionService) Send(r *http.Request, args *SendArgs, reply *SendResponse) error {
	log.Println("Received Transaction.Send request.")

	from := common.HexToAddress(args.From)
	to := common.HexToAddress(args.To)
	amount := new(big.Int)
	amount.SetString(args.Amount, 0)

	req := types.TransactionRequest{
		Transaction: args.Transaction,
		From:        from,
		To:          to,
		Amount:      amount,
	}

	ch := make(chan types.TransactionRequest)
	t.TxChan <- ch
	ch <- req
	res := <-ch
	close(ch)

	if res.Response.Error != nil {
		return res.Response.Error
	}

	*reply = SendResponse{
		Transaction: res.Response.Transaction,
	}

	return nil
}
