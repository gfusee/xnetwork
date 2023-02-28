package export

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-go/state"
)

type formatterPlainJson struct {
}

func (f *formatterPlainJson) toText(accounts []*state.UserAccountData, args formatterArgs) (string, error) {
	records := make([]plainBalance, 0, len(accounts))

	for _, account := range accounts {
		address := addressConverter.Encode(account.Address)
		balance := account.Balance.String()

		records = append(records, plainBalance{
			Address: address,
			Balance: balance,
		})
	}

	recordsJson, err := json.MarshalIndent(records, "", fourSpaces)
	if err != nil {
		return "", err
	}

	return string(recordsJson), nil
}

func (f *formatterPlainJson) getFileExtension() string {
	return "json"
}
