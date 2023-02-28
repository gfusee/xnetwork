package main

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/tidwall/gjson"
)

const (
	multipleSearchBulk = 10000
	accountsEsdtIndex  = "accountsesdt"
)

type extraTokensChecker struct {
	nftBalancesGetter tokenBalancesGetter
	elasticClient     elasticMultiSearchClient
}

func newExtraTokensCrossChecker(client elasticMultiSearchClient, nftBalancesGetter tokenBalancesGetter) (crossTokenChecker, error) {
	if client == nil {
		return nil, errNilElasticClient
	}
	if nftBalancesGetter == nil {
		return nil, errNilTokenBalancesGetter
	}

	return &extraTokensChecker{
		nftBalancesGetter: nftBalancesGetter,
		elasticClient:     client,
	}, nil
}

func (etc *extraTokensChecker) crossCheckExtraTokens(tokens map[string]struct{}) ([]string, error) {
	numTokens := len(tokens)
	log.Info("starting to cross-check", "num of tokens", numTokens)

	bulkSize := core.MinInt(multipleSearchBulk, numTokens)
	tokensThatStillExist := make([]string, 0)
	requests := make([]string, 0, bulkSize)
	currentTokenIdx := 0
	counterRequests := 0
	for token := range tokens {
		currentTokenIdx++
		requests = append(requests, createRequest(token))

		notEnoughRequests := len(requests) < bulkSize
		notLastBulk := currentTokenIdx != numTokens
		if notEnoughRequests && notLastBulk {
			continue
		}

		respBytes, err := etc.elasticClient.GetMultiple(accountsEsdtIndex, requests)
		if err != nil {
			log.Error("elasticClient.GetMultiple(accountsesdt, requests)",
				"error", err,
				"requests", requests)
			return nil, err
		}

		responses := gjson.Get(string(respBytes), "responses").Array()
		tokensThatStillExistFromRequest, err := etc.checkIndexerResponse(requests, responses)
		if err != nil {
			return nil, err
		}

		go printProgress(numTokens, currentTokenIdx)

		counterRequests += len(requests)
		requests = make([]string, 0, bulkSize)
		tokensThatStillExist = append(tokensThatStillExist, tokensThatStillExistFromRequest...)
	}

	log.Info("finished cross-checking",
		"total num of tokens", numTokens,
		"total num of tokens cross-checked", currentTokenIdx,
		"total num of tokens requests in indexer", counterRequests)

	if numTokens != currentTokenIdx || numTokens != counterRequests {
		return nil, errors.New("failed to cross check all tokens, check logs")
	}

	return tokensThatStillExist, nil
}

func (etc *extraTokensChecker) checkIndexerResponse(requests []string, responses []gjson.Result) ([]string, error) {
	tokensThatStillExist := make([]string, 0)
	for idxRequestedToken, res := range responses {
		hits := res.Get("hits.hits").Array()
		if len(hits) != 0 {
			token := gjson.Get(requests[idxRequestedToken], "query.match.identifier.query").String()
			log.Debug("found token in indexer with hits/accounts",
				"token", token,
				"num hits/accounts", len(hits))

			tokenExists, err := etc.crossCheckToken(hits, token)
			if err != nil {
				return nil, err
			}

			if tokenExists {
				tokensThatStillExist = append(tokensThatStillExist, token)
			}
		}
		idxRequestedToken++
	}

	return tokensThatStillExist, nil
}

func (etc *extraTokensChecker) crossCheckToken(hits []gjson.Result, token string) (bool, error) {
	tokenExists := false
	for _, hit := range hits {
		address := hit.Get("_source.address").String()
		balance, err := etc.nftBalancesGetter.GetBalance(address, token)
		if err != nil {
			return false, err
		}

		log.Debug("checking gateway if token still exists in trie",
			"token", token,
			"address", address)

		if balance != "0" {
			tokenExists = true
			log.Error("cross-check failed; found token which is still in other address",
				"token", token,
				"balance", balance,
				"address", address)
			break
		}

		log.Warn("possible indexer problem",
			"token", token,
			"hit in address", address,
			"found in trie", false)
	}

	return tokenExists, nil
}

func createRequest(token string) string {
	return `{"query" : {"match" : { "identifier": {"query":"` + token + `","operator":"and"}}}}`
}

func printProgress(numTokens, numTokensCrossChecked int) {
	log.Info("status",
		"num cross checked tokens", numTokensCrossChecked,
		"remaining num of tokens to check", numTokens-numTokensCrossChecked,
		"progress(%)", (100*numTokensCrossChecked)/numTokens) // this should not panic with div by zero, since func is only called if numTokens > 0
}
