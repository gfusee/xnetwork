package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTokenBalanceBodyResponse(balance, error string) string {
	return fmt.Sprintf(` {
		"data": {
			"tokenData": {
				"balance": %s,
			}
		},
		"error": "%s",
	}`, balance, error)
}

func TestTokenBalanceGetter_GetBalance(t *testing.T) {
	proxyUrl := "proxy"
	address := "erd1abc"
	token := "token-rand-0f"
	expectedBalance := "44112"

	t.Run("should work without errors", func(t *testing.T) {
		getFunc := func(url string) (resp *http.Response, err error) {
			require.Equal(t, fmt.Sprintf("%s/address/%s/nft/%s/nonce/%d", proxyUrl, address, "token-rand", 15), url)

			body := createTokenBalanceBodyResponse(expectedBalance, "")
			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(body)),
			}, nil
		}

		tbg := newTokenBalanceGetter(proxyUrl, getFunc)

		balance, err := tbg.GetBalance(address, token)
		require.Nil(t, err)
		require.Equal(t, expectedBalance, balance)
	})

	t.Run("should work after few retrials", func(t *testing.T) {
		body := ""
		errorCt := 0
		getFunc := func(url string) (resp *http.Response, err error) {
			require.Equal(t, fmt.Sprintf("%s/address/%s/nft/%s/nonce/%d", proxyUrl, address, "token-rand", 15), url)

			errorCt++
			switch errorCt {
			case 1:
				return nil, errors.New("first error")
			case 2:
				body = createTokenBalanceBodyResponse("0", "second error")
			case 3:
				return &http.Response{Body: nil}, nil
			case 4:
				body = createTokenBalanceBodyResponse(expectedBalance, "")
			default:
				require.Fail(t, "should not do another retrial")
			}

			return &http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(body)),
			}, nil
		}

		tbg := newTokenBalanceGetter(proxyUrl, getFunc)

		balance, err := tbg.GetBalance(address, token)
		require.Nil(t, err)
		require.Equal(t, expectedBalance, balance)
	})

	t.Run("could not get balance after max retrials", func(t *testing.T) {
		errorCt := 0
		getFunc := func(url string) (resp *http.Response, err error) {
			require.Equal(t, fmt.Sprintf("%s/address/%s/nft/%s/nonce/%d", proxyUrl, address, "token-rand", 15), url)
			errorCt++
			return nil, errors.New("local error")
		}

		tbg := newTokenBalanceGetter(proxyUrl, getFunc)

		balance, err := tbg.GetBalance(address, token)
		require.Empty(t, balance)
		require.Equal(t, maxRequestsRetrial, errorCt)

		require.NotNil(t, err)
		require.ErrorIs(t, err, errCouldNotGetBalance)
		require.True(t, strings.Contains(err.Error(), address))
		require.True(t, strings.Contains(err.Error(), fmt.Sprintf("%d", maxRequestsRetrial)))
	})
}
