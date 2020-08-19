package fs

import (
	"os"
	"time"

	"gopkg.in/djherbis/times.v1"
)

func isExists(pathName string) bool {
	_, err := os.Stat(pathName)
	if err != nil {
		return false
	}

	return true
}

func isRegular(pathName string) (bool, error) {
	fileStat, err := os.Stat(pathName)
	if err != nil {
		return false, err
	}

	if isRegular := fileStat.Mode().IsRegular(); !isRegular {
		return false, nil
	}

	return true, nil
}

func isDirectory(pathName string) (bool, error) {
	fileStat, err := os.Stat(pathName)
	if err != nil {
		return false, err
	}

	if isDir := fileStat.Mode().IsDir(); !isDir {
		return false, nil
	}

	return true, nil
}

func fileTimes(pathName string) (atime time.Time, mtime time.Time, err error) {
	fstat, err := times.Stat(pathName)
	if err != nil {
		return
	}

	atime = fstat.AccessTime()
	mtime = fstat.ModTime()

	return
}

func setFileTimes(pathName string, atime, mtime time.Time) {
	_ = os.Chtimes(pathName, atime, mtime)
}

func isFileLocked(pathName string) bool {
	file, err := os.Open(pathName)
	if err != nil {
		return true
	}
	file.Close()

	if err := os.Rename(pathName, pathName); err != nil {
		return true
	}

	return false
}
