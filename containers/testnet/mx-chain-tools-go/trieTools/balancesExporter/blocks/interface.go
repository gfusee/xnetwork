package blocks

type trieWrapper interface {
	IsRootHashAvailable(rootHash []byte) bool
}
