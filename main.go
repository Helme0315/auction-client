package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	// Importing the general purpose Cosmos blockchain client
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
)

func main() {
	// Prefix to use for account addresses.
	// The address prefix was assigned to the auction blockchain
	// using the `--address-prefix` flag during scaffolding.
	addressPrefix := "auction"

	// Create a Cosmos client instance
	cosmos, err := cosmosclient.New(
		context.Background(),
		cosmosclient.WithAddressPrefix(addressPrefix),
	)

	if err != nil {
		log.Fatal(err)
	}
	for {
		if cosmos.WaitForNextBlock(context.Background()) != nil {
			continue
		}

		latestHeight, err := cosmos.LatestBlockHeight(context.Background())
		if err != nil {
			continue
		}

		fmt.Println("Last Block Number: ", latestHeight)

		txs, err := cosmos.GetBlockTXs(context.Background(), latestHeight)
		if err != nil {
			continue
		}

		for _, tx := range txs {
			events, _ := tx.GetEvents()
			fmt.Println(events)

		}

		fmt.Println(strconv.Itoa(int(latestHeight)))
		resp, err := http.Get("http://0.0.0.0:26657/block_results?height=" + strconv.Itoa(int(latestHeight)))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(resp)
	}

}
