package mocks

// TokenBalanceGetterStub -
type TokenBalanceGetterStub struct {
	GetBalanceCalled func(address, token string) (string, error)
}

// GetBalance -
func (tbg *TokenBalanceGetterStub) GetBalance(address, token string) (string, error) {
	if tbg.GetBalanceCalled != nil {
		return tbg.GetBalanceCalled(address, token)
	}

	return "", nil
}
