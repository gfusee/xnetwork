package trieToolsCommon

// This is not concurrent safe
type addressTokensMap struct {
	internalMap map[string]map[string]struct{}
}

// NewAddressTokensMap creates a new map<address, tokens> handler
func NewAddressTokensMap() AddressTokensMap {
	return &addressTokensMap{
		internalMap: make(map[string]map[string]struct{}),
	}
}

// Add will add all provided tokens to the corresponding address
func (atm *addressTokensMap) Add(address string, tokens map[string]struct{}) {
	_, addressExists := atm.internalMap[address]
	if !addressExists {
		atm.internalMap[address] = make(map[string]struct{})
	}

	atm.addTokens(address, tokens)
}

func (atm *addressTokensMap) addTokens(address string, tokens map[string]struct{}) {
	for token := range tokens {
		atm.internalMap[address][token] = struct{}{}
	}
}

func copyTokens(tokens map[string]struct{}) map[string]struct{} {
	ret := make(map[string]struct{})
	for token := range tokens {
		ret[token] = struct{}{}
	}

	return ret
}

// HasAddress checks if the address is in map
func (atm *addressTokensMap) HasAddress(address string) bool {
	_, found := atm.internalMap[address]
	return found
}

// NumAddresses returns the num of addresses in map
func (atm *addressTokensMap) NumAddresses() uint64 {
	return uint64(len(atm.internalMap))
}

// NumTokens returns the num of tokens in map for all addresses
func (atm *addressTokensMap) NumTokens() uint64 {
	numTokens := uint64(0)
	for _, tokens := range atm.internalMap {
		numTokens += uint64(len(tokens))
	}

	return numTokens
}

// GetMapCopy returns an internal copy map
func (atm *addressTokensMap) GetMapCopy() map[string]map[string]struct{} {
	addressTokensMapCopy := make(map[string]map[string]struct{})

	for address, tokens := range atm.internalMap {
		addressTokensMapCopy[address] = make(map[string]struct{})
		for token := range tokens {
			addressTokensMapCopy[address][token] = struct{}{}
		}
	}

	return addressTokensMapCopy
}

// GetTokens returns all tokens of the provided address
func (atm *addressTokensMap) GetTokens(address string) map[string]struct{} {
	return copyTokens(atm.internalMap[address])
}

// Delete deletes the map entry for the provided address
func (atm *addressTokensMap) Delete(address string) {
	delete(atm.internalMap, address)
}

// GetAllTokens returns all tokens from all addresses
func (atm *addressTokensMap) GetAllTokens() map[string]struct{} {
	allTokens := make(map[string]struct{})
	for _, tokens := range atm.internalMap {
		for token := range tokens {
			allTokens[token] = struct{}{}
		}
	}

	return allTokens
}

// Clone returns a shallow clone of the current object
func (atm *addressTokensMap) Clone() AddressTokensMap {
	mapCopy := atm.GetMapCopy()
	return &addressTokensMap{
		internalMap: mapCopy,
	}
}
