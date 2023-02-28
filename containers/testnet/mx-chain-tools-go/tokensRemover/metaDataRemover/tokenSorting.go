package main

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

type interval struct {
	start uint64
	end   uint64
}

func (i *interval) split(index uint64) (*interval, *interval) {
	first := &interval{
		start: i.start,
		end:   i.start + index - 1,
	}

	second := &interval{
		start: first.end + 1,
		end:   i.end,
	}

	return first, second
}

type tokenWithInterval struct {
	tokenID  string
	interval *interval
}

type tokenData struct {
	tokenID   string
	intervals []*interval
}

func sortTokensIDByNonce(tokens map[string]struct{}) (map[string][]uint64, error) {
	ret := make(map[string][]uint64)
	for token := range tokens {
		splits := strings.Split(token, "-")
		if len(splits) != 3 {
			return nil, fmt.Errorf("found %w = %s; expected format = [ticker-randSequence-nonce]", errInvalidTokenFormat, token)
		}

		tokenID := splits[0] + "-" + splits[1] // ticker-randSequence
		nonceStr := splits[2]
		nonceBI := big.NewInt(0)
		_, ok := nonceBI.SetString(nonceStr, 16)
		if !ok {
			return nil, fmt.Errorf("%w; token = %s, nonce string = %s", errCouldNotConvertNonceToBigInt, tokenID, nonceStr)
		}

		ret[tokenID] = append(ret[tokenID], nonceBI.Uint64())
	}

	log.Info("found", "num of tokensID", len(ret))
	for _, nonces := range ret {
		sort.SliceStable(nonces, func(i, j int) bool {
			return nonces[i] < nonces[j]
		})
	}

	return ret, nil
}

func groupTokensByIntervals(tokens map[string][]uint64) map[string][]*interval {
	ret := make(map[string][]*interval)
	for token, nonces := range tokens {
		ret[token] = getIntervals(nonces)
	}

	return ret
}

func getIntervals(nonces []uint64) []*interval {
	numNonces := len(nonces)
	intervals := make([]*interval, 0)
	for idx := 0; idx < numNonces; idx++ {
		nonce := nonces[idx]
		if idx+1 >= numNonces {
			intervals = append(intervals, &interval{
				start: nonce,
				end:   nonce,
			})
			break
		}

		currInterval := &interval{start: nonce}
		numConsecutiveNonces := uint64(0)
		for idx < numNonces-1 {
			currNonce := nonces[idx]
			nextNonce := nonces[idx+1]
			if nextNonce-currNonce > 1 {
				break
			}

			numConsecutiveNonces++
			idx++
		}

		currInterval.end = currInterval.start + numConsecutiveNonces
		intervals = append(intervals, currInterval)
	}

	return intervals
}

func sortTokenIntervalsByMaxConsecutiveNonces(tokens map[string][]*interval) []*tokenWithInterval {
	ret := make([]*tokenWithInterval, 0)
	for token, intervals := range tokens {
		for _, currInterval := range intervals {
			ret = append(ret, &tokenWithInterval{
				tokenID:  token,
				interval: currInterval,
			})

		}
	}

	sort.SliceStable(ret, func(i, j int) bool {
		consecutiveNonces1 := ret[i].interval.end - ret[i].interval.start + 1
		consecutiveNonces2 := ret[j].interval.end - ret[j].interval.start + 1

		if consecutiveNonces1 == consecutiveNonces2 {
			return ret[i].tokenID < ret[j].tokenID
		}

		return consecutiveNonces1 > consecutiveNonces2
	})

	return ret
}

func groupTokenIntervalsInBulks(tokens []*tokenWithInterval, bulkSize uint64) [][]*tokenData {
	bulks := make([][]*tokenData, 0, bulkSize)
	currBulk := make(map[string][]*interval, 0)
	numNoncesInBulk := uint64(0)

	tokensCopy := make([]*tokenWithInterval, len(tokens))
	copy(tokensCopy, tokens)

	index := 0
	for index < len(tokensCopy) {
		currTokenData := tokensCopy[index]
		currInterval := currTokenData.interval
		currTokenID := currTokenData.tokenID

		noncesInInterval := currInterval.end - currInterval.start + 1
		availableSlots := bulkSize - numNoncesInBulk
		if availableSlots >= noncesInInterval {
			numNoncesInBulk += noncesInInterval
			currBulk[currTokenID] = append(currBulk[currTokenID], currInterval)

		} else {
			first, second := currInterval.split(availableSlots)

			numNoncesInBulk += availableSlots
			currBulk[currTokenID] = append(currBulk[currTokenID], first)
			tokensCopy = insert(tokensCopy, index+1, &tokenWithInterval{tokenID: currTokenID, interval: second})
		}

		bulkFull := numNoncesInBulk == bulkSize
		lastInterval := index == len(tokensCopy)-1
		if bulkFull || lastInterval {
			bulks = append(bulks, tokensMapToOrderedArray(currBulk))

			currBulk = make(map[string][]*interval, 0)
			numNoncesInBulk = 0
		}

		index++
	}

	return bulks
}

func insert(tokens []*tokenWithInterval, index int, token *tokenWithInterval) []*tokenWithInterval {
	if len(tokens) <= index {
		return append(tokens, token)
	}

	tokens = append(tokens[:index+1], tokens[index:]...)
	tokens[index] = token
	return tokens
}

func tokensMapToOrderedArray(tokens map[string][]*interval) []*tokenData {
	ret := make([]*tokenData, 0, len(tokens))

	for token, intervals := range tokens {
		ret = append(ret, &tokenData{
			tokenID:   token,
			intervals: intervals,
		})
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].tokenID < ret[j].tokenID
	})

	return ret
}
