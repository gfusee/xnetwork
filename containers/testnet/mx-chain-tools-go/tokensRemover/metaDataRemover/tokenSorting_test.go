package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortTokensIDByNonce(t *testing.T) {
	t.Parallel()

	t.Run("invalid token format, should error", func(t *testing.T) {
		t.Parallel()

		tokens := map[string]struct{}{
			"token1-rand1-0f": {},
			"token1-rand1":    {},
		}

		sortedTokens, err := sortTokensIDByNonce(tokens)
		require.Nil(t, sortedTokens)
		require.ErrorIs(t, err, errInvalidTokenFormat)
		require.True(t, strings.Contains(err.Error(), "token1-rand1"))
	})

	t.Run("invalid nonce format, should error", func(t *testing.T) {
		t.Parallel()

		tokens := map[string]struct{}{
			"token1-rand1-0f": {},
			"token1-rand1-zx": {},
		}

		sortedTokens, err := sortTokensIDByNonce(tokens)
		require.Nil(t, sortedTokens)
		require.ErrorIs(t, err, errCouldNotConvertNonceToBigInt)
		require.True(t, strings.Contains(err.Error(), "zx"))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		tokens := map[string]struct{}{
			"token1-rand1-0f": {},
			"token1-rand1-01": {},
			"token1-rand1-0a": {},
			"token1-rand1-0b": {},

			"token2-rand2-04": {},

			"token3-rand3-04": {},
			"token3-rand3-08": {},
		}

		sortedTokens, err := sortTokensIDByNonce(tokens)
		require.Nil(t, err)
		require.Equal(t, sortedTokens, map[string][]uint64{
			"token1-rand1": {1, 10, 11, 15},
			"token2-rand2": {4},
			"token3-rand3": {4, 8},
		})
	})
}

func TestGroupTokensByIntervals(t *testing.T) {
	tokens := map[string][]uint64{
		"token1": {1, 2, 3, 8, 9, 10},
		"token2": {1},
		"token3": {3, 9},
		"token4": {11, 12},
		"token5": {10, 100, 101, 102, 111},
		"token6": {4, 5, 6, 7},
	}

	sortedTokens := groupTokensByIntervals(tokens)
	require.Equal(t, sortedTokens,
		map[string][]*interval{
			"token1": {
				{
					start: 1,
					end:   3,
				},
				{
					start: 8,
					end:   10,
				},
			},
			"token2": {
				{
					start: 1,
					end:   1,
				},
			},
			"token3": {
				{
					start: 3,
					end:   3,
				},
				{
					start: 9,
					end:   9,
				},
			},
			"token4": {
				{
					start: 11,
					end:   12,
				},
			},
			"token5": {
				{
					start: 10,
					end:   10,
				},
				{
					start: 100,
					end:   102,
				},
				{
					start: 111,
					end:   111,
				},
			},
			"token6": {
				{
					start: 4,
					end:   7,
				},
			},
		},
	)
}

func TestSortTokenIntervalsByMaxConsecutiveNonces(t *testing.T) {
	tokensIntervals := map[string][]*interval{
		"token1": {
			{
				start: 0,
				end:   0,
			},
			{
				start: 1,
				end:   1,
			},
			{
				start: 2,
				end:   3,
			},
			{
				start: 4,
				end:   8,
			},
		},
		"token2": {
			{
				start: 1,
				end:   5,
			},
		},
		"token3": {
			{
				start: 0,
				end:   0,
			},
			{
				start: 1,
				end:   4,
			},
			{
				start: 6,
				end:   7,
			},
		},
	}

	sortedTokens := sortTokenIntervalsByMaxConsecutiveNonces(tokensIntervals)
	expectedOutput := []*tokenWithInterval{
		{
			tokenID: "token1",
			interval: &interval{
				start: 4,
				end:   8,
			},
		},
		{
			tokenID: "token2",
			interval: &interval{
				start: 1,
				end:   5,
			},
		},
		{
			tokenID: "token3",
			interval: &interval{
				start: 1,
				end:   4,
			},
		},
		{
			tokenID: "token1",
			interval: &interval{
				start: 2,
				end:   3,
			},
		},
		{
			tokenID: "token3",
			interval: &interval{
				start: 6,
				end:   7,
			},
		},
		{
			tokenID: "token1",
			interval: &interval{
				start: 0,
				end:   0,
			},
		},
		{
			tokenID: "token1",
			interval: &interval{
				start: 1,
				end:   1,
			},
		},
		{
			tokenID: "token3",
			interval: &interval{
				start: 0,
				end:   0,
			},
		},
	}

	require.Equal(t, expectedOutput, sortedTokens)
}

func TestGroupTokenIntervalsInBulks(t *testing.T) {
	tokensIntervals := []*tokenWithInterval{
		{
			tokenID: "token1",
			interval: &interval{
				start: 4,
				end:   8,
			},
		},
		{
			tokenID: "token2",
			interval: &interval{
				start: 1,
				end:   5,
			},
		},
		{
			tokenID: "token3",
			interval: &interval{
				start: 1,
				end:   4,
			},
		},
		{
			tokenID: "token1",
			interval: &interval{
				start: 2,
				end:   3,
			},
		},
		{
			tokenID: "token3",
			interval: &interval{
				start: 6,
				end:   7,
			},
		},
		{
			tokenID: "token1",
			interval: &interval{
				start: 0,
				end:   0,
			},
		},
		{
			tokenID: "token1",
			interval: &interval{
				start: 1,
				end:   1,
			},
		},
		{
			tokenID: "token3",
			interval: &interval{
				start: 0,
				end:   0,
			},
		},
	}

	output := groupTokenIntervalsInBulks(tokensIntervals, 6)
	bulk1 := []*tokenData{
		{
			tokenID: "token1",
			intervals: []*interval{
				{
					start: 4,
					end:   8,
				},
			},
		},
		{
			tokenID: "token2",
			intervals: []*interval{
				{
					start: 1,
					end:   1,
				},
			},
		},
	}
	bulk2 := []*tokenData{
		{
			tokenID: "token2",
			intervals: []*interval{
				{
					start: 2,
					end:   5,
				},
			},
		},
		{
			tokenID: "token3",
			intervals: []*interval{
				{
					start: 1,
					end:   2,
				},
			},
		},
	}
	bulk3 := []*tokenData{
		{
			tokenID: "token1",
			intervals: []*interval{
				{
					start: 2,
					end:   3,
				},
			},
		},
		{
			tokenID: "token3",
			intervals: []*interval{
				{
					start: 3,
					end:   4,
				},
				{
					start: 6,
					end:   7,
				},
			},
		},
	}
	bulk4 := []*tokenData{
		{
			tokenID: "token1",
			intervals: []*interval{
				{
					start: 0,
					end:   0,
				},
				{
					start: 1,
					end:   1,
				},
			},
		},
		{
			tokenID: "token3",
			intervals: []*interval{
				{
					start: 0,
					end:   0,
				},
			},
		},
	}
	expectedOutput := [][]*tokenData{bulk1, bulk2, bulk3, bulk4}
	require.Equal(t, expectedOutput, output)
}
