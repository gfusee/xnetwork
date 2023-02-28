package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadTokensInput(t *testing.T) {
	tokensMap, err := readTokensInput("tokensTestData/tokens.json")
	require.Nil(t, err)
	expectedMap := map[uint32]map[string]struct{}{
		0: {
			"ZZZ0-c5aa13-01":  {},
			"ZZZ1-c5aa13-01":  {},
			"ZZZ10-eb20df-01": {},
		},
		1: {
			"AAA0-f1fac9-01": {},
			"AAA0-f1fac9-03": {},
			"ZZZ9-ae1fa4-01": {},
		},
		2: {
			"BBB0-adde72-01": {},
			"BBB1-adde72-01": {},
		},
	}
	require.Equal(t, expectedMap, tokensMap)
}
