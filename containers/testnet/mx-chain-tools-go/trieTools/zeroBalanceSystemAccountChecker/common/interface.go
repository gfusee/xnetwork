package common

import "io"

// FileInfo should provide basic information about a file
type FileInfo interface {
	Name() string
	IsDir() bool
}

// FileHandler defines what a sys file handler should do (e.g. read directories, get working dir)
type FileHandler interface {
	Open(name string) (io.Reader, error)
	ReadAll(r io.Reader) ([]byte, error)
	Getwd() (dir string, err error)
	ReadDir(dirname string) ([]FileInfo, error)
}
