package fs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isExists(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	assert.True(t, isExists(file.Name()))
	assert.False(t, isExists(file.Name()+"_non_existent"))
}

func Test_isRegular(t *testing.T) {
	dirName, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dirName)

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	regular, err := isRegular(file.Name())
	assert.True(t, regular)
	assert.Nil(t, err)

	regular, err = isRegular(dirName)
	assert.False(t, regular)
	assert.Nil(t, err)

	regular, err = isRegular(file.Name() + "_non_existent")
	assert.False(t, regular)
	assert.NotNil(t, err)
}

func Test_isDirectory(t *testing.T) {
	dirName, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dirName)

	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	dir, err := isDirectory(dirName)
	assert.True(t, dir)
	assert.Nil(t, err)

	dir, err = isDirectory(file.Name())
	assert.False(t, dir)
	assert.Nil(t, err)

	dir, err = isDirectory(dirName + "_non_existent")
	assert.False(t, dir)
	assert.NotNil(t, err)
}

func Test_isFileLocked(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	assert.True(t, isFileLocked(file.Name()))

	file.Close()

	assert.False(t, isFileLocked(file.Name()))
}
