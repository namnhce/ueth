package wallets

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/namnhce/ueth/pkg/services"
)

func DoSend(ctx services.CLIContext) error {
	privateKey := ctx.String("private-key")
	if privateKey == "" {
		return errors.New("must provide private key")
	}

	value := ctx.Float64("value")
	if value == 0 {
		return errors.New("must provide ETH value")
	}

	inputFileName := ctx.String("input")
	if inputFileName == "" {
		return errors.New("must provide csv filepath")
	}

	err := send(privateKey, value, inputFileName)
	if err != nil {
		return err
	}

	fmt.Println("Transactions sent successfully!")
	return nil
}

func send(privKey string, inValue float64, csvFilePath string) error {
	// Connect to Ethereum network
	client, err := ethclient.Dial("https://goerli.base.org")
	if err != nil {
		log.Fatal(err)
	}

	// Private key of your wallet
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		log.Fatal(err)
	}

	// Address of your wallet
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Nonce of your wallet
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(int64(inValue * 1e18))

	// Gas limit for the transaction
	gasLimit := uint64(21000)

	// Gas price for the transaction
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Read the recipient addresses from the CSV file
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Create and sign each transaction
	for _, record := range records {
		// Convert recipient address to Ethereum common.Address
		toAddress := common.HexToAddress(record[0])

		// Create a new transaction
		tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

		chainID := big.NewInt(84531)
		// Sign the transaction
		signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatal(err)
		}

		// Send the transaction to the Ethereum network
		err = client.SendTransaction(context.Background(), signedTx)
		if err != nil {
			log.Fatal(err)
		}

		// Increment nonce for next transaction
		nonce++
	}

	return nil
}
