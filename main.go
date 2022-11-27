package main

import (
	"auction/x/auction/types"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	// Importing the general purpose Cosmos blockchain client
)

func main() {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 30)

	publicKey := privateKey.PublicKey

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

	accountName := "bob"

	account, err := cosmos.Account(accountName)
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

		fmt.Println("Transactions: ", txs)

		for _, tx := range txs {
			events, _ := tx.GetEvents()
			fmt.Println(events)

		}

		var bidAmount = "10"
		var encryptBidString = EncryptString(bidAmount, publicKey)
		// broadcast keyshare message
		broadcastMsg := &types.MsgBidAuction{
			AuctionIndex: "1",
			BidAmount:    encryptBidString,
		}
		broadcastResp, err := cosmos.BroadcastTx(context.Background(), account, broadcastMsg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Sent keyshare: ", broadcastResp)
	}

}

func EncryptString(secretMessage string,
	key rsa.PublicKey) string {
	rng := rand.Reader
	ciphertext, err := rsa.EncryptPKCS1v15(rng, &key,
		[]byte(secretMessage))

	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(ciphertext)
}
