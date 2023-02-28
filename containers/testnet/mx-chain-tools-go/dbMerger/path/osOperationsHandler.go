package path

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("path")

const dirPermMode = os.FileMode(0755)

type osOperationsHandler struct {
}

// NewOsOperationsHandler returns a new instance of a handler that deals with the OS-level functions
func NewOsOperationsHandler() *osOperationsHandler {
	return &osOperationsHandler{}
}

// CopyDirectory is able to recursively copy the contents of one directory to another
func (handler *osOperationsHandler) CopyDirectory(destination string, source string) error {
	log.Debug("copying raw data", "source", source, "destination", destination)
	entries, errReadDir := os.ReadDir(source)
	if errReadDir != nil {
		return errReadDir
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(destination, entry.Name())

		fileInfo, errStat := os.Stat(sourcePath)
		if errStat != nil {
			return errStat
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			err := createIfNotExists(destPath, dirPermMode)
			if err != nil {
				return err
			}

			err = handler.CopyDirectory(destPath, sourcePath)
			if err != nil {
				return err
			}
		case os.ModeSymlink:
			err := copySymLink(destPath, sourcePath)
			if err != nil {
				return err
			}
		default:
			err := copyFile(destPath, sourcePath)
			if err != nil {
				return err
			}
		}

		err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid))
		if err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			errChmod := os.Chmod(destPath, fInfo.Mode())
			if errChmod != nil {
				return errChmod
			}
		}
	}
	return nil
}

func copyFile(dstFile string, srcFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer func() {
		errClose := out.Close()
		log.LogIfError(errClose)
	}()

	in, err := os.Open(srcFile)
	defer func() {
		errClose := in.Close()
		log.LogIfError(errClose)
	}()

	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	err := os.MkdirAll(dir, perm)
	if err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func copySymLink(dest string, source string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}

	return os.Symlink(link, dest)
}

// CheckIfDirectoryIsEmpty returns true if the directory exists and is empty
func (handler *osOperationsHandler) CheckIfDirectoryIsEmpty(directory string) error {
	entries, errReadDir := os.ReadDir(directory)
	if errReadDir != nil {
		return fmt.Errorf("%w while reading the directory %s", errReadDir, directory)
	}

	if len(entries) > 0 {
		return fmt.Errorf("%w %s", errDirectoryIsNotEmpty, directory)
	}

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *osOperationsHandler) IsInterfaceNil() bool {
	return handler == nil
}
