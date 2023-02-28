package trieToolsCommon

// AddressTokensMap should handle a map<address, tokens>
type AddressTokensMap interface {
	Add(addr string, tokens map[string]struct{})
	Delete(address string)
	GetAllTokens() map[string]struct{}
	GetTokens(address string) map[string]struct{}
	GetMapCopy() map[string]map[string]struct{}
	Clone() AddressTokensMap
	HasAddress(addr string) bool
	NumAddresses() uint64
	NumTokens() uint64
}
