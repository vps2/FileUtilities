package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

//File абстракция над файлом файловой системы
type File struct {
	PathName string
}

//Name возвращает имя файла
func (f *File) Name() string {
	return filepath.Base(f.PathName)
}

//AbsolutePath возвращает полный путь к файлу
func (f *File) AbsolutePath() string {
	return f.PathName
}

func (f *File) String() string {
	return f.AbsolutePath()
}

//ModTime возвращает время последней модификации файла
func (f *File) ModTime() (time.Time, error) {
	if err := f.validate(); err != nil {
		return time.Time{}, err
	}

	_, modTime, err := fileTimes(f.AbsolutePath())
	return modTime, err
}

//AccessTime возвращает время последнего доступа к файлу
func (f *File) AccessTime() (time.Time, error) {
	if err := f.validate(); err != nil {
		return time.Time{}, err
	}

	accessTime, _, err := fileTimes(f.AbsolutePath())
	return accessTime, err
}

//Delete удаляет файл
func (f *File) Delete() error {
	if err := f.validate(); err != nil {
		return err
	}

	if err := os.Remove(f.AbsolutePath()); err != nil {
		return err
	}

	return nil
}

//CopyTo копирует файл в новое расположение. Если операция копирования прошла удачно, то возвращается указатель на новый файл.
func (f *File) CopyTo(path string) (*File, error) {
	dstFile := &File{PathName: filepath.Join(path, f.Name())}

	//модифицируем время модификации и доступа в новом файле, на такие же значения, как в оригинальном
	defer func() {
		if atime, mtime, err := fileTimes(f.AbsolutePath()); err == nil {
			setFileTimes(dstFile.AbsolutePath(), atime, mtime)
		}
	}()

	if err := f.validate(); err != nil {
		return nil, err
	}

	if exists := isExists(dstFile.AbsolutePath()); exists {
		return nil, fmt.Errorf("file '%s' already exists: %w", dstFile.AbsolutePath(), ErrAlreadyExists)
	}
	if locked := isFileLocked(f.AbsolutePath()); locked {
		return nil, fmt.Errorf("file '%s' is blocked: %w", f.AbsolutePath(), ErrBlocked)
	}

	source, err := os.Open(f.AbsolutePath())
	if err != nil {
		return nil, fmt.Errorf("can not open a file '%s': %w", f.AbsolutePath(), ErrBlocked)
	}
	defer source.Close()

	destination, err := os.Create(dstFile.AbsolutePath())
	if err != nil {
		return nil, err
	}
	defer destination.Close()

	err = copy(destination, source)
	if err != nil {
		destination.Close()
		dstFile.Delete()

		return nil, err
	}

	return dstFile, nil
}

func copy(dst io.Writer, src io.Reader) error {
	bufferSize := 100 << 20 //100Mb
	buf := make([]byte, bufferSize)

	if _, err := io.CopyBuffer(dst, src, buf); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrCopy)
	}

	return nil
}

//MoveTo перемещает файл в новое расположение
func (f *File) MoveTo(path string) error {
	targetFile, err := f.CopyTo(path)
	if err != nil {
		return err
	}

	originalFile := *f
	f.PathName = targetFile.AbsolutePath()

	return originalFile.Delete()
}

func (f *File) validate() error {
	pathName := f.AbsolutePath()
	if exists := isExists(pathName); !exists {
		return fmt.Errorf("file '%s' is not exists: %w", pathName, ErrNotExists)
	}
	if regular, _ := isRegular(pathName); !regular {
		return fmt.Errorf("'%s' is not a regular file: %w", pathName, ErrNotRegular)
	}

	return nil
}
