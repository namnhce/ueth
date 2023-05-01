package wallets

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/namnhce/ueth/pkg/services"
)

type Wallet struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privatekey"`
}

func generateEthWallets(mnemonic string, numWallets int) ([]Wallet, error) {
	// Create an Ethereum HD wallet from the mnemonic seed phrase
	blackList, err := getBlacklist()
	if err != nil {
		return nil, err
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	// Generate multiple Ethereum accounts from the HD wallet
	wallets := make([]Wallet, 0)
	for i := 0; i < numWallets; i++ {

		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		account, err := wallet.Derive(path, false)
		if err != nil {
			return nil, err
		}

		privateKey, err := wallet.PrivateKeyHex(account)
		if err != nil {
			return nil, err
		}

		address := account.Address.Hex()
		if _, ok := blackList[address]; ok {
			numWallets++
			continue
		}

		wallets = append(wallets, Wallet{
			Address:    address,
			PrivateKey: privateKey,
		})
	}

	return wallets, nil
}

func exportCsv(wallets []Wallet, filename string) (string, error) {
	// Create the CSV file
	csvFile, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer csvFile.Close()

	// Create the CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write the header row
	header := []string{"Address", "Private Key"}
	writer.Write(header)

	// Write the data rows
	for _, wallet := range wallets {
		row := []string{wallet.Address, wallet.PrivateKey}
		writer.Write(row)
	}

	return csvFile.Name(), nil
}

func DoGenerateWallet(ctx services.CLIContext) error {
	mnemonic := ctx.String("mnemonic")
	if mnemonic == "" {
		return errors.New("must provide mnemonic")
	}

	numWallets := ctx.Int("num-wallets")
	if numWallets == 0 {
		return errors.New("must provide number of accounts")
	}

	if numWallets > 1000 {
		return errors.New("number of accounts must be less than 1000")
	}

	outputFileName := ctx.String("output")
	if outputFileName == "" {
		return errors.New("must provide outputFileName")
	}

	wallets, err := generateEthWallets(mnemonic, numWallets)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	_, err = exportCsv(wallets, outputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Generate wallet successfully")
	return nil
}

func getBlacklist() (map[string]string, error) {
	file, err := os.Open("blacklist.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	defer file.Close()
	rs, err := parseCSVFile(file)
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func parseCSVFile(file io.Reader) (map[string]string, error) {
	csvReader := csv.NewReader(file)
	_, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	rowNumber := 1
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error when read content csv at row: %d", rowNumber)
			continue
		}
		record := parseField(row[0])
		if record != "" {
			result[record] = record
		}
	}

	return result, nil
}

func parseField(field string) string {
	field = strings.TrimSpace(field)
	if isCSVNull(field) {
		field = ""
	}

	return field
}

func isCSVNull(value string) bool {
	loweredVal := strings.ToLower(strings.TrimSpace(value))
	return loweredVal == "(null)" || loweredVal == "null"
}
