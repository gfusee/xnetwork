package export

import (
	"fmt"
	"strings"

	"github.com/multiversx/mx-chain-go/state"
)

type plainBalance struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

type formatterPlainText struct {
}

func (f *formatterPlainText) toText(accounts []*state.UserAccountData, args formatterArgs) (string, error) {
	var builder strings.Builder

	for _, account := range accounts {
		address := addressConverter.Encode(account.Address)
		balance := account.Balance.String()
		line := fmt.Sprintf("%s %s %s\n", address, balance, args.currency)
		_, err := builder.WriteString(line)
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}

func (f *formatterPlainText) getFileExtension() string {
	return "txt"
}
