package export

import (
	"encoding/json"

	"github.com/multiversx/mx-chain-go/state"
)

type rosettaBalance struct {
	AccountIdentifier rosettaAccountIdentifier `json:"account_identifier"`
	Currency          *rosettaCurrency         `json:"currency"`
	Value             string                   `json:"value"`
}

type rosettaAccountIdentifier struct {
	Address string `json:"address"`
}

type rosettaCurrency struct {
	Symbol   string `json:"symbol"`
	Decimals uint   `json:"decimals"`
}

type formatterRosettaJson struct {
}

func (f *formatterRosettaJson) toText(accounts []*state.UserAccountData, args formatterArgs) (string, error) {
	records := make([]rosettaBalance, 0, len(accounts))

	currency := &rosettaCurrency{
		Symbol:   args.currency,
		Decimals: args.currencyDecimals,
	}

	for _, account := range accounts {
		address := addressConverter.Encode(account.Address)
		balance := account.Balance.String()

		records = append(records, rosettaBalance{
			AccountIdentifier: rosettaAccountIdentifier{
				Address: address,
			},
			Currency: currency,
			Value:    balance,
		})
	}

	recordsJson, err := json.MarshalIndent(records, "", fourSpaces)
	if err != nil {
		return "", err
	}

	return string(recordsJson), nil
}

func (f *formatterRosettaJson) getFileExtension() string {
	return "rosetta.json"
}
