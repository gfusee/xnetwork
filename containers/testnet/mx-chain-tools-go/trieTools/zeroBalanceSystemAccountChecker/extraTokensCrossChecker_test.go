package main

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/mocks"
	"github.com/stretchr/testify/require"
)

type source struct {
	Address string `json:"address"`
}

type addressHits struct {
	Source source `json:"_source"`
}

type hits struct {
	Hits []addressHits `json:"hits"`
}

type response struct {
	Hits hits `json:"hits"`
}

type indexerResponse struct {
	Responses []response `json:"responses"`
}

func requireSameSliceDifferentOrder(t *testing.T, s1, s2 []string) {
	require.Equal(t, len(s1), len(s2))

	for _, elemInS1 := range s1 {
		require.Contains(t, s2, elemInS1)
	}
}

func requireIndexerRequestsAreCorrect(t *testing.T, requests []string, index string, tokens map[string]struct{}) {
	expectedRequests := make([]string, 0, len(tokens))
	for token := range tokens {
		expectedRequests = append(expectedRequests, `{"query" : {"match" : { "identifier": {"query":"`+token+`","operator":"and"}}}}`)
	}

	require.Equal(t, accountsEsdtIndex, index)
	requireSameSliceDifferentOrder(t, expectedRequests, requests)
}

func TestNewExtraTokensCrossChecker(t *testing.T) {
	t.Parallel()

	t.Run("nil elastic indexer, should error", func(t *testing.T) {
		t.Parallel()

		etc, err := newExtraTokensCrossChecker(nil, &mocks.TokenBalanceGetterStub{})
		require.Nil(t, etc)
		require.Equal(t, errNilElasticClient, err)
	})

	t.Run("nil token balances getter, should error", func(t *testing.T) {
		t.Parallel()

		etc, err := newExtraTokensCrossChecker(&mocks.ElasticClientStub{}, nil)
		require.Nil(t, etc)
		require.Equal(t, errNilTokenBalancesGetter, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		etc, err := newExtraTokensCrossChecker(&mocks.ElasticClientStub{}, &mocks.TokenBalanceGetterStub{})
		require.NotNil(t, etc)
		require.Nil(t, err)
	})
}

func TestExtraTokensChecker_CrossCheckExtraTokensErrorCases(t *testing.T) {
	t.Parallel()

	tokens := map[string]struct{}{
		"token": {},
	}

	t.Run("error requesting elastic client multiple search", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("err elastic client")
		elasticClient := &mocks.ElasticClientStub{
			GetMultipleCalled: func(index string, requests []string) ([]byte, error) {
				return nil, expectedErr
			},
		}

		getBalanceCalled := false
		balanceGetter := &mocks.TokenBalanceGetterStub{
			GetBalanceCalled: func(address, token string) (string, error) {
				getBalanceCalled = true
				return "", nil
			},
		}

		checker, _ := newExtraTokensCrossChecker(elasticClient, balanceGetter)
		extraTokens, err := checker.crossCheckExtraTokens(tokens)
		require.Nil(t, extraTokens)
		require.Equal(t, expectedErr, err)
		require.False(t, getBalanceCalled)
	})

	t.Run("error requesting balance", func(t *testing.T) {
		t.Parallel()

		getMultipleCalledCt := 0
		elasticClient := &mocks.ElasticClientStub{
			GetMultipleCalled: func(index string, requests []string) ([]byte, error) {
				getMultipleCalledCt++
				requireIndexerRequestsAreCorrect(t, requests, index, tokens)

				resp := &indexerResponse{
					Responses: []response{
						{
							Hits: hits{
								Hits: []addressHits{{Source: source{Address: "address"}}},
							},
						},
					},
				}

				responseBytes, _ := json.Marshal(resp)
				return responseBytes, nil
			},
		}

		expectedErr := errors.New("err get balance")
		balanceGetter := &mocks.TokenBalanceGetterStub{
			GetBalanceCalled: func(address, token string) (string, error) {
				return "", expectedErr
			},
		}

		checker, _ := newExtraTokensCrossChecker(elasticClient, balanceGetter)
		extraTokens, err := checker.crossCheckExtraTokens(tokens)
		require.Nil(t, extraTokens)
		require.Equal(t, expectedErr, err)
		require.Equal(t, 1, getMultipleCalledCt)
	})
}

func TestExtraTokensChecker_CrossCheckExtraTokens(t *testing.T) {
	t.Parallel()

	address1 := "address1"
	address2 := "address2"
	address3 := "address3"

	token1 := "token1"
	token2 := "token2"
	tokens := map[string]struct{}{
		token1: {},
		token2: {},
	}

	t.Run("tokens not found in any other address when querying indexer, should not check gateway balances", func(t *testing.T) {
		t.Parallel()

		getMultipleCalledCt := 0
		elasticClient := &mocks.ElasticClientStub{
			GetMultipleCalled: func(index string, requests []string) ([]byte, error) {
				getMultipleCalledCt++
				requireIndexerRequestsAreCorrect(t, requests, index, tokens)

				resp := &indexerResponse{
					Responses: []response{
						{
							Hits: hits{
								Hits: []addressHits{},
							},
						},
						{
							Hits: hits{
								Hits: []addressHits{},
							},
						},
					},
				}

				responseBytes, _ := json.Marshal(resp)
				return responseBytes, nil
			},
		}

		getBalanceCalled := false
		balanceGetter := &mocks.TokenBalanceGetterStub{
			GetBalanceCalled: func(address, token string) (string, error) {
				getBalanceCalled = true
				return "", nil
			},
		}

		checker, _ := newExtraTokensCrossChecker(elasticClient, balanceGetter)
		extraTokens, err := checker.crossCheckExtraTokens(tokens)
		require.Nil(t, err)
		require.Empty(t, extraTokens)
		require.False(t, getBalanceCalled)
		require.Equal(t, 1, getMultipleCalledCt)
	})

	t.Run("tokens found in indexer, but not found in trie with balance => possible indexing problem", func(t *testing.T) {
		t.Parallel()

		getMultipleCalledCt := 0
		elasticClient := &mocks.ElasticClientStub{
			GetMultipleCalled: func(index string, requests []string) ([]byte, error) {
				getMultipleCalledCt++
				requireIndexerRequestsAreCorrect(t, requests, index, tokens)

				resp := &indexerResponse{
					Responses: []response{
						{
							Hits: hits{
								Hits: []addressHits{
									{Source: source{Address: address1}},
									{Source: source{Address: address2}},
								},
							},
						},
						{
							Hits: hits{
								Hits: []addressHits{
									{Source: source{Address: address1}},
								},
							},
						},
					},
				}

				responseBytes, _ := json.Marshal(resp)
				return responseBytes, nil
			},
		}

		getBalanceCalledCt := 0
		balanceGetter := &mocks.TokenBalanceGetterStub{
			GetBalanceCalled: func(address, token string) (string, error) {
				getBalanceCalledCt++
				require.True(t, address == address1 || address == address2)
				return "0", nil

			},
		}

		checker, _ := newExtraTokensCrossChecker(elasticClient, balanceGetter)
		extraTokens, err := checker.crossCheckExtraTokens(tokens)
		require.Nil(t, err)
		require.Empty(t, extraTokens)
		require.Equal(t, 3, getBalanceCalledCt)
		require.Equal(t, 1, getMultipleCalledCt)
	})

	t.Run("token1 found in indexer and trie with balance, should be output as token that still exists", func(t *testing.T) {
		t.Parallel()

		getMultipleCalledCt := 0
		elasticClient := &mocks.ElasticClientStub{
			GetMultipleCalled: func(index string, requests []string) ([]byte, error) {
				getMultipleCalledCt++
				requireIndexerRequestsAreCorrect(t, requests, index, tokens)

				resp := &indexerResponse{
					Responses: []response{
						{
							Hits: hits{
								Hits: []addressHits{},
							},
						},
						{
							Hits: hits{
								Hits: []addressHits{},
							},
						},
					},
				}

				token1Hits := hits{
					Hits: []addressHits{
						{Source: source{Address: address1}},
						{Source: source{Address: address2}},
						{Source: source{Address: address3}},
					},
				}

				if strings.Contains(requests[0], "token1") {
					resp.Responses[0].Hits = token1Hits
				} else {
					resp.Responses[1].Hits = token1Hits
				}

				responseBytes, _ := json.Marshal(resp)
				return responseBytes, nil
			},
		}

		getBalanceCalledCt := 0
		balanceGetter := &mocks.TokenBalanceGetterStub{
			GetBalanceCalled: func(address, token string) (string, error) {
				getBalanceCalledCt++
				switch address {
				case address1:
					// address1 will generate a warn(possible indexer problem), since token1 is found in indexer with owner, but not in trie(has 0 balance)
					return "0", nil
				case address2:
					// address2 will have a hit in indexer + balance(123) in trie
					return "123", nil
				default:
					// address3 should not be checked , since the loop will break once the first hit with balance is found
					require.Fail(t, " should not have queried balance for another address")
				}

				return "0", nil
			},
		}

		checker, _ := newExtraTokensCrossChecker(elasticClient, balanceGetter)
		extraTokens, err := checker.crossCheckExtraTokens(tokens)
		require.Nil(t, err)
		require.Equal(t, []string{token1}, extraTokens)
		require.Equal(t, 2, getBalanceCalledCt)
		require.Equal(t, 1, getMultipleCalledCt)
	})

}
