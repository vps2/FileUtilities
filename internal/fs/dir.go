package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//FilterFunc если функция возвращает true, то данный экземпляр содержимого каталога должен содержаться в выборке.
type FilterFunc func(fileInfo os.FileInfo) bool

var defaultFilterFunc = func(fileInfo os.FileInfo) bool {
	return true
}

//DirReader представляет собой просмотрщик содержимого каталога, указанного в поле Path.
type DirReader struct {
	path   string
	filter FilterFunc
}

//NewDirReader возвращает настроенный экземпляр DirReader
func NewDirReader(path string) DirReader {
	return DirReader{
		path:   path,
		filter: defaultFilterFunc,
	}
}

//NewDirReaderWithFilter возвращает настроенный экземпляр DirReader
func NewDirReaderWithFilter(path string, filter FilterFunc) DirReader {
	return DirReader{
		path:   path,
		filter: filter,
	}
}

//ReadWithFilter возвращает содержимое каталога.
func (r DirReader) Read() (res []*File, err error) {
	if err = r.validate(); err != nil {
		return
	}
	entries, err := ioutil.ReadDir(r.path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if ok := r.filter(entry); ok {
			entryPathName := filepath.Join(r.path, entry.Name())
			res = append(res, &File{entryPathName})
		}
	}

	return
}

func (r DirReader) validate() error {
	if exists := isExists(r.path); !exists {
		return fmt.Errorf("directory '%s' is not exists: %w", r.path, ErrNotExists)
	}
	if regular, _ := isDirectory(r.path); !regular {
		return fmt.Errorf("'%s' is not a directory: %w", r.path, ErrNotDirectory)
	}

	return nil
}
