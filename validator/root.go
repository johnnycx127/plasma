package validator

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	encoding_json "encoding/json"

	"github.com/kyokan/plasma/chain"
	"github.com/kyokan/plasma/db"
	"github.com/kyokan/plasma/eth"
	"github.com/kyokan/plasma/userclient"
	"github.com/kyokan/plasma/util"
)

// clientResponse represents a JSON-RPC response returned to a client.
type ClientResponse struct {
	Result *encoding_json.RawMessage `json:"result"`
	Error  interface{}               `json:"error"`
	Id     uint64                    `json:"id"`
}

func RootNodeListener(rootUrl string, storage db.PlasmaStorage, client eth.Client, userAddress string) {
	rootClient := userclient.NewRootClient(rootUrl)
	for {
		log.Println("Watching root node...")

		block, err := storage.LatestBlock()

		if err != nil {
			log.Fatalf("Failed to get latest block: %v", err)
		}

		var blockNum uint64

		if block == nil { // first block in the plasma chain is genesis
			blockNum = 1
		} else {
			blockNum = block.Header.Number + 1
		}

		log.Printf("Looking for block number: %d\n", blockNum)

		response := rootClient.GetBlock(blockNum)

		if response != nil {
			log.Printf("Found block number: %d\n", blockNum)
			plasmaBlock := response.Block

			// Block number for the contract is off by one
			contractBlock := plasma.GetBlock(util.NewUint64(blockNum))

			if IsValidBlock(plasmaBlock, contractBlock) {
				log.Println("Block is valid, saving locally.")
				storage.SaveBlock(response.Block)
			} else {
				_, err := storage.GetInvalidBlock(response.Block.BlockHash)

				if err != nil && err.Error() == "leveldb: not found" {
					log.Println("Block is not valid, starting exit of utxos.")
					// Exit all utxos, because we suspect root node is dishonest
					ExitUTXOs(rootUrl, plasma, userAddress)

					log.Println("Saving invalid block...")

					storage.SaveInvalidBlock(response.Block)
				} else {
					log.Println("We already invalidated this block")
				}
			}
		}

		// Need to wait longer, because we need to wait for block to be submitted.
		time.Sleep(10 * time.Second)
	}
}

func ExitUTXOs(rootUrl string, plasma *eth.PlasmaClient, userAddress string) {
	rootClient := userclient.NewRootClient(rootUrl)
	res := rootClient.GetUTXOs(userAddress)
	txs := res.Transactions

	type UTXO struct {
		BlkNum    uint64
		TxIdx     uint32
		OutputIdx int
	}

	utxosByBlock := make(map[uint64][]UTXO)

	for _, tx := range txs {
		utxos := utxosByBlock[tx.BlkNum]

		if utxos == nil {
			utxos = []UTXO{}
		}

		// Collect a list of outputs because technically both can belong to the user.
		var outputIdxs []int

		if tx.Output0.NewOwner.String() == userAddress {
			outputIdxs = append(outputIdxs, 0)
		} else if tx.Output1.NewOwner.String() == userAddress {
			outputIdxs = append(outputIdxs, 1)
		} else {
			log.Fatalf("Transaction must have at least one output that belongs to address: %s\n", userAddress)
		}

		for _, outputIdx := range outputIdxs {
			utxos = append(
				utxos,
				UTXO{
					tx.BlkNum,
					tx.TxIdx,
					outputIdx,
				},
			)
		}

		utxosByBlock[tx.BlkNum] = utxos
	}

	for blkNum, utxos := range utxosByBlock {
		res2 := rootClient.GetBlock(blkNum)

		for _, utxo := range utxos {
			log.Printf("Exiting block: %d, tx: %d, output: %d\n", utxo.BlkNum, utxo.TxIdx, utxo.OutputIdx)

			plasma.StartExit(
				res2.Block,
				res2.Transactions,
				util.NewUint64(utxo.BlkNum),
				util.NewUint32(utxo.TxIdx),
				util.NewInt(utxo.OutputIdx),
			)

			time.Sleep(3 * time.Second)
		}
	}
}

func IsValidBlock(block *chain.Block, plasmaBlock eth.Block) bool {
	fmt.Println(block.Header.Number)
	fmt.Println(hex.EncodeToString(block.Header.RLPMerkleRoot))
	fmt.Println(hex.EncodeToString(plasmaBlock.Root))
	return bytes.Equal(block.Header.RLPMerkleRoot, plasmaBlock.Root)
}
