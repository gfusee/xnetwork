package main

type crossTokenChecker interface {
	crossCheckExtraTokens(tokens map[string]struct{}) ([]string, error)
}

type tokenBalancesGetter interface {
	GetBalance(address, token string) (string, error)
}

type elasticMultiSearchClient interface {
	GetMultiple(index string, requests []string) ([]byte, error)
}
