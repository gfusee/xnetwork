package process

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeIntervals(t *testing.T) {
	start := int64(1653974028)
	end := int64(1653981227)

	res, err := computeIntervals(start, end, 3)
	require.Nil(t, err)
	require.Equal(t, []*interval{
		{
			start: 1653974028,
			stop:  1653976427,
		},
		{
			start: 1653976427,
			stop:  1653978826,
		},
		{
			start: 1653978826,
			stop:  1653981227,
		},
	}, res)
}
