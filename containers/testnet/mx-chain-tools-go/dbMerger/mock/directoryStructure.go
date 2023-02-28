package mock

import (
	"errors"
	"fmt"
	"os"
	"strings"

	logger "github.com/multiversx/mx-chain-logger-go"
)

const pathSeparator = "/"

var log = logger.GetOrCreate("mock")

type pathPart struct {
	name        string
	isDirectory bool
	children    []*pathPart
}

func newPathPart(name string) *pathPart {
	dirs := strings.Split(name, pathSeparator)

	return &pathPart{
		name:        dirs[0],
		isDirectory: true,
	}
}

func (part *pathPart) addChild(name string, isDirectory bool) {
	if !part.isDirectory {
		log.Error("can not add child to a non-directory part",
			"name", part.name, "child", name)
		return
	}

	dirs := strings.Split(name, pathSeparator)
	currentPathPart := part
	for i := 0; i < len(dirs)-1; i++ {
		currentPathPart = currentPathPart.getOrCreateChild(dirs[i], isDirectory)
	}

	currentPathPart.getOrCreateChild(dirs[len(dirs)-1], isDirectory)
}

func (part *pathPart) getOrCreateChild(name string, isDirectory bool) *pathPart {
	if part.name == name {
		return part
	}

	for _, child := range part.children {
		if child.name == name {
			return child
		}
	}

	if !part.isDirectory {
		log.Error("can not add a child to a non-directory path part", "name", part.name)
		return part
	}

	child := &pathPart{
		name:        name,
		isDirectory: isDirectory,
	}
	part.children = append(part.children, child)

	return child
}

func (part *pathPart) string(prefixForName string, hasSiblings bool) string {
	strs := prefixForName + "+- " + part.name + "\n"

	prefixForChild := prefixForName
	if hasSiblings {
		prefixForChild += "|  "
	} else {
		prefixForChild += "   "
	}

	if len(part.children) > 0 {
		strs += prefixForChild + "|\n"
	}

	for i := 0; i < len(part.children); i++ {
		child := part.children[i]

		hasSiblings = len(part.children) > 1
		isLastElement := i == len(part.children)-1

		strs += child.string(prefixForChild, hasSiblings && !isLastElement)
	}

	if len(part.children) > 0 {
		strs += prefixForChild + "\n"
	}
	return strs
}

func (part *pathPart) listDirectory(paths []string) ([]os.DirEntry, error) {
	if len(paths) == 0 {
		return nil, errors.New("empty path to list directories from")
	}

	if paths[0] != part.name {
		return nil, fmt.Errorf("can not traverse to the '%s` path", paths[0])
	}
	if len(paths) == 1 {
		dirEntries := make([]os.DirEntry, 0, len(paths))
		for _, child := range part.children {
			dirEntries = append(dirEntries,
				&DirEntryStub{
					NameValue:  child.name,
					IsDirValue: child.isDirectory,
				})
		}

		return dirEntries, nil
	}

	for _, child := range part.children {
		if child.name == paths[1] {
			return child.listDirectory(paths[1:])
		}
	}

	return nil, fmt.Errorf("child %s not found", paths[1])
}

// DirectoryStructure holds the root and is able to manipulate sub-paths
type DirectoryStructure struct {
	root    *pathPart
	crtPath *pathPart
}

// AddPath will add a known path
func (ds *DirectoryStructure) AddPath(path string, isDirectory bool) {
	if ds.root == nil {
		ds.root = newPathPart(path)
		ds.crtPath = ds.root
	}

	ds.root.addChild(path, isDirectory)
}

// ListDirectory will list the directory by traversing the path parts all the way from the root instance
func (ds *DirectoryStructure) ListDirectory(name string) ([]os.DirEntry, error) {
	if ds.root == nil {
		return nil, errors.New("nil root path, no directories added")
	}

	return ds.root.listDirectory(strings.Split(name, pathSeparator))
}

// String returns the directory struct as a string
func (ds *DirectoryStructure) String() string {
	return ds.root.string("", false)
}
