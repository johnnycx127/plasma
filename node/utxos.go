package node

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kyokan/plasma/db"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/urfave/cli.v1"
)

// TODO: migrate to root userclient.
func PrintUTXOs(c *cli.Context) {
	db, storage, err := db.CreateStorage(c.GlobalString("db"), nil)

	if err != nil {
		log.Panic("Failed to establish connection with database:", err)
	}

	defer db.Close()

	addrStr := c.String("addr")

	if addrStr == "" {
		log.Panic("Addr is required.")
	}

	addr := common.HexToAddress(c.String("addr"))
	txs, err := storage.UTXOs(&addr)

	if err != nil {
		log.Panic("Failed to get UTXOs: ", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Hash", "Amount", "Block Number", "Tx Index"})
	for _, tx := range txs {
		table.Append([]string{
			common.ToHex(tx.Hash()),
			tx.OutputFor(&addr).Amount.String(),
			fmt.Sprint(tx.BlkNum),
			fmt.Sprint(tx.TxIdx),
		})
	}

	table.Render()
}
