package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileInfo(t *testing.T) {

	dir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "testfile")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	entry, err := os.ReadDir(dir)
	assert.NoError(t, err)

	fileInfo, err := getFileInfo(dir+"/"+entry[0].Name(), entry[0])
	assert.NoError(t, err)
	assert.Equal(t, entry[0].Name(), fileInfo.Name)
	assert.Equal(t, "file", fileInfo.Type)
	assert.False(t, fileInfo.IsSymlink)
}

func TestListRootFiles(t *testing.T) {

	dir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	file, err := os.CreateTemp(dir, "testfile")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	jsonData, err := listRootFiles(dir, "")
	assert.NoError(t, err)
	assert.Contains(t, jsonData, "testfile")
}

func TestAddTrailingSlash(t *testing.T) {
	assert.Equal(t, "/path/", AddTrailingSlash("/path"))
	assert.Equal(t, "/path/", AddTrailingSlash("/path/"))
}

func TestAddSlash(t *testing.T) {
	assert.Equal(t, "/path", AddSlash("path"))
	assert.Equal(t, "/path", AddSlash("/path"))
}

func TestRemoveDoubleSlashes(t *testing.T) {
	assert.Equal(t, "/path/to/file", RemoveDoubleSlashes("/path//to///file"))
}
